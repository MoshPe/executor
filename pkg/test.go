package pkg

import (
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"io/ioutil"
	"log"
	"os"
)

func Test() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	yfile, err := ioutil.ReadFile("docker-compose.yml")

	if err != nil {

		log.Fatal(err)
	}
	cfg := types.ConfigDetails{
		Version:     "",
		WorkingDir:  workingDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yml", Content: yfile},
		},
		Environment: nil,
	}
	project, _ := loader.Load(cfg, func(options *loader.Options) {
		options.SkipConsistencyCheck = true
		options.SkipNormalization = true
		options.ResolvePaths = true
		options.Name = "docker-compose parser"
	})
	fmt.Println(project.Name)
	for _, service := range project.Services {
		fmt.Println(service.Environment)
	}
}