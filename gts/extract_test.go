/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"testing"
)

// TestExtractID_BasicEntityID tests extracting entity ID from JSON content
func TestExtractID_BasicEntityID(t *testing.T) {
	tests := []struct {
		name          string
		content       map[string]any
		expectedID    string
		expectedField string
	}{
		{
			name: "Extract from gtsId field",
			content: map[string]any{
				"gtsId": "gts.vendor.package.namespace.type.v0",
				"name":  "Test Entity",
			},
			expectedID:    "gts.vendor.package.namespace.type.v0",
			expectedField: "gtsId",
		},
		{
			name: "Extract from $id field",
			content: map[string]any{
				"$id":  "gts.vendor.package.namespace.type.v1",
				"name": "Test Entity",
			},
			expectedID:    "gts.vendor.package.namespace.type.v1",
			expectedField: "$id",
		},
		{
			name: "Extract from id field (fallback)",
			content: map[string]any{
				"id":   "gts.vendor.package.namespace.type.v2",
				"name": "Test Entity",
			},
			expectedID:    "gts.vendor.package.namespace.type.v2",
			expectedField: "id",
		},
		{
			name: "Priority: gtsId over id",
			content: map[string]any{
				"gtsId": "gts.vendor1.package.namespace.type.v0",
				"id":    "gts.vendor2.package.namespace.type.v0",
			},
			expectedID:    "gts.vendor1.package.namespace.type.v0",
			expectedField: "gtsId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractID(tt.content, nil)
			if result.ID != tt.expectedID {
				t.Errorf("Expected ID %q, got %q", tt.expectedID, result.ID)
			}
			if result.SelectedEntityField == nil || *result.SelectedEntityField != tt.expectedField {
				var got string
				if result.SelectedEntityField != nil {
					got = *result.SelectedEntityField
				}
				t.Errorf("Expected field %q, got %q", tt.expectedField, got)
			}
		})
	}
}

// TestExtractID_SchemaID tests extracting schema ID from JSON content
func TestExtractID_SchemaID(t *testing.T) {
	tests := []struct {
		name                string
		content             map[string]any
		expectedSchemaID    string
		expectedSchemaField string
	}{
		{
			name: "Extract from $schema field",
			content: map[string]any{
				"gtsId":   "gts.vendor.package.namespace.type.v0.1",
				"$schema": "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchemaID:    "gts.vendor.package.namespace.type.v0~",
			expectedSchemaField: "$schema",
		},
		{
			name: "Extract from gtsTid field",
			content: map[string]any{
				"gtsId":  "gts.vendor.package.namespace.type.v0.1",
				"gtsTid": "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchemaID:    "gts.vendor.package.namespace.type.v0~",
			expectedSchemaField: "gtsTid",
		},
		{
			name: "Derive from entity ID with tilde",
			content: map[string]any{
				"gtsId": "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchemaID:    "gts.vendor.package.namespace.type.v0~",
			expectedSchemaField: "", // Not set - entity ID itself is a type
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractID(tt.content, nil)
			if result.SchemaID != tt.expectedSchemaID {
				t.Errorf("Expected SchemaID %q, got %q", tt.expectedSchemaID, result.SchemaID)
			}

			// Handle both empty string expectation and actual value
			var got string
			if result.SelectedSchemaIDField != nil {
				got = *result.SelectedSchemaIDField
			}
			if got != tt.expectedSchemaField {
				t.Errorf("Expected schema field %q, got %q", tt.expectedSchemaField, got)
			}
		})
	}
}

// TestExtractID_IsSchema tests detecting JSON Schema documents
func TestExtractID_IsSchema(t *testing.T) {
	tests := []struct {
		name           string
		content        map[string]any
		expectedSchema bool
	}{
		{
			name: "JSON Schema with http://json-schema.org/",
			content: map[string]any{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"gtsId":   "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchema: true,
		},
		{
			name: "JSON Schema with https://json-schema.org/",
			content: map[string]any{
				"$schema": "https://json-schema.org/draft/2020-12/schema",
				"gtsId":   "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchema: true,
		},
		{
			name: "GTS Schema with gts:// prefix",
			content: map[string]any{
				"$schema": "gts://vendor.package.namespace.type.v0~",
			},
			expectedSchema: true,
		},
		{
			name: "GTS Schema with gts. prefix",
			content: map[string]any{
				"$schema": "gts.vendor.package.namespace.type.v0~",
			},
			expectedSchema: true,
		},
		{
			name: "Regular entity (not a schema)",
			content: map[string]any{
				"gtsId": "gts.vendor.package.namespace.type.v0.1",
				"name":  "Test Entity",
			},
			expectedSchema: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractID(tt.content, nil)
			if result.IsSchema != tt.expectedSchema {
				t.Errorf("Expected IsSchema %v, got %v", tt.expectedSchema, result.IsSchema)
			}
		})
	}
}

// TestExtractID_CustomConfig tests using custom configuration
func TestExtractID_CustomConfig(t *testing.T) {
	customCfg := &GtsConfig{
		EntityIDFields: []string{"customId", "id"},
		SchemaIDFields: []string{"customType", "type"},
	}

	content := map[string]any{
		"customId":   "gts.vendor.package.namespace.type.v0",
		"customType": "gts.vendor.package.namespace.type.v0~",
	}

	result := ExtractID(content, customCfg)

	if result.ID != "gts.vendor.package.namespace.type.v0" {
		t.Errorf("Expected ID from customId field, got %q", result.ID)
	}
	if result.SelectedEntityField == nil || *result.SelectedEntityField != "customId" {
		var got string
		if result.SelectedEntityField != nil {
			got = *result.SelectedEntityField
		}
		t.Errorf("Expected customId field, got %q", got)
	}
	if result.SchemaID != "gts.vendor.package.namespace.type.v0~" {
		t.Errorf("Expected SchemaID from customType field, got %q", result.SchemaID)
	}
	if result.SelectedSchemaIDField == nil || *result.SelectedSchemaIDField != "customType" {
		var got string
		if result.SelectedSchemaIDField != nil {
			got = *result.SelectedSchemaIDField
		}
		t.Errorf("Expected customType field, got %q", got)
	}
}

// TestExtractID_NoValidID tests extraction when no valid GTS ID is found
func TestExtractID_NoValidID(t *testing.T) {
	content := map[string]any{
		"name":        "Test Entity",
		"description": "No GTS ID here",
	}

	result := ExtractID(content, nil)

	if result.ID != "" {
		t.Errorf("Expected empty ID, got %q", result.ID)
	}
	if result.SelectedEntityField != nil {
		t.Errorf("Expected nil SelectedEntityField, got %q", *result.SelectedEntityField)
	}
}

// TestExtractID_InvalidIDInField tests extraction when field contains invalid GTS ID
func TestExtractID_InvalidIDInField(t *testing.T) {
	content := map[string]any{
		"gtsId": "not-a-valid-gts-id",
		"id":    "gts.vendor.package.namespace.type.v0",
	}

	result := ExtractID(content, nil)

	// Should fallback to the "id" field which has a valid GTS ID
	if result.ID != "gts.vendor.package.namespace.type.v0" {
		t.Errorf("Expected fallback to valid ID, got %q", result.ID)
	}
}

// TestExtractID_SchemaIDFallback tests schema ID extraction using entity ID as fallback
func TestExtractID_SchemaIDFallback(t *testing.T) {
	content := map[string]any{
		"$schema": "gts.vendor.package.namespace.type.v0~",
	}

	result := ExtractID(content, nil)

	// When no entity ID field is found, schema ID should be used as entity ID too
	if result.ID != "gts.vendor.package.namespace.type.v0~" {
		t.Errorf("Expected ID to fallback to schema ID, got %q", result.ID)
	}
	if result.SchemaID != "gts.vendor.package.namespace.type.v0~" {
		t.Errorf("Expected SchemaID, got %q", result.SchemaID)
	}
}

// =============================================================================
// Tests for URI prefix "gts://" in JSON Schema $id field
// The gts:// prefix is used in JSON Schema for URI compatibility.
// GtsEntity strips it when parsing so the GtsID works with normal "gts." format.
// =============================================================================

// TestExtractID_GtsURIPrefix_InDollarIdField tests that gts:// prefix is stripped from $id field
func TestExtractID_GtsURIPrefix_InDollarIdField(t *testing.T) {
	content := map[string]any{
		"$id":     "gts://gts.vendor.package.namespace.type.v1.0~",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type":    "object",
	}

	result := ExtractID(content, nil)

	// The gts:// prefix should be stripped from the $id field
	if result.ID != "gts.vendor.package.namespace.type.v1.0~" {
		t.Errorf("Expected ID without gts:// prefix %q, got %q", "gts.vendor.package.namespace.type.v1.0~", result.ID)
	}
	if !result.IsSchema {
		t.Errorf("Expected IsSchema to be true")
	}
}

// TestExtractID_GtsURIPrefix_NotStrippedFromOtherFields tests that gts:// prefix is NOT stripped from non-$id fields
func TestExtractID_GtsURIPrefix_NotStrippedFromOtherFields(t *testing.T) {
	// gts:// prefix in non-$id field should NOT be stripped (and results in invalid GTS ID)
	content := map[string]any{
		"id": "gts://gts.vendor.package.namespace.type.v1.0",
	}

	result := ExtractID(content, nil)

	// The "id" field is not $id, so the gts:// prefix is NOT stripped
	// The value "gts://gts.vendor..." is not a valid GTS ID
	if result.ID != "" {
		t.Errorf("Expected empty ID (gts:// prefix in 'id' field should not be stripped), got %q", result.ID)
	}
}

// TestExtractID_GtsColonPrefix_NotValid tests that gts: (without //) is NOT a valid prefix
func TestExtractID_GtsColonPrefix_NotValid(t *testing.T) {
	// "gts:" (without //) is NOT a valid prefix - only "gts://" is valid
	content := map[string]any{
		"$id":     "gts:gts.vendor.package.namespace.type.v1.0~",
		"$schema": "http://json-schema.org/draft-07/schema#",
	}

	result := ExtractID(content, nil)

	// With "gts:" prefix (not "gts://"), the ID is not stripped and won't be valid
	// The entity should NOT have a valid GTS ID
	if result.ID != "" {
		t.Errorf("Expected empty ID (gts: prefix without // should not be stripped), got %q", result.ID)
	}
}

// TestExtractID_GtsURIPrefix_WithoutPrefix tests that IDs without prefix still work
func TestExtractID_GtsURIPrefix_WithoutPrefix(t *testing.T) {
	// IDs without gts:// prefix should work as before
	content := map[string]any{
		"$id":     "gts.vendor.package.namespace.type.v1.0~",
		"$schema": "http://json-schema.org/draft-07/schema#",
	}

	result := ExtractID(content, nil)

	if result.ID != "gts.vendor.package.namespace.type.v1.0~" {
		t.Errorf("Expected ID %q, got %q", "gts.vendor.package.namespace.type.v1.0~", result.ID)
	}
}
