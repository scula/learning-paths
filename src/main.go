package main

import (
	"os"

	"github.com/amazemind-io/learning-paths/learning_paths"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cmd := argsWithoutProg[0]
	path_file := argsWithoutProg[1]

	learning_paths.ValidateCmd(cmd, path_file, argsWithoutProg)

}
