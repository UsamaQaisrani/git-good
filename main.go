package main

import (
	"fmt"
	"os"
	"usamaqaisrani/git-good/plumbing"
	"usamaqaisrani/git-good/porcelain"
)

var commands = map[string]struct{}{
	"init":        {},
	"hash-object": {},
	"add":         {},
	"ls-files": {},
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command provided.")
		return
	}

	if _, ok := commands[args[0]]; !ok {
		fmt.Printf("%s command not defined\n", args[0])
		return
	}

	switch args[0] {
	case "init":
		porcelain.Init()
	case "hash-object":
		hashObject(args)
	case "add":
		stage(args)
	case "ls-files":
		readIndex()
	}
}

func hashObject(args []string) {
	if len(args) < 2 {
		fmt.Print("Missing path to file name.")
		return
	}
	fmt.Println("Creating hash of the file")
	content, err := plumbing.ReadFile(args[1])
	if err != nil {
		fmt.Printf("Error occured while reading contents of %s: %s", args[1], err)
		return
	}
	plumbing.HashFile(content)
}

func stage(args []string) {
	if len(args) < 2 {
		fmt.Print("Missing path to file name.")
		return
	}
	fmt.Println("Creating hash of the file")
	porcelain.Stage(args[1])
}

func readIndex() {
	plumbing.ReadIndex()
}
