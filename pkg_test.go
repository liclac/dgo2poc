package dgo2poc

import (
	"os"
	"testing"

	"golang.org/x/oauth2"

	_ "github.com/joho/godotenv/autoload"
)

// Test token used for tests. Specify with DGO_TEST_TOKEN.
var TestToken *oauth2.Token

func init() {
	if t := os.Getenv("DGO_TEST_TOKEN"); t != "" {
		TestToken = BotToken(t)
	}
}

func requireTestToken(t *testing.T) {
	if TestToken == nil {
		t.Error("DGO_TEST_TOKEN env var not set")
		t.FailNow()
	}
}
