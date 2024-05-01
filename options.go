package minienv

import "os"

func WithOverrides(overrides map[string]string) Option {
	return func(m map[string]string) error {
		for k, v := range overrides {
			m[k] = v
		}

		return nil
	}
}

func WithFile(canFail bool, files ...string) Option {
	return func(m map[string]string) error {
		if len(files) == 0 || files == nil {
			files = []string{".env"}
		}

		for _, file := range files {
			envs, err := parseEnvFile(file)
			if err != nil {
				if !canFail {
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
}

func parseEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return nil, nil
}
