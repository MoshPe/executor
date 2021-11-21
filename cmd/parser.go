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
	"errors"
	"executor/pkg"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// parserCmd represents the parser command
var parserCmd = &cobra.Command{
	Use:   "parser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1{
			return errors.New("accepts 1 arg(s)")
		}
		return nil
	},
	Example: `executor ./resources/{file name}.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fileName string
			filePath string
		)
		filePath = args[0]
		fileExists, err := ConfigFile.FileExists(filePath)
		if err != nil {
			fmt.Println(err)
		}
		if fileExists {
			fileName = filepath.Base(filePath)
			if err != nil {
				fmt.Println(err.Error())

			}
		} else {
			fmt.Printf("File %v doest not Exists", filePath)
			return
		}
		fileName = ConfigFile.FilenameWithoutExtension(fileName)
		fmt.Println(fileName)
		fmt.Println(filePath)
		Project,err = pkg.ParserYaml(fileName)
		if err != nil {
			log.Fatalln("Couldn't parse the file")
		}
	},
}

func init() {
	RootCmd.AddCommand(parserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
