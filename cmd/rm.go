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
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var ContainerRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			RemoveVolumes: cmd.Flag("force").Changed,
			RemoveLinks:   cmd.Flag("link").Changed,
			Force:         cmd.Flag("volumes").Changed,
		})
		if err != nil {
			panic(err)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("required a container ID / name")
		}
		return nil
	},
}

func init() {
	ContainerRmCmd.Flags().BoolP("force", "f", false, "Force the removal of a running container (uses SIGKILL)")
	ContainerRmCmd.Flags().BoolP("link", "l", false, "Remove the specified link")
	ContainerRmCmd.Flags().BoolP("volumes", "v", false, "Remove anonymous volumes associated with the container")
	ContainerCmd.AddCommand(ContainerRmCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
