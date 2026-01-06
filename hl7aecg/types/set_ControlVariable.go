package types

// SetObservationCode sets the code for the related observation.
//
// Parameters:
//   - code: The observation code (e.g., "21612-7" for "Reported Age")
//   - codeSystem: The code system OID (e.g., LOINC_OID)
//   - displayName: Optional display name for the code
//   - codeSystemName: Optional code system name (e.g., "LOINC")
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetObservationCode(code string, codeSystem CodeSystemOID, displayName, codeSystemName string) *RelatedObservation {
	if ro.Code == nil {
		ro.Code = &Code[string, CodeSystemOID]{}
	}
	ro.Code.SetCode(code, codeSystem, displayName, codeSystemName)
	return ro
}

// SetText sets the textual description of the observation.
//
// Parameters:
//   - text: The observation description
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetText(text string) *RelatedObservation {
	ro.Text = &text
	return ro
}

// SetValue sets the observed value as a physical quantity.
//
// Parameters:
//   - value: The numeric value (e.g., "34")
//   - unit: The unit of measurement (e.g., "a" for years, "mo" for months)
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetValue(value, unit string) *RelatedObservation {
	ro.Value = &PhysicalQuantity{
		XsiType: "PQ",
		Value:   value,
		Unit:    unit,
	}
	return ro
}

// SetAuthorPerson sets the observation author as a person.
//
// Parameters:
//   - authorID: Optional author ID root
//   - authorExtension: Optional author ID extension
//   - personName: The name of the person (can be initials or full name)
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetAuthorPerson(authorID, authorExtension, personName string) *RelatedObservation {
	if ro.Author == nil {
		ro.Author = &ObservationAuthor{}
	}

	// Set author ID if provided
	if authorID != "" || authorExtension != "" {
		if ro.Author.AssignedEntity.ID == nil {
			ro.Author.AssignedEntity.ID = &ID{}
		}
		ro.Author.AssignedEntity.ID.SetID(authorID, authorExtension)
	}

	// Set person
	if ro.Author.AssignedEntity.AssignedAuthorType == nil {
		ro.Author.AssignedEntity.AssignedAuthorType = &AssignedAuthorType{}
	}
	ro.Author.AssignedEntity.AssignedAuthorType.AssignedPerson = &ObservationAssignedPerson{
		Name: &personName,
	}

	return ro
}

// SetAuthorDevice sets the observation author as a device.
//
// Parameters:
//   - deviceID: Optional device ID root
//   - deviceExtension: Optional device ID extension
//   - deviceCode: Optional device type code
//   - modelName: Optional device model name
//   - softwareName: Optional software name and version
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetAuthorDevice(deviceID, deviceExtension string, deviceCode DeviceTypeCode, modelName, softwareName string) *RelatedObservation {
	if ro.Author == nil {
		ro.Author = &ObservationAuthor{}
	}

	// Set author ID if provided
	if deviceID != "" || deviceExtension != "" {
		if ro.Author.AssignedEntity.ID == nil {
			ro.Author.AssignedEntity.ID = &ID{}
		}
		ro.Author.AssignedEntity.ID.SetID(deviceID, deviceExtension)
	}

	// Set device
	if ro.Author.AssignedEntity.AssignedAuthorType == nil {
		ro.Author.AssignedEntity.AssignedAuthorType = &AssignedAuthorType{}
	}

	device := &ObservationAssignedDevice{}

	// Set device ID
	if deviceID != "" || deviceExtension != "" {
		device.ID = &ID{}
		device.ID.SetID(deviceID, deviceExtension)
	}

	// Set device code
	if deviceCode != "" {
		device.Code = &Code[DeviceTypeCode, CodeSystemOID]{}
		device.Code.SetCode(deviceCode, CodeSystemOID(""), "", "")
	}

	// Set model and software
	if modelName != "" {
		device.ManufacturerModelName = &modelName
	}
	if softwareName != "" {
		device.SoftwareName = &softwareName
	}

	ro.Author.AssignedEntity.AssignedAuthorType.AssignedDevice = device

	return ro
}

// SetAuthoringOrganization sets the organization responsible for the observation author.
//
// Parameters:
//   - orgID: Optional organization globally known ID root
//   - orgName: Organization name
//   - roleID: Optional role-specific ID root
//   - roleExtension: Optional role-specific ID extension
//
// Returns the RelatedObservation for method chaining.
func (ro *RelatedObservation) SetAuthoringOrganization(orgID, orgName, roleID, roleExtension string) *RelatedObservation {
	if ro.Author == nil {
		ro.Author = &ObservationAuthor{}
	}

	org := &RepresentedAuthoringOrganization{}

	// Set organization ID
	if orgID != "" {
		org.ID = &ID{Root: orgID}
	}

	// Set organization name
	if orgName != "" {
		org.Name = &orgName
	}

	// Set role-specific identification
	if roleID != "" || roleExtension != "" {
		org.Identification = &OrganizationIdentification{
			ID: &ID{},
		}
		org.Identification.ID.SetID(roleID, roleExtension)
	}

	ro.Author.AssignedEntity.RepresentedAuthoringOrganization = org

	return ro
}

// SetDeviceManufacturer sets the manufacturer information for an author device.
//
// Parameters:
//   - manufacturerID: Manufacturer organization ID
//   - manufacturerName: Manufacturer organization name
//
// Returns the ObservationAssignedDevice for method chaining.
func (dev *ObservationAssignedDevice) SetDeviceManufacturer(manufacturerID, manufacturerName string) *ObservationAssignedDevice {
	if dev.PlayedManufacturedDevice == nil {
		dev.PlayedManufacturedDevice = &PlayedManufacturedDevice{}
	}

	mfg := &ObservationManufacturingOrganization{}

	if manufacturerID != "" {
		mfg.ID = &ID{Root: manufacturerID}
	}

	if manufacturerName != "" {
		mfg.Name = &manufacturerName
	}

	dev.PlayedManufacturedDevice.ManufacturingOrganization = mfg

	return dev
}
