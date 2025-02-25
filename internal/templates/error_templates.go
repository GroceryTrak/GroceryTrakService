package templates

// General error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"An error occurred"`
}

// 400 Bad Request error
type BadRequestResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

// 401 Unauthorized error
type UnauthorizedResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

// 403 Forbidden error
type ForbiddenResponse struct {
	Error string `json:"error" example:"Access denied"`
}

// 404 Not Found error
type NotFoundResponse struct {
	Error string `json:"error" example:"Resource not found"`
}

// 409 Conflict error
type ConflictResponse struct {
	Error string `json:"error" example:"User already exists"`
}

// 500 Internal Server Error
type InternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
