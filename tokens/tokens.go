package tokens

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var keys struct {
	publicKey ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func CreateKeys() {
	var err error
	keys.publicKey, keys.privateKey, err = ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("error generating key %v\n", err)
		panic(err)
	}
}

type customClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func CreateToken() string {
	

	// Create the Claims
	claims := customClaims{
		[]string{"user"},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject: "lol",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	ss, sErr := token.SignedString(keys.privateKey)
	if sErr != nil {
		fmt.Printf("error signing key %v\n", sErr)
	}
	
	return ss
}

type returnedClaims struct {
	User string
	Roles []string
}

func ValidateToken(tokenString string) (claims returnedClaims, ok bool) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if vErr := token.Claims.Valid(); vErr != nil {
			return nil, fmt.Errorf("error validating claims")
		}

		return keys.publicKey, nil
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