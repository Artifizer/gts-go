/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package gts

// IDValidationResult represents the result of GTS ID validation
type IDValidationResult struct {
	ID    string `json:"id"`
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidateGtsID validates a GTS identifier and returns a result
func ValidateGtsID(gtsID string) *IDValidationResult {
	_, err := NewGtsID(gtsID)
	if err != nil {
		return &IDValidationResult{
			ID:    gtsID,
			Valid: false,
			Error: err.Error(),
		}
	}
	return &IDValidationResult{
		ID:    gtsID,
		Valid: true,
		Error: "",
	}
}

// ExtractGtsID extracts GTS ID from JSON content
func ExtractGtsID(content map[string]any, cfg *GtsConfig) *ExtractIDResult {
	return ExtractID(content, cfg)
}

// ParseGtsID parses a GTS identifier into its components
func ParseGtsID(gtsID string) ParseIDResult {
	return ParseID(gtsID)
}

// UUIDResult represents the result of GTS ID to UUID conversion
type UUIDResult struct {
	ID    string `json:"id"`
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

// IDToUUID converts a GTS ID to a UUID
func IDToUUID(gtsID string) *UUIDResult {
	id, err := NewGtsID(gtsID)
	if err != nil {
		return &UUIDResult{
			ID:    gtsID,
			UUID:  "",
			Error: err.Error(),
		}
	}

	uuid := id.ToUUID()
	return &UUIDResult{
		ID:    gtsID,
		UUID:  uuid.String(),
		Error: "",
	}
}
