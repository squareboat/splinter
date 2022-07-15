package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
)

func GetMigrationFileNames() []string {
	var migrationFileNames []string
	files, err := ioutil.ReadDir(viper.GetString(constants.SPLINTER_PATH))
	if err != nil {
		logger.Log.Error(err)
		return migrationFileNames
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), constants.FILE_EXTENSION) {
			migrationFileNames = append(migrationFileNames, file.Name())
		}
	}
	return migrationFileNames
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
	var parsed []string
	var currIdx int = strings.Index(text, ";")
	for {
		if currIdx == -1 {
			tempStr := text + *remainingString
			remainingString = &tempStr
			break
		}
		str := text[:currIdx+1]
		if remainingString != nil && *remainingString != "" {
			tempStr := *remainingString + " " + str
			remainingString = &tempStr
		}
		parsed = append(parsed, str)
		text = text[currIdx+1:]
		currIdx = strings.Index(text, ";")
		if remainingString != nil {
			tempStr := ""
			remainingString = &tempStr
		}
	}
	return parsed
}

func ParseFile(filename string, mode string) ([]string, error) {
	filePath := fmt.Sprintf("%s/%s", viper.GetString(constants.SPLINTER_PATH), filename)
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
		filename := fmt.Sprintf("%s/%d_%s%s", viper.GetString(constants.SPLINTER_PATH), time.Now().UnixMicro(), name, constants.FILE_EXTENSION)
		file, err := os.Create(filename)
		file.Write([]byte("[up]\n\n"))
		file.Write([]byte("[down]\n\n"))
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
}
