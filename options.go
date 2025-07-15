package minienv

import (
	"bufio"
	"log"
	"maps"
	"os"
	"regexp"
)

type Option func(*LoadConfig) error

// Supply a map of values that will be used as fallback values if no
// matching environment variable was found.
// The keys are case-sensitive.
func WithFallbackValues(values map[string]string) Option {
	return func(c *LoadConfig) error {
		maps.Copy(c.Values, values)

		return nil
	}
}

// Supply a prefix that will be added to all environment variables and fallback values.
func WithPrefix(prefix string) Option {
	return func(c *LoadConfig) error {
		c.Prefix = prefix
		return nil
	}
}

// Supply a list of files to load environment variables from that will be
// uses as fallback values in case no matching env variable was found.
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
