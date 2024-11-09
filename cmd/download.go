package cmd

import (
	"context"
	"fmt"
	"github.com/ides15/todoist"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a backup in JSON format from Todoist",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := todoist.NewClient(Config.TOKEN)
		if err != nil {
			panic(err)
		}
		projects, _, err := client.Projects.List(context.Background(), "")
		if err != nil {
			panic(err)
		}

		for _, p := range projects {
			fmt.Println(p.ID, p.Name)
		}
	},
}
