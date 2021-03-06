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
	"code.cloudfoundry.org/bytefmt"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"time"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List images",
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

		images, err := cli.ImageList(ctx, types.ImageListOptions{})
		if err != nil {
			panic(err)
		}

		if len(images) == 0 {
			fmt.Println("There are no images available")
		}

		fmt.Printf("Image ID\tSize\t\tCreated\t\t\tTags\n")
		for _, image := range images {
			fmt.Printf("%.12s\t%s\t\t%s\t\t%s\t\n", image.ID[7:], bytefmt.ByteSize(uint64(image.Size)), time.Unix(image.Created, 0).Format(time.Stamp), image.RepoTags[0])

		}
	},
}

func init() {
	ImageCmd.AddCommand(lsCmd)
	var images = *lsCmd
	images.Use = "images"
	RootCmd.AddCommand(&images)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
