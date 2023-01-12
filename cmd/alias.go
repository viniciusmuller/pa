package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

type PaAlias struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type PaProject struct {
	Aliases []PaAlias `json:"aliases"`
}

type PaData struct {
	Projects map[string]PaProject
}

var alias string
var command string
var deleteAlias bool

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "the alias to use")
	versionCmd.PersistentFlags().BoolVarP(&deleteAlias, "delete", "d", false, "whether to delete the alias")
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

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current directory: %s", err)
		}

		var dataFilePath = path.Join(dataDirectory, "data.json")
		content, err := ioutil.ReadFile(dataFilePath)
		if err != nil {
			log.Fatalf("error opening data file: %s", err)
		}

		var data PaData
		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Fatalf("error parsing data file: %s", err)
		}

		project, ok := data.Projects[currentDir]
		if !ok {
			data.Projects = make(map[string]PaProject)
			project = PaProject{}
			data.Projects[currentDir] = project
		}

		if !deleteAlias {
			alias := PaAlias{
				Name:    alias,
				Command: command,
			}

			project, err = addProjectAlias(project, alias)
			if err != nil {
				log.Fatalf("could not add alias to project: %s", err)
			}
		}

		if deleteAlias {
			project = deleteProjectAlias(project, alias)
		}

		data.Projects[currentDir] = project

		file, _ := json.MarshalIndent(data, "", "  ")
		_ = ioutil.WriteFile(dataFilePath, file, 0644)
	},
}

func addProjectAlias(project PaProject, alias PaAlias) (PaProject, error) {
	for _, existingAlias := range project.Aliases {
		if alias.Name == existingAlias.Name {
			return PaProject{}, fmt.Errorf("alias '%s' already exists in project", alias.Name)
		}
	}

	project.Aliases = append(project.Aliases, alias)
	return project, nil
}

func deleteProjectAlias(project PaProject, targetAlias string) PaProject {
	var filteredAliases []PaAlias
	for _, alias := range project.Aliases {
		if alias.Name != targetAlias {
			filteredAliases = append(filteredAliases, alias)
		}
	}

	project.Aliases = filteredAliases
	return project
}
