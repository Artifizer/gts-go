/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdCompatibility = &Command{
	UsageLine: "compatibility -old <old-schema-id> -new <new-schema-id>",
	Short:     "check compatibility between two schemas",
	Long: `
Compatibility checks whether two schema versions are compatible.

The -old flag specifies the old schema GTS ID.
The -new flag specifies the new schema GTS ID.
Requires -path to be set to load entities.

Example:

	gts -path ./examples compatibility -old gts.vendor.pkg.ns.type.v1~ -new gts.vendor.pkg.ns.type.v2~
	`,
}

var (
	compatOld string
	compatNew string
)

func init() {
	cmdCompatibility.Run = runCompatibility
	cmdCompatibility.Flag.StringVar(&compatOld, "old", "", "old schema GTS ID")
	cmdCompatibility.Flag.StringVar(&compatNew, "new", "", "new schema GTS ID")
}

func runCompatibility(cmd *Command, args []string) {
	if compatOld == "" || compatNew == "" {
		cmd.Usage()
	}

	store := newStore()
	result := store.CheckCompatibility(compatOld, compatNew)
	writeJSON(result)
}
