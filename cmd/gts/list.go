/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdList = &Command{
	UsageLine: "list [-limit n]",
	Short:     "list all entities",
	Long: `
List displays all entities in the store.

The -limit flag limits the number of results (default: 100).
Requires -path to be set to load entities.

Example:

	gts -path ./examples list -limit 50
	`,
}

var (
	listLimit int
)

func init() {
	cmdList.Run = runList
	cmdList.Flag.IntVar(&listLimit, "limit", 100, "maximum number of results")
}

func runList(cmd *Command, args []string) {
	store := newStore()
	result := store.List(listLimit)
	writeJSON(result)
}
