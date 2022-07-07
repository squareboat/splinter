package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/the-e3n/migrator/constants"
)

func MigParser(fileName string) ([]string, []string, error) {
	var upArr []string
	var downArr []string
	var isUp bool
	file, err := os.Open(fileName)
	if err != nil {
		return upArr, downArr, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentString string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "[down]" {
			isUp = false
			continue
		} else if strings.ToLower(text) == "[up]" {
			isUp = true
			continue
		}
		text = strings.Trim(text, " \n")
		if text == "" {
			continue
		}
		idx := strings.Index(text, ";")
		if idx == -1 {
			currentString += text
		} else {
			currentString += text + " "
			if isUp {
				upArr = append(upArr, currentString)
			} else {
				downArr = append(downArr, currentString)
			}
			currentString = ""
		}
	}
	return upArr, downArr, nil
}

type QueriesToRun struct {
	Up   []string
	Down []string
}

func ParseAllMigrations() []QueriesToRun {
	querys := []QueriesToRun{}
	migrationPath := viper.GetString(constants.SPLINTER_PATH)
	files, _ := ioutil.ReadDir(migrationPath)
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for _, file := range files {
		fmt.Printf("Parsing :- %s\n", file.Name())
		if file.IsDir() {
			continue
		}
		query := QueriesToRun{}
		up, down, err := MigParser(migrationPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		query.Up = append(query.Up, up...)
		query.Down = append(query.Down, down...)
		querys = append(querys, query)
	}
	return querys
}
func ParseRollbackMigration() ([]string, error) {
	migrationPath := viper.GetString(constants.SPLINTER_PATH)
	files, _ := ioutil.ReadDir(migrationPath)
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	filePath := fmt.Sprintf("%s/%s", migrationPath, files[len(files)-1].Name())
	_, down, err := MigParser(filePath)
	return down, err
}

func CreateMigrationFile(names []string) {
	for _, name := range names {
		filename := fmt.Sprintf("%s/%d_%s.sql", viper.GetString(constants.SPLINTER_PATH), time.Now().UnixMicro(), name)
		file, err := os.Create(filename)
		file.Write([]byte("[up]\n\n"))
		file.Write([]byte("[down]\n\n"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
}
