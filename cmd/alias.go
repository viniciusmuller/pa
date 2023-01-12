package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

type PaDataAliases struct {
	Name    string
	Command string
}

type PaData struct {
	Aliases []PaDataAliases
}

var alias string
var command string

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "the alias to use")
	versionCmd.PersistentFlags().StringVarP(&command, "command", "c", "", "the command to run")
}

var versionCmd = &cobra.Command{
	Use:   "alias",
	Short: "Creates an alias",
	Long:  "Aliases are project-specific and are designed for saving you keystrokes",
	Run: func(cmd *cobra.Command, args []string) {
		var dataDirectory = path.Join(xdg.DataHome, "/pa")
		if _, err := os.Stat(dataDirectory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dataDirectory, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		// TODO: Read config file, update it based on user input and write it
		// back again
		var configFile = path.Join(dataDirectory, "data.json")
		data := PaData{
			Aliases: []PaDataAliases{
				{
					Name:    "test",
					Command: "mix test",
				},
			},
		}

		file, _ := json.MarshalIndent(data, "", "  ")
		_ = ioutil.WriteFile(configFile, file, 0644)
	},
}
