package info

// Logo type is a struct contains logo and information useful for output
type Logo struct {
	// Logo is a logo string
	Logo string

	// Lines is a logo lines count
	Lines int

	// MaxLength is a len on longest string in logo
	MaxLength int

	// AccentColor is a accent logo color
	AccentColor string
}
