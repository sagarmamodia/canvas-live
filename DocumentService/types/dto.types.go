package types

import (
	"document-service/model"
)

// Dtos
type AllDocumentsDto struct {
	OwnedDocuments  []model.Document `json:"ownedDocuments"`
	SharedDocuments []model.Document `json:"sharedDocuments"`
}

type CreatedResponse struct {
	ID string `json:"id"`
}

type ShareDocumentPostData struct {
	CollaboratorUserID string `json:"collaboratorUserId"`
	DocumentID         string `json:"documentId"`
	AccessType         string `json:"accessType"`
}

type DeleteDocumentPostData struct {
	DocumentID string `json:"documentId"`
}
