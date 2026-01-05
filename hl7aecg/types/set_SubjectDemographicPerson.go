package types

// SetName sets the subject's name (often just initials for privacy).
//
// Example: "BDB" (initials only)
func (s *SubjectDemographicPerson) SetName(name string) *SubjectDemographicPerson {
	s.Name = &name
	return s
}

// SetGender sets the administrative gender code.
//
// Valid codes (from HL7 AdministrativeGender vocabulary):
//   - GENDER_MALE ("M")
//   - GENDER_FEMALE ("F")
//   - GENDER_UNDIFFERENTIATED ("UN")
//
// Example: SetGender(GENDER_MALE)
func (s *SubjectDemographicPerson) SetGender(gender GenderCode, codeSystem CodeSystemOID) *SubjectDemographicPerson {
	if s.AdministrativeGenderCode == nil {
		s.AdministrativeGenderCode = &Code[GenderCode, CodeSystemOID]{}
	}
	s.AdministrativeGenderCode.SetCode(gender, codeSystem, "", "")
	return s
}

// SetBirthDate sets the subject's date of birth.
//
// Format: YYYYMMDD (e.g., "19530508" for May 8, 1953)
func (s *SubjectDemographicPerson) SetBirthDate(birthDate string) *SubjectDemographicPerson {
	if s.BirthTime == nil {
		s.BirthTime = &Time{}
	}
	s.BirthTime = &Time{Value: birthDate}
	return s
}

// SetRace sets the race code.
//
// Valid codes (from HL7 Race vocabulary):
//   - RACE_NATIVE_AMERICAN ("1002-5")
//   - RACE_ASIAN ("2028-9")
//   - RACE_BLACK_OR_AFRICAN_AMERICAN ("2054-5")
//   - RACE_HAWAIIAN_OR_PACIFIC_ISLAND ("2076-8")
//   - RACE_WHITE ("2106-3")
//   - RACE_OTHER ("2131-1")
//
// Example: SetRace(RACE_WHITE)
func (s *SubjectDemographicPerson) SetRace(raceCode RaceCode, codeSystem CodeSystemOID, codeSystemName, display string) *SubjectDemographicPerson {
	if s.RaceCode == nil {
		s.RaceCode = &Code[RaceCode, CodeSystemOID]{}
	}
	if codeSystem == "" {
		codeSystem = HL7_Race_OID
	}
	s.RaceCode.SetCode(raceCode, codeSystem, codeSystemName, display)
	return s
}

// SetPatientID sets the primary patient identifier.
//
// This is typically the hospital or institution's internal patient ID.
//
// Example: SetPatientID("25060897140")
func (s *SubjectDemographicPerson) SetPatientID(patientID string) *SubjectDemographicPerson {
	s.PatientID = patientID
	return s
}

// SetSecondPatientID sets the optional secondary patient identifier.
//
// Used when the patient has multiple identification numbers
// (e.g., different hospital systems).
//
// Example: SetSecondPatientID("ALT-123456")
func (s *SubjectDemographicPerson) SetSecondPatientID(secondPatientID string) *SubjectDemographicPerson {
	s.SecondPatientID = secondPatientID
	return s
}

// SetAge sets the subject's age at the time of ECG acquisition.
//
// Can be represented as a number (years) or other format.
//
// Example: SetAge("65")
func (s *SubjectDemographicPerson) SetAge(age string) *SubjectDemographicPerson {
	s.Age = age
	return s
}

// SetPaced sets whether the patient has a cardiac pacemaker.
//
// Example: SetPaced(true)
func (s *SubjectDemographicPerson) SetPaced(paced bool) *SubjectDemographicPerson {
	s.Paced = paced
	return s
}

// AddMedication adds a medication to the patient's medication list.
//
// Each call adds one medication. Call multiple times to add multiple medications.
//
// Example:
//
//	AddMedication("Aspirin 100mg").
//	AddMedication("Metoprolol 50mg")
func (s *SubjectDemographicPerson) AddMedication(medication string) *SubjectDemographicPerson {
	// Remove empty default entry if present
	if len(s.Medications.Medication) == 1 && s.Medications.Medication[0] == "" {
		s.Medications.Medication = []string{}
	}
	s.Medications.Medication = append(s.Medications.Medication, medication)
	return s
}

// SetMedications initializes or resets the medications list.
//
// Use this to start a fresh medication list, then chain with AddMedication().
//
// Example: SetMedications().AddMedication("Aspirin")
func (s *SubjectDemographicPerson) SetMedications() *SubjectDemographicPerson {
	s.Medications = Medications{
		Medication: []string{},
	}
	// Add one empty entry by default
	s.Medications.Medication = append(s.Medications.Medication, "")
	return s
}

// AddClinicalClassification adds a clinical classification to the patient's record.
//
// Each call adds one classification. Call multiple times to add multiple classifications.
//
// Example:
//
//	AddClinicalClassification("Hypertension").
//	AddClinicalClassification("Diabetes Type 2")
func (s *SubjectDemographicPerson) AddClinicalClassification(classification string) *SubjectDemographicPerson {
	// Remove empty default entry if present
	if len(s.ClinicalClassifications.ClinicalClassification) == 1 && s.ClinicalClassifications.ClinicalClassification[0] == "" {
		s.ClinicalClassifications.ClinicalClassification = []string{}
	}
	s.ClinicalClassifications.ClinicalClassification = append(s.ClinicalClassifications.ClinicalClassification, classification)
	return s
}

// SetClinicalClassifications initializes or resets the clinical classifications list.
//
// Use this to start a fresh classifications list, then chain with AddClinicalClassification().
//
// Example: SetClinicalClassifications().AddClinicalClassification("Hypertension")
func (s *SubjectDemographicPerson) SetClinicalClassifications() *SubjectDemographicPerson {
	s.ClinicalClassifications = ClinicalClassifications{
		ClinicalClassification: []string{},
	}
	s.ClinicalClassifications.ClinicalClassification = append(s.ClinicalClassifications.ClinicalClassification, "")
	return s
}

// SetEmptyClinicalClassifications initializes the classifications list with a specific number of empty entries.
//
// This is useful to match XML structures that require empty <ClinicalClassification/> elements.
//
// Example: SetEmptyClinicalClassifications(2) generates:
//
//	<ClinicalClassifications>
//	  <ClinicalClassification/>
//	  <ClinicalClassification/>
//	</ClinicalClassifications>
func (s *SubjectDemographicPerson) SetEmptyClinicalClassifications(count int) *SubjectDemographicPerson {
	s.ClinicalClassifications = ClinicalClassifications{
		ClinicalClassification: make([]string, count),
	}
	return s
}

// SetBed sets the patient's bed location within the facility.
//
// Example: SetBed("12A")
func (s *SubjectDemographicPerson) SetBed(bed string) *SubjectDemographicPerson {
	s.Bed = bed
	return s
}

// SetRoom sets the patient's room number or identifier.
//
// Example: SetRoom("302")
func (s *SubjectDemographicPerson) SetRoom(room string) *SubjectDemographicPerson {
	s.Room = room
	return s
}

// SetPointOfCare sets the care unit or department.
//
// Example: SetPointOfCare("Cardiology ICU")
func (s *SubjectDemographicPerson) SetPointOfCare(pointOfCare string) *SubjectDemographicPerson {
	s.PointOfCare = pointOfCare
	return s
}
