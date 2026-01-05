package types

import (
	"fmt"
	"sync"
)

// Create Singloton rootID stringvar rootID string

type InstanceID struct {
	ID        string
	Extension string
}

var (
	instanceID *InstanceID
	once       sync.Once
)

// SetID sets the ID of the instance.
//
// Parameters:
//   - id: The root ID value (OID or custom identifier)
//   - extension: Optional extension to the ID
//
// Note: No automatic UUID generation. The caller must provide a valid ID.

func (i *ID) SetID(id, extension string) {
	if id == "" {
		if instanceID != nil {
			i.Root = instanceID.ID
		}
	} else {
		i.Root = id

	}
	if extension != "" {
		i.Extension = extension
	}
}

func SetRootID(id, extension string) *InstanceID {
	once.Do(func() {
		fmt.Println("Setting singleton InstanceID")
		instanceID = &InstanceID{ID: id, Extension: extension}
	})
	return instanceID
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
