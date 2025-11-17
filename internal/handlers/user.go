package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tilshansanoj/golangrestapi/internal/auth"
	"github.com/tilshansanoj/golangrestapi/internal/dto"
	"github.com/tilshansanoj/golangrestapi/internal/errorhandler"
	"github.com/tilshansanoj/golangrestapi/internal/middlewares"
	"github.com/tilshansanoj/golangrestapi/internal/store"
	"github.com/tilshansanoj/golangrestapi/internal/utils"
	"github.com/tilshansanoj/golangrestapi/internal/validation"
)

// CreateUser godoc
// @Summary Create an user
// @Description Create a new user by providing username, email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Create User Request"
// @Success 201 {string} string
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 409 {object} errorhandler.ErrorResponse
// @Failure 500 {object} errorhandler.ErrorResponse
// @Router /users/register [post]
func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for creating a user
		ctx := r.Context()
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			slog.Error("Invalid request payload", "error", err)
			return
		}

		// Validate the request
		if err := validation.Validate(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
			slog.Error("Validation error", "error", err)
			return
		}

		// Start a transaction
		tx, err := h.DB.BeginTx(ctx, nil)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to start transaction")
			slog.Error("Failed to start transaction", "error", err)
			return
		}
		defer tx.Rollback()

		// Check if username or email already exists
		_, err = h.Queries.GetUserByUsername(ctx, req.Username)
		if err == nil {
			errorhandler.ResponseWithError(w, http.StatusConflict, "Username already exists")
			slog.Error("Username already exists", "username", req.Username)
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to hash password")
			slog.Error("Failed to hash password", "error", err)
			return
		}
		
		_, err = h.Queries.CreateUser(ctx, store.CreateUserParams{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		})
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to create user")
			slog.Error("Failed to create user", "error", err)
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
			slog.Error("Failed to commit transaction", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusCreated, "User created successfully", req.Username)
		slog.Info("User created successfully", "username", req.Username, "email", req.Email)
	}
}

func (h *Handler) GetUserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for retrieving user profile
		ctx := r.Context()

		var req dto.GetUserProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			slog.Error("Invalid request payload", "error", err)
			return
		}

		//validate the request
		if err := validation.Validate(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
			slog.Error("Validation error", "error", err)
			return
		}

		user, err := h.Queries.GetUserByID(ctx, req.UserID)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to retrieve user profile")
			slog.Error("Failed to retrieve user profile", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "User profile retrieved successfully", user)
		slog.Info("User profile retrieved successfully", "user", user)
	}
}

// GetUsers godoc
// @Summary Get users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 404 {object} errorhandler.ErrorResponse
// @Router /users [get]
func (h *Handler) GetAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for retrieving all users
		ctx := r.Context()

		users, err := h.Queries.ListUsers(ctx)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
			slog.Error("Failed to retrieve users", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "Users retrieved successfully", users)
		slog.Info("Users retrieved successfully", "users", users)
	}
}

// GetUserToken godoc
// @Summary Get user token
// @Description Get user token
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.LoginUserRequest true "Login User Request"
// @Success 200 {string} utils.ResponseWithSuccess
// @Failure 400 {object} errorhandler.ErrorResponse
// @Router /users/login [post]
func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for user login
		ctx := r.Context()
		var req dto.LoginUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			slog.Error("Invalid request payload", "error", err)
			return
		}

		//validate the request
		if err := validation.Validate(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
			slog.Error("Validation error", "error", err)
			return
		}

		// Retrieve user by username
		user, err := h.Queries.GetUserByUsername(ctx, req.Username)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Invalid credentials")
			slog.Error("Invalid credentials", "error", err)
			return
		}

		// Verify password
		if !utils.CheckPasswordHash(user.Password, req.Password) {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Invalid credentials")
			slog.Error("Invalid credentials for user", "username", req.Username)
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, err := auth.GenerateJWT(int64(user.ID), user.Username, user.Email, jwtKey)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token")
			slog.Error("Failed to generate token", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "Login successful", map[string]string{"token": token})
		slog.Info("User logged in successfully", "username", req.Username)
	}
}

// GetUserProfile godoc
// @Summary Get user by JWT token
// @Description Get user by providing JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} models.User
// @Failure 400 {object} errorhandler.ErrorResponse
// @Router /users/me [get]
func (h *Handler) GetUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler logic for retrieving users
		ctx := r.Context()

		claims, ok := r.Context().Value(middlewares.UserClaimKey).(*auth.Claims)
		if !ok {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized access, Login to continue")
			slog.Error("Unauthorized access attempt, Login to continue")
			return
		}
		
		userID := claims.UserID

		// check the redis instance first
		cacheKey := fmt.Sprintf("user: %d", userID)
		if cached, err := h.Redis.Get(ctx, cacheKey).Result(); err == nil{
			var user store.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil{
				utils.ResponseWithSuccess(w, http.StatusOK, "Success (from cache/redis)", user)
				return 
			}
		}

		// Fallback to the Database
		user, err := h.Queries.GetUserByID(ctx, int32(userID))
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusNotFound, "User not found")
			slog.Error("Failed to retrieve user", "error", err)
			return
		}

		// Set to redis
		userJson, _ := json.Marshal(user)
		h.Redis.Set(r.Context(), cacheKey, userJson, 5*time.Minute)

		utils.ResponseWithSuccess(w, http.StatusOK, "User retrieved successfully", user)
		slog.Info("User retrieved successfully", "user", user)
	}
}

// UpdateUser godoc
// @Summary Update an user
// @Description Update an existing user by providing username and email
// @Tags users
// @Param id path int true "User ID"
// @Accept json
// @Produce json
// @Param user body dto.UpdateUserRequest true "Update User Request"
// @Success 201 {object} models.User
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 500 {object} errorhandler.ErrorResponse
// @Router /users/{id} [patch]
func (h *Handler)UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()
		
		ID := r.PathValue("id")
		userID, _ := strconv.Atoi(ID)
		if ID == "" {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Missing User ID")
			slog.Error("Missing User ID")
		}

		var req dto.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			slog.Error("Invalid request payload", "error", err)
			return
		}

		// Start a transaction
		tx, err := h.DB.BeginTx(ctx, nil)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to start transaction")
			slog.Error("Failed to start transaction", "error", err)
			return
		}
		defer tx.Rollback()

		_, err = h.Queries.UpdateUser(ctx, store.UpdateUserParams{
			ID: int32(userID),
			Username: req.Username,
			Email: req.Email,
		})
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to update user")
			slog.Error("Failed to create user", "error", err)
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
			slog.Error("Failed to commit transaction", "error", err)
			return
		}

		utils.ResponseWithSuccess(w, http.StatusCreated, "User updated successfully", req)
		slog.Info("User updated successfully", "username", req.Username, "email", req.Email)

	}
}

// DeleteUser godoc
// @Summary Delete an user
// @Description Delete an existing user by providing user ID
// @Tags users
// @Param id path int true "User ID"
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 500 {object} errorhandler.ErrorResponse
// @Router /users/{id} [delete]
func (h *Handler)DeleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()

		UserID := r.PathValue("id")

		if UserID == "" {
			utils.ResponseWithError(w, http.StatusBadRequest, "Missing User ID")
			slog.Error("Failed to retrieve user", "error", "Missing User ID")
			return
		}

		id, err := strconv.Atoi(UserID)
		if err	!= nil {
			slog.Error("Failed to parse UserID into Integer", "error", err)
			utils.ResponseWithError(w, http.StatusBadRequest, "Faled to parse UserID into Integer")
		}

		// Start a transaction
		tx, err := h.DB.BeginTx(ctx, nil)
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to start transaction")
			slog.Error("Failed to start transaction", "error", err)
			return
		}
		defer tx.Rollback()

		_, err = h.Queries.DeleteUser(ctx, int32(id))
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

		utils.ResponseWithSuccess(w, http.StatusOK, "User deleted successfully", id)
		slog.Info("User deleted successfully", "id", id )
		
	}
}

// logout user handler
func (h *Handler) LogoutUserHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		// extract the jwt claims from the context
		ctx := r.Context()

		claims, ok := r.Context().Value(middlewares.UserClaimKey).(*auth.Claims)
		if !ok {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized access, Login to continue")
			slog.Error("Unauthorized access attempt, Login to continue")
			return
		}

		// extract the token from the authorization header
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			errorhandler.ResponseWithError(w, http.StatusBadRequest, "Missing token in Authorization header")
			slog.Error("Missing token in Authorization header")
			return
		}

		// conver the expiration time.Time
		expirationTime := time.Unix(claims.ExpiresAt, 0)
		now := time.Now()
		ttl := expirationTime.Sub(now)

		if ttl <= 0 {
			ttl = time.Minute * 5 // fallback ttl
		}

		//blacklist the token in redis
		err := h.Redis.Set(ctx, tokenString, "blacklisted", ttl).Err()
		if err != nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Failed to logout user")
			slog.Error("Failed to blacklist token in redis", "error", err)
			return
		}

		//clear any cached user data in redis
		userIDstr := fmt.Sprintf("user: %d", claims.UserID)
		if err := h.clearUserSession(userIDstr); err != nil {
			slog.Error("Failed to clear user session from redis", "error", err)
		}
		utils.ResponseWithSuccess(w, http.StatusOK, "User logged out successfully", nil)
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")	
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func (h *Handler) clearUserSession(userID string) error {
	pattern := fmt.Sprintf("session:%s:*", userID)

	// Background context for Redis operations
	ctx := context.Background()

	// Use Redis SCAN to find keys matching the pattern
	iter := h.Redis.Scan(ctx, 0, pattern, 0).Iterator()

	// Delete each matching key
	for iter.Next(ctx) {
		if err := h.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			slog.Error("Failed to delete user session from redis", "error", err)
		}
	}
	return iter.Err()
}
