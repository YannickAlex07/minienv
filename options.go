package minienv

import (
	"bufio"
	"os"
	"regexp"
)

// Supply a map of overrides that will take precedence over
// any other environment variables.
// The keys are case-sensitive.
func WithOverrides(overrides map[string]string) Option {
	return func(m map[string]string) error {
		for k, v := range overrides {
			m[k] = v
		}

		return nil
	}
}

// Supply a list of files to load environment variables from.
// If any error occures it is ignored. Use WithRequiredFile if you want
// to fail in case any error is raised.
func WithFile(files ...string) Option {
	return func(m map[string]string) error {
		parseFiles(false, m, files...)
		return nil
	}
}

// Supply a list of files to load environment variables from.
// Will raise an error if any error occures.
func WithRequiredFile(files ...string) Option {
	return func(m map[string]string) error {
		err := parseFiles(true, m, files...)
		if err != nil {
			return err
		}

		return nil
	}
}

func parseFiles(raiseError bool, m map[string]string, files ...string) error {
	if len(files) == 0 || files == nil {
		files = []string{".env"}
	}

	for _, file := range files {
		envs, err := parseSingleEnvFile(file)
		if err != nil {
			if raiseError {
				return err
			}

			continue
		}

		for k, v := range envs {
			m[k] = v
		}
	}

	return nil
}

func parseSingleEnvFile(path string) (map[string]string, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
