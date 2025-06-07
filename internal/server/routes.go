package server

import (
	"encoding/json"
	"log"
	"net/http"
	"url_razor/internal/shortener"
	"url_razor/internal/store"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/create_short_url", s.CreateShortUrlHandler)
	mux.HandleFunc("/{shortUrl}", s.RedirectHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request UrlCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.LongUrl == "" || request.UserId == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	shortUrl := shortener.GenerateShortLink(request.LongUrl, request.UserId)
	store.SaveUrlMapping(shortUrl, request.LongUrl, request.UserId)

	response := map[string]string{"short_url": shortUrl}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("shortUrl")
	originalUrl := store.RetrieveInitialUrl(shortUrl)

	if originalUrl == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusFound)
}
