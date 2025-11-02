package repository

import (
	"context"
	"document-service/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentRepository struct {
	collection                *mongo.Collection
	sharedDocRecordCollection *mongo.Collection
}

func NewDocumentRepository(client *mongo.Client, database string, collection string, sharedDocCollectionName string) *DocumentRepository {
	coll := client.Database(database).Collection(collection)
	shared := client.Database(sharedDocCollectionName).Collection(sharedDocCollectionName)
	return &DocumentRepository{
		collection:                coll,
		sharedDocRecordCollection: shared,
	}
}

func (r *DocumentRepository) CreateNewDocument(ctx context.Context, title string, ownerId string) (model.Document, error) {

	// Create a Document
	emptyDocument := model.Document{
		Title:   title,
		OwnerID: ownerId,
		Slides: []model.Slide{
			{
				Background: "#FFFFFF",
				Objects:    []model.DocObject{},
			},
		},
	}

	// Insert Document
	result, err := r.collection.InsertOne(ctx, emptyDocument)
	if err != nil {
		return model.Document{}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		emptyDocument.ID = oid
	}

	return emptyDocument, nil
}

func (r *DocumentRepository) DeleteDocument(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("[DocumentRepository] Invalid document id: %v\n", err)
		return err
	}

	filter := bson.M{"_id": objectId}

	// Execute Deletion
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("[DocumentRepository] Error deleting document: %v\n", err)
		return err
	}

	if result.DeletedCount == 1 {
		fmt.Printf("[DocumentRepository] Successfully deleted 1 document with ID: %s\n", id)
	} else {
		fmt.Printf("[DocumentRepository] No document found with ID: %s\n", id)
	}

	return nil
}

func (r *DocumentRepository) FindOwnedDocuments(ctx context.Context, userId string) ([]model.Document, error) {

	filter := bson.M{"ownerId": userId}
	// Execute the query
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[DocumentRepository][FindOwnedDocuments] Error retrieving documents: %v\n", err)
		return []model.Document{}, err
	}
	defer cursor.Close(ctx)

	// Decode all Documents in documents slice
	documents := []model.Document{}
	if err = cursor.All(ctx, &documents); err != nil {
		fmt.Printf("[DocumentRepository][FindOwnedDocuments] Error decoding documents: %v\n", err)
		return []model.Document{}, err
	}

	return documents, nil
}

func (r *DocumentRepository) FindSharedDocuments(ctx context.Context, userId string) ([]model.Document, error) {

	filter := bson.M{"userId": userId}

	// Get IDs of documents shared with the current user
	cursor, err := r.sharedDocRecordCollection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[DocumentRepository][FindSharedDocuments] Error retrieving shared document records: %v\n", err)
		return []model.Document{}, err
	}
	defer cursor.Close(ctx)

	var sharedDocRecords []model.SharedDocRecord
	if err = cursor.All(ctx, &sharedDocRecords); err != nil {
		fmt.Printf("[DocumentRepository][FindSharedDocuments] Error decoding shared document records: %v\n", err)
		return []model.Document{}, err
	}

	var ids []primitive.ObjectID
	for _, record := range sharedDocRecords {
		objectId, err := primitive.ObjectIDFromHex(record.DocumentID)
		if err != nil {
			continue
		}
		ids = append(ids, objectId)
	}

	// Get documents
	// if ids is empty return empty slice
	if len(ids) == 0 {
		return []model.Document{}, nil
	}

	filter = bson.M{
		"_id": bson.M{"$in": ids},
	}

	cursor, err = r.collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[DocumentRepository][FindSharedDocuments] Error retrieving documents: %v\n", err)
		return []model.Document{}, err
	}
	defer cursor.Close(ctx)

	var documents []model.Document
	if err = cursor.All(ctx, &documents); err != nil {
		fmt.Printf("[DocumentRepository][FindSharedDocuments] Error decoding documents: %v\n", err)
		return []model.Document{}, nil
	}

	return documents, nil
}

func (r *DocumentRepository) CreateSharedDocRecord(ctx context.Context, collaboratorUserId string, documentId, accessType string) (model.SharedDocRecord, error) {

	// Create shared document record object
	sharedDocRecord := model.SharedDocRecord{
		UserID:     collaboratorUserId,
		DocumentID: documentId,
		AccessType: accessType,
	}

	// Execute the query
	result, err := r.sharedDocRecordCollection.InsertOne(ctx, sharedDocRecord)
	if err != nil {
		fmt.Printf("[DocumentRepository] Error creating sharing record: %v\n", err)
		return model.SharedDocRecord{}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		sharedDocRecord.ID = oid
	}

	return sharedDocRecord, nil
}
