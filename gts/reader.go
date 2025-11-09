/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

// GtsReader is an interface for reading JSON entities from various sources
type GtsReader interface {
	// Next returns the next JsonEntity or nil when exhausted
	Next() *JsonEntity

	// ReadByID reads a JsonEntity by its ID
	// Returns nil if the entity is not found
	ReadByID(entityID string) *JsonEntity

	// Reset resets the iterator to start from the beginning
	Reset()
}
