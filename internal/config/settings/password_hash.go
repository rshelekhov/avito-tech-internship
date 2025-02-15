package settings

import "github.com/rshelekhov/avito-tech-internship/internal/domain/service/token"

type PasswordHash struct {
	Pepper     string `mapstructure:"PASSWORD_HASH_PEPPER"`
	BcryptCost int    `mapstructure:"PASSWORD_HASH_BCRYPT_COST" envDefault:"10"`
}

func ToPasswordHashConfig(params PasswordHash) token.PasswordHash {
	return token.PasswordHash{
		Pepper:     params.Pepper,
		BcryptCost: params.BcryptCost,
	}
}
