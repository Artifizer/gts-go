/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
)

// StoreGtsObjectNotFoundError is returned when a GTS entity is not found in the store
type StoreGtsObjectNotFoundError struct {
	EntityID string
}

func (e *StoreGtsObjectNotFoundError) Error() string {
	return fmt.Sprintf("JSON object with GTS ID '%s' not found in store", e.EntityID)
}

// StoreGtsSchemaNotFoundError is returned when a GTS schema is not found in the store
type StoreGtsSchemaNotFoundError struct {
	EntityID string
}

func (e *StoreGtsSchemaNotFoundError) Error() string {
	return fmt.Sprintf("JSON schema with GTS ID '%s' not found in store", e.EntityID)
}

// StoreGtsSchemaForInstanceNotFoundError is returned when a schema ID cannot be determined for an instance
type StoreGtsSchemaForInstanceNotFoundError struct {
	EntityID string
}

func (e *StoreGtsSchemaForInstanceNotFoundError) Error() string {
	return fmt.Sprintf("Can't determine JSON schema ID for instance with GTS ID '%s'", e.EntityID)
}

// GtsStore manages a collection of JSON entities and schemas
type GtsStore struct {
	byID   map[string]*JsonEntity
	reader GtsReader
}

// NewGtsStore creates a new GtsStore, optionally populating it from a reader
func NewGtsStore(reader GtsReader) *GtsStore {
	store := &GtsStore{
		byID:   make(map[string]*JsonEntity),
		reader: reader,
	}

	// Populate from reader if provided
	if reader != nil {
		store.populateFromReader()
	}

	return store
}

// populateFromReader loads all entities from the reader into the store
func (s *GtsStore) populateFromReader() {
	if s.reader == nil {
		return
	}

	for {
		entity := s.reader.Next()
		if entity == nil {
			break
		}
		if entity.GtsID != nil && entity.GtsID.ID != "" {
			s.byID[entity.GtsID.ID] = entity
		}
	}
}

// Register adds a JsonEntity to the store
func (s *GtsStore) Register(entity *JsonEntity) error {
	if entity.GtsID == nil || entity.GtsID.ID == "" {
		return fmt.Errorf("entity must have a valid gts_id")
	}
	s.byID[entity.GtsID.ID] = entity
	return nil
}

// RegisterSchema registers a schema with the given type ID
// This is a legacy method for backward compatibility
func (s *GtsStore) RegisterSchema(typeID string, schema map[string]any) error {
	if typeID[len(typeID)-1] != '~' {
		return fmt.Errorf("schema type_id must end with '~'")
	}

	// Parse to validate
	gtsID, err := NewGtsID(typeID)
	if err != nil {
		return err
	}

	entity := &JsonEntity{
		GtsID:    gtsID,
		Content:  schema,
		IsSchema: true,
	}

	s.byID[typeID] = entity
	return nil
}

// Get retrieves a JsonEntity by its ID
// If not found in cache, attempts to fetch from reader
func (s *GtsStore) Get(entityID string) *JsonEntity {
	// Check cache first
	if entity, ok := s.byID[entityID]; ok {
		return entity
	}

	// Try to fetch from reader
	if s.reader != nil {
		entity := s.reader.ReadByID(entityID)
		if entity != nil {
			s.byID[entityID] = entity
			return entity
		}
	}

	return nil
}

// GetSchemaContent retrieves schema content as a map (legacy method)
func (s *GtsStore) GetSchemaContent(typeID string) (map[string]any, error) {
	entity := s.Get(typeID)
	if entity == nil {
		return nil, fmt.Errorf("schema not found: %s", typeID)
	}
	if !entity.IsSchema {
		return nil, fmt.Errorf("entity is not a schema: %s", typeID)
	}
	return entity.Content, nil
}

// Items returns all entity ID and entity pairs
func (s *GtsStore) Items() map[string]*JsonEntity {
	return s.byID
}

// Count returns the number of entities in the store
func (s *GtsStore) Count() int {
	return len(s.byID)
}

// EntityInfo represents basic information about an entity
type EntityInfo struct {
	ID       string `json:"id"`
	SchemaID string `json:"schema_id"`
	IsSchema bool   `json:"is_schema"`
}

// ListResult represents the result of listing entities
type ListResult struct {
	Entities []EntityInfo `json:"entities"`
	Count    int          `json:"count"`
	Total    int          `json:"total"`
}

// List returns a list of entities up to the specified limit
func (s *GtsStore) List(limit int) *ListResult {
	total := len(s.byID)
	entities := []EntityInfo{}

	count := 0
	for id, entity := range s.byID {
		if count >= limit {
			break
		}
		entities = append(entities, EntityInfo{
			ID:       id,
			SchemaID: entity.SchemaID,
			IsSchema: entity.IsSchema,
		})
		count++
	}

	return &ListResult{
		Entities: entities,
		Count:    count,
		Total:    total,
	}
}
