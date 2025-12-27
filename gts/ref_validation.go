/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
	"strings"
)

// RefValidationError represents a validation error for $ref values
type RefValidationError struct {
	FieldPath string
	RefValue  string
	Reason    string
}

func (e *RefValidationError) Error() string {
	return fmt.Sprintf("$ref validation failed for field '%s': %s", e.FieldPath, e.Reason)
}

// RefValidator validates $ref constraints in GTS schemas
type RefValidator struct {
}

// NewRefValidator creates a new $ref validator
func NewRefValidator() *RefValidator {
	return &RefValidator{}
}

// ValidateSchemaRefs validates all $ref values in a schema
func (v *RefValidator) ValidateSchemaRefs(schema map[string]interface{}, schemaPath string) []*RefValidationError {
	var errors []*RefValidationError
	v.visitSchemaForRefs(schema, schemaPath, &errors)
	return errors
}

// visitSchemaForRefs recursively visits schema nodes to find and validate $ref values
func (v *RefValidator) visitSchemaForRefs(schema map[string]interface{}, path string, errors *[]*RefValidationError) {
	if schema == nil {
		return
	}

	// Check for $ref field
	if refValue, hasRef := schema["$ref"]; hasRef {
		refPath := "$ref"
		if path != "" {
			refPath = path + "/$ref"
		}
		if err := v.validateRef(refValue, refPath); err != nil {
			*errors = append(*errors, err)
		}
	}

	// Recurse into nested structures
	for key, value := range schema {
		if key == "$ref" {
			continue // Already processed above
		}

		nestedPath := key
		if path != "" {
			nestedPath = path + "/" + key
		}

		switch val := value.(type) {
		case map[string]interface{}:
			v.visitSchemaForRefs(val, nestedPath, errors)
		case []interface{}:
			for idx, item := range val {
				if itemMap, ok := item.(map[string]interface{}); ok {
					v.visitSchemaForRefs(itemMap, fmt.Sprintf("%s[%d]", nestedPath, idx), errors)
				}
			}
		}
	}
}

// validateRef validates a single $ref value according to GTS specification
func (v *RefValidator) validateRef(refValue interface{}, fieldPath string) *RefValidationError {
	refStr, ok := refValue.(string)
	if !ok {
		return &RefValidationError{
			FieldPath: fieldPath,
			RefValue:  fmt.Sprintf("%v", refValue),
			Reason:    fmt.Sprintf("$ref value must be a string, got %T", refValue),
		}
	}

	refStr = strings.TrimSpace(refStr)
	if refStr == "" {
		return &RefValidationError{
			FieldPath: fieldPath,
			RefValue:  refStr,
			Reason:    "$ref value cannot be empty",
		}
	}

	// $ref must use gts:// URI format for GTS references

	// Case 1: Local refs (JSON Pointer) - must start with #
	if strings.HasPrefix(refStr, "#") {
		return nil // Valid local reference
	}

	// Case 2: GTS refs - must use gts:// prefix with valid GTS ID
	if strings.HasPrefix(refStr, "gts://") {
		// Strip prefix and validate the GTS ID
		gtsID := strings.TrimPrefix(refStr, GtsURIPrefix)
		if !IsValidGtsID(gtsID) {
			return &RefValidationError{
				FieldPath: fieldPath,
				RefValue:  refStr,
				Reason:    fmt.Sprintf("contains invalid GTS identifier '%s'", gtsID),
			}
		}
		return nil // Valid GTS URI reference
	}

	// Case 3: Invalid formats

	// Bare GTS ID (missing gts:// prefix)
	if strings.HasPrefix(refStr, "gts.") && IsValidGtsID(refStr) {
		return &RefValidationError{
			FieldPath: fieldPath,
			RefValue:  refStr,
			Reason:    "must be a local ref (starting with '#') or a GTS URI (starting with 'gts://')",
		}
	}

	// HTTP/HTTPS URIs
	if strings.HasPrefix(refStr, "http://") || strings.HasPrefix(refStr, "https://") {
		return &RefValidationError{
			FieldPath: fieldPath,
			RefValue:  refStr,
			Reason:    "must be a local ref (starting with '#') or a GTS URI (starting with 'gts://')",
		}
	}

	// Any other format
	return &RefValidationError{
		FieldPath: fieldPath,
		RefValue:  refStr,
		Reason:    "must be a local ref (starting with '#') or a GTS URI (starting with 'gts://')",
	}
}
