package repository

import (
	"DocumentUpdatesConsumer/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentRepository struct {
	collection *mongo.Collection
}

func NewDocumentRepository(client *mongo.Client, database string, collection string) *DocumentRepository {
	coll := client.Database(database).Collection(collection)
	return &DocumentRepository{
		collection: coll,
	}
}

func (r *DocumentRepository) AddNewSlide(ctx context.Context, documentId string, slideId string) error {
	objectId, err := primitive.ObjectIDFromHex(documentId)
	if err != nil {
		fmt.Printf("[DocumentRepository] Invalid document id: %v\n", err)
		return err
	}

	// check if document exists or not
	filter := bson.M{"_id": objectId}
	var doc model.Document
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		fmt.Printf("[DocumentRepository][FindOwnedDocuments] Error decoding documents: %v\n", err)
		return err
	}

	// document exists
	// create new slide
	newSlide := model.Slide{
		ID:         slideId,
		Background: "#fff",
		Objects:    make([]model.Object, 0, 1),
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "slides", Value: newSlide},
		}},
	}

	// Execute the UpdateOne
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// 4. Check the result
	if result.ModifiedCount == 1 {
		fmt.Println("Successfully pushed new slide to the document list.")
	} else if result.MatchedCount == 0 {
		return fmt.Errorf("document not found with ID: %s", documentId)
	}

	return nil
}

func (r *DocumentRepository) RemoveSlide(ctx context.Context, docId string, slideId string) error {

	// --- 1. Top-Level FILTER: Find the Document ---
	docObjectID, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return fmt.Errorf("invalid Document ID format: %w", err)
	}
	docFilter := bson.M{"_id": docObjectID}

	// --- 2. Construct the $pull Update
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			// Key: The name of the array field to pull from ("slides")
			// Value: The query that identifies the element(s) to remove.
			{Key: "slides", Value: bson.M{"_id": slideId}},
		}},
	}

	// --- 3. Execute UpdateOne (No Array Filters Required) ---
	// We pass nil for the options since arrayFilters is not needed.
	result, err := r.collection.UpdateOne(
		ctx,
		docFilter,
		update,
		// options.Update() is optional here, as no complex options are used
	)

	if err != nil {
		return fmt.Errorf("[Repository][RemoveSlide] database update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("[Repository][RemoveSlide] Slide was not found or document ID is incorrect")
	}

	fmt.Printf("[Repository][RemoveSlide] Successfully deleted slide %s. Modified: %d\n", slideId, result.ModifiedCount)
	return nil
}

func (r *DocumentRepository) UpdateElement(ctx context.Context, docId string, slideId string, elementId string, updatedFields map[string]interface{}) error {

	// --- 1. Top-Level FILTER: Find the Document ---
	docObjectID, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return fmt.Errorf("invalid Document ID format: %w", err)
	}
	docFilter := bson.M{"_id": docObjectID}

	// --- 2. ARRAY FILTERS: Target the Slide and the Element ---
	// Array Filters are defined using a slice of BSON documents.
	arrayFilters := bson.A{
		// Filter 1 (for the Slides array): Find the slide that matches the slideID.
		// The identifier 'elem' can be used later in the $set path.
		bson.M{"elem._id": slideId},

		// Filter 2 (for the Objects array inside the matched slide): Find the element that matches the elementID.
		// The identifier 'obj' can be used later in the $set path.
		bson.M{"obj._id": elementId},
	}

	// --- 3. Construct the $SET Update ---
	// We use bson.D for the $set operator because order matters.
	// The $set value itself is built dynamically from the map[string]interface{}

	// Create the $set stage
	setStage := bson.D{}

	// CRITICAL STEP: Build the full path for the update
	// "slides.$[elem].objects.$[obj].<field>"
	// - $[elem]: Targets the slide found by Filter 1.
	// - objects.$[obj]: Targets the object found by Filter 2.

	for key, value := range updatedFields {
		fullPath := fmt.Sprintf("slides.$[elem].objects.$[obj].attributes.%s", key)
		setStage = append(setStage, bson.E{Key: fullPath, Value: value})
	}

	update := bson.D{
		{Key: "$set", Value: setStage},
	}

	// --- 4. Execute UpdateOne with Array Filters ---
	result, err := r.collection.UpdateOne(
		ctx,
		docFilter,
		update,
		options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters}),
	)

	if err != nil {
		return fmt.Errorf("[Repository][UpdateElement] database update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("[Repository][UpdateElement] no element was found or modified (IDs may be incorrect)")
	}

	fmt.Printf("[Repository][UpdateElement] Successfully updated 1 element. Matched: %d, Modified: %d\n",
		result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *DocumentRepository) CreateElement(ctx context.Context, docId string, slideId string, newElementData model.Object) error {
	docObjectId, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		fmt.Printf("[DocumentRepository][CreateElement] Invalid document id: %v\n", err)
		return err
	}

	// --- 1. Top-Level Filter: Find the Document ---
	// Match the main document by its ID.
	docFilter := bson.M{"_id": docObjectId}

	// --- 2. ARRAY FILTERS: Target the Slide ---
	// Define a filter to find the correct slide within the "slides" array.
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			// The identifier 'elem' will point to the matching slide sub-document.
			bson.M{"elem._id": slideId},
		},
	}

	// --- 3. Construct the $PUSH Update ---
	// We use bson.D for ordered operators.

	// CRITICAL PATH: slides.$[elem].elements
	// This path targets the 'elements' array inside the slide where elem._id matches slideID.
	updatePath := "slides.$[elem].objects"

	update := bson.D{
		{Key: "$push", Value: bson.D{
			// $push to the specific path defined by the positional filtered identifier '$[elem]'
			{Key: updatePath, Value: newElementData},
		}},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		docFilter,
		update,
		options.Update().SetArrayFilters(arrayFilters),
	)

	if err != nil {
		return fmt.Errorf("[Repository][CreateElement] database update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("[Repository][CreateElement] no element was created (IDs may be incorrect)")
	}

	fmt.Printf("[Repository][CreateElement] Successfully created 1 element. Matched: %d, Modified: %d\n",
		result.MatchedCount, result.ModifiedCount)

	return nil
}

func (r *DocumentRepository) DeleteElement(ctx context.Context, docId string, slideId string, elementId string) error {
	docObjectId, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		fmt.Printf("[DocumentRepository][CreateElement] Invalid document id: %v\n", err)
		return err
	}

	// --- 1. Top-Level Filter: Find the Document ---
	docFilter := bson.M{"_id": docObjectId}

	// --- 2. ARRAY FILTERS: Target the Slide ---
	// We use the identifier 'elem' to find the specific slide based on its ID.
	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			// Filter 1: Find the slide that matches the slideID.
			bson.M{"elem._id": slideId},
		},
	}

	// --- 3. Construct the $PULL Update ---
	// The $pull operator removes elements from an array that match a specified query.
	// We use the positional filtered identifier $[elem] to target the correct slide.

	// CRITICAL PATH: slides.$[elem].objects
	// This path targets the 'objects' array inside the slide where elem._id matches slideID.
	updatePath := "slides.$[elem].objects"

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			// $pull from the target array field (updatePath)
			{Key: updatePath, Value: bson.M{"_id": elementId}},
		}},
	}

	// --- 4. Execute UpdateOne with Array Filters ---
	result, err := r.collection.UpdateOne(
		ctx,
		docFilter,
		update,
		options.Update().SetArrayFilters(arrayFilters),
	)

	if err != nil {
		return fmt.Errorf("database $pull update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		// This means either the document, slide, or element wasn't found/deleted.
		return fmt.Errorf("element not found or deleted (Element ID: %s)", elementId)
	}

	fmt.Printf("Successfully deleted element %s from slide %s.\n", elementId, slideId)
	return nil
}
