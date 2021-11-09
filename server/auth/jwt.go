package auth

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func JWT(audience string, subjID uuid.UUID, key []byte) string {
	sig, _ := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: key},
		(&jose.EncrypterOptions{}).WithType("JWT"))

	cl := jwt.Claims{
		Subject:  subjID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Audience: []string{audience},
	}

	raw, _ := jwt.Encrypted(sig).Claims(cl).CompactSerialize()
	return raw
}
