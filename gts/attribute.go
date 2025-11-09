/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
	"strconv"
	"strings"
)

// AttributeResult represents the result of attribute path resolution
type AttributeResult struct {
	GtsID           string   `json:"gts_id"`
	Path            string   `json:"path"`
	Value           any      `json:"value,omitempty"`
	Resolved        bool     `json:"resolved"`
	Error           string   `json:"error,omitempty"`
	AvailableFields []string `json:"available_fields,omitempty"`
}

// GetAttribute retrieves an attribute value from an entity using a path selector
// Format: "gts_id@path.to.field" or "gts_id@array[0].field"
// see gts-python ops.py attr method
func (s *GtsStore) GetAttribute(gtsWithPath string) *AttributeResult {
	// Split GTS ID from attribute path
	gtsID, path := splitAtPath(gtsWithPath)

	// Check if @ symbol was provided
	if path == "" {
		return &AttributeResult{
			GtsID:    gtsID,
			Path:     "",
			Resolved: false,
			Error:    "Attribute selector requires '@path' in the identifier",
		}
	}

	// Get entity from store
	entity := s.Get(gtsID)
	if entity == nil {
		return &AttributeResult{
			GtsID:    gtsID,
			Path:     path,
			Resolved: false,
			Error:    fmt.Sprintf("Entity not found: %s", gtsID),
		}
	}

	// Resolve path in entity content
	return resolveAttributePath(gtsID, path, entity.Content)
}

// splitAtPath splits a GTS ID with path into GTS ID and attribute path
// see gts-python gts.py GtsID.split_at_path method
func splitAtPath(gtsWithPath string) (string, string) {
	if !strings.Contains(gtsWithPath, "@") {
		return gtsWithPath, ""
	}

	parts := strings.SplitN(gtsWithPath, "@", 2)
	gtsID := parts[0]
	path := ""
	if len(parts) == 2 {
		path = parts[1]
	}

	return gtsID, path
}

// resolveAttributePath resolves an attribute path in content
// see gts-python path_resolver.py JsonPathResolver.resolve method
func resolveAttributePath(gtsID, path string, content map[string]any) *AttributeResult {
	result := &AttributeResult{
		GtsID:           gtsID,
		Path:            path,
		Resolved:        false,
		AvailableFields: []string{},
	}

	// Parse path into parts
	parts := parsePath(path)

	// Traverse content following the path
	var current any = content
	for _, part := range parts {
		switch node := current.(type) {
		case map[string]any:
			// Expect field name, not array index
			if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
				result.Error = fmt.Sprintf("Path not found at segment '%s' in '%s', see available fields", part, path)
				result.AvailableFields = collectAvailableFields(node, "")
				return result
			}

			// Check if field exists
			val, exists := node[part]
			if !exists {
				result.Error = fmt.Sprintf("Path not found at segment '%s' in '%s', see available fields", part, path)
				result.AvailableFields = collectAvailableFields(node, "")
				return result
			}

			current = val

		case []any:
			// Expect array index
			var idx int
			var err error

			if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
				// Parse [N] format
				idxStr := part[1 : len(part)-1]
				idx, err = strconv.Atoi(idxStr)
				if err != nil {
					result.Error = fmt.Sprintf("Expected list index at segment '%s'", part)
					result.AvailableFields = collectAvailableFieldsFromArray(node, "")
					return result
				}
			} else {
				// Try to parse as integer
				idx, err = strconv.Atoi(part)
				if err != nil {
					result.Error = fmt.Sprintf("Expected list index at segment '%s'", part)
					result.AvailableFields = collectAvailableFieldsFromArray(node, "")
					return result
				}
			}

			// Check bounds
			if idx < 0 || idx >= len(node) {
				result.Error = fmt.Sprintf("Index out of range at segment '%s'", part)
				result.AvailableFields = collectAvailableFieldsFromArray(node, "")
				return result
			}

			current = node[idx]

		default:
			result.Error = fmt.Sprintf("Cannot descend into %T at segment '%s'", current, part)
			return result
		}
	}

	// Successfully resolved
	result.Value = current
	result.Resolved = true
	return result
}

// parsePath parses an attribute path into parts, handling array indices
// see gts-python path_resolver.py JsonPathResolver._parts method
func parsePath(path string) []string {
	// Normalize path (replace / with .)
	normalized := strings.ReplaceAll(path, "/", ".")

	// Split by dots but preserve array indices
	rawParts := []string{}
	for _, seg := range strings.Split(normalized, ".") {
		if seg != "" {
			rawParts = append(rawParts, seg)
		}
	}

	// Parse each segment to handle array indices
	parts := []string{}
	for _, seg := range rawParts {
		parts = append(parts, parsePathSegment(seg)...)
	}

	return parts
}

// parsePathSegment parses a path segment into sub-parts, extracting array indices
// see gts-python path_resolver.py JsonPathResolver._parse_part method
func parsePathSegment(seg string) []string {
	out := []string{}
	buf := ""
	i := 0

	for i < len(seg) {
		ch := seg[i]
		if ch == '[' {
			// Save any accumulated buffer
			if buf != "" {
				out = append(out, buf)
				buf = ""
			}

			// Find closing bracket
			j := strings.Index(seg[i+1:], "]")
			if j == -1 {
				// No closing bracket, treat rest as literal
				buf += seg[i:]
				break
			}

			// Extract [index] as a part
			j += i + 1 // Adjust for offset
			out = append(out, seg[i:j+1])
			i = j + 1
		} else {
			buf += string(ch)
			i++
		}
	}

	if buf != "" {
		out = append(out, buf)
	}

	return out
}

// collectAvailableFields recursively collects available fields from a map
// see gts-python path_resolver.py JsonPathResolver._list_available method
func collectAvailableFields(node map[string]any, prefix string) []string {
	fields := []string{}

	for key, val := range node {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		fields = append(fields, path)

		// Recurse into nested structures
		switch v := val.(type) {
		case map[string]any:
			fields = append(fields, collectAvailableFields(v, path)...)
		case []any:
			fields = append(fields, collectAvailableFieldsFromArray(v, path)...)
		}
	}

	return fields
}

// collectAvailableFieldsFromArray collects available indices and nested fields from an array
// see gts-python path_resolver.py JsonPathResolver._list_available method
func collectAvailableFieldsFromArray(node []any, prefix string) []string {
	fields := []string{}

	for i, val := range node {
		path := fmt.Sprintf("[%d]", i)
		if prefix != "" {
			path = prefix + path
		}
		fields = append(fields, path)

		// Recurse into nested structures
		switch v := val.(type) {
		case map[string]any:
			fields = append(fields, collectAvailableFields(v, path)...)
		case []any:
			fields = append(fields, collectAvailableFieldsFromArray(v, path)...)
		}
	}

	return fields
}
