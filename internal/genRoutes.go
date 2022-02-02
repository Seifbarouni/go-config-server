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

func GenerateRoutes(app *fiber.App) {
	configFiles:= getConfigFiles()
	for _, file := range configFiles {
		if file.IsDir() {
			continue
		}
		fname:=file.Name()
		// check if the file is a yaml file
		if !strings.HasSuffix(fname,".yaml") && !strings.HasSuffix(fname,".yml") {
			fmt.Printf("%s is not a yaml file\n",fname)
			continue
		}
		// delete .yaml or .yml from the file name
		fname = strings.Replace(fname,".yaml","",-1)
		fname = strings.Replace(fname,".yml","",-1)

		// check if - exists only once in the file name
		if strings.Count(fname, "-") != 1 {
			fmt.Printf("%s is not a valid file name\nThe file name should be [appName]-[branchName].yaml",fname)
			continue
		}
		
		// check if - exist in the end or the beginning of the file name
		if strings.HasPrefix(fname,"-") || strings.HasSuffix(fname,"-") {
			fmt.Printf("%s is not a valid file name\nThe file name should be [appName]-[branchName].yaml",fname)
			continue
		}
		
		// split the file name by -
		appAndBranch:=strings.Split(fname,"-")
		
		appName := appAndBranch[0]
		branchName := appAndBranch[1]
		
		yamlFile, err := ioutil.ReadFile("./config/"+file.Name())
		if err != nil {
			panic(err)
		}
		i := make(map[string]interface{})
		err = yaml.Unmarshal(yamlFile, &i)
		if err != nil {
			panic(err)
		}

		app.Get(fmt.Sprintf("/%s/%s",appName,branchName),func(c *fiber.Ctx) error {
			return c.JSON(i)
		})

	}
}