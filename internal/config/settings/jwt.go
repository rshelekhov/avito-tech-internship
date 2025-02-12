package settings

import (
	"time"

	"github.com/rshelekhov/avito-tech-internship/internal/domain/service/token"
)

type JWT struct {
	Secret string        `mapstructure:"JWT_SECRET" envDefault:"secret"`
	TTL    time.Duration `mapstructure:"JWT_TTL" envDefault:"24h"`
}

func ToJWTConfig(params JWT) token.JWT {
	return token.JWT{
		Secret: params.Secret,
		TTL:    params.TTL,
	}
}
