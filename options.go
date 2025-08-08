package minienv

import (
	"bufio"
	"log"
	"maps"
	"os"
	"regexp"
)

// The Option func can be used to configure the loading behavior of minienv.
type Option func(*LoadConfig) error

// WithFallbackValues allows you to set fallback values for environment variables.
// These will be applied if no environment variable is found.
// If no fallback is sepcified either, the default value will be used.
func WithFallbackValues(values map[string]string) Option {
	return func(c *LoadConfig) error {
		maps.Copy(c.Values, values)

		return nil
	}
}

// WithPrefix allows you to set a prefix for the environment variables.
// Each environment variable will be prefixed with this value.
// If a specified value within the struct tag already has the specified prefix,
// it will not be prefixed again.
func WithPrefix(prefix string) Option {
	return func(c *LoadConfig) error {
		c.Prefix = prefix
		return nil
	}
}

// WithEnvFile can be used to read environment values from an .env file.
// If required is set to true, the file must exist. If it is set to false,
// the file is optional and will not cause an error if it does not exist.
func WithEnvFile(file string, required bool) Option {
	return func(c *LoadConfig) error {
		envs, err := parseEnvFile(file)
		if err != nil {
			if os.IsNotExist(err) && !required {
				return nil
			}

			return err
		}

		maps.Copy(c.Values, envs)

		return nil
	}
}

func parseEnvFile(path string) (map[string]string, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close env file %s: %v", path, err)
		}
	}()

	overrides := map[string]string{}

	// scan file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// compile regex
	r, err := regexp.Compile(`^(?P<key>\w+)=["']?(?P<value>[^'"]*)['"]?.*$`)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines
		if len(line) == 0 {
			continue
		}

		// check if line is a valid env line
		matches := r.FindStringSubmatch(line)
		if len(matches) == 0 || matches == nil {
			continue
		}

		overrides[matches[r.SubexpIndex("key")]] = matches[r.SubexpIndex("value")]
	}

	return overrides, nil
}
