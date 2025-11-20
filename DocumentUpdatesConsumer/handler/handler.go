package handler

import (
	"DocumentUpdatesConsumer/model"
	"DocumentUpdatesConsumer/repository"
	"DocumentUpdatesConsumer/types"
	"context"
	"encoding/json"
	"fmt"
)

func DocumentUpdatesHandler(ctx context.Context, r *repository.DocumentRepository, msg types.Message) {

	var actionMsg map[string]interface{}
	err := json.Unmarshal([]byte(msg.Body), &actionMsg)
	if err != nil {
		fmt.Printf("[DocumentUpdatesHandler] error unmarshalling message body")
		return
	}

	// fmt.Printf("\n ============ Action Msg ============= \n %v\n", actionMsg)

	actVal := actionMsg["action"].(string) // it is always possible as only validated data is pushed to kafka
	if actVal == "add_slide" {
		fmt.Printf("[DocumentUpdatesHandler] AddSlide message received by consumer")
		slideId, ok := actionMsg["slideId"].(string)
		if !ok {
			fmt.Printf("[DocumentUpdatesHandler] slideId missing")
			return
		}

		err := r.AddNewSlide(ctx, msg.DocumentID, slideId)
		if err != nil {
			fmt.Printf("[DocumentUpdatesHandler] Error adding new slide")
			return
		}

	} else if actVal == "remove_slide" {
		fmt.Printf("[DocumentUpdatesHandler] RemoveSlide message received by consumer")
		slideId, ok := actionMsg["slideId"].(string)
		if !ok {
			fmt.Printf("[DocumentUpdatesHandler] slideId missing")
			return
		}

		err := r.RemoveSlide(ctx, msg.DocumentID, slideId)
		if err != nil {
			fmt.Printf("[DocumentUpdatesHandler] Error adding new slide")
			return
		}

	} else if actVal == "delete" {
		fmt.Printf("[DocumentUpdatesHandler] Delete message received by consumer")
		// msg contains the docId; the actionMsg must contain slideId and objectId
		docId := msg.DocumentID
		slideId := actionMsg["slideId"].(string)
		objectId := actionMsg["objectId"].(string)
		err := r.DeleteElement(ctx, docId, slideId, objectId)
		if err != nil {
			fmt.Printf("[DocumentUpdatesHandler] Error deleting object")
			return
		}

	} else if actVal == "update" {
		fmt.Printf("[DocumentUpdatesHandler] Update message received by consumer")
		// msg contains the docId; the actionMsg must contain slideId and objectId
		docId := msg.DocumentID
		slideId := actionMsg["slideId"].(string)
		objectId := actionMsg["objectId"].(string)

		// updated fields actionMsg["updatedAttributes"] is of type interface it need to be converted to map[string]interface
		updatedFields, ok := actionMsg["updatedAttributes"].(map[string]interface{})
		if !ok {
			fmt.Printf("[DocumentUpdatesHandler] Error converting updatedAttributes to map[string]interface{}: %s\n", err)
			return
		}

		err := r.UpdateElement(ctx, docId, slideId, objectId, updatedFields)
		if err != nil {
			fmt.Printf("[DocumentUpdatesHandler] Error updating object: %s\n", err)
			return
		}

	} else if actVal == "create" {
		fmt.Printf("[DocumentUpdatesHandler] Create message received by consumer")
		// msg contains the docId; the actionMsg must contain slideId and objectId
		docId := msg.DocumentID
		slideId := actionMsg["slideId"].(string)
		objectId := actionMsg["objectId"].(string)
		objectType := actionMsg["objectType"].(string)

		// updated fields actionMsg["updatedAttributes"] is of type interface it need to be converted to map[string]interface
		attr, ok := actionMsg["attributes"].(map[string]interface{})
		if !ok {
			fmt.Printf("[DocumentUpdatesHandler] Error converting updatedAttributes to map[string]interface{}:- %s\n", err)
			return
		}

		// create model.Object
		obj := model.Object{
			ID:         objectId,
			Type:       objectType,
			Attributes: attr,
		}

		err := r.CreateElement(ctx, docId, slideId, obj)
		if err != nil {
			fmt.Printf("[DocumentUpdatesHandler] Error creating object:- %s\n", err)
			return
		}
	} else {
		fmt.Printf("[DocumentUpdatesHandler] Unknown message received by consumer")
	}
}
