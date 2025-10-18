package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/puspo/basicnewproject/internal/models"
	"github.com/puspo/basicnewproject/internal/services"
)

// PostHandlers groups HTTP handlers for posts.
type PostHandlers struct {
	service *services.PostService
}

// NewPostHandlers creates a new PostHandlers instance.
func NewPostHandlers(service *services.PostService) *PostHandlers {
	return &PostHandlers{
		service: service,
	}
}

// Register wires the handlers onto a ServeMux.
func (h *PostHandlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/posts", h.handlePosts)
	mux.HandleFunc("/posts/", h.handlePostByID)
}

// handleHealth handles the health check endpoint.
func (h *PostHandlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// handlePosts services /posts for GET (list) and POST (create).
func (h *PostHandlers) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listPosts(w, r)
	case http.MethodPost:
		h.createPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handlePostByID services /posts/{id} for GET (get), PUT (update), DELETE (delete).
func (h *PostHandlers) handlePostByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/posts/"), "/")
	idStr := parts[0]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getPost(w, r, id)
	case http.MethodPut:
		h.updatePost(w, r, id)
	case http.MethodDelete:
		h.deletePost(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PostHandlers) createPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	p, err := decodePostFromBody(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := h.service.Create(ctx, p)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *PostHandlers) listPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Basic defaults; you could parse query params here.
	posts, err := h.service.List(ctx, 0, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Ensure we never return null; JSON null would break some clients.
	if posts == nil {
		posts = []models.Post{}
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *PostHandlers) getPost(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	post, err := h.service.Get(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if post == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandlers) updatePost(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	p, err := decodePostFromBody(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	ok, err := h.service.Update(ctx, id, p)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"updated": true})
}

func (h *PostHandlers) deletePost(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	ok, err := h.service.Delete(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

// decodePostFromBody parses JSON request body into a Post struct.
func decodePostFromBody(r *http.Request) (*models.Post, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB limit
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var p models.Post
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// writeJSON is a small helper to write JSON responses consistently.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
