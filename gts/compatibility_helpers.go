/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
	"sort"
	"strings"
)

// Helper functions for compatibility checking

// getPropertiesMap safely extracts properties as map[string]any
func getPropertiesMap(schema map[string]any) map[string]any {
	if props, ok := schema["properties"].(map[string]any); ok {
		return props
	}
	return make(map[string]any)
}

// getRequiredSet safely extracts required fields as a set
func getRequiredSet(schema map[string]any) map[string]bool {
	set := make(map[string]bool)
	if req, ok := schema["required"].([]any); ok {
		for _, item := range req {
			if str, ok := item.(string); ok {
				set[str] = true
			}
		}
	}
	return set
}

// getString safely extracts a string value from map
func getString(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getMap safely extracts a map value
func getMap(m map[string]any, key string) map[string]any {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]any); ok {
			return mapVal
		}
	}
	return nil
}

// getNumber safely extracts a number (int or float) value
func getNumber(m map[string]any, key string) *float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return &v
		case int:
			f := float64(v)
			return &f
		case int64:
			f := float64(v)
			return &f
		}
	}
	return nil
}

// getStringSlice safely extracts a string slice from enum
func getStringSlice(m map[string]any, key string) []string {
	result := []string{}
	if val, ok := m[key]; ok {
		if slice, ok := val.([]any); ok {
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
		}
	}
	return result
}

// getKeys returns all keys from a map as a set
func getKeys(m map[string]any) map[string]bool {
	keys := make(map[string]bool)
	for k := range m {
		keys[k] = true
	}
	return keys
}

// setDifference returns elements in a that are not in b
func setDifference(a, b map[string]bool) []string {
	diff := []string{}
	for k := range a {
		if !b[k] {
			diff = append(diff, k)
		}
	}
	sort.Strings(diff)
	return diff
}

// setIntersection returns elements that exist in both a and b
func setIntersection(a, b map[string]bool) []string {
	intersection := []string{}
	for k := range a {
		if b[k] {
			intersection = append(intersection, k)
		}
	}
	sort.Strings(intersection)
	return intersection
}

// joinStrings joins string slice with comma separator
func joinStrings(strs []string) string {
	return strings.Join(strs, ", ")
}

// stringSliceToSet converts string slice to set
func stringSliceToSet(slice []string) map[string]bool {
	set := make(map[string]bool)
	for _, s := range slice {
		set[s] = true
	}
	return set
}

// setToString converts a set to a sorted comma-separated string
func setToString(set map[string]bool) string {
	keys := []string{}
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

// floatToString converts float to string
func floatToString(f float64) string {
	// Remove trailing zeros and decimal point if integer
	s := fmt.Sprintf("%.10f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
