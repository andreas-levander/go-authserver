package tokens

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Keys struct {
	publicKey *ed25519.PublicKey
	privateKey *ed25519.PrivateKey
}

func CreateKeys() *Keys{
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("error generating key %v\n", err)
		panic(err)
	}
	return &Keys{&publicKey, &privateKey}
}

type customClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (keys *Keys) Create(username string, roles []string) string {
	// Create the Claims
	claims := customClaims{
		roles,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-authserver",
			Subject: username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	ss, sErr := token.SignedString(*keys.privateKey)

	if sErr != nil {
		fmt.Printf("error signing key %v\n", sErr)
	}
	
	return ss
}

type returnedClaims struct {
	User string
	Roles []string
}

func (keys *Keys) Validate(tokenString string) (claims returnedClaims, ok bool) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if vErr := token.Claims.Valid(); vErr != nil {
			return nil, fmt.Errorf("error validating claims")
		}

		return *keys.publicKey, nil
	})
	if err != nil {
		fmt.Printf("error verifying key: %v\n", err)
		return returnedClaims{}, false
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return returnedClaims{claims.Subject, claims.Roles}, true
	}
	return returnedClaims{}, false
}