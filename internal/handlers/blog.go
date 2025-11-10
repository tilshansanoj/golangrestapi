package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/tilshansanoj/golangrestapi/internal/dto"
	"github.com/tilshansanoj/golangrestapi/internal/errorhandler"
	"github.com/tilshansanoj/golangrestapi/internal/store"
	"github.com/tilshansanoj/golangrestapi/internal/utils"
)

func (h *Handler) CreateBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for creating a blog post
		ctx := r.Context()

		var req dto.CreateBlogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			slog.Error("Failed to decode create blog request", "error", err)
			return
		}

		_, err := h.Queries.CreateBlog(ctx, store.CreateBlogParams{
			Title:   req.Title,
			Content: req.Content,
			UserID:  sql.NullInt32{Int32: req.UserID, Valid: true},
		})
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to create blog post")
			slog.Error("Failed to create blog post", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusCreated, "Blog post created successfully", req.Title)
		slog.Info("Blog post created successfully", "title", req.Title, "user_id", req.UserID)
	}
}

func (h *Handler) ListBlogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for retrieving blog posts
		ctx := r.Context()
		
		blogs, err := h.Queries.ListBlogs(ctx)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to retrieve blog posts")
			slog.Error("Failed to retrieve blog posts", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "Blog posts retrieved successfully", blogs)
		slog.Info("Blog posts retrieved successfully", "blogs", blogs)
	}
}

func (h *Handler) GetBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		
		BlogID := r.PathValue("id")
		
		if BlogID == "" {
			utils.ResponseWithError(w, http.StatusBadRequest, "Missing Blog ID")
			slog.Error("Failed to retrieve blog post", "error", "Missing Blog ID")
			return
		}

		ID, err := strconv.Atoi(BlogID)
		if err !=  nil {
			slog.Error("Failed to parse BlogID into Integer", "error", err)
			utils.ResponseWithError(w, http.StatusBadRequest, "Faled to parse BlogID into Integer")
		}

		// Call service to get blog by ID
		blog, err := h.Queries.GetBlogbyId(ctx, int32(ID))
		if err != nil {
			slog.Error("Failed to get blog", "error", err)
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to get blog")
			return
		}

		if blog.ID == 0 {
			utils.ResponseWithError(w, http.StatusNotFound, "Blog not found")
			return
		}

		slog.Info("Blog retrieved:", "blog", blog)
		utils.ResponseWithSuccess(w, http.StatusOK, " Blog Retrieved:", blog)
	}
}

func (h *Handler)UpdateBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()

		BlogID := r.PathValue("id")
		if BlogID == "" {
			utils.ResponseWithError(w, http.StatusBadRequest, "Missing Blog ID")
			slog.Error("Failed to retrieve blog post", "error", "Missing Blog ID")
			return
		}

		ID, _ := strconv.Atoi(BlogID)
		if BlogID ==  "" {
			slog.Error("Failed to parse BlogID into Integer")
			utils.ResponseWithError(w, http.StatusBadRequest, "Faled to parse BlogID into Integer")
		}
		var req dto.UpdateBlogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Invalid request Payload")
			slog.Error("Invalid request payload", "error", err)
			return 
		}

		// Start a trasaction
		tx, err := h.DB.BeginTx(ctx,nil)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to start transaction")
			slog.Error("Failed to start transaction", "error", err)
			return
		}
		defer tx.Rollback()

		_, err = h.Queries.UpdateBlog(ctx, store.UpdateBlogParams{
			ID : int32(ID),
			Content: req.Content,
			Title: req.Title,
		})
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to update user")
			slog.Error("Failed to update blog", "error", err)
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
			slog.Error("Failed to commit transaction", "error", err)
			return
		}
		utils.ResponseWithSuccess(w, http.StatusCreated, "Blog updated successfully", req)
		slog.Info("Blog updated successfully", "data", req)

	}
}

func (h *Handler)DeleteBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		BlogID := r.PathValue("id")

		ID, err := strconv.Atoi(BlogID)
		if BlogID ==  "" {
			slog.Error("Failed to parse BlogID into Integer", "error", err)
			utils.ResponseWithError(w, http.StatusBadRequest, "Faled to parse BlogID into Integer")
		}
		
		// Start a trasaction
		tx, err := h.DB.BeginTx(ctx,nil)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to start transaction")
			slog.Error("Failed to start transaction", "error", err)
			return
		}
		defer tx.Rollback()

		_, err = h.Queries.DeleteBlog(ctx, int32(ID))
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to update user")
			slog.Error("Failed to delete user", "error", err)
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
			slog.Error("Failed to commit transaction", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "Blog deleted successfully", ID)
		slog.Info("Blog deleted successfully", "id", ID )
	}
}
 