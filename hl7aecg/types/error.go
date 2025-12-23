package types

import "fmt"

// =============================================================================
// Error Types
// =============================================================================

// ValidationError represents a validation error for HL7 aECG structures.
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Description of the validation failure
	Value   string // Optional: the invalid value
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("validation error on field %s:\n- %s (value: %s)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("validation error on field %s:\n- %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// NewValidationErrorWithValue creates a new validation error with a value.
func NewValidationErrorWithValue(field, message, value string) *ValidationError {
	return &ValidationError{Field: field, Message: message, Value: value}
}

// =============================================================================
// Common Validation Errors
// =============================================================================

var (
	// ErrMissingID indicates a required ID is missing
	ErrMissingID = NewValidationError("ID", "ID.Root is required")

	// ErrInvalidID indicates ID format is invalid
	ErrInvalidID = NewValidationError("ID", "ID.Root must be a valid UUID or OID")

	// ErrMissingCode indicates a required Code is missing
	ErrMissingCode = NewValidationError("Code", "Code is required")

	// ErrConfientialityCode indicates confidentiality code is invalid
	ErrConfidentialityCode = NewValidationError("ConfidentialityCode", "ConfidentialityCode must be one of \n\t- 'S'\n\t- 'I'\n\t- 'B'\n\t- 'C'")

	// ErrReasonCode indicates reason code is invalid
	ErrReasonCode = NewValidationError("ReasonCode", "ReasonCode must be a valid code from the appropriate code system: \n\t- PER_PROTOCOL\n\t- NOT_IN_PROTOCOL\n\t- IN_PROTOCAL_WRONG_EVENT")

	// ErrMissingCodeSystem indicates a required CodeSystem is missing
	ErrMissingCodeSystem = NewValidationError("CodeSystem", "CodeSystem is required")

	// ErrMissingEffectiveTime indicates effective time is missing
	ErrMissingEffectiveTime = NewValidationError("EffectiveTime", "EffectiveTime is required")

	// ErrInvalidEffectiveTime indicates effective time has invalid values
	ErrInvalidEffectiveTime = NewValidationError("EffectiveTime", "EffectiveTime must have at least Low or High or Center")

	// ErrInvalidTimeFormat indicates time format is invalid
	ErrInvalidTimeFormat = NewValidationError("Time", "Time format must be unix timestamp")

	// ErrInvalidOID indicates OID format is invalid
	ErrInvalidOID = NewValidationError("OID", "OID must be in dot-separated format (e.g., 2.16.840.1.113883.3.1)")

	// ErrInvalidUUID indicates UUID format is invalid
	ErrInvalidUUID = NewValidationError("UUID", "UUID must be in format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")

	// ErrMissingClinicalTrial indicates clinical trial information is missing
	ErrMissingClinicalTrial = NewValidationError("ClinicalTrial", "ClinicalTrial information is required")

	// ErrTrialSubjectCode indicates trial subject code is invalid
	ErrTrialSubjectCode = NewValidationError("TrialSubjectCode", "TrialSubjectCode must be one of \n\t- 'SCREENING'\n\t- 'ENROLLED'")

	// ErrMissingSubject indicates subject information is missing
	ErrMissingSubject = NewValidationError("Subject", "Subject information is required")

	// ErrRaceCode indicates
	ErrRaceCode = NewValidationError("RaceCode", "RaceCode must be one of \n\t- 'F'\n\t- 'M'\n\t- 'U'")

	// ErrSequenceLengthMismatch indicates sequences in a set have different lengths
	ErrSequenceLengthMismatch = NewValidationError("SequenceSet", "All sequences in a SequenceSet must have the same length")

	// ErrMissingTimeSequence indicates time sequence is missing
	ErrMissingTimeSequence = NewValidationError("SequenceSet", "SequenceSet must have at least one time sequence (TIME_ABSOLUTE or TIME_RELATIVE)")

	// ErrMissingLeadSequence indicates no lead sequences found
	ErrMissingLeadSequence = NewValidationError("SequenceSet", "SequenceSet must have at least one lead sequence")

	// ErrInvalidDigits indicates digits format is invalid
	ErrInvalidDigits = NewValidationError("Digits", "Digits must be space-separated integers")

	// ErrInvalidIncrement indicates increment value is invalid
	ErrInvalidIncrement = NewValidationError("Increment", "Increment value must be a positive number")

	// ErrInvalidScale indicates scale value is invalid
	ErrInvalidScale = NewValidationError("Scale", "Scale value must be a non-zero number")

	// ErrInvalidVoltageRange indicates voltage values are out of reasonable range
	ErrInvalidVoltageRange = NewValidationError("Voltage", "Voltage values out of reasonable range (-10mV to +10mV)")

	// ErrInvalidHeartRate indicates heart rate is out of reasonable range
	ErrInvalidHeartRate = NewValidationError("HeartRate", "Heart rate out of reasonable range (30-250 bpm)")
)

// =============================================================================
// Validation Context
// =============================================================================

// ValidationContext provides context for validation operations.
type ValidationContext struct {
	StrictMode bool     // If true, apply stricter validation rules
	Warnings   []string // Non-fatal warnings collected during validation
	Errors     []error  // Errors collected during validation
}

// NewValidationContext creates a new validation context.
func NewValidationContext(strictMode bool) *ValidationContext {
	return &ValidationContext{
		StrictMode: strictMode,
		Warnings:   make([]string, 0),
		Errors:     make([]error, 0),
	}
}

// AddWarning adds a warning to the context.
func (ctx *ValidationContext) AddWarning(warning string) {
	ctx.Warnings = append(ctx.Warnings, warning)
}

// AddError adds an error to the context.
func (ctx *ValidationContext) AddError(err error) {
	ctx.Errors = append(ctx.Errors, err)
}

// HasErrors returns true if any errors were collected.
func (ctx *ValidationContext) HasErrors() bool {
	return len(ctx.Errors) > 0
}

// HasWarnings returns true if any warnings were collected.
func (ctx *ValidationContext) HasWarnings() bool {
	return len(ctx.Warnings) > 0
}

// GetError returns a combined error if any errors exist.
func (ctx *ValidationContext) GetError() error {
	if !ctx.HasErrors() {
		return nil
	}
	if len(ctx.Errors) == 1 {
		return ctx.Errors[0]
	}
	return fmt.Errorf("multiple validation errors: %v", ctx.Errors)
}
