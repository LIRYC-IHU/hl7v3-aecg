package types

type SetCode[T ~string, U ~string] interface {
	SetCode(code T, codeSystem U, display string)
	GetCode() *Code[T, U]
	IsEmpty() bool
	String() string
}

func (c *Code[T, U]) SetCode(code T, codeSystem U, display string) {
	newCode := NewCode(code, codeSystem, display)
	*c = *newCode
}

func (c *Code[T, U]) GetCode() *Code[T, U] {
	return c
}

// IsEmpty returns true if the Code is empty (i.e., has no Code value).
func (c Code[T, U]) IsEmpty() bool {
	return c.Code == ""
}

// String returns the string representation of the Code.
func (c Code[T, U]) String() string {
	if c.DisplayName != "" {
		return c.DisplayName
	}
	return string(c.Code)
}
