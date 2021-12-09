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
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var containerLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		//can add filtering options
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		if len(containers) == 0 {
			fmt.Println("There are no containers available")
		}

		for _, container := range containers {
			fmt.Printf("Container ID: %s\nImage: %s\nCreated: %s"+
				"\nCommand: \"%s\"\n"+
				"Status: %s\n", container.ID, container.Image, time.Unix(container.Created, 0).Format(time.RFC850), container.Command,
				container.Status)
			fmt.Printf("Ports:\n")
			for _, port := range container.Ports {
				fmt.Printf("\t %d:%d", port.PrivatePort, port.PublicPort)
			}
			fmt.Printf("Names:")
			for _, name := range container.Names {
				fmt.Printf("\t %s\n", name[1:])
			}
		}
	},
}

func init() {
	ContainerCmd.AddCommand(containerLsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
