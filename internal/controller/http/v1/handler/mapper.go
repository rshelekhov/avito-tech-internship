package handler

import "github.com/rshelekhov/avito-tech-internship/internal/domain/entity"

func toUserCredentials(request *AuthRequest) entity.UserCredentials {
	return entity.UserCredentials{
		Username: request.username,
		Password: request.password,
	}
}
