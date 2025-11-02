package handler

import (
	"document-service/model"
	"document-service/repository"
	"encoding/json"
	"net/http"
)

// Dtos
type AllDocumentsDto struct {
	OwnedDocuments  []model.Document `json:"ownedDocuments"`
	SharedDocuments []model.Document `json:"sharedDocuments"`
}

type CreatedResponse struct {
	ID string `json:"id"`
}

type DocumentHandler struct {
	DocumentRepository *repository.DocumentRepository
}

func (h DocumentHandler) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request allowed", http.StatusBadRequest)
		return
	}

	// Retrieve user data
	userId := r.Header.Get("X-User-ID")
	if userId == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	// Get all owned documents
	ownedDocuments, err := h.DocumentRepository.FindOwnedDocuments(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error retrieving owned documents", http.StatusInternalServerError)
		return
	}

	// Get all shared documents
	sharedDocuments, err := h.DocumentRepository.FindSharedDocuments(r.Context(), userId)
	if err != nil {
		http.Error(w, "Error retrieving shared documents", http.StatusInternalServerError)
		return
	}

	result := AllDocumentsDto{OwnedDocuments: ownedDocuments, SharedDocuments: sharedDocuments}
	// Json response
	json.NewEncoder(w).Encode(result)

}

func (h DocumentHandler) CreateNewDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request allowed", http.StatusBadRequest)
		return
	}

	// Retrieve user data
	userId := r.Header.Get("X-User-ID")
	if userId == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	// Create document
	createdDoc, err := h.DocumentRepository.CreateNewDocument(r.Context(), "Untitled", userId)
	if err != nil {
		http.Error(w, "Error creating document", http.StatusInternalServerError)
		return
	}

	response := CreatedResponse{ID: createdDoc.ID.Hex()}

	json.NewEncoder(w).Encode(response)
}
