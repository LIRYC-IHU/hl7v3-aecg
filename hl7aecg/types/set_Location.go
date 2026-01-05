package types

// SetName sets the name of the trial site location.
//
// Example: SetName("1st Clinic of Milwaukee")
func (s *SiteLocation) SetName(name string) *SiteLocation {
	s.Name = &name
	return s
}

// SetAddress initializes the address for the site location.
//
// Use this to start configuring the address, then chain with SetCity, SetState, SetCountry.
//
// Example: SetAddress().SetCity("Milwaukee").SetState("WI").SetCountry("USA")
func (s *SiteLocation) SetAddress() *SiteLocation {
	if s.Addr == nil {
		s.Addr = &Address{}
	}
	return s
}

// SetCity sets the city for the site address.
//
// Automatically initializes the address if needed.
//
// Example: SetCity("Milwaukee")
func (s *SiteLocation) SetCity(city string) *SiteLocation {
	if s.Addr == nil {
		s.Addr = &Address{}
	}
	s.Addr.City = &city
	return s
}

// SetState sets the state or province for the site address.
//
// Automatically initializes the address if needed.
//
// Example: SetState("WI")
func (s *SiteLocation) SetState(state string) *SiteLocation {
	if s.Addr == nil {
		s.Addr = &Address{}
	}
	s.Addr.State = &state
	return s
}

// SetCountry sets the country for the site address.
//
// Automatically initializes the address if needed.
//
// Example: SetCountry("USA")
func (s *SiteLocation) SetCountry(country string) *SiteLocation {
	if s.Addr == nil {
		s.Addr = &Address{}
	}
	s.Addr.Country = &country
	return s
}

// SetFullAddress sets the complete address in one call.
//
// Example: SetFullAddress("Milwaukee", "WI", "USA")
func (s *SiteLocation) SetFullAddress(city, state, country string) *SiteLocation {
	s.Addr = &Address{
		City:    &city,
		State:   &state,
		Country: &country,
	}
	return s
}

// =============================================================================
// Address setters
// =============================================================================

// SetCity sets the city for the address.
//
// Example: SetCity("Milwaukee")
func (a *Address) SetCity(city string) *Address {
	a.City = &city
	return a
}

// SetState sets the state or province for the address.
//
// Example: SetState("WI")
func (a *Address) SetState(state string) *Address {
	a.State = &state
	return a
}

// SetCountry sets the country for the address.
//
// Example: SetCountry("USA")
func (a *Address) SetCountry(country string) *Address {
	a.Country = &country
	return a
}
