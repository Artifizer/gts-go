/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdQuery = &Command{
	UsageLine: "query -expr <expression> [-limit n]",
	Short:     "query entities using an expression",
	Long: `
Query filters entities using a GTS query expression.

The -expr flag specifies the query expression.
The -limit flag limits the number of results (default: 100).
Requires -path to be set to load entities.

Example:

	gts -path ./examples query -expr "gts.vendor.pkg.*" -limit 10
	`,
}

var (
	queryExpr  string
	queryLimit int
)

func init() {
	cmdQuery.Run = runQuery
	cmdQuery.Flag.StringVar(&queryExpr, "expr", "", "query expression")
	cmdQuery.Flag.IntVar(&queryLimit, "limit", 100, "maximum number of results")
}

func runQuery(cmd *Command, args []string) {
	if queryExpr == "" {
		cmd.Usage()
	}

	store := newStore()
	result := store.Query(queryExpr, queryLimit)
	writeJSON(result)
}
