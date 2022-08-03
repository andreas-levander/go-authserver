package tokens

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
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


func (keys *Keys) Create(username string, roles []string, token_ttl int) string {
	privateKey := *keys.privateKey

	token, err := jwt.NewBuilder().
				Claim("roles", roles).
				Claim("kid", privateKey.KeyID()).
				Issuer("go-authserver").
				IssuedAt(time.Now()).
				Expiration(time.Now().Add(time.Minute * time.Duration(token_ttl))).
				Subject(username).
				Build()

	if err != nil {
		panic("failed creating token" + err.Error())
	}

	  // Sign a JWT!
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.EdDSA, privateKey))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return ""
	}
	
	return string(signed)
}

type returnedClaims struct {
	User string `json:"sub"`
	Roles []string `json:"roles"`
}

func (keys *Keys) Validate(tokenString string) (claims returnedClaims, ok bool) {
	verifiedToken, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.EdDSA, *keys.publicKey), jwt.WithValidate(true))
	if err != nil {
		fmt.Printf("failed to verify JWS: %s\n", err)
		return returnedClaims{}, false
	}
	jsonm, err := json.Marshal(verifiedToken)
	if err != nil {
		fmt.Println("error marshaling token: " + err.Error())
		return returnedClaims{}, false
	}
	var Rclaims returnedClaims
	if err := json.Unmarshal(jsonm, &Rclaims); err != nil {
		fmt.Println("error unmarshaling token: " + err.Error())
		return returnedClaims{}, false
	}
	
	return Rclaims, true
}

func (keys *Keys) PublicKey() *jwk.Key {
	return keys.publicKey
}