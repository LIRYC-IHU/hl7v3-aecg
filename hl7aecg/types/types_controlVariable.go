package types

// ControlVariable captures related information about the subject or ECG collection conditions.
//
// Control variables capture related information that is not part of the core ECG data,
// such as the subject's age, fasting status, or other relevant observations.
//
// XML Structure:
//
//	<controlVariable>
//	  <relatedObservation>
//	    <code code="21612-7" codeSystem="2.16.840.1.113883.6.1" displayName="Reported Age"/>
//	    <value xsi:type="PQ" value="34" unit="a"/>
//	  </relatedObservation>
//	</controlVariable>
//
// Cardinality: Optional (0..* in Series)
type ControlVariable struct {
	// RelatedObservation contains the observation details.
	//
	// XML Tag: <relatedObservation>...</relatedObservation>
	// Cardinality: Required (within ControlVariable)
	RelatedObservation RelatedObservation `xml:"relatedObservation"`
}

// RelatedObservation represents an observation related to the ECG.
//
// This can capture various types of observations such as age, fasting status,
// or other clinical information relevant to the ECG interpretation.
//
// XML Structure:
//
//	<relatedObservation>
//	  <code code="21612-7" codeSystem="2.16.840.1.113883.6.1" displayName="Reported Age"/>
//	  <text>Subject age at time of ECG</text>
//	  <value xsi:type="PQ" value="34" unit="a"/>
//	  <author>...</author>
//	</relatedObservation>
//
// Cardinality: Required (within ControlVariable)
type RelatedObservation struct {
	// Code identifies the type of observation.
	//
	// Common codes include:
	//   - LOINC "21612-7": Reported Age
	//   - LOINC "49541-6": Fasting status
	//
	// XML Tag: <code code="..." codeSystem="..." displayName="..."/>
	// Cardinality: Optional
	Code *Code[string, CodeSystemOID] `xml:"code,omitempty"`

	// Text provides a textual description of the observation.
	//
	// This can be used when the observation is not a simple quantity
	// that can be put into the value element.
	//
	// XML Tag: <text>...</text>
	// Cardinality: Optional
	Text *string `xml:"text,omitempty"`

	// Value contains the observed value.
	//
	// For example, if the subject's age is observed, this would contain
	// the age value with appropriate units (e.g., "34 years").
	//
	// Note: The xsi:type attribute is typically "PQ" (Physical Quantity)
	//
	// XML Tag: <value xsi:type="PQ" value="..." unit="..."/>
	// Cardinality: Optional
	Value *PhysicalQuantity `xml:"value,omitempty"`

	// Author identifies who or what device recorded this observation.
	//
	// XML Tag: <author>...</author>
	// Cardinality: Optional
	Author *ObservationAuthor `xml:"author,omitempty"`
}

// ObservationAuthor identifies who or what device recorded the observation.
//
// XML Structure:
//
//	<author>
//	  <assignedEntity>...</assignedEntity>
//	</author>
//
// Cardinality: Optional (within RelatedObservation)
type ObservationAuthor struct {
	// AssignedEntity contains the author's identification and organization.
	//
	// XML Tag: <assignedEntity>...</assignedEntity>
	// Cardinality: Required (within ObservationAuthor)
	AssignedEntity AssignedEntity `xml:"assignedEntity"`
}

// AssignedEntity represents the person or device that authored an observation.
//
// XML Structure:
//
//	<assignedEntity>
//	  <id root="2.16.840.1.113883.3.5" extension="TECH_23"/>
//	  <assignedAuthorType>...</assignedAuthorType>
//	  <representedAuthoringOrganization>...</representedAuthoringOrganization>
//	</assignedEntity>
//
// Cardinality: Required (within ObservationAuthor)
type AssignedEntity struct {
	// ID is the unique identifier of the observation author.
	//
	// This is the ID assigned by the RepresentedAuthoringOrganization.
	// For example, if the authoring organization is the trial site,
	// this would be the ID assigned by that site.
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// AssignedAuthorType specifies whether the author is a person or device.
	//
	// XML Tag: <assignedAuthorType>...</assignedAuthorType>
	// Cardinality: Optional
	AssignedAuthorType *AssignedAuthorType `xml:"assignedAuthorType,omitempty"`

	// RepresentedAuthoringOrganization is the organization responsible for the author.
	//
	// XML Tag: <representedAuthoringOrganization>...</representedAuthoringOrganization>
	// Cardinality: Optional
	RepresentedAuthoringOrganization *RepresentedAuthoringOrganization `xml:"representedAuthoringOrganization,omitempty"`
}

// AssignedAuthorType specifies whether the author is a person or device.
//
// Only one of AssignedPerson or AssignedDevice should be set.
//
// XML Structure:
//
//	<assignedAuthorType>
//	  <assignedPerson>
//	    <name>JMK</name>
//	  </assignedPerson>
//	</assignedAuthorType>
//
// Cardinality: Optional (within AssignedEntity)
type AssignedAuthorType struct {
	// AssignedPerson identifies a person who authored the observation.
	//
	// XML Tag: <assignedPerson>...</assignedPerson>
	// Cardinality: Optional (mutually exclusive with AssignedDevice)
	AssignedPerson *ObservationAssignedPerson `xml:"assignedPerson,omitempty"`

	// AssignedDevice identifies a device that authored the observation.
	//
	// XML Tag: <assignedDevice>...</assignedDevice>
	// Cardinality: Optional (mutually exclusive with AssignedPerson)
	AssignedDevice *ObservationAssignedDevice `xml:"assignedDevice,omitempty"`
}

// ObservationAssignedPerson represents a person who authored an observation.
//
// XML Structure:
//
//	<assignedPerson>
//	  <name>JMK</name>
//	</assignedPerson>
//
// Cardinality: Optional (within AssignedAuthorType)
type ObservationAssignedPerson struct {
	// Name is the name of the person who authored the observation.
	//
	// Can be initials, full name, or structured PersonName.
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}

// ObservationAssignedDevice represents a device that authored an observation.
//
// XML Structure:
//
//	<assignedDevice>
//	  <id root="1.3.6.1.4.1.57054" extension="SN234"/>
//	  <code code="DEVICE_TYPE" codeSystem=""/>
//	  <manufacturerModelName>Device Model</manufacturerModelName>
//	  <softwareName>v1.0</softwareName>
//	  <playedManufacturedDevice>...</playedManufacturedDevice>
//	</assignedDevice>
//
// Cardinality: Optional (within AssignedAuthorType)
type ObservationAssignedDevice struct {
	// ID is the unique identifier of the device.
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Code identifies the type of device.
	//
	// Note: At the time of guide publishing, no formal vocabulary existed.
	//
	// XML Tag: <code code="..." codeSystem=""/>
	// Cardinality: Optional
	Code *Code[DeviceTypeCode, CodeSystemOID] `xml:"code,omitempty"`

	// ManufacturerModelName is the model name of the device.
	//
	// XML Tag: <manufacturerModelName>...</manufacturerModelName>
	// Cardinality: Optional
	ManufacturerModelName *string `xml:"manufacturerModelName,omitempty"`

	// SoftwareName is the name and version of software running in the device.
	//
	// XML Tag: <softwareName>...</softwareName>
	// Cardinality: Optional
	SoftwareName *string `xml:"softwareName,omitempty"`

	// PlayedManufacturedDevice contains the device manufacturer information.
	//
	// XML Tag: <playedManufacturedDevice>...</playedManufacturedDevice>
	// Cardinality: Optional
	PlayedManufacturedDevice *PlayedManufacturedDevice `xml:"playedManufacturedDevice,omitempty"`
}

// PlayedManufacturedDevice contains information about the device manufacturer.
//
// XML Structure:
//
//	<playedManufacturedDevice>
//	  <manufacturingOrganization>
//	    <id root="1.3.6.1.4.1.57054"/>
//	    <name>Device Manufacturer Inc.</name>
//	  </manufacturingOrganization>
//	</playedManufacturedDevice>
//
// Cardinality: Optional (within ObservationAssignedDevice)
type PlayedManufacturedDevice struct {
	// ManufacturingOrganization identifies the organization that manufactured the device.
	//
	// XML Tag: <manufacturingOrganization>...</manufacturingOrganization>
	// Cardinality: Optional
	ManufacturingOrganization *ObservationManufacturingOrganization `xml:"manufacturingOrganization,omitempty"`
}

// ObservationManufacturingOrganization represents the manufacturer of a device.
//
// XML Structure:
//
//	<manufacturingOrganization>
//	  <id root="1.3.6.1.4.1.57054"/>
//	  <name>Device Manufacturer Inc.</name>
//	</manufacturingOrganization>
//
// Cardinality: Optional (within PlayedManufacturedDevice)
type ObservationManufacturingOrganization struct {
	// ID is the unique identifier of the manufacturing organization.
	//
	// XML Tag: <id root="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Name is the name of the manufacturing organization.
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}

// RepresentedAuthoringOrganization is the organization responsible for the person or device.
//
// XML Structure:
//
//	<representedAuthoringOrganization>
//	  <name>1st Clinic of Milwaukee</name>
//	  <identification>
//	    <id root="2.16.840.1.113883.3.5" extension="SITE_1"/>
//	  </identification>
//	</representedAuthoringOrganization>
//
// Cardinality: Optional (within AssignedEntity)
type RepresentedAuthoringOrganization struct {
	// ID is the globally known unique identifier of the organization.
	//
	// This is an ID that would be globally known outside of the
	// organization's role in the trial.
	//
	// XML Tag: <id root="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Name is the name of the organization.
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`

	// Identification contains the role-specific ID of the organization.
	//
	// This is mostly an ID assigned by the trial or the sponsoring organization.
	//
	// XML Tag: <identification>...</identification>
	// Cardinality: Optional
	Identification *OrganizationIdentification `xml:"identification,omitempty"`
}

// OrganizationIdentification contains the role-specific ID of an organization.
//
// XML Structure:
//
//	<identification>
//	  <id root="2.16.840.1.113883.3.5" extension="SITE_1"/>
//	</identification>
//
// Cardinality: Optional (within RepresentedAuthoringOrganization)
type OrganizationIdentification struct {
	// ID is the role-specific identifier of the organization.
	//
	// This is most likely an ID assigned by the trial or sponsoring organization.
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`
}
