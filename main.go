package main

import "github.com/viniciusmuller/pa/cmd"

type Alias struct {
	Name            string
	OriginalCommand string
}

func main() {
	cmd.Execute()
}
