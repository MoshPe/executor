package pkg

import (
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"io/ioutil"
	"log"
	"os"
)

func ParserYaml(filename string) (*types.Project, error){

	workingDir, err := os.Getwd()
	if err != nil {
		return nil,err
	}

	ymlFile, err := ioutil.ReadFile(filename+".yml")
	if err != nil {
		return nil,err
	}

	cfg := types.ConfigDetails{
		Version:     "",
		WorkingDir:  workingDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: filename, Content: ymlFile},
		},
		Environment: nil,
	}

	project, err := loader.Load(cfg, func(options *loader.Options) {
		options.SkipConsistencyCheck = true
		options.SkipNormalization = true
		options.ResolvePaths = true
		options.Name = filename
	})
	if err != nil {
		return nil,err
	}
	log.Println("Config file "+filename+" has been loaded!")
	return project,nil
}