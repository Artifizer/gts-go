/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

// ParseIDSegment represents a parsed segment component from a GTS identifier
type ParseIDSegment struct {
	Vendor    string
	Package   string
	Namespace string
	Type      string
	VerMajor  int
	VerMinor  *int
	IsType    bool
}

// ParseIDResult represents the result of parsing a GTS identifier
type ParseIDResult struct {
	ID       string
	OK       bool
	Segments []ParseIDSegment
	Error    string
}

// ParseID decomposes a GTS identifier into its constituent parts
// Returns a ParseIDResult with OK=true and populated Segments on success,
// or OK=false with an Error message on failure
func ParseID(gtsID string) ParseIDResult {
	id, err := NewGtsID(gtsID)
	if err != nil {
		return ParseIDResult{
			ID:       gtsID,
			OK:       false,
			Segments: nil,
			Error:    err.Error(),
		}
	}

	segments := make([]ParseIDSegment, len(id.Segments))
	for i, seg := range id.Segments {
		segments[i] = ParseIDSegment{
			Vendor:    seg.Vendor,
			Package:   seg.Package,
			Namespace: seg.Namespace,
			Type:      seg.Type,
			VerMajor:  seg.VerMajor,
			VerMinor:  seg.VerMinor,
			IsType:    seg.IsType,
		}
	}

	return ParseIDResult{
		ID:       gtsID,
		OK:       true,
		Segments: segments,
		Error:    "",
	}
}
