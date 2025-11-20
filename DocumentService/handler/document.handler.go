package handler

import (
	"document-service/repository"
	"document-service/types"
	"encoding/json"
	"fmt"
	"net/http"
)

// ===========================================

type DocumentHandler struct {
	DocumentRepository *repository.DocumentRepository
}

// ====================== Get all documents handler =======================================

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

	result := types.AllDocumentsDto{OwnedDocuments: ownedDocuments, SharedDocuments: sharedDocuments}
	// Json response
	json.NewEncoder(w).Encode(result)

}

// ================================ Create New Empty Document Handler ===========================

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

	response := types.CreatedResponse{ID: createdDoc.ID.Hex()}

	json.NewEncoder(w).Encode(response)
}

// ================================= Share Document Handler ==============================

func (h DocumentHandler) ShareDocument(w http.ResponseWriter, r *http.Request) {
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

	// Decode data from request
	var data types.ShareDocumentPostData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Check if the user actually own the document
	isUserOwner, err := h.DocumentRepository.IsDocumentOwnedByUser(r.Context(), userId, data.DocumentID)
	if err != nil {
		http.Error(w, "Error verifying ownership of the document", http.StatusInternalServerError)
		return
	}

	if !isUserOwner {
		http.Error(w, "Only the owner can share documents with other users", http.StatusBadRequest)
		return
	}

	// Create sharing record
	_, err = h.DocumentRepository.CreateCollaborationRecord(r.Context(), data.CollaboratorUserID, data.DocumentID, data.AccessType)
	if err != nil {
		http.Error(w, "Error creating a collaboration record", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success")

}

// ================================= Delete Document Handler ==============================

func (h DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
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

	// Decode data from request
	var data types.DeleteDocumentPostData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Check if the user actually own the document
	isUserOwner, err := h.DocumentRepository.IsDocumentOwnedByUser(r.Context(), userId, data.DocumentID)
	if err != nil {
		http.Error(w, "Error verifying ownership of the document", http.StatusInternalServerError)
		return
	}

	if !isUserOwner {
		http.Error(w, "Only the owner can share documents with other users", http.StatusBadRequest)
		return
	}

	// Delete document
	err = h.DocumentRepository.DeleteDocument(r.Context(), data.DocumentID)
	if err != nil {
		http.Error(w, "Error deleting document", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Success")
}
