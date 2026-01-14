package types

import (
	"encoding/xml"
	"testing"
)

// TestID_SetID tests the SetID method for ID type
func TestID_SetID(t *testing.T) {
	tests := []struct {
		name          string
		root          string
		extension     string
		wantRootEmpty bool // if root is empty, UUID should be generated
	}{
		{
			name:          "Set ID with both root and extension",
			root:          "2.16.840.1.113883.3.1234",
			extension:     "TEST-001",
			wantRootEmpty: false,
		},
		{
			name:          "Set ID with extension only (root stays empty without singleton)",
			root:          "",
			extension:     "TEST-002",
			wantRootEmpty: true, // No UUID generation, root stays empty
		},
		{
			name:          "Set ID with root only",
			root:          "2.16.840.1.113883.3.1234",
			extension:     "",
			wantRootEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ID{}
			id.SetID(tt.root, tt.extension)

			if id.Root == "" && !tt.wantRootEmpty {
				t.Errorf("SetID() Root is empty, expected non-empty")
			}

			if tt.extension != "" && id.Extension != tt.extension {
				t.Errorf("SetID() Extension = %v, want %v", id.Extension, tt.extension)
			}

			// If root was provided, it should be preserved
			if tt.root != "" && id.Root != tt.root {
				t.Errorf("SetID() Root = %v, want %v", id.Root, tt.root)
			}
		})
	}
}

// TestID_IsEmpty tests the IsEmpty method for ID type
func TestID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want bool
	}{
		{
			name: "Empty ID",
			id:   ID{},
			want: true,
		},
		{
			name: "ID with root only",
			id:   ID{Root: "2.16.840.1.113883.3.1234"},
			want: false,
		},
		{
			name: "ID with extension only",
			id:   ID{Extension: "TEST-001"},
			want: true, // Empty because Root is required
		},
		{
			name: "ID with both root and extension",
			id:   ID{Root: "2.16.840.1.113883.3.1234", Extension: "TEST-001"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsEmpty(); got != tt.want {
				t.Errorf("ID.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestID_String tests the String method for ID type
func TestID_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "ID with both root and extension",
			id:   ID{Root: "2.16.840.1.113883.3.1234", Extension: "TEST-001"},
			want: "2.16.840.1.113883.3.1234^TEST-001",
		},
		{
			name: "ID with root only",
			id:   ID{Root: "2.16.840.1.113883.3.1234"},
			want: "2.16.840.1.113883.3.1234",
		},
		{
			name: "Empty ID",
			id:   ID{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("ID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCode_NewCode tests the NewCode constructor function
func TestCode_NewCode(t *testing.T) {
	tests := []struct {
		name               string
		code               string
		codeSystem         string
		displayName        string
		codeSystemName     string
		wantCode           string
		wantCodeSystem     string
		wantCodeSystemName string
		wantDisplayName    string
	}{
		{
			name:               "Create code with all fields",
			code:               "93000",
			codeSystem:         "2.16.840.1.113883.6.12",
			codeSystemName:     "CPT-4",
			displayName:        "Routine ECG",
			wantCode:           "93000",
			wantCodeSystem:     "2.16.840.1.113883.6.12",
			wantCodeSystemName: "CPT-4",
			wantDisplayName:    "Routine ECG",
		},
		{
			name:               "Create code without display name",
			code:               "RHYTHM",
			codeSystem:         "2.16.840.1.113883.5.4",
			codeSystemName:     "CPT-4",
			displayName:        "",
			wantCode:           "RHYTHM",
			wantCodeSystem:     "2.16.840.1.113883.5.4",
			wantCodeSystemName: "CPT-4",
			wantDisplayName:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCode[string, string](tt.code, tt.codeSystem, tt.codeSystemName, tt.displayName)

			if got.Code != tt.wantCode {
				t.Errorf("NewCode() Code = %v, want %v", got.Code, tt.wantCode)
			}
			if got.CodeSystem != tt.wantCodeSystem {
				t.Errorf("NewCode() CodeSystem = %v, want %v", got.CodeSystem, tt.wantCodeSystem)
			}
			if got.DisplayName != tt.wantDisplayName {
				t.Errorf("NewCode() DisplayName = %v, want %v", got.DisplayName, tt.wantDisplayName)
			}
			if got.CodeSystemName != tt.wantCodeSystemName {
				t.Errorf("NewCode() CodeSystemName = %v, want %v", got.CodeSystemName, tt.wantCodeSystemName)
			}
		})
	}
}

// TestCode_SetCode tests the SetCode method
func TestCode_SetCode(t *testing.T) {
	tests := []struct {
		name               string
		initialCode        *Code[string, string]
		code               string
		codeSystem         string
		codeSystemName     string
		displayName        string
		wantCode           string
		wantCodeSystem     string
		wantCodeSystemName string
		wantDisplayName    string
	}{
		{
			name:               "Set code on empty code",
			initialCode:        &Code[string, string]{},
			code:               "M",
			codeSystem:         "2.16.840.1.113883.5.1",
			codeSystemName:     "CPT-4",
			displayName:        "Male",
			wantCode:           "M",
			wantCodeSystem:     "2.16.840.1.113883.5.1",
			wantCodeSystemName: "CPT-4",
			wantDisplayName:    "Male",
		},
		{
			name: "Update existing code",
			initialCode: &Code[string, string]{
				Code:           "OLD",
				CodeSystem:     "OLD_SYSTEM",
				DisplayName:    "Old Display",
				CodeSystemName: "OLD_NAME",
			},
			code:               "NEW",
			codeSystem:         "NEW_SYSTEM",
			displayName:        "New Display",
			codeSystemName:     "NEW_NAME",
			wantCode:           "NEW",
			wantCodeSystem:     "NEW_SYSTEM",
			wantCodeSystemName: "NEW_NAME",
			wantDisplayName:    "New Display",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialCode.SetCode(tt.code, tt.codeSystem, tt.codeSystemName, tt.displayName)

			if tt.initialCode.Code != tt.wantCode {
				t.Errorf("SetCode() Code = %v, want %v", tt.initialCode.Code, tt.wantCode)
			}
			if tt.initialCode.CodeSystem != tt.wantCodeSystem {
				t.Errorf("SetCode() CodeSystem = %v, want %v", tt.initialCode.CodeSystem, tt.wantCodeSystem)
			}
			if tt.initialCode.DisplayName != tt.wantDisplayName {
				t.Errorf("SetCode() DisplayName = %v, want %v", tt.initialCode.DisplayName, tt.wantDisplayName)
			}
			if tt.initialCode.CodeSystemName != tt.wantCodeSystemName {
				t.Errorf("SetCode() CodeSystemName = %v, want %v", tt.initialCode.CodeSystemName, tt.wantCodeSystemName)
			}
		})
	}
}

// TestCode_IsEmpty tests the IsEmpty method for Code type
func TestCode_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		code Code[string, string]
		want bool
	}{
		{
			name: "Empty code",
			code: Code[string, string]{},
			want: true,
		},
		{
			name: "Code with code value only",
			code: Code[string, string]{Code: "M"},
			want: false,
		},
		{
			name: "Code with all fields",
			code: Code[string, string]{
				Code:        "M",
				CodeSystem:  "2.16.840.1.113883.5.1",
				DisplayName: "Male",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.IsEmpty(); got != tt.want {
				t.Errorf("Code.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCode_String tests the String method for Code type
func TestCode_String(t *testing.T) {
	tests := []struct {
		name string
		code Code[string, string]
		want string
	}{
		{
			name: "Code with display name",
			code: Code[string, string]{
				Code:        "M",
				CodeSystem:  "2.16.840.1.113883.5.1",
				DisplayName: "Male",
			},
			want: "Male",
		},
		{
			name: "Code without display name",
			code: Code[string, string]{
				Code:       "M",
				CodeSystem: "2.16.840.1.113883.5.1",
			},
			want: "M",
		},
		{
			name: "Empty code",
			code: Code[string, string]{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.want {
				t.Errorf("Code.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEffectiveTime_IsEmpty tests the IsEmpty method for EffectiveTime
func TestEffectiveTime_IsEmpty(t *testing.T) {
	tests := []struct {
		name          string
		effectiveTime EffectiveTime
		want          bool
	}{
		{
			name:          "Empty effective time",
			effectiveTime: EffectiveTime{},
			want:          true,
		},
		{
			name: "Effective time with low only",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "20231223120000"},
			},
			want: false,
		},
		{
			name: "Effective time with high only",
			effectiveTime: EffectiveTime{
				High: Time{Value: "20231223120000"},
			},
			want: false,
		},
		{
			name: "Effective time with both low and high",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "20231223120000"},
				High: Time{Value: "20231223120010"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.effectiveTime.IsEmpty(); got != tt.want {
				t.Errorf("EffectiveTime.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestID_XMLMarshal tests XML marshaling for ID type
func TestID_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		wantXML string
	}{
		{
			name: "ID with both root and extension",
			id: ID{
				Root:      "2.16.840.1.113883.3.1234",
				Extension: "TEST-001",
			},
			wantXML: `<ID root="2.16.840.1.113883.3.1234" extension="TEST-001"></ID>`,
		},
		{
			name: "ID with root only",
			id: ID{
				Root: "2.16.840.1.113883.3.1234",
			},
			wantXML: `<ID root="2.16.840.1.113883.3.1234"></ID>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.id)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			got := string(data)
			if got != tt.wantXML {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.wantXML)
			}
		})
	}
}

// TestID_XMLUnmarshal tests XML unmarshaling for ID type
func TestID_XMLUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		xmlData   string
		wantRoot  string
		wantExt   string
		wantError bool
	}{
		{
			name:      "Unmarshal ID with both attributes",
			xmlData:   `<id root="2.16.840.1.113883.3.1234" extension="TEST-001"></id>`,
			wantRoot:  "2.16.840.1.113883.3.1234",
			wantExt:   "TEST-001",
			wantError: false,
		},
		{
			name:      "Unmarshal ID with root only",
			xmlData:   `<id root="2.16.840.1.113883.3.1234"></id>`,
			wantRoot:  "2.16.840.1.113883.3.1234",
			wantExt:   "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := xml.Unmarshal([]byte(tt.xmlData), &id)

			if (err != nil) != tt.wantError {
				t.Errorf("xml.Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if id.Root != tt.wantRoot {
					t.Errorf("xml.Unmarshal() Root = %v, want %v", id.Root, tt.wantRoot)
				}
				if id.Extension != tt.wantExt {
					t.Errorf("xml.Unmarshal() Extension = %v, want %v", id.Extension, tt.wantExt)
				}
			}
		})
	}
}

// TestCode_XMLMarshal tests XML marshaling for Code type
func TestCode_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		code    Code[string, string]
		wantXML string
	}{
		{
			name: "Code with all attributes",
			code: Code[string, string]{
				Code:        "93000",
				CodeSystem:  "2.16.840.1.113883.6.12",
				DisplayName: "Routine ECG",
			},
			wantXML: `<Code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"></Code>`,
		},
		{
			name: "Code without display name",
			code: Code[string, string]{
				Code:       "RHYTHM",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			wantXML: `<Code code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"></Code>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.code)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			got := string(data)
			if got != tt.wantXML {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.wantXML)
			}
		})
	}
}

// TestTime_XMLMarshal tests XML marshaling for Time type
func TestTime_XMLMarshal(t *testing.T) {
	tr := true
	f := false
	tests := []struct {
		name    string
		time    Time
		wantXML string
	}{
		{
			name:    "Time with HL7 timestamp",
			time:    Time{Value: "20231223120000", Inclusive: &tr},
			wantXML: `<Time value="20231223120000" inclusive="true"></Time>`,
		},
		{
			name:    "Time with milliseconds",
			time:    Time{Value: "20231223120000.123", Inclusive: &f},
			wantXML: `<Time value="20231223120000.123" inclusive="false"></Time>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.time)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			got := string(data)
			if got != tt.wantXML {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.wantXML)
			}
		})
	}
}

// TestCode_XMLUnmarshal tests XML unmarshaling for Code type with generic parameters
func TestCode_XMLUnmarshal(t *testing.T) {
	tests := []struct {
		name               string
		xmlData            string
		wantCode           string
		wantCodeSystem     string
		wantCodeSystemName string
		wantDisplayName    string
		wantError          bool
	}{
		{
			name:               "Unmarshal Code with all attributes",
			xmlData:            `<code code="93000" codeSystem="2.16.840.1.113883.6.12" codeSystemName="CPT-4" displayName="Routine ECG"></code>`,
			wantCode:           "93000",
			wantCodeSystem:     "2.16.840.1.113883.6.12",
			wantCodeSystemName: "CPT-4",
			wantDisplayName:    "Routine ECG",
			wantError:          false,
		},
		{
			name:               "Unmarshal Code with required attributes only",
			xmlData:            `<code code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"></code>`,
			wantCode:           "RHYTHM",
			wantCodeSystem:     "2.16.840.1.113883.5.4",
			wantCodeSystemName: "",
			wantDisplayName:    "",
			wantError:          false,
		},
		{
			name:               "Unmarshal Code with empty codeSystem (valid for informal vocabularies)",
			xmlData:            `<code code="B" codeSystem="" displayName="Blinded"></code>`,
			wantCode:           "B",
			wantCodeSystem:     "",
			wantCodeSystemName: "",
			wantDisplayName:    "Blinded",
			wantError:          false,
		},
		{
			name:               "Unmarshal Code self-closing tag",
			xmlData:            `<code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>`,
			wantCode:           "MDC_ECG_LEAD_I",
			wantCodeSystem:     "2.16.840.1.113883.6.24",
			wantCodeSystemName: "",
			wantDisplayName:    "",
			wantError:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var code Code[string, string]
			err := xml.Unmarshal([]byte(tt.xmlData), &code)

			if (err != nil) != tt.wantError {
				t.Errorf("xml.Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if code.Code != tt.wantCode {
					t.Errorf("xml.Unmarshal() Code = %v, want %v", code.Code, tt.wantCode)
				}
				if code.CodeSystem != tt.wantCodeSystem {
					t.Errorf("xml.Unmarshal() CodeSystem = %v, want %v", code.CodeSystem, tt.wantCodeSystem)
				}
				if code.CodeSystemName != tt.wantCodeSystemName {
					t.Errorf("xml.Unmarshal() CodeSystemName = %v, want %v", code.CodeSystemName, tt.wantCodeSystemName)
				}
				if code.DisplayName != tt.wantDisplayName {
					t.Errorf("xml.Unmarshal() DisplayName = %v, want %v", code.DisplayName, tt.wantDisplayName)
				}
			}
		})
	}
}

// TestCode_XMLUnmarshal_TypedGenerics tests XML unmarshaling with specialized generic types
func TestCode_XMLUnmarshal_TypedGenerics(t *testing.T) {
	tests := []struct {
		name           string
		xmlData        string
		wantCode       CPT_CODE
		wantCodeSystem CodeSystemOID
		wantError      bool
	}{
		{
			name:           "Unmarshal into Code[CPT_CODE, CodeSystemOID]",
			xmlData:        `<code code="93000" codeSystem="2.16.840.1.113883.6.12"></code>`,
			wantCode:       CPT_CODE("93000"),
			wantCodeSystem: CodeSystemOID("2.16.840.1.113883.6.12"),
			wantError:      false,
		},
		{
			name:           "Unmarshal CPT code with display name",
			xmlData:        `<code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"></code>`,
			wantCode:       CPT_CODE_ECG_Routine,
			wantCodeSystem: CPT_OID,
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var code Code[CPT_CODE, CodeSystemOID]
			err := xml.Unmarshal([]byte(tt.xmlData), &code)

			if (err != nil) != tt.wantError {
				t.Errorf("xml.Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if code.Code != tt.wantCode {
					t.Errorf("xml.Unmarshal() Code = %v, want %v", code.Code, tt.wantCode)
				}
				if code.CodeSystem != tt.wantCodeSystem {
					t.Errorf("xml.Unmarshal() CodeSystem = %v, want %v", code.CodeSystem, tt.wantCodeSystem)
				}
			}
		})
	}
}

// TestCode_XMLUnmarshal_InHL7AEcgContext tests Code unmarshaling within HL7AEcg document
func TestCode_XMLUnmarshal_InHL7AEcgContext(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="728989ec-b8bc-49cd-9a5a-30be5ade1db5"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"/>
	<text>Test ECG</text>
	<confidentialityCode code="B" codeSystem="" displayName="Blinded"/>
	<reasonCode code="PER_PROTOCOL" codeSystem=""/>
</AnnotatedECG>`

	var ecg HL7AEcg
	err := xml.Unmarshal([]byte(xmlData), &ecg)
	if err != nil {
		t.Fatalf("xml.Unmarshal() error = %v", err)
	}

	// Verify main Code field (Code[CPT_CODE, CodeSystemOID])
	if ecg.Code == nil {
		t.Fatal("HL7AEcg.Code is nil, expected non-nil")
	}
	if ecg.Code.Code != CPT_CODE("93000") {
		t.Errorf("Code.Code = %v, want %v", ecg.Code.Code, CPT_CODE("93000"))
	}
	if ecg.Code.CodeSystem != CodeSystemOID("2.16.840.1.113883.6.12") {
		t.Errorf("Code.CodeSystem = %v, want %v", ecg.Code.CodeSystem, CodeSystemOID("2.16.840.1.113883.6.12"))
	}
	if ecg.Code.DisplayName != "Routine ECG" {
		t.Errorf("Code.DisplayName = %v, want %v", ecg.Code.DisplayName, "Routine ECG")
	}

	// Verify ConfidentialityCode (Code[ConfidentialityCode, string])
	if ecg.ConfidentialityCode == nil {
		t.Fatal("HL7AEcg.ConfidentialityCode is nil, expected non-nil")
	}
	if ecg.ConfidentialityCode.Code != ConfidentialityCode("B") {
		t.Errorf("ConfidentialityCode.Code = %v, want %v", ecg.ConfidentialityCode.Code, ConfidentialityCode("B"))
	}
	if ecg.ConfidentialityCode.DisplayName != "Blinded" {
		t.Errorf("ConfidentialityCode.DisplayName = %v, want %v", ecg.ConfidentialityCode.DisplayName, "Blinded")
	}

	// Verify ReasonCode (Code[ReasonCode, string])
	if ecg.ReasonCode == nil {
		t.Fatal("HL7AEcg.ReasonCode is nil, expected non-nil")
	}
	if ecg.ReasonCode.Code != ReasonCode("PER_PROTOCOL") {
		t.Errorf("ReasonCode.Code = %v, want %v", ecg.ReasonCode.Code, ReasonCode("PER_PROTOCOL"))
	}

	// Verify other fields
	if ecg.Text != "Test ECG" {
		t.Errorf("Text = %v, want %v", ecg.Text, "Test ECG")
	}
	if ecg.ID == nil || ecg.ID.Root != "728989ec-b8bc-49cd-9a5a-30be5ade1db5" {
		t.Errorf("ID.Root = %v, want %v", ecg.ID.Root, "728989ec-b8bc-49cd-9a5a-30be5ade1db5")
	}
}

// TestCode_XMLRoundTrip tests marshal → unmarshal produces equivalent struct
func TestCode_XMLRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		code Code[string, string]
	}{
		{
			name: "Round-trip with all attributes",
			code: Code[string, string]{
				Code:           "93000",
				CodeSystem:     "2.16.840.1.113883.6.12",
				CodeSystemName: "CPT-4",
				DisplayName:    "Routine ECG",
			},
		},
		{
			name: "Round-trip with required attributes only",
			code: Code[string, string]{
				Code:       "RHYTHM",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
		},
		{
			name: "Round-trip with empty codeSystem",
			code: Code[string, string]{
				Code:        "B",
				CodeSystem:  "",
				DisplayName: "Blinded",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to XML
			data, err := xml.Marshal(&tt.code)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			// Unmarshal back to struct
			var unmarshalled Code[string, string]
			err = xml.Unmarshal(data, &unmarshalled)
			if err != nil {
				t.Fatalf("xml.Unmarshal() error = %v", err)
			}

			// Compare
			if unmarshalled.Code != tt.code.Code {
				t.Errorf("Round-trip Code = %v, want %v", unmarshalled.Code, tt.code.Code)
			}
			if unmarshalled.CodeSystem != tt.code.CodeSystem {
				t.Errorf("Round-trip CodeSystem = %v, want %v", unmarshalled.CodeSystem, tt.code.CodeSystem)
			}
			if unmarshalled.CodeSystemName != tt.code.CodeSystemName {
				t.Errorf("Round-trip CodeSystemName = %v, want %v", unmarshalled.CodeSystemName, tt.code.CodeSystemName)
			}
			if unmarshalled.DisplayName != tt.code.DisplayName {
				t.Errorf("Round-trip DisplayName = %v, want %v", unmarshalled.DisplayName, tt.code.DisplayName)
			}
		})
	}
}
