/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
	"strings"
)

// JsonFile represents a JSON file containing one or more entities
type JsonFile struct {
	Path    string
	Name    string
	Content any
}

// JsonEntity represents a JSON object with extracted GTS identifiers
type JsonEntity struct {
	GtsID                 *GtsID
	SchemaID              string
	SelectedEntityField   string
	SelectedSchemaIDField string
	IsSchema              bool
	Content               map[string]any
	File                  *JsonFile
	ListSequence          *int
	Label                 string
}

// ExtractIDResult holds the result of extracting ID information from JSON content
type ExtractIDResult struct {
	ID                    string
	SchemaID              string
	SelectedEntityField   string
	SelectedSchemaIDField string
	IsSchema              bool
}

// NewJsonEntity creates a JsonEntity from JSON content using the provided config
func NewJsonEntity(content map[string]any, cfg *GtsConfig) *JsonEntity {
	return NewJsonEntityWithFile(content, cfg, nil, nil)
}

// NewJsonEntityWithFile creates a JsonEntity with file and sequence information
func NewJsonEntityWithFile(content map[string]any, cfg *GtsConfig, file *JsonFile, listSequence *int) *JsonEntity {
	if cfg == nil {
		cfg = DefaultGtsConfig()
	}

	entity := &JsonEntity{
		Content:      content,
		IsSchema:     isJSONSchema(content),
		File:         file,
		ListSequence: listSequence,
	}

	// Extract entity ID
	entityIDValue := entity.calcJSONEntityID(cfg)

	// Extract schema ID
	entity.SchemaID = entity.calcJSONSchemaID(cfg, entityIDValue)

	// If no valid GTS ID found in entity fields, use schema ID as fallback
	if entityIDValue == "" || !IsValidGtsID(entityIDValue) {
		if entity.SchemaID != "" && IsValidGtsID(entity.SchemaID) {
			entityIDValue = entity.SchemaID
		}
	}

	// Create GtsID if valid
	if entityIDValue != "" && IsValidGtsID(entityIDValue) {
		gtsID, _ := NewGtsID(entityIDValue)
		entity.GtsID = gtsID
	}

	// Set label
	entity.setLabel()

	return entity
}

// setLabel sets the entity's label based on file, sequence, or GTS ID
func (e *JsonEntity) setLabel() {
	if e.File != nil && e.ListSequence != nil {
		e.Label = fmt.Sprintf("%s#%d", e.File.Name, *e.ListSequence)
	} else if e.File != nil {
		e.Label = e.File.Name
	} else if e.GtsID != nil {
		e.Label = e.GtsID.ID
	} else {
		e.Label = ""
	}
}

// isJSONSchema checks if the content represents a JSON Schema
func isJSONSchema(content map[string]any) bool {
	if content == nil {
		return false
	}

	schemaURL, ok := content["$schema"]
	if !ok {
		return false
	}

	schemaStr, ok := schemaURL.(string)
	if !ok {
		return false
	}

	return strings.HasPrefix(schemaStr, "http://json-schema.org/") ||
		strings.HasPrefix(schemaStr, "https://json-schema.org/") ||
		strings.HasPrefix(schemaStr, "gts://") ||
		strings.HasPrefix(schemaStr, "gts.")
}

// getFieldValue retrieves a string value from content field
func (e *JsonEntity) getFieldValue(field string) string {
	if e.Content == nil {
		return ""
	}

	val, ok := e.Content[field]
	if !ok {
		return ""
	}

	strVal, ok := val.(string)
	if !ok {
		return ""
	}

	return strings.TrimSpace(strVal)
}

// firstNonEmptyField finds the first non-empty field, preferring valid GTS IDs
func (e *JsonEntity) firstNonEmptyField(fields []string) (string, string) {
	// First pass: look for valid GTS IDs
	for _, field := range fields {
		val := e.getFieldValue(field)
		if val != "" && IsValidGtsID(val) {
			return field, val
		}
	}

	// Second pass: any non-empty string
	for _, field := range fields {
		val := e.getFieldValue(field)
		if val != "" {
			return field, val
		}
	}

	return "", ""
}

// calcJSONEntityID extracts the entity ID from JSON content
func (e *JsonEntity) calcJSONEntityID(cfg *GtsConfig) string {
	field, value := e.firstNonEmptyField(cfg.EntityIDFields)
	e.SelectedEntityField = field
	return value
}

// calcJSONSchemaID extracts the schema ID from JSON content
func (e *JsonEntity) calcJSONSchemaID(cfg *GtsConfig, entityIDValue string) string {
	field, value := e.firstNonEmptyField(cfg.SchemaIDFields)
	if value != "" {
		e.SelectedSchemaIDField = field
		return value
	}

	// If no schema ID field found, try to derive from entity ID
	if entityIDValue != "" && IsValidGtsID(entityIDValue) {
		// If entity ID ends with ~, it's already a type ID
		if strings.HasSuffix(entityIDValue, "~") {
			e.SelectedSchemaIDField = e.SelectedEntityField
			return entityIDValue
		}

		// Find last ~ and return everything up to and including it
		lastTilde := strings.LastIndex(entityIDValue, "~")
		if lastTilde > 0 {
			e.SelectedSchemaIDField = e.SelectedEntityField
			return entityIDValue[:lastTilde+1]
		}
	}

	return ""
}

// ExtractID extracts GTS ID information from JSON content
func ExtractID(content map[string]any, cfg *GtsConfig) *ExtractIDResult {
	entity := NewJsonEntity(content, cfg)

	result := &ExtractIDResult{
		SchemaID:              entity.SchemaID,
		SelectedEntityField:   entity.SelectedEntityField,
		SelectedSchemaIDField: entity.SelectedSchemaIDField,
		IsSchema:              entity.IsSchema,
	}

	if entity.GtsID != nil {
		result.ID = entity.GtsID.ID
	}

	return result
}
