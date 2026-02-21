package domain

import "errors"

// Match domain errors.
var (
	ErrMatchNotFound       = errors.New("match not found")
	ErrMatchResultNotFound = errors.New("match result not found")
	ErrResultAlreadyExists = errors.New("match result already reported")
	ErrSameTeam            = errors.New("home team and away team cannot be the same")
)
