package dtos

type ErrorResponse struct {
	Error string `json:"error" example:"An error occurred"`
}

type BadRequestResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

type UnauthorizedResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

type ForbiddenResponse struct {
	Error string `json:"error" example:"Access denied"`
}

type NotFoundResponse struct {
	Error string `json:"error" example:"Resource not found"`
}

type ConflictResponse struct {
	Error string `json:"error" example:"User already exists"`
}

type InternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
