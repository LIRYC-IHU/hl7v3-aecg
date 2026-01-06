package types

import (
	"context"
	"sync"
	"testing"
)

// TestSeries_Validate tests Series validation
func TestSeries_Validate(t *testing.T) {
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
		series    Series
		wantError *ValidationError
	}{
		{
			name: "Valid series - minimal",
			series: Series{
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low:  Time{Value: "20021122091000"},
					High: Time{Value: "20021122091010"},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid series with ID",
			series: Series{
				ID: &ID{Root: "b65deea0-078e-11d9-9669-0800200c9a66"},
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid series with author",
			series: Series{
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       REPRESENTATIVE_BEAT_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
				Author: &Author{SeriesAuthor: SeriesAuthor{
					ManufacturedSeriesDevice: ManufacturedSeriesDevice{
						ID: &ID{Root: "2.16.840.1.113883.3.1234.5", Extension: "ECG-001"},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Missing code",
			series: Series{
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
			},
			wantError: ErrMissingCode,
		},
		{
			name: "Invalid effective time",
			series: Series{
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "invalid-time"},
				},
			},
			wantError: ErrInvalidTimeFormat,
		},
		{
			name: "Invalid series ID",
			series: Series{
				ID: &ID{Root: ""},
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.series.Validate(ctx, vctx)

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

// TestSeriesAuthor_Validate tests SeriesAuthor validation
func TestSeriesAuthor_Validate(t *testing.T) {
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
		name         string
		seriesAuthor SeriesAuthor
		wantError    *ValidationError
	}{
		{
			name: "Valid author - minimal",
			seriesAuthor: SeriesAuthor{
				ManufacturedSeriesDevice: ManufacturedSeriesDevice{},
			},
			wantError: nil,
		},
		{
			name: "Valid author with ID",
			seriesAuthor: SeriesAuthor{
				ID: &ID{Root: "2.16.840.1.113883.3.5", Extension: "45"},
				ManufacturedSeriesDevice: ManufacturedSeriesDevice{
					ID: &ID{Root: "1.3.6.1.4.1.57054", Extension: "SN234-AR9-102993"},
				},
			},
			wantError: nil,
		},
		{
			name: "Valid author with manufacturer organization",
			seriesAuthor: SeriesAuthor{
				ManufacturedSeriesDevice: ManufacturedSeriesDevice{
					ID: &ID{Root: "1.3.6.1.4.1.57054"},
				},
				ManufacturerOrganization: &ManufacturerOrganization{
					ID:   &ID{Root: "1.3.6.1.4.1.57054"},
					Name: stringPtr("ECG Devices By Smith, Inc."),
				},
			},
			wantError: nil,
		},
		{
			name: "Invalid author ID",
			seriesAuthor: SeriesAuthor{
				ID: &ID{Root: ""},
				ManufacturedSeriesDevice: ManufacturedSeriesDevice{},
			},
			wantError: ErrMissingID,
		},
		{
			name: "Invalid device ID",
			seriesAuthor: SeriesAuthor{
				ManufacturedSeriesDevice: ManufacturedSeriesDevice{
					ID: &ID{Root: ""},
				},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.seriesAuthor.Validate(ctx, vctx)

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

// TestManufacturedSeriesDevice_Validate tests ManufacturedSeriesDevice validation
func TestManufacturedSeriesDevice_Validate(t *testing.T) {
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
		name   string
		device ManufacturedSeriesDevice
		wantError    *ValidationError
	}{
		{
			name:      "Valid device - minimal",
			device:    ManufacturedSeriesDevice{},
			wantError: nil,
		},
		{
			name: "Valid device with ID",
			device: ManufacturedSeriesDevice{
				ID: &ID{Root: "1.3.6.1.4.1.57054", Extension: "SN234-AR9-102993"},
			},
			wantError: nil,
		},
		{
			name: "Valid device with code",
			device: ManufacturedSeriesDevice{
				Code: &Code[DeviceTypeCode, CodeSystemOID]{
					Code:       DEVICE_12LEAD_ECG,
					CodeSystem: CodeSystemOID(""),
				},
			},
			wantError: nil,
		},
		{
			name: "Valid device - complete",
			device: ManufacturedSeriesDevice{
				ID: &ID{Root: "1.3.6.1.4.1.57054", Extension: "SN234-AR9-102993"},
				Code: &Code[DeviceTypeCode, CodeSystemOID]{
					Code:       DEVICE_12LEAD_ECG,
					CodeSystem: CodeSystemOID(""),
				},
				ManufacturerModelName: stringPtr("CardioMax Pro 3000"),
				SoftwareName:          stringPtr("v2.4.1"),
			},
			wantError: nil,
		},
		{
			name: "Invalid device ID",
			device: ManufacturedSeriesDevice{
				ID: &ID{Root: ""},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.device.Validate(ctx, vctx)

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

// TestManufacturerOrganization_Validate tests ManufacturerOrganization validation
func TestManufacturerOrganization_Validate(t *testing.T) {
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
		name         string
		organization ManufacturerOrganization
		wantError    *ValidationError
	}{
		{
			name:         "Valid organization - minimal",
			organization: ManufacturerOrganization{},
			wantError:    nil,
		},
		{
			name: "Valid organization with ID",
			organization: ManufacturerOrganization{
				ID: &ID{Root: "1.3.6.1.4.1.57054"},
			},
			wantError: nil,
		},
		{
			name: "Valid organization with name",
			organization: ManufacturerOrganization{
				Name: stringPtr("ECG Devices By Smith, Inc."),
			},
			wantError: nil,
		},
		{
			name: "Valid organization - complete",
			organization: ManufacturerOrganization{
				ID:   &ID{Root: "1.3.6.1.4.1.57054"},
				Name: stringPtr("Philips Healthcare"),
			},
			wantError: nil,
		},
		{
			name: "Invalid organization ID",
			organization: ManufacturerOrganization{
				ID: &ID{Root: ""},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.organization.Validate(ctx, vctx)

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

// TestSeries_Validate_WithCompleteAuthor tests Series validation with complete author hierarchy
func TestSeries_Validate_WithCompleteAuthor(t *testing.T) {
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
		series    Series
		wantError *ValidationError
	}{
		{
			name: "Valid series with complete author hierarchy",
			series: Series{
				ID: &ID{Root: "b65deea0-078e-11d9-9669-0800200c9a66"},
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low:  Time{Value: "20021122091000"},
					High: Time{Value: "20021122091010"},
				},
				Author: &Author{SeriesAuthor: SeriesAuthor{
					ID: &ID{Root: "2.16.840.1.113883.3.5", Extension: "45"},
					ManufacturedSeriesDevice: ManufacturedSeriesDevice{
						ID: &ID{Root: "1.3.6.1.4.1.57054", Extension: "SN234-AR9-102993"},
						Code: &Code[DeviceTypeCode, CodeSystemOID]{
							Code:       DEVICE_12LEAD_ECG,
							CodeSystem: CodeSystemOID(""),
						},
						ManufacturerModelName: stringPtr("CardioMax Pro 3000"),
						SoftwareName:          stringPtr("v2.4.1"),
					},
					ManufacturerOrganization: &ManufacturerOrganization{
						ID:   &ID{Root: "1.3.6.1.4.1.57054"},
						Name: stringPtr("CardioTech Medical Devices Inc."),
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Invalid - author with missing device ID",
			series: Series{
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
				Author: &Author{SeriesAuthor: SeriesAuthor{
					ManufacturedSeriesDevice: ManufacturedSeriesDevice{
						ID: &ID{Root: ""},
					},
				},
			},
			wantError: ErrMissingID,
		},
		{
			name: "Invalid - author with missing manufacturer org ID",
			series: Series{
				Code: &Code[SeriesTypeCode, CodeSystemOID]{
					Code:       RHYTHM_CODE,
					CodeSystem: HL7_ActCode_OID,
				},
				EffectiveTime: EffectiveTime{
					Low: Time{Value: "20021122091000"},
				},
				Author: &Author{SeriesAuthor: SeriesAuthor{
					ManufacturedSeriesDevice: ManufacturedSeriesDevice{},
					ManufacturerOrganization: &ManufacturerOrganization{
						ID: &ID{Root: ""},
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

			tt.series.Validate(ctx, vctx)

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
