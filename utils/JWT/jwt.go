package JWT

import (
	"github.com/dgrijalva/jwt-go"
	"encoding/hex"
)

var (
	key = []byte("TSL DARK SECRET")
)

func GenerateToken(data map[string] interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	for k,v := range data {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err  := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString([]byte(tokenStr)),nil
}


func ParseToken(tokenStr string) (data map[string] interface{}, err error) {
	realToken, err := hex.DecodeString(tokenStr)
	tokenStr = string(realToken)
	claims := make(jwt.MapClaims)
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return data, err
	}

	return claims, nil
}