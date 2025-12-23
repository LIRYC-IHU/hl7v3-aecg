package types

import "github.com/google/uuid"

type Identifiable interface {
	SetID(id, extension string) Identifiable
	GetID() ID
	IsEmpty() bool
	String() string
}

// SetID sets the ID of the TrailSubject instance.
func (i *ID) SetID(id, extension string) {
	if id == "" {
		id = uuid.New().String()
	}
	if extension != "" {
		i.Extension = extension
	}
	i.Root = id
}

func (i ID) GetID() ID {
	return i
}

// IsEmpty returns true if the ID is empty (i.e., has no Root).
func (id ID) IsEmpty() bool {
	return id.Root == ""
}

// String returns the string representation of the ID in the format "Root^Extension".
func (id ID) String() string {
	if id.Extension != "" {
		return id.Root + "^" + id.Extension
	}
	return id.Root
}
