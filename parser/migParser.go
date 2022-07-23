package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/squareboat/splinter/config"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/logger"
)

func GetMigrationFileNames() ([]string, error) {
	var migrationFileNames []string
	files, err := ioutil.ReadDir(config.GetMigrationsPath())
	if err != nil {
		logger.Log.Error(err)
		return migrationFileNames, err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), constants.DEFAULT_FILE_EXTENSION) {
			migrationFileNames = append(migrationFileNames, file.Name())
		}
	}
	return migrationFileNames, nil
}

func fileParser(lines []string, down bool) []string {
	var upArr []string
	var downArr []string
	var isUp bool
	var parts []string
	var remainingString string
	for _, text := range lines {
		if strings.ToLower(text) == constants.MIGRATION_UP_IDENTIFIER {
			isUp = true
			continue
		} else if strings.ToLower(text) == constants.MIGRATION_DOWN_IDENTIFIER {
			if !down {
				return upArr
			}
			isUp = false
			continue
		}
		text = strings.Trim(text, " \n")
		if text == "" {
			continue
		}
		parts = stringParser(text, &remainingString)
		if isUp {
			upArr = append(upArr, parts...)
		} else {
			downArr = append(downArr, parts...)
		}

	}
	if down {
		return downArr
	}
	return upArr
}

func stringParser(text string, remainingString *string) []string {
	if remainingString == nil {
		logger.Log.Warn("remainingString is nil")
		return []string{}
	}
	var parsed []string
	var currIdx int = strings.Index(text, ";")
	for {

		if currIdx == -1 {
			// TODO Find a better way to do this if possible
			*remainingString = *remainingString + " " + text
			return parsed
		}

		str := text[:currIdx+1]
		if remainingString != nil && *remainingString != "" {
			str = *remainingString + " " + str
		}

		parsed = append(parsed, str)
		text = text[currIdx+1:]
		currIdx = strings.Index(text, ";")
		if remainingString != nil {
			// TODO Find a better way to do this if possible
			*remainingString = ""
		}
	}
}

func ParseFile(filename string, mode string) ([]string, error) {
	filePath := fmt.Sprintf("%s/%s", config.GetMigrationsPath(), filename)
	file, err := os.ReadFile(filePath)
	if err != nil {
		logger.Log.WithError(err).Error("Error reading file")
		return nil, err
	}
	strs := strings.Split(string(file), "\n")

	queries := fileParser(strs, mode == constants.MIGRATION_DOWN)
	return queries, nil
}

func CreateMigrationFile(names []string) {
	for _, name := range names {
		filename := fmt.Sprintf("%s/%d_%s%s", config.GetMigrationsPath(), time.Now().UnixMicro(), name, constants.DEFAULT_FILE_EXTENSION)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		file.Write([]byte("[up]\n\n"))
		file.Write([]byte("[down]\n\n"))
		file.Close()
	}
}
