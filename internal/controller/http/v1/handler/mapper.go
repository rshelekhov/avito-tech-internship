package handler

import "github.com/rshelekhov/merch-store/internal/domain/entity"

func toUserCredentials(request *AuthRequest) entity.UserCredentials {
	return entity.UserCredentials{
		Username: request.Username,
		Password: request.Password,
	}
}
