package backend

import (
	"bufio"
	"os"
	"strings"
)

type Section struct {
	Name      string
	Variables map[string]string
}

func ReadEnvFile(filename string) ([]Section, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sections []Section
	var currentSection *Section
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a section header
		if strings.HasPrefix(line, "#") && strings.HasSuffix(line, "CONFIGURATION") {
			// Create a new section
			sectionName := strings.TrimSuffix(strings.TrimPrefix(line, "#"), "CONFIGURATION")
			currentSection = &Section{
				Name:      strings.TrimSpace(sectionName),
				Variables: make(map[string]string),
			}
			sections = append(sections, *currentSection)
		} else if strings.Contains(line, "=") {
			// Parse the variable and its default value
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				// Skip lines that don't look like key=value
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Add the variable to the current section
			if currentSection != nil {
				currentSection.Variables[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}
