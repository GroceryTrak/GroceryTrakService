package repository

import (
	"context"

	"github.com/GroceryTrak/GroceryTrakService/internal/clients"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AuthRepository interface {
	RegisterUser(req dtos.RegisterRequest, role string) (dtos.RegisterResponse, error)
	LoginUser(req dtos.LoginRequest) (dtos.LoginResponse, error)
	ConfirmUser(req dtos.ConfirmRequest) (dtos.ConfirmResponse, error)
	ResendCode(req dtos.ResendRequest) (dtos.ResendResponse, error)
}

type AuthRepositoryImpl struct{}

func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}

func (r *AuthRepositoryImpl) RegisterUser(req dtos.RegisterRequest, role string) (dtos.RegisterResponse, error) {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: &clients.ClientID,
		Username: &req.Username,
		Password: &req.Password,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Username),
			},
		},
	}

	_, err := clients.CognitoClient.SignUp(context.TODO(), input)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	return dtos.RegisterResponse{Message: "User registered successfully. Please confirm your email."}, nil
}

func (r *AuthRepositoryImpl) LoginUser(req dtos.LoginRequest) (dtos.LoginResponse, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &clients.ClientID,
		AuthParameters: map[string]string{
			"USERNAME": req.Username,
			"PASSWORD": req.Password,
		},
	}

	resp, err := clients.CognitoClient.InitiateAuth(context.TODO(), input)
	if err != nil {
		return dtos.LoginResponse{}, err
	}

	authResult := resp.AuthenticationResult
	return dtos.LoginResponse{
		AccessToken:  *authResult.AccessToken,
		IdToken:      *authResult.IdToken,
		RefreshToken: *authResult.RefreshToken,
		ExpiresIn:    authResult.ExpiresIn,
		TokenType:    *authResult.TokenType,
	}, nil
}

func (r *AuthRepositoryImpl) ConfirmUser(req dtos.ConfirmRequest) (dtos.ConfirmResponse, error) {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(clients.ClientID),
		Username:         aws.String(req.Username),
		ConfirmationCode: aws.String(req.Code),
	}

	_, err := clients.CognitoClient.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		return dtos.ConfirmResponse{}, err
	}

	return dtos.ConfirmResponse{Message: "User confirmed successfully."}, nil
}

func (r *AuthRepositoryImpl) ResendCode(req dtos.ResendRequest) (dtos.ResendResponse, error) {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(clients.ClientID),
		Username: aws.String(req.Username),
	}

	_, err := clients.CognitoClient.ResendConfirmationCode(context.TODO(), input)
	if err != nil {
		return dtos.ResendResponse{}, err
	}

	return dtos.ResendResponse{Message: "Confirmation code resent successfully."}, nil
}
