/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdValidate = &Command{
	UsageLine: "validate -id <gts-id>",
	Short:     "validate an instance against its schema",
	Long: `
Validate checks an instance against its corresponding schema.

The -id flag specifies the GTS ID of the instance.
Requires -path to be set to load entities.

Example:

	gts -path ./examples validate -id gts.vendor.pkg.ns.type.v1.0
	`,
}

var (
	validateInstance string
)

func init() {
	cmdValidate.Run = runValidate
	cmdValidate.Flag.StringVar(&validateInstance, "id", "", "GTS ID of the instance")
}

func runValidate(cmd *Command, args []string) {
	if validateInstance == "" {
		cmd.Usage()
	}

	store := newStore()
	result := store.ValidateInstance(validateInstance)
	writeJSON(result)
}
