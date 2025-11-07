package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required" min=3,max=20"`
	Email    string `json:"email" validate:"required" format="email"`
	Password string `json:"password" validate:"required" min=6,max=100"`
}



type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required" min=3,max=100"`
	Content string `json:"content" validate:"required" min=10"`
	UserID  int32  `json:"user_id"`
}

type GetUserProfileRequest struct {
	UserID int32 `json:"user_id" validate:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email string `json: "email,omitempty"`
}

type UpdateBlogRequest struct {
	Title string `json: "title,omitempty"`
	Content string `json: "content, omitempty`
}

