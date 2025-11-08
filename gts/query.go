/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"errors"
	"fmt"
	"strings"
)

// QueryResult represents the result of a GTS query execution
type QueryResult struct {
	Error   string           `json:"error"`
	Count   int              `json:"count"`
	Limit   int              `json:"limit"`
	Results []map[string]any `json:"results,omitempty"`
}

// Query filters entities by a GTS query expression
// Supports:
// - Exact match: "gts.x.core.events.event.v1~"
// - Wildcard match: "gts.x.core.events.*"
// - With filters: "gts.x.core.events.event.v1~[status=active]"
// - Wildcard with filters: "gts.x.core.*[status=active]"
// - Wildcard filter values: "gts.x.core.*[status=active, category=*]"
// see gts-python store.py query method
func (s *GtsStore) Query(expr string, limit int) *QueryResult {
	if limit <= 0 {
		limit = 100 // Default limit
	}

	result := &QueryResult{
		Error:   "",
		Count:   0,
		Limit:   limit,
		Results: []map[string]any{},
	}

	// Parse the query expression to extract base pattern and filters
	basePattern, filters, err := s.parseQueryExpression(expr)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Determine if pattern is wildcard
	isWildcard := strings.Contains(basePattern, "*")

	// Validate the pattern
	if err := s.validateQueryPattern(basePattern, isWildcard); err != nil {
		result.Error = err.Error()
		return result
	}

	// Filter entities
	for _, entity := range s.byID {
		if len(result.Results) >= limit {
			break
		}

		// Skip entities without valid content or GTS ID
		if len(entity.Content) == 0 || entity.GtsID == nil {
			continue
		}

		// Check if ID matches the pattern
		if !s.matchesIDPattern(entity.GtsID, basePattern, isWildcard) {
			continue
		}

		// Check filters
		if !s.matchesFilters(entity.Content, filters) {
			continue
		}

		result.Results = append(result.Results, entity.Content)
	}

	result.Count = len(result.Results)
	return result
}

// parseQueryExpression parses the query expression into base pattern and filters
// see gts-python store.py query method
func (s *GtsStore) parseQueryExpression(expr string) (string, map[string]string, error) {
	// Split by '[' to separate base pattern from filters
	parts := strings.SplitN(expr, "[", 2)
	basePattern := strings.TrimSpace(parts[0])

	filters := make(map[string]string)
	if len(parts) == 2 {
		// Extract filter string (remove trailing ])
		filterStr := strings.TrimSpace(parts[1])
		if !strings.HasSuffix(filterStr, "]") {
			return "", nil, errors.New("invalid query: missing closing bracket ']'")
		}
		filterStr = strings.TrimSuffix(filterStr, "]")

		// Parse filters
		filters = s.parseQueryFilters(filterStr)
	}

	return basePattern, filters, nil
}

// parseQueryFilters parses filter expressions from query string
// see gts-python store.py _parse_query_filters method
func (s *GtsStore) parseQueryFilters(filterStr string) map[string]string {
	filters := make(map[string]string)
	if filterStr == "" {
		return filters
	}

	// Split by comma to handle multiple filters
	parts := strings.Split(filterStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			// Remove quotes from value if present
			value = strings.Trim(value, `"'`)

			filters[key] = value
		}
	}

	return filters
}

// validateQueryPattern validates the query pattern
// see gts-python store.py _validate_query_pattern method
func (s *GtsStore) validateQueryPattern(basePattern string, isWildcard bool) error {
	if isWildcard {
		// Wildcard pattern must end with .* or ~*
		if !strings.HasSuffix(basePattern, ".*") && !strings.HasSuffix(basePattern, "~*") {
			return errors.New("invalid query: wildcard patterns must end with .* or ~*")
		}

		// Validate as wildcard pattern
		_, err := validateWildcard(basePattern)
		if err != nil {
			return fmt.Errorf("invalid query: %w", err)
		}
	} else {
		// Non-wildcard pattern must be a complete valid GTS ID
		gtsID, err := NewGtsID(basePattern)
		if err != nil {
			return fmt.Errorf("invalid query: %w", err)
		}

		// Must have at least one valid segment
		if len(gtsID.Segments) == 0 {
			return errors.New("invalid query: GTS ID has no valid segments")
		}
	}

	return nil
}

// matchesIDPattern checks if entity ID matches the query pattern
// see gts-python store.py _matches_id_pattern method
func (s *GtsStore) matchesIDPattern(entityID *GtsID, basePattern string, isWildcard bool) bool {
	if entityID == nil {
		return false
	}

	// Use the existing MatchIDPattern function
	matchResult := MatchIDPattern(entityID.ID, basePattern)
	return matchResult.Match
}

// matchesFilters checks if entity content matches all filter criteria
// see gts-python store.py _matches_filters method
func (s *GtsStore) matchesFilters(entityContent map[string]any, filters map[string]string) bool {
	if len(filters) == 0 {
		return true
	}

	for key, value := range filters {
		entityValue := fmt.Sprintf("%v", entityContent[key])

		// Support wildcard in filter values
		if value == "*" {
			// Wildcard matches any non-empty value
			if entityValue == "" || entityValue == "<nil>" {
				return false
			}
		} else if entityValue != value {
			return false
		}
	}

	return true
}
