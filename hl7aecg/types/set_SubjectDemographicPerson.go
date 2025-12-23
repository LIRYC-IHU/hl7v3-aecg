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
	s.AdministrativeGenderCode.SetCode(gender, codeSystem, "")
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
func (s *SubjectDemographicPerson) SetRace(raceCode RaceCode, codeSystem CodeSystemOID, display string) *SubjectDemographicPerson {
	if s.RaceCode == nil {
		s.RaceCode = &Code[RaceCode, CodeSystemOID]{}
	}
	if codeSystem == "" {
		codeSystem = HL7_Race_OID
	}
	s.RaceCode.SetCode(raceCode, codeSystem, display)
	return s
}
