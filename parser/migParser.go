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
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/logger"
)

func MigParser(filepath string,) ([]string, []string, error) {
	var upArr []string
	var downArr []string
	var isUp bool
	file, err := os.Open(filepath)
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

func ParseAllMigrations() map[string]QueriesToRun {
	querys := map[string]QueriesToRun{}
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
		querys[file.Name()] = query
	}
	return querys
}

func ParseRollbackMigration(migration_name string) ([]string, error) {
	filePath := findMigrationFilePath(migration_name)
	exists, _ := os.Stat(filePath)
	if exists != nil {
		_, down, err := MigParser(filePath)
		return down, err
	}
	logger.Log.Fatal("Migration file not found")
	os.Exit(1)
	return nil, nil
}

func findMigrationFilePath(name string) string {
	migrationPath := viper.GetString(constants.SPLINTER_PATH)
	files, _ := ioutil.ReadDir(migrationPath)
	for _, file := range files {
		if file.Name() == name {
			return fmt.Sprintf("%s/%s", migrationPath, file.Name())
		}
	}
	return ""
}

func CreateMigrationFile(names []string) {
	viper.AddConfigPath("./")
	viper.SetConfigFile("test.json")
	logger.Log.Info(viper.AllSettings())
	for _, name := range names {
		filename := fmt.Sprintf("%s/%d_%s.sql", viper.GetString(constants.SPLINTER_PATH), time.Now().UnixMicro(), name)
		file, err := os.Create(filename)
		file.Write([]byte("[up]\n\n"))
		file.Write([]byte("[down]\n\n"))
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
}
