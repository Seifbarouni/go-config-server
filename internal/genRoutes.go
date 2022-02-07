package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

func getConfigFiles() []os.FileInfo {
	f, err := os.Open("./config")
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func isValidYamlFileName(fname string) (bool,string,string) {
	// check if the file is a yaml file
	if !strings.HasSuffix(fname, ".yaml") && !strings.HasSuffix(fname, ".yml") {
		return false, "", ""
	}

	// delete .yaml or .yml from the file name
	fname = strings.Replace(fname, ".yaml", "", -1)
	fname = strings.Replace(fname, ".yml", "", -1)

	// check if - exists  in the file name
	dashCount:=strings.Count(fname, "-")

	if  dashCount == 0 {
		return false, "", ""
	}

	// check if - exist in the end or the beginning of the file name
	if strings.HasPrefix(fname, "-") || strings.HasSuffix(fname, "-") {
		return false, "", ""
	}

	// get the last index of -
	lastDashIndex := strings.LastIndex(fname, "-")
	appName:=fname[:lastDashIndex]
	branchName:=fname[lastDashIndex+1:]
	
	
	return true, appName, branchName
}

func GenerateRoutes(app *fiber.App) {
	configFiles := getConfigFiles()
	for _, file := range configFiles {
		if file.IsDir() {
			continue
		}
		fname := file.Name()

		valid, appName, branchName := isValidYamlFileName(fname)

		if !valid {
			continue
		}
		

		yamlFile, err := ioutil.ReadFile("./config/" + file.Name())
		if err != nil {
			panic(err)
		}
		i := make(map[string]interface{})
		err = yaml.Unmarshal(yamlFile, &i)
		if err != nil {
			panic(err)
		}

		app.Get(fmt.Sprintf("/%s/%s", appName, branchName), func(c *fiber.Ctx) error {
			return c.JSON(i)
		})
	}
}
