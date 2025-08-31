package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/edipretoro/boot.dev/go_web_server/internal/auth"
	"github.com/edipretoro/boot.dev/go_web_server/internal/database"
)

func sanitizeChirp(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range strings.Fields(chirp) {
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirp = strings.ReplaceAll(chirp, word, "****")
		}
	}
	return chirp
}

func apiHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)("OK"))
}

func validateSizeChirp(chirp string) bool {
	return len(chirp) <= 140
}

// func validateChirp(w http.ResponseWriter, req *http.Request) {
// 	type chirpBody struct {
// 		Chirp string `json:"body"`
// 	}
// 	type chirpResponse struct {
// 		ErrorMessage string `json:"error,omitempty"`
// 		CleanedBody  string `json:"cleaned_body,omitempty"`
// 		// extra        string `json:"extra,omitempty"`
// 	}
// 	decoder := json.NewDecoder(req.Body)
// 	chirp := chirpBody{}
// 	response := chirpResponse{}
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	if err := decoder.Decode(&chirp); err != nil {
// 		response.ErrorMessage = "Error decoding chirp body"
// 		dat, err := json.Marshal(response)
// 		if err != nil {
// 			log.Printf("Error marshalling response: %v", err)
// 		}
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(dat)
// 		return
// 	}

// 	if !validateSizeChirp(chirp.Chirp) {
// 		response.ErrorMessage = "Chirp is too long"
// 		w.WriteHeader(http.StatusBadRequest)
// 	} else {
// 		response.ErrorMessage = ""
// 		response.CleanedBody = sanitizeChirp(chirp.Chirp)
// 	}
// 	dat, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("Error marshalling response: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(dat)
// }

func addUser(w http.ResponseWriter, req *http.Request) {
	type userBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type userResponse struct {
		ID           uuid.UUID `json:"id,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
		Email        string    `json:"email,omitempty"`
		ErrorMessage string    `json:"error,omitempty"`
	}
	decoder := json.NewDecoder(req.Body)
	user := userBody{}
	response := userResponse{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := decoder.Decode(&user); err != nil {
		response.ErrorMessage = "Error decoding user body"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		response.ErrorMessage = "Error processing password"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}
	newUser, err := apiCfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		response.ErrorMessage = "Error creating user"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.ErrorMessage = ""
		response.ID = newUser.ID
		response.CreatedAt = newUser.CreatedAt
		response.UpdatedAt = newUser.UpdatedAt
		response.Email = newUser.Email

	}
	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func addChirp(w http.ResponseWriter, req *http.Request) {
	type chirpBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type chirpResponse struct {
		ChirpID      uuid.UUID `json:"chirp_id,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
		Body         string    `json:"body,omitempty"`
		UserID       uuid.UUID `json:"user_id,omitempty"`
		ErrorMessage string    `json:"error,omitempty"`
	}
	decoder := json.NewDecoder(req.Body)
	chirp := chirpBody{}
	response := chirpResponse{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := decoder.Decode(&chirp); err != nil {
		response.ErrorMessage = "Error decoding chirp body"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	if !validateSizeChirp(chirp.Body) {
		response.ErrorMessage = "Chirp is too long"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}
	chirp.Body = sanitizeChirp(chirp.Body)

	newChirp, err := apiCfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   chirp.Body,
		UserID: chirp.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		response.ErrorMessage = "Error creating chirp"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.ErrorMessage = ""
		response.ChirpID = newChirp.ID
		response.CreatedAt = newChirp.CreatedAt
		response.UpdatedAt = newChirp.UpdatedAt
		response.Body = newChirp.Body
		response.UserID = newChirp.UserID
		w.WriteHeader(http.StatusCreated)
	}
	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
	}
	w.Write(dat)
}

func getAllChirps(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	type chirpsResponse struct {
		ErrorMessage string `json:"error,omitempty"`
	}
	response := chirpsResponse{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	chirpsDb, err := apiCfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		log.Printf("Error getting all chirps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.Write(dat)
		return
	}
	chirps := make([]chirp, 0, len(chirpsDb))
	for _, c := range chirpsDb {
		chirps = append(chirps, chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling chirps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func getChirpByID(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing chirp ID (%s): %v", req.PathValue("chirpID"), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirp, err := apiCfg.dbQueries.GetChirpByID(req.Context(), userID)
	if err != nil {
		log.Printf("Error getting chirp by ID: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func loginUser(w http.ResponseWriter, req *http.Request) {
	type userBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type userResponse struct {
		ID           uuid.UUID `json:"id,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
		UpdatedAt    time.Time `json:"updated_at,omitempty"`
		Email        string    `json:"email,omitempty"`
		ErrorMessage string    `json:"error,omitempty"`
	}
	decoder := json.NewDecoder(req.Body)
	user := userBody{}
	response := userResponse{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := decoder.Decode(&user); err != nil {
		response.ErrorMessage = "Error decoding user body"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	existingUser, err := apiCfg.dbQueries.GetUserByEmail(req.Context(), user.Email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		response.ErrorMessage = "Error getting user"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	if auth.CheckPasswordHash(user.Password, existingUser.HashedPassword) != nil {
		response.ErrorMessage = "Invalid password"
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(dat)
		return
	}

	response.ErrorMessage = ""
	response.ID = existingUser.ID
	response.CreatedAt = existingUser.CreatedAt
	response.UpdatedAt = existingUser.UpdatedAt
	response.Email = existingUser.Email

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
