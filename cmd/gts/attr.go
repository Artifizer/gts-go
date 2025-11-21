/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdAttr = &Command{
	UsageLine: "attr -path <gts-id@path>",
	Short:     "get attribute value from a GTS entity",
	Long: `
Attr retrieves an attribute value from a GTS entity using path notation.

The -path flag specifies the GTS ID with attribute path (e.g., gts.x.y.z.v1.0@field.subfield).
Requires the global -path flag to be set to load entities.

Example:

	gts -path ./examples attr -path gts.vendor.pkg.ns.type.v1.0@name
	`,
}

var (
	attrPath string
)

func init() {
	cmdAttr.Run = runAttr
	cmdAttr.Flag.StringVar(&attrPath, "path", "", "GTS ID with attribute path")
}

func runAttr(cmd *Command, args []string) {
	if attrPath == "" {
		cmd.Usage()
	}

	store := newStore()
	result := store.GetAttribute(attrPath)
	writeJSON(result)
}
