package dtos

type RegisterRequest struct {
	Username string `json:"username" example:"abc@gmail.com"`
	Password string `json:"password" example:"Password@123"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully."`
}

type LoginRequest struct {
	Username string `json:"username" example:"abc@gmail.com"`
	Password string `json:"password" example:"Password@123"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR..."`
	IdToken      string `json:"idToken" example:"eyJhbGciOiJIUzI1NiIsInR..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR..."`
	ExpiresIn    int32  `json:"expiresIn" example:"3600"`
	TokenType    string `json:"tokenType" example:"Bearer"`
}

type ConfirmRequest struct {
	Username string `json:"username" example:"abc@gmail.com"`
	Code     string `json:"code" example:"123456"`
}

type ConfirmResponse struct {
	Message string `json:"message" example:"User confirmed successfully."`
}

type ResendRequest struct {
	Username string `json:"username" example:"abc@gmail.com"`
}

type ResendResponse struct {
	Message string `json:"message" example:"Confirmation code resent successfully."`
}
