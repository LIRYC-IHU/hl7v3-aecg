package types

import "testing"

// TestResponsibleParty_SetInvestigatorID tests setting investigator ID
func TestResponsibleParty_SetInvestigatorID(t *testing.T) {
	tests := []struct {
		name      string
		root      string
		extension string
	}{
		{
			name:      "Set ID with root and extension",
			root:      "2.16.840.1.113883.3.6",
			extension: "INV_001",
		},
		{
			name:      "Set ID with extension only (uses singleton)",
			root:      "",
			extension: "trialInvestigator",
		},
		{
			name:      "Set ID with root only",
			root:      "2.16.840.1.113883.3.6",
			extension: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResponsibleParty{
				TrialInvestigator: TrialInvestigator{
					ID: ID{},
				},
			}

			result := rp.SetInvestigatorID(tt.root, tt.extension)

			// Check method chaining
			if result != rp {
				t.Error("SetInvestigatorID() should return the same ResponsibleParty for chaining")
			}

			// Check extension is set correctly
			if tt.extension != "" && rp.TrialInvestigator.ID.Extension != tt.extension {
				t.Errorf("Extension = %v, want %v", rp.TrialInvestigator.ID.Extension, tt.extension)
			}

			// Check root is set correctly (when provided)
			if tt.root != "" && rp.TrialInvestigator.ID.Root != tt.root {
				t.Errorf("Root = %v, want %v", rp.TrialInvestigator.ID.Root, tt.root)
			}
		})
	}
}

// TestResponsibleParty_SetInvestigatorName tests setting investigator name
func TestResponsibleParty_SetInvestigatorName(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		given      string
		family     string
		suffix     string
		wantPrefix bool
		wantGiven  bool
		wantFamily bool
		wantSuffix bool
	}{
		{
			name:       "Set all name components",
			prefix:     "Dr.",
			given:      "John",
			family:     "Smith",
			suffix:     "MD",
			wantPrefix: true,
			wantGiven:  true,
			wantFamily: true,
			wantSuffix: true,
		},
		{
			name:       "Set only given and family",
			prefix:     "",
			given:      "Jane",
			family:     "Doe",
			suffix:     "",
			wantPrefix: false,
			wantGiven:  true,
			wantFamily: true,
			wantSuffix: false,
		},
		{
			name:       "Set only prefix and suffix",
			prefix:     "Prof.",
			given:      "",
			family:     "",
			suffix:     "PhD",
			wantPrefix: true,
			wantGiven:  false,
			wantFamily: false,
			wantSuffix: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ResponsibleParty{}

			result := rp.SetInvestigatorName(tt.prefix, tt.given, tt.family, tt.suffix)

			// Check method chaining
			if result != rp {
				t.Error("SetInvestigatorName() should return the same ResponsibleParty for chaining")
			}

			// Check InvestigatorPerson is initialized
			if rp.TrialInvestigator.InvestigatorPerson == nil {
				t.Error("InvestigatorPerson should be initialized")
				return
			}

			// Check Name is initialized
			if rp.TrialInvestigator.InvestigatorPerson.Name == nil {
				t.Error("Name should be initialized")
				return
			}

			name := rp.TrialInvestigator.InvestigatorPerson.Name

			// Check prefix
			if tt.wantPrefix {
				if name.Prefix == nil || *name.Prefix != tt.prefix {
					t.Errorf("Prefix = %v, want %v", name.Prefix, tt.prefix)
				}
			} else if name.Prefix != nil {
				t.Errorf("Prefix should be nil, got %v", *name.Prefix)
			}

			// Check given
			if tt.wantGiven {
				if name.Given == nil || *name.Given != tt.given {
					t.Errorf("Given = %v, want %v", name.Given, tt.given)
				}
			} else if name.Given != nil {
				t.Errorf("Given should be nil, got %v", *name.Given)
			}

			// Check family
			if tt.wantFamily {
				if name.Family == nil || *name.Family != tt.family {
					t.Errorf("Family = %v, want %v", name.Family, tt.family)
				}
			} else if name.Family != nil {
				t.Errorf("Family should be nil, got %v", *name.Family)
			}

			// Check suffix
			if tt.wantSuffix {
				if name.Suffix == nil || *name.Suffix != tt.suffix {
					t.Errorf("Suffix = %v, want %v", name.Suffix, tt.suffix)
				}
			} else if name.Suffix != nil {
				t.Errorf("Suffix should be nil, got %v", *name.Suffix)
			}
		})
	}
}

// TestResponsibleParty_SetEmptyInvestigatorName tests setting empty name
func TestResponsibleParty_SetEmptyInvestigatorName(t *testing.T) {
	rp := &ResponsibleParty{}

	result := rp.SetEmptyInvestigatorName()

	// Check method chaining
	if result != rp {
		t.Error("SetEmptyInvestigatorName() should return the same ResponsibleParty for chaining")
	}

	// Check InvestigatorPerson is initialized
	if rp.TrialInvestigator.InvestigatorPerson == nil {
		t.Error("InvestigatorPerson should be initialized")
		return
	}

	// Check Name is initialized
	if rp.TrialInvestigator.InvestigatorPerson.Name == nil {
		t.Error("Name should be initialized")
		return
	}

	// Check all name components are nil (empty)
	name := rp.TrialInvestigator.InvestigatorPerson.Name
	if name.Prefix != nil {
		t.Errorf("Prefix should be nil, got %v", *name.Prefix)
	}
	if name.Given != nil {
		t.Errorf("Given should be nil, got %v", *name.Given)
	}
	if name.Family != nil {
		t.Errorf("Family should be nil, got %v", *name.Family)
	}
	if name.Suffix != nil {
		t.Errorf("Suffix should be nil, got %v", *name.Suffix)
	}
}

// TestResponsibleParty_MethodChaining tests method chaining
func TestResponsibleParty_MethodChaining(t *testing.T) {
	rp := &ResponsibleParty{
		TrialInvestigator: TrialInvestigator{
			ID: ID{},
		},
	}

	// Test chaining multiple methods
	result := rp.
		SetInvestigatorID("2.16.840.1.113883.3.6", "INV_001").
		SetInvestigatorName("Dr.", "John", "Smith", "MD")

	if result != rp {
		t.Error("Method chaining should return the same ResponsibleParty instance")
	}

	// Verify all values are set
	if rp.TrialInvestigator.ID.Root != "2.16.840.1.113883.3.6" {
		t.Errorf("Root not set correctly: %v", rp.TrialInvestigator.ID.Root)
	}
	if rp.TrialInvestigator.ID.Extension != "INV_001" {
		t.Errorf("Extension not set correctly: %v", rp.TrialInvestigator.ID.Extension)
	}
	if rp.TrialInvestigator.InvestigatorPerson == nil || rp.TrialInvestigator.InvestigatorPerson.Name == nil {
		t.Error("Name should be initialized")
		return
	}

	name := rp.TrialInvestigator.InvestigatorPerson.Name
	if name.Prefix == nil || *name.Prefix != "Dr." {
		t.Error("Prefix not set correctly")
	}
	if name.Given == nil || *name.Given != "John" {
		t.Error("Given not set correctly")
	}
	if name.Family == nil || *name.Family != "Smith" {
		t.Error("Family not set correctly")
	}
	if name.Suffix == nil || *name.Suffix != "MD" {
		t.Error("Suffix not set correctly")
	}
}
