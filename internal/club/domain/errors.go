package domain

import "errors"

// Team domain errors.
var (
	ErrTeamNotFound = errors.New("team not found")
)

// Player domain errors.
var (
	ErrPlayerNotFound    = errors.New("player not found")
	ErrJerseyNumberTaken = errors.New("jersey number already taken in this team")
)
