package gitlab

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Define static errors.
var (
	ErrInvalidConfigFormat = errors.New("invalid config format")
	ErrNoRepositoriesFound = errors.New("no repositories found")
)

// ExtractRepositories parses the config.js file and extracts the repositories array.
func ExtractRepositories(configPath string) ([]string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer file.Close()

	var repositories []string

	inRepositoriesSection := false

	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if isRepositoriesArrayStart(line) {
			inRepositoriesSection = true

			continue
		}

		if inRepositoriesSection {
			if err := processRepositoryLine(line, &repositories, &inRepositoriesSection); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if len(repositories) == 0 {
		return nil, ErrNoRepositoriesFound
	}

	return repositories, nil
}

// processRepositoryLine processes a single line within the repositories section.
func processRepositoryLine(line string, repositories *[]string, inRepositoriesSection *bool) error {
	if isRepositoriesArrayEnd(line) {
		*inRepositoriesSection = false

		return nil
	}

	if isInvalidFormatAfterEnd(line) {
		return ErrInvalidConfigFormat
	}

	if isCommentLine(line) {
		return nil
	}

	line = trimRepositoryLine(line)
	if line != "" && line != "[" {
		*repositories = append(*repositories, line)
		logrus.WithField("repository", line).Debug("Found repository")
	}

	return nil
}

// isRepositoriesArrayStart checks if a line indicates the start of the repositories array.
func isRepositoriesArrayStart(line string) bool {
	return strings.HasPrefix(line, "repositories:")
}

// isRepositoriesArrayEnd checks if the line marks the end of the repositories array.
func isRepositoriesArrayEnd(line string) bool {
	return strings.Contains(line, "]")
}

// isInvalidFormatAfterEnd checks if a line indicates an invalid format after the end of the repositories array.
func isInvalidFormatAfterEnd(line string) bool {
	return strings.HasSuffix(line, "};")
}

// isCommentLine checks if a line is a comment.
func isCommentLine(line string) bool {
	return strings.HasPrefix(line, "//")
}

// trimRepositoryLine removes surrounding quotes and commas from repository entries and trims spaces.
func trimRepositoryLine(line string) string {
	line = strings.Trim(line, `",`)
	line = strings.TrimSpace(line)

	return line
}
