package platform

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenRejectsWrongSecretAndExpired(t *testing.T) {
	secret := "test-secret-at-least-32-characters"
	token, e := Token(secret, "client", "Client", "CLIENT")
	if e != nil {
		t.Fatal(e)
	}

	c, e := Parse(secret, token)
	if e != nil || c.Subject != "client" {
		t.Fatal("valid token rejected")
	}

	if _, e = Parse(secret+"wrong", token); e == nil {
		t.Fatal("wrong signature accepted")
	}

	expired, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{Role: "CLIENT", RegisteredClaims: jwt.RegisteredClaims{Subject: "client", Issuer: "robust", Audience: jwt.ClaimStrings{"credit"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour))}}).SignedString([]byte(secret))
	if _, e = Parse(secret, expired); e == nil {
		t.Fatal("expired token accepted")
	}
}

func TestErrorHidesInternalDetails(t *testing.T) {
	b, _ := json.Marshal(HTTPError(RPCError(Invalid("bad input"))))
	if string(b) == "" || HTTPError(RPCError(Conflict("conflict"))).Status != 409 {
		t.Fatal("invalid error mapping")
	}
}
