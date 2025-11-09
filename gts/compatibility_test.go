/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"testing"
)

func TestCheckCompatibility_BackwardCompatible_AddOptional(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.event.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "timestamp", "userId"},
		"properties": map[string]any{
			"eventId":   map[string]any{"type": "string"},
			"timestamp": map[string]any{"type": "string", "format": "date-time"},
			"userId":    map[string]any{"type": "string"},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (adds optional field with default)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.event.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "timestamp", "userId"},
		"properties": map[string]any{
			"eventId":   map[string]any{"type": "string"},
			"timestamp": map[string]any{"type": "string", "format": "date-time"},
			"userId":    map[string]any{"type": "string"},
			"metadata": map[string]any{
				"type":    "object",
				"default": map[string]any{},
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility
	result := store.CheckCompatibility("gts.x.core.compat.event.v1.0~", "gts.x.core.compat.event.v1.1~")

	if !result.IsBackwardCompatible {
		t.Errorf("Expected backward compatible, got false. Errors: %v", result.BackwardErrors)
	}
	if result.OldID != "gts.x.core.compat.event.v1.0~" {
		t.Errorf("Expected old ID, got: %s", result.OldID)
	}
	if result.NewID != "gts.x.core.compat.event.v1.1~" {
		t.Errorf("Expected new ID, got: %s", result.NewID)
	}
}

func TestCheckCompatibility_BackwardIncompatible_AddRequired(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.breaking.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (adds required field - breaking!)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.breaking.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "newRequiredField"},
		"properties": map[string]any{
			"eventId":          map[string]any{"type": "string"},
			"newRequiredField": map[string]any{"type": "string"},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: should NOT be backward compatible
	result := store.CheckCompatibility("gts.x.core.compat.breaking.v1.0~", "gts.x.core.compat.breaking.v1.1~")

	if result.IsBackwardCompatible {
		t.Error("Expected backward incompatible, got true")
	}
	if len(result.BackwardErrors) == 0 {
		t.Error("Expected backward errors, got none")
	}
}

func TestCheckCompatibility_ForwardCompatible_OpenModel(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema with additionalProperties: true
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.forward.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
		},
		"additionalProperties": true,
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (adds new field)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.forward.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "newField"},
		"properties": map[string]any{
			"eventId":  map[string]any{"type": "string"},
			"newField": map[string]any{"type": "string"},
		},
		"additionalProperties": true,
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: should be forward compatible
	result := store.CheckCompatibility("gts.x.core.compat.forward.v1.0~", "gts.x.core.compat.forward.v1.1~")

	if !result.IsForwardCompatible {
		t.Errorf("Expected forward compatible, got false. Errors: %v", result.ForwardErrors)
	}
}

func TestCheckCompatibility_ForwardIncompatible_RemoveRequired(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.fwd_break.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "importantField"},
		"properties": map[string]any{
			"eventId":        map[string]any{"type": "string"},
			"importantField": map[string]any{"type": "string"},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (removes required field)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.fwd_break.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: should NOT be forward compatible
	result := store.CheckCompatibility("gts.x.core.compat.fwd_break.v1.0~", "gts.x.core.compat.fwd_break.v1.1~")

	if result.IsForwardCompatible {
		t.Error("Expected forward incompatible, got true")
	}
	if len(result.ForwardErrors) == 0 {
		t.Error("Expected forward errors, got none")
	}
}

func TestCheckCompatibility_FullyCompatible(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema (open model)
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.full.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
		},
		"additionalProperties": true,
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (adds optional field with default)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.full.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
			"optionalField": map[string]any{
				"type":    "string",
				"default": "default_value",
			},
		},
		"additionalProperties": true,
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: should be fully compatible
	result := store.CheckCompatibility("gts.x.core.compat.full.v1.0~", "gts.x.core.compat.full.v1.1~")

	if !result.IsBackwardCompatible {
		t.Errorf("Expected backward compatible, got false. Errors: %v", result.BackwardErrors)
	}
	if !result.IsForwardCompatible {
		t.Errorf("Expected forward compatible, got false. Errors: %v", result.ForwardErrors)
	}
	if !result.IsFullyCompatible {
		t.Error("Expected fully compatible, got false")
	}
}

func TestCheckCompatibility_TypeChange(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.typechange.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "count"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
			"count":   map[string]any{"type": "number"},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (changes count type from number to string)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.typechange.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "count"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
			"count":   map[string]any{"type": "string"},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: should be incompatible both ways
	result := store.CheckCompatibility("gts.x.core.compat.typechange.v1.0~", "gts.x.core.compat.typechange.v1.1~")

	if result.IsBackwardCompatible {
		t.Error("Expected backward incompatible due to type change")
	}
	if result.IsForwardCompatible {
		t.Error("Expected forward incompatible due to type change")
	}
	if result.IsFullyCompatible {
		t.Error("Expected fully incompatible due to type change")
	}
}

func TestCheckCompatibility_EnumExpansion(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 schema with enum
	v10Schema := map[string]any{
		"$id":      "gts.x.core.compat.enum.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "status"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
			"status": map[string]any{
				"type": "string",
				"enum": []any{"active", "inactive"},
			},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 schema (adds enum value)
	v11Schema := map[string]any{
		"$id":      "gts.x.core.compat.enum.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"eventId", "status"},
		"properties": map[string]any{
			"eventId": map[string]any{"type": "string"},
			"status": map[string]any{
				"type": "string",
				"enum": []any{"active", "inactive", "pending"},
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility: forward compatible, not backward
	result := store.CheckCompatibility("gts.x.core.compat.enum.v1.0~", "gts.x.core.compat.enum.v1.1~")

	if result.IsBackwardCompatible {
		t.Error("Expected backward incompatible due to enum expansion")
	}
	if !result.IsForwardCompatible {
		t.Errorf("Expected forward compatible, got false. Errors: %v", result.ForwardErrors)
	}
	if result.IsFullyCompatible {
		t.Error("Expected not fully compatible")
	}
}

func TestCheckCompatibility_NestedObjectChanges(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 with nested object
	v10Schema := map[string]any{
		"$id":      "gts.x.core.nested_compat.order.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"orderId", "customer"},
		"properties": map[string]any{
			"orderId": map[string]any{"type": "string"},
			"customer": map[string]any{
				"type":     "object",
				"required": []any{"customerId", "name"},
				"properties": map[string]any{
					"customerId": map[string]any{"type": "string"},
					"name":       map[string]any{"type": "string"},
				},
			},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 with additional nested field
	v11Schema := map[string]any{
		"$id":      "gts.x.core.nested_compat.order.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"orderId", "customer"},
		"properties": map[string]any{
			"orderId": map[string]any{"type": "string"},
			"customer": map[string]any{
				"type":     "object",
				"required": []any{"customerId", "name"},
				"properties": map[string]any{
					"customerId": map[string]any{"type": "string"},
					"name":       map[string]any{"type": "string"},
					"email":      map[string]any{"type": "string"},
				},
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility
	result := store.CheckCompatibility("gts.x.core.nested_compat.order.v1.0~", "gts.x.core.nested_compat.order.v1.1~")

	if !result.IsBackwardCompatible {
		t.Errorf("Expected backward compatible for nested optional field. Errors: %v", result.BackwardErrors)
	}
}

func TestCheckCompatibility_ConstraintRelaxation(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 with strict constraints
	v10Schema := map[string]any{
		"$id":      "gts.x.core.constraints.product.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"productId", "price"},
		"properties": map[string]any{
			"productId": map[string]any{"type": "string"},
			"price": map[string]any{
				"type":    "number",
				"minimum": 0,
				"maximum": 1000,
			},
			"name": map[string]any{
				"type":      "string",
				"minLength": 3,
				"maxLength": 50,
			},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 with relaxed constraints
	v11Schema := map[string]any{
		"$id":      "gts.x.core.constraints.product.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"productId", "price"},
		"properties": map[string]any{
			"productId": map[string]any{"type": "string"},
			"price": map[string]any{
				"type":    "number",
				"minimum": 0,
				"maximum": 10000,
			},
			"name": map[string]any{
				"type":      "string",
				"minLength": 1,
				"maxLength": 100,
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility - should be backward compatible
	result := store.CheckCompatibility("gts.x.core.constraints.product.v1.0~", "gts.x.core.constraints.product.v1.1~")

	if !result.IsBackwardCompatible {
		t.Errorf("Expected backward compatible for constraint relaxation. Errors: %v", result.BackwardErrors)
	}
}

func TestCheckCompatibility_ConstraintTightening(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 with loose constraints
	v10Schema := map[string]any{
		"$id":      "gts.x.core.tight.item.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"itemId", "quantity"},
		"properties": map[string]any{
			"itemId": map[string]any{"type": "string"},
			"quantity": map[string]any{
				"type":    "integer",
				"minimum": 1,
				"maximum": 1000,
			},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 with tighter constraints
	v11Schema := map[string]any{
		"$id":      "gts.x.core.tight.item.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"itemId", "quantity"},
		"properties": map[string]any{
			"itemId": map[string]any{"type": "string"},
			"quantity": map[string]any{
				"type":    "integer",
				"minimum": 1,
				"maximum": 100,
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility - should NOT be backward compatible
	result := store.CheckCompatibility("gts.x.core.tight.item.v1.0~", "gts.x.core.tight.item.v1.1~")

	if result.IsBackwardCompatible {
		t.Error("Expected backward incompatible for constraint tightening")
	}
	if len(result.BackwardErrors) == 0 {
		t.Error("Expected backward errors for constraint tightening")
	}
}

func TestCheckCompatibility_ArrayItemSchemaChange(t *testing.T) {
	store := NewGtsStore(nil)

	// Register v1.0 with simple array items
	v10Schema := map[string]any{
		"$id":      "gts.x.core.array_compat.list.v1.0~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"listId", "items"},
		"properties": map[string]any{
			"listId": map[string]any{"type": "string"},
			"items": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":     "object",
					"required": []any{"id", "value"},
					"properties": map[string]any{
						"id":    map[string]any{"type": "string"},
						"value": map[string]any{"type": "number"},
					},
				},
			},
		},
	}
	v10Entity := NewJsonEntity(v10Schema, DefaultGtsConfig())
	if err := store.Register(v10Entity); err != nil {
		t.Fatalf("Failed to register v1.0 schema: %v", err)
	}

	// Register v1.1 with additional array item field
	v11Schema := map[string]any{
		"$id":      "gts.x.core.array_compat.list.v1.1~",
		"$schema":  "http://json-schema.org/draft-07/schema#",
		"type":     "object",
		"required": []any{"listId", "items"},
		"properties": map[string]any{
			"listId": map[string]any{"type": "string"},
			"items": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":     "object",
					"required": []any{"id", "value"},
					"properties": map[string]any{
						"id":    map[string]any{"type": "string"},
						"value": map[string]any{"type": "number"},
						"label": map[string]any{"type": "string"},
					},
				},
			},
		},
	}
	v11Entity := NewJsonEntity(v11Schema, DefaultGtsConfig())
	if err := store.Register(v11Entity); err != nil {
		t.Fatalf("Failed to register v1.1 schema: %v", err)
	}

	// Check compatibility
	result := store.CheckCompatibility("gts.x.core.array_compat.list.v1.0~", "gts.x.core.array_compat.list.v1.1~")

	if !result.IsBackwardCompatible {
		t.Errorf("Expected backward compatible for array item optional field. Errors: %v", result.BackwardErrors)
	}
}

func TestCheckCompatibility_EntityNotFound(t *testing.T) {
	store := NewGtsStore(nil)

	// Check compatibility with non-existent schemas
	result := store.CheckCompatibility("gts.x.nonexistent.schema.v1.0~", "gts.x.nonexistent.schema.v1.1~")

	if result.IsBackwardCompatible || result.IsForwardCompatible {
		t.Error("Expected incompatible for non-existent schemas")
	}
	if len(result.BackwardErrors) == 0 {
		t.Error("Expected backward errors for non-existent schemas")
	}
	if result.BackwardErrors[0] != "Schema not found" {
		t.Errorf("Expected 'Schema not found' error, got: %s", result.BackwardErrors[0])
	}
}

func TestInferDirection(t *testing.T) {
	tests := []struct {
		name     string
		fromID   string
		toID     string
		expected string
	}{
		{
			name:     "Up direction (v1.0 to v1.1)",
			fromID:   "gts.x.core.schema.test.v1.0~",
			toID:     "gts.x.core.schema.test.v1.1~",
			expected: "up",
		},
		{
			name:     "Down direction (v1.5 to v1.2)",
			fromID:   "gts.x.core.schema.test.v1.5~",
			toID:     "gts.x.core.schema.test.v1.2~",
			expected: "down",
		},
		{
			name:     "None direction (same version)",
			fromID:   "gts.x.core.schema.test.v1.0~",
			toID:     "gts.x.core.schema.test.v1.0~",
			expected: "none",
		},
		{
			name:     "Unknown direction (invalid ID)",
			fromID:   "invalid",
			toID:     "also-invalid",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferDirection(tt.fromID, tt.toID)
			if result != tt.expected {
				t.Errorf("Expected direction %s, got %s", tt.expected, result)
			}
		})
	}
}
