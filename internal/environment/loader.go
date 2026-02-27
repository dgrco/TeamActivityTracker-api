package environment

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables and ensure all necessary variables are set.
// Returns a pointer to a filled Environment object on success.
// On failure, exits the program with an explanation.
func Load() *Environment {
	env := &Environment{}

	// Load .env file if it exists
	_ = godotenv.Load()

	// Read variables into corresponding Environment fields
	env.DatabaseURL = readVar("DATABASE_URL")
	env.JWTSecret = readVar("JWT_SECRET")
	env.Port = readVar("PORT")

	return env
}

// Returns a string if the variable exists.
// If the variable does not exist, exit and print an error.
func readVar(envVar string) string {
	val, ok := os.LookupEnv(envVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "ERROR: required environment variable \"%s\" not set\n", envVar)
		os.Exit(1)
	}

	return val
}
