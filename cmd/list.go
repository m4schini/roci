/*
Copyright © 2023 github.com/m4schini
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all containers managed by roci",
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			format = MustGetString(cmd, "format")
		)

		fmt.Println("list called")
		containers, err := confs.List()
		if err != nil {
			return err
		}

		switch format {
		case "table":
			w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
			fmt.Fprint(w, "ID\tPID\tSTATUS\tBUNDLE\n")
			for _, item := range containers {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
					item.ID,
					item.Pid,
					item.Status,
					item.Bundle)
			}
			if err := w.Flush(); err != nil {
				return err
			}
		case "json":
			if err := json.NewEncoder(os.Stdout).Encode(containers); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid format option")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().String("format", "table", "Possible values: table, json")
}
