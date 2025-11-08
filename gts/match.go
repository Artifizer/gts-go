/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

import (
	"fmt"
	"strings"
)

// MatchIDResult represents the result of matching a GTS identifier against a pattern
type MatchIDResult struct {
	Candidate string `json:"candidate"`
	Pattern   string `json:"pattern"`
	Match     bool   `json:"match"`
	Error     string `json:"error"`
}

// InvalidWildcardError represents an error when a wildcard pattern is invalid
type InvalidWildcardError struct {
	Pattern string
	Cause   string
}

func (e *InvalidWildcardError) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("Invalid GTS wildcard pattern: %s: %s", e.Pattern, e.Cause)
	}
	return fmt.Sprintf("Invalid GTS wildcard pattern: %s", e.Pattern)
}

// MatchIDPattern matches a candidate GTS identifier against a pattern with wildcards
// Returns a MatchIDResult with Match=true if the candidate matches the pattern,
// or Match=false with an optional Error message on failure or mismatch
func MatchIDPattern(candidate, pattern string) MatchIDResult {
	// Parse candidate
	candidateID, err := NewGtsID(candidate)
	if err != nil {
		return MatchIDResult{
			Candidate: candidate,
			Pattern:   pattern,
			Match:     false,
			Error:     err.Error(),
		}
	}

	// Validate and parse pattern
	patternID, err := validateWildcard(pattern)
	if err != nil {
		return MatchIDResult{
			Candidate: candidate,
			Pattern:   pattern,
			Match:     false,
			Error:     err.Error(),
		}
	}

	// Perform matching
	match := wildcardMatch(candidateID, patternID)

	return MatchIDResult{
		Candidate: candidate,
		Pattern:   pattern,
		Match:     match,
		Error:     "",
	}
}

// validateWildcard validates a wildcard pattern and returns a parsed GtsID
func validateWildcard(pattern string) (*GtsID, error) {
	p := strings.TrimSpace(pattern)

	// Must start with gts.
	if !strings.HasPrefix(p, GtsPrefix) {
		return nil, &InvalidWildcardError{
			Pattern: pattern,
			Cause:   fmt.Sprintf("Does not start with '%s'", GtsPrefix),
		}
	}

	// Count wildcards
	wildcardCount := strings.Count(p, "*")
	if wildcardCount > 1 {
		return nil, &InvalidWildcardError{
			Pattern: pattern,
			Cause:   "The wildcard '*' token is allowed only once",
		}
	}

	// If wildcard exists, must be at the end
	if wildcardCount == 1 {
		if !strings.HasSuffix(p, ".*") && !strings.HasSuffix(p, "~*") {
			return nil, &InvalidWildcardError{
				Pattern: pattern,
				Cause:   "The wildcard '*' token is allowed only at the end of the pattern",
			}
		}
	}

	// Try to parse as a GtsID
	id, err := NewGtsID(p)
	if err != nil {
		return nil, &InvalidWildcardError{
			Pattern: pattern,
			Cause:   err.Error(),
		}
	}

	return id, nil
}

// wildcardMatch performs the actual matching between candidate and pattern
func wildcardMatch(candidate, pattern *GtsID) bool {
	if candidate == nil || pattern == nil {
		return false
	}

	// If no wildcard in pattern, perform exact match with version flexibility
	if !strings.Contains(pattern.ID, "*") {
		return matchSegments(pattern.Segments, candidate.Segments)
	}

	// Wildcard case
	if strings.Count(pattern.ID, "*") > 1 || !strings.HasSuffix(pattern.ID, "*") {
		return false
	}

	// Use segment matching for wildcard patterns too
	return matchSegments(pattern.Segments, candidate.Segments)
}

// matchSegments matches pattern segments against candidate segments
func matchSegments(patternSegs, candidateSegs []*GtsIDSegment) bool {
	// If pattern is longer than candidate, no match
	if len(patternSegs) > len(candidateSegs) {
		return false
	}

	for i, pSeg := range patternSegs {
		cSeg := candidateSegs[i]

		// If pattern segment is a wildcard, check non-wildcard fields first
		if pSeg.IsWildcard {
			// Check the fields that are set (non-empty) in the wildcard pattern
			if pSeg.Vendor != "" && pSeg.Vendor != cSeg.Vendor {
				return false
			}
			if pSeg.Package != "" && pSeg.Package != cSeg.Package {
				return false
			}
			if pSeg.Namespace != "" && pSeg.Namespace != cSeg.Namespace {
				return false
			}
			if pSeg.Type != "" && pSeg.Type != cSeg.Type {
				return false
			}
			// Check version fields if they are set in the pattern
			if pSeg.VerMajor != 0 && pSeg.VerMajor != cSeg.VerMajor {
				return false
			}
			if pSeg.VerMinor != nil && (cSeg.VerMinor == nil || *pSeg.VerMinor != *cSeg.VerMinor) {
				return false
			}
			// Check is_type flag if set
			if pSeg.IsType && pSeg.IsType != cSeg.IsType {
				return false
			}
			// Wildcard matches - accept anything after this point
			return true
		}

		// Non-wildcard segment - all fields must match
		if pSeg.Vendor != cSeg.Vendor {
			return false
		}
		if pSeg.Package != cSeg.Package {
			return false
		}
		if pSeg.Namespace != cSeg.Namespace {
			return false
		}
		if pSeg.Type != cSeg.Type {
			return false
		}

		// Check version matching
		// Major version must match
		if pSeg.VerMajor != cSeg.VerMajor {
			return false
		}

		// Minor version: if pattern has no minor version, accept any minor in candidate
		// If pattern has minor version, it must match exactly
		if pSeg.VerMinor != nil {
			if cSeg.VerMinor == nil || *pSeg.VerMinor != *cSeg.VerMinor {
				return false
			}
		}
		// else: pattern has no minor version, so any minor version in candidate is OK

		// Check is_type flag matches
		if pSeg.IsType != cSeg.IsType {
			return false
		}
	}

	// If we've matched all pattern segments, it's a match
	return true
}
