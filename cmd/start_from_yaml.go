// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"executor/pkg"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:        "start",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Example:           "",
	ValidArgs:         nil,
	ValidArgsFunction: nil,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("accepts 1 arg(s)")
		}
		return nil
	},
	ArgAliases:             nil,
	BashCompletionFunction: "",
	Deprecated:             "",
	Annotations:            nil,
	Version:                "",
	PersistentPreRun:       nil,
	PersistentPreRunE:      nil,
	PreRun:                 nil,
	PreRunE: 				nil,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fileName string
			filePath string
		)
		filePath = args[0]
		fileExists, err := ConfigFile.FileExists(filePath)
		if err != nil {
			panic(err)
		}
		if fileExists {
			fileName = filepath.Base(filePath)
		} else {
			errMsg := "File "+filePath+" doesn't Exists"
			panic(errMsg)
		}
		fileName = ConfigFile.FilenameWithoutExtension(fileName)
		Project, err = pkg.ParserYaml(fileName)
		if err != nil {
			log.Fatalln(err)
		}
		if err = startServices(); err != nil {
			panic(err)
		}
	},
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startServices() error{
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	authConfig := types.AuthConfig{
		Username: "moshpe",
		Password: "MoshPe2969999",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	for _, service := range Project.Services {
		fmt.Println("image = ",service.Image)
		reader, err := cli.ImagePull(ctx, service.Image, types.ImagePullOptions{
			RegistryAuth: authStr,
		})
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return err
		}

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: service.Image,
			Cmd:   strslice.StrSlice(service.Command),
			Tty:   service.Tty,
			OpenStdin: service.StdinOpen,
		}, nil, nil, nil, service.ContainerName)
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		//out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		//if err != nil {
		//	panic(err)
		//}
		//
		//stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}
	return nil
}
