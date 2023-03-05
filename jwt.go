package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtMake 制作jwt串
func JwtMake(id, secret, issuer, subject, audience string, maxAge int) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  []string{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(maxAge) * time.Second)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtStr, err := token.SignedString(Str2SliceByte(secret))
	if err != nil {
		return "", NewSecurityError(err.Error())
	}
	return jwtStr, nil
}

// JwtParse 解析jwt
func JwtParse(jwtStr, secret string) (*jwt.RegisteredClaims, error) {
	tok, err := jwt.ParseWithClaims(jwtStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return Str2SliceByte(secret), nil
	})
	if err != nil {
		return nil, NewSecurityError("token 不可用!")
	}
	if claims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, NewSecurityError("token 不可用!")
}
