package token

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenGen struct {
	issuer  string
	nowFunc func() time.Time
	privkey *rsa.PrivateKey
}

func NewJWTTokenGen(issuer string, privkey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:  issuer,
		nowFunc: time.Now,
		privkey: privkey,
	}
}

func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	now := t.nowFunc()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  now.Unix(),
		Subject:   accountID,
		ExpiresAt: now.Add(expire).Unix(),
	})

	return token.SignedString(t.privkey)
}
