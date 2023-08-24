package im_auth

import (
	"fmt"
	"testing"
)

func TestBuildClaims(t *testing.T) {
	InitSecret("lalala")
	token, err := NewToken("u1", 0, 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("token: %s\n", token)
	claims, err := GetClaimFromToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(claims)
}
