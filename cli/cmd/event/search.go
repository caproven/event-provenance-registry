// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package event

import (
	"encoding/json"
	"fmt"

	"github.com/sassoftware/event-provenance-registry/cli/cmd/common"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Searches for event objects",
	Long:    `Searches for event objects`,
	PreRunE: common.BindFlagsE,
	RunE:    runSearchEvent,
}

// runSearchEvent runs the search and return error
func runSearchEvent(_ *cobra.Command, _ []string) error {
	dryrun := viper.GetBool("dry-run")
	noindent := viper.GetBool("no-indent")

	url := viper.GetString("url")
	c, err := common.GetClient(url)
	if err != nil {
		return err
	}

	params := make(map[string]interface{})

	id := viper.GetString("id")
	if id != "" {
		params["id"] = id
	}

	fields, err := common.ProcessSearchFields(viper.GetStringSlice("fields"), &storage.Event{})
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("ID: %s", id)
		fmt.Printf("Fields: %v", fields)
		return nil
	}

	events, err := c.SearchEvents(params, fields)
	if err != nil {
		return err
	}

	if noindent {
		content, err := json.Marshal(events)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", content)
		return nil
	}
	content, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", content)
	return nil
}

// NewSearchCmd returns a new search command
func NewSearchCmd() *cobra.Command {
	searchCmd.Flags().String("id", "", "Id for the event")
	searchCmd.Flags().String("fields", "id name version multipass", "Space delimited list of fields, or 'all' for all user fields")
	searchCmd.Flags().String("jsonpath", "", "JSONPath expression to apply to output")
	searchCmd.Flags().String("url", "http://localhost:8042", "EPR base url")
	searchCmd.Flags().Bool("dry-run", false, "do a dry run of the command")
	searchCmd.Flags().Bool("no-indent", false, "do not indent the JSON output")
	return searchCmd
}
