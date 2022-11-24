package db

import (
	crand "crypto/rand"
	"testing"

	"github.com/sethvargo/go-password/password"
)

func TestGenApiKey(t *testing.T) {

	var err error

	keyGeneratorInput = &password.GeneratorInput{
		LowerLetters: password.LowerLetters,
		UpperLetters: password.UpperLetters,
		Digits:       password.Digits,
		Symbols:      "~@#%^&*()_+-={}[]:;<>?,./", // Some characters removed from the original, to be easily usable in Bash
		Reader:       crand.Reader,
	}

	keyGenerator, err = password.NewGenerator(keyGeneratorInput)
	if err != nil {
		t.Fatalf("FAIL: failed to create key generator: %s", err)
	}

	key, err := genAPIKey()
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	t.Logf("KEY: %s\n", key)
}

func BenchmarkGenApiKey(b *testing.B) {

	var err error

	keyGeneratorInput = &password.GeneratorInput{
		LowerLetters: password.LowerLetters,
		UpperLetters: password.UpperLetters,
		Digits:       password.Digits,
		Symbols:      "~@#%^&*()_+-={}[]:;<>?,./", // Some characters removed from the original, to be easily usable in Bash
		Reader:       crand.Reader,
	}

	keyGenerator, err = password.NewGenerator(keyGeneratorInput)
	if err != nil {
		b.Fatalf("FAIL: failed to create key generator: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := genAPIKey()
		if err != nil {
			b.Fatalf("FAIL: %s\n", err)
		}
	}
}
