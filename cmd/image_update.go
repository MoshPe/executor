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
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
const updateSuccess = "Updated image successfully"

var containerID []string
var imageID []string

type containersToUpdate struct {
	containers      []types.Container
	containerConfig []types.ContainerJSON
}

var updateCmd = &cobra.Command{
	Use:   "update IMAGE_FROM IMAGE_TO",
	Short: "Update all running containers with the same image but newer version - USE WITH CARE!!!",
	Long: `The application stops and rms all containers with the exact image, pull an image
and starting a new containers again - USE WITH CARE!!!.`,
	Run: func(cmd *cobra.Command, args []string) {
		var con containersToUpdate
		containerID = make([]string, 1)
		imageIdFrom := args[0]
		imageIdTo := args[1]
		imageID = append(imageID, imageIdTo)
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
		//running the pull image command
		PullCmd.Run(cmd, imageID)

		if len(containers) == 0 {
			fmt.Println("No containers available - " + updateSuccess + imageIdTo)
		}

		//search for the containers with the exact image to update
		for _, container := range containers {
			if container.Image == imageIdFrom {
				con.containers = append(con.containers, container)
			}
		}

		con.containerConfig = make([]types.ContainerJSON, len(con.containers))
		//saving the container configuration
		for i, container := range con.containers {
			containerSpec, err := cli.ContainerInspect(ctx, container.ID)
			if err != nil {
				panic(err)
			}
			con.containerConfig[i] = containerSpec
		}

		//stopping the containers
		for _, container := range con.containers {
			containerID[0] = container.ID
			StopCmd.Run(cmd, containerID)
		}

		//rming the containers
		for _, container := range con.containers {
			err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				panic(err)
			}
		}

		//restarting the containers with the new image
		for i := range con.containers {
			containerConfig := con.containerConfig[i].HostConfig
			net := network.NetworkingConfig{EndpointsConfig: con.containerConfig[i].NetworkSettings.Networks}
			//setting the image
			con.containerConfig[i].Config.Image = imageIdTo
			response, err := cli.ContainerCreate(ctx, con.containerConfig[i].Config, containerConfig, &net,
				nil, con.containerConfig[i].Name)

			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, response.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			fmt.Println(response.ID)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("required an image name")
		}
		return nil
	},
}

func init() {
	ImageCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
