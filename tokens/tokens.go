package tokens

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Tokens interface {
	Create(username string, roles []string, token_ttl int) string
	Validate(tokenString string) (claims returnedClaims, ok bool)
	PublicKey() *jwk.Key
}

type Keys struct {
	publicKey *jwk.Key
	privateKey *jwk.Key
}

func CreateKeys() *Keys{
	kid := uuid.New()

	rawPublicKey, rawPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Println("failed to create keys")
	}
	privateKey, err := jwk.FromRaw(rawPrivateKey)
	if err != nil {
		panic("failed to import key from raw: " + err.Error())
	}
	publicKey, err := jwk.FromRaw(rawPublicKey)
	if err != nil {
		panic("failed to import key from raw: " + err.Error())
	}
	privateKey.Set(jwk.KeyIDKey, kid.String())
	publicKey.Set(jwk.KeyIDKey, kid.String())

	return &Keys{&publicKey, &privateKey}
}

type customClaims struct {
	Roles []string `json:"roles"`
	Kid string `json:"kid"`
	jwt.RegisteredClaims
}

func (keys *Keys) Create(username string, roles []string, token_ttl int) string {
	privateKey := *keys.privateKey
	var rawPrivateKey ed25519.PrivateKey
	if err := privateKey.Raw(&rawPrivateKey); err != nil {
		panic("failed to get raw key" + err.Error())
	}
	// Create the Claims
	claims := customClaims{
		roles,
		privateKey.KeyID(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(token_ttl))),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-authserver",
			Subject: username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	ss, sErr := token.SignedString(rawPrivateKey)

	if sErr != nil {
		panic("error sigining token: " + sErr.Error())
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

func (keys *Keys) PublicKey() *jwk.Key {
	return keys.publicKey
}