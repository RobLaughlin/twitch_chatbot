package env

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Validate returns (map[string]string, error) where map is all of the valid env keys to values.
// envFile is the filepath for the .env file, and validKeys is a list of keys to expect in the .env file.
func Validate(envFile string, validKeys []string) (map[string]string, error) {
	env, err := godotenv.Read(envFile)

	if err != nil {
		return map[string]string{}, err
	}

	for _, key := range validKeys {
		if _, ok := env[key]; !ok {
			return map[string]string{}, fmt.Errorf("Missing env variable: %s", key)
		}
	}

	return env, nil
}
