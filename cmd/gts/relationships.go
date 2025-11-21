/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

var cmdRelationships = &Command{
	UsageLine: "relationships -id <gts-id>",
	Short:     "resolve relationships for an entity",
	Long: `
Relationships builds a graph of schema relationships for an entity.

The -id flag specifies the GTS ID of the entity.
Requires -path to be set to load entities.

Example:

	gts -path ./examples relationships -id gts.vendor.pkg.ns.type.v1~
	`,
}

var (
	relationshipsID string
)

func init() {
	cmdRelationships.Run = runRelationships
	cmdRelationships.Flag.StringVar(&relationshipsID, "id", "", "GTS ID of the entity")
}

func runRelationships(cmd *Command, args []string) {
	if relationshipsID == "" {
		cmd.Usage()
	}

	store := newStore()
	result := store.BuildSchemaGraph(relationshipsID)
	writeJSON(result)
}
