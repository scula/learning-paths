package main

import (
	"os"

	"amazemind.io/learning-paths/validator"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cmd := argsWithoutProg[0]
	path_file := argsWithoutProg[1]

	validator.ValidateCmd(cmd, path_file, argsWithoutProg)

}
