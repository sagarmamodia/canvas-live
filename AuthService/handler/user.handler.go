// package handler

// import (
// 	"auth-service/repository"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"
// )

// type UserDto struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// }

// type UserHandler struct {
// 	UserRepository *repository.UserRepository
// }

// func (h UserHandler) RetrieveSearchedUsers(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Only GET request allowed", http.StatusBadRequest)
// 		return
// 	}

// 	params := r.URL.Query()
// 	q := params.Get("q")
// 	if q == "" {
// 		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
// 		defer cancel()
// 		users, err := h.UserRepository.FindAll(ctx)
// 		if err != nil {
// 			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
// 			return
// 		}

// 		userDtos := []UserDto{}
// 		for _, user := range users {
// 			userDto := UserDto{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}
// 			userDtos = append(userDtos, userDto)
// 		}

// 		json.NewEncoder(w).Encode(userDtos)
// 	}
// }
package handler

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func (h UserHandler) RetrieveSearchedUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Validations
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request allowed", http.StatusBadRequest)
		return
	}

	// 2. Set JSON Header
	w.Header().Set("Content-Type", "application/json")

	// 3. Setup Context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 4. Get Query Params
	params := r.URL.Query()
	q := params.Get("q")

	var users []model.User
	var err error

	// 5. Logic Branching
	if q == "" {
		// If query is empty, get all users
		users, err = h.UserRepository.FindAll(ctx)
	} else {
		// If query exists, search for them
		users, err = h.UserRepository.FindByQuery(ctx, q)
	}

	// 6. Error Handling
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	// 7. Convert to DTOs
	userDtos := []UserDto{}
	
	// Ensure users is not nil before looping
	if users != nil {
		for _, user := range users {
			userDto := UserDto{
				ID:       user.ID.Hex(),
				Username: user.Username,
				Email:    user.Email,
			}
			userDtos = append(userDtos, userDto)
		}
	}

	// 8. Send Response (Handles [] case automatically)
	json.NewEncoder(w).Encode(userDtos)
}