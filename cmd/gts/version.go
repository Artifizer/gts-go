/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

import (
	"fmt"
	"runtime/debug"
)

var cmdVersion = &Command{
	UsageLine: "version",
	Short:     "print GTS version",
	Long:      `Version prints the GTS version.`,
}

func init() {
	cmdVersion.Run = runVersion
}

func runVersion(cmd *Command, args []string) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("gts version unknown")
		return
	}

	fmt.Printf("gts version %s\n", info.Main.Version)
	if verbose > 0 {
		fmt.Printf("go version %s\n", info.GoVersion)
		fmt.Printf("path %s\n", info.Path)
	}
}
