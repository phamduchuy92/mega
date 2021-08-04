package test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestFastHTTPClientGet give an example of using fasthttp for sending request
func TestBcrypt(t *testing.T) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("admin"), 15)
	t.Logf("%s", string(bytes))
}
