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

// WithFile allows you to specify a list of environment files that should be read.
func WithFile(required bool, files ...string) Option {
	return func(c *LoadConfig) error {
		values, err := readEnvFiles(required, files...)
		if err != nil {
			return err
		}

		maps.Copy(c.Values, values)

		return nil
	}
}

// Reads a list of env-files and sets them in the load config
func readEnvFiles(shouldRaiseError bool, files ...string) (map[string]string, error) {
	values := make(map[string]string)

	if len(files) == 0 || files == nil {
		files = []string{".env"}
	}

	for _, file := range files {
		envs, err := parseEnvFile(file)
		if err != nil {
			if shouldRaiseError {
				return nil, err
			}

			continue
		}

		maps.Copy(values, envs)
	}

	return values, nil
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
