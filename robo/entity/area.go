package entity

// Area is a rectangular shaped area to clean.
type Area struct {
	Name  string `json:"name,omitempty"`   // Each cleaning area should definitely have a name to make reports nicer.
	SizeX int    `json:"size_x,omitempty"` // X side size in millimeters.
	SizeY int    `json:"size_y,omitempty"` // Y side size in millimeters.

	// Number of grid square passes needed before the square can be considered clean.
	PassesNeeded int `json:"passes_needed,omitempty"`

	Common
}

// NewArea creates a new area, e.g. a room or a corridor.
func NewArea(name string, sizeX, sizeY, passesNeeded int) *Area {
	return &Area{
		Name:         name,
		SizeX:        sizeX,
		SizeY:        sizeY,
		PassesNeeded: passesNeeded,
		Common: Common{
			UID:       "_:" + AreaUID,
			DType:     []string{"Area"},
			CreatedAt: now(),
		},
	}
}
