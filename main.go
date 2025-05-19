package main

import (
	"fmt"
	"os"

	"github.com/EmreKb/gpw/pkg/analyzer"
	"github.com/EmreKb/gpw/pkg/fshelper"
	"github.com/EmreKb/gpw/pkg/graphview"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gpw <path>")
		os.Exit(1)
	}

	pathArg := os.Args[1]
	root, err := fshelper.SubTree(pathArg)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	anlyze, err := analyzer.Analyze(root)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = graphview.View(anlyze, "output.png")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
