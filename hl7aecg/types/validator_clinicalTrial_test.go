package types

import (
	"context"
	"sync"
	"testing"
)

// TestLocation_Validate tests Location validation
func TestLocation_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name      string
		location  Location
		wantError *ValidationError
	}{
		{
			name: "Valid location with trial site",
			location: Location{
				TrialSite: TrialSite{
					ID: ID{Root: "2.16.840.1.113883.3.5"},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid location with site details",
			location: Location{
				TrialSite: TrialSite{
					ID: ID{Root: "2.16.840.1.113883.3.5", Extension: "SITE_1"},
					Location: &SiteLocation{
						Name: stringPtr("Test Site"),
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Missing trial site ID",
			location: Location{
				TrialSite: TrialSite{
					ID: ID{Root: ""},
				},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.location.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestTrialSite_Validate tests TrialSite validation
func TestTrialSite_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name      string
		trialSite TrialSite
		wantError *ValidationError
	}{
		{
			name: "Valid trial site with ID only",
			trialSite: TrialSite{
				ID: ID{Root: "2.16.840.1.113883.3.5", Extension: "SITE_1"},
			},
			wantError: nil,
		},
		{
			name: "Valid trial site with location details",
			trialSite: TrialSite{
				ID: ID{Root: "2.16.840.1.113883.3.5", Extension: "SITE_1"},
				Location: &SiteLocation{
					Name: stringPtr("1st Clinic of Milwaukee"),
					Addr: &Address{
						City:    stringPtr("Milwaukee"),
						State:   stringPtr("WI"),
						Country: stringPtr("USA"),
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid trial site with responsible party",
			trialSite: TrialSite{
				ID: ID{Root: "2.16.840.1.113883.3.5", Extension: "SITE_1"},
				ResponsibleParty: &ResponsibleParty{
					TrialInvestigator: TrialInvestigator{
						ID: ID{Root: "2.16.840.1.113883.3.6", Extension: "INV_001"},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Missing trial site ID",
			trialSite: TrialSite{
				ID: ID{Root: ""},
			},
			wantError: ErrMissingID,
		},
		{
			name: "Invalid responsible party (missing investigator ID)",
			trialSite: TrialSite{
				ID: ID{Root: "2.16.840.1.113883.3.5"},
				ResponsibleParty: &ResponsibleParty{
					TrialInvestigator: TrialInvestigator{
						ID: ID{Root: ""},
					},
				},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.trialSite.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestResponsibleParty_Validate tests ResponsibleParty validation
func TestResponsibleParty_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name             string
		responsibleParty ResponsibleParty
		wantError        *ValidationError
	}{
		{
			name: "Valid responsible party with ID only",
			responsibleParty: ResponsibleParty{
				TrialInvestigator: TrialInvestigator{
					ID: ID{Root: "2.16.840.1.113883.3.6", Extension: "INV_001"},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid responsible party with investigator details",
			responsibleParty: ResponsibleParty{
				TrialInvestigator: TrialInvestigator{
					ID: ID{Root: "2.16.840.1.113883.3.6", Extension: "INV_001"},
					InvestigatorPerson: &InvestigatorPerson{
						Name: &PersonName{
							Given:  stringPtr("John"),
							Family: stringPtr("Smith"),
						},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Missing investigator ID",
			responsibleParty: ResponsibleParty{
				TrialInvestigator: TrialInvestigator{
					ID: ID{Root: ""},
				},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.responsibleParty.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestTrialInvestigator_Validate tests TrialInvestigator validation
func TestTrialInvestigator_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name              string
		trialInvestigator TrialInvestigator
		wantError         *ValidationError
	}{
		{
			name: "Valid investigator with ID only",
			trialInvestigator: TrialInvestigator{
				ID: ID{Root: "2.16.840.1.113883.3.6", Extension: "INV_001"},
			},
			wantError: nil,
		},
		{
			name: "Valid investigator with person details",
			trialInvestigator: TrialInvestigator{
				ID: ID{Root: "2.16.840.1.113883.3.6", Extension: "INV_001"},
				InvestigatorPerson: &InvestigatorPerson{
					Name: &PersonName{
						Prefix: stringPtr("Dr."),
						Given:  stringPtr("John"),
						Family: stringPtr("Smith"),
						Suffix: stringPtr("MD"),
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid investigator with OID",
			trialInvestigator: TrialInvestigator{
				ID: ID{Root: "2.16.840.1.113883.3.6"},
			},
			wantError: nil,
		},
		{
			name: "Missing investigator ID",
			trialInvestigator: TrialInvestigator{
				ID: ID{Root: ""},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.trialInvestigator.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestClinicalTrial_Validate_WithLocation tests ClinicalTrial validation with Location
func TestClinicalTrial_Validate_WithLocation(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name          string
		clinicalTrial ClinicalTrial
		wantError     *ValidationError
	}{
		{
			name: "Valid trial with location",
			clinicalTrial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.4"},
				Location: &Location{
					TrialSite: TrialSite{
						ID: ID{Root: "2.16.840.1.113883.3.5"},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid trial with complete location hierarchy",
			clinicalTrial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.4"},
				Location: &Location{
					TrialSite: TrialSite{
						ID: ID{Root: "2.16.840.1.113883.3.5"},
						ResponsibleParty: &ResponsibleParty{
							TrialInvestigator: TrialInvestigator{
								ID: ID{Root: "2.16.840.1.113883.3.6"},
							},
						},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Invalid - missing trial site ID in location",
			clinicalTrial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.4"},
				Location: &Location{
					TrialSite: TrialSite{
						ID: ID{Root: ""},
					},
				},
			},
			wantError: ErrMissingID,
		},
		{
			name: "Valid trial with activity time",
			clinicalTrial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.4"},
				ActivityTime: EffectiveTime{
					Low:  Time{Value: "20010509"},
					High: Time{Value: "20020316"},
				},
			},
			wantError: nil,
		},
		{
			name: "Invalid - bad activity time format",
			clinicalTrial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.4"},
				ActivityTime: EffectiveTime{
					Low: Time{Value: "invalid-date"},
				},
			},
			wantError: ErrInvalidTimeFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.clinicalTrial.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}
