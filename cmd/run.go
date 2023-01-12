package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an existing alias",
	Long:  "Runs an alias. Only aliases defined in the current directory are considered",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requestedAlias := args[0]

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current directory: %s", err)
		}

		data, err := readDataFile()
		if err != nil {
			log.Fatalf("could not read data file: %s", err)
		}

		paAlias, err := findProjectAlias(data.Projects[currentDir], requestedAlias)
		if err != nil {
			log.Fatalf("could not find alias '%s': %s", requestedAlias, err)
		}

		splittedCommand := strings.Split(paAlias.Command, " ")
		command, args = splittedCommand[0], splittedCommand[1:]

		shellCommand := exec.Command(command)
		shellCommand.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		shellCommand.Stdout = &out

		err = shellCommand.Run()
		if err != nil {
			log.Fatalf("error running command: %s", err)
		}
	},
}

func findProjectAlias(project PaProject, targetAlias string) (PaAlias, error) {
	for _, alias := range project.Aliases {
		if alias.Name == targetAlias {
			return alias, nil
		}
	}
	return PaAlias{}, fmt.Errorf("alias '%s' not found", targetAlias)
}
