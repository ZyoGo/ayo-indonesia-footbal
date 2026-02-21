package domain

import (
	"strings"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/ulid"
)

type Position int

const (
	PositionGK  Position = iota + 1 // Goalkeeper
	PositionCB                       // Center Back
	PositionLB                       // Left Back
	PositionRB                       // Right Back
	PositionLWB                      // Left Wing Back
	PositionRWB                      // Right Wing Back
	PositionCDM                      // Central Defensive Midfielder
	PositionCM                       // Central Midfielder
	PositionCAM                      // Central Attacking Midfielder
	PositionLM                       // Left Midfielder
	PositionRM                       // Right Midfielder
	PositionLW                       // Left Winger
	PositionRW                       // Right Winger
	PositionCF                       // Center Forward
	PositionST                       // Striker
	PositionSS                       // Second Striker
)

var positionNames = map[Position]string{
	PositionGK:  "GK",
	PositionCB:  "CB",
	PositionLB:  "LB",
	PositionRB:  "RB",
	PositionLWB: "LWB",
	PositionRWB: "RWB",
	PositionCDM: "CDM",
	PositionCM:  "CM",
	PositionCAM: "CAM",
	PositionLM:  "LM",
	PositionRM:  "RM",
	PositionLW:  "LW",
	PositionRW:  "RW",
	PositionCF:  "CF",
	PositionST:  "ST",
	PositionSS:  "SS",
}

var positionValues = map[string]Position{
	"GK":  PositionGK,
	"CB":  PositionCB,
	"LB":  PositionLB,
	"RB":  PositionRB,
	"LWB": PositionLWB,
	"RWB": PositionRWB,
	"CDM": PositionCDM,
	"CM":  PositionCM,
	"CAM": PositionCAM,
	"LM":  PositionLM,
	"RM":  PositionRM,
	"LW":  PositionLW,
	"RW":  PositionRW,
	"CF":  PositionCF,
	"ST":  PositionST,
	"SS":  PositionSS,
}

func (p Position) String() string {
	if name, ok := positionNames[p]; ok {
		return name
	}
	return ""
}

func ParsePosition(s string) (Position, bool) {
	p, ok := positionValues[strings.ToUpper(strings.TrimSpace(s))]
	return p, ok
}

func (p Position) IsValid() bool {
	_, ok := positionNames[p]
	return ok
}

const (
	maxPlayerNameLength = 255
	maxHeight           = 300.0 // cm
	maxWeight           = 300.0 // kg
)

type Player struct {
	ID           string
	TeamID       string
	Name         string
	Height       float64
	Weight       float64
	Position     Position
	JerseyNumber int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewPlayer(teamID, name string, height, weight float64, position Position, jerseyNumber int) (*Player, error) {
	name = strings.TrimSpace(name)

	if teamID == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "team ID is required")
	}
	if name == "" {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "player name is required")
	}
	if len(name) > maxPlayerNameLength {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "player name must not exceed %d characters", maxPlayerNameLength)
	}
	if height < 0 || height > maxHeight {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "height must be between 0 and %.0f cm", maxHeight)
	}
	if weight < 0 || weight > maxWeight {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "weight must be between 0 and %.0f kg", maxWeight)
	}
	if !position.IsValid() {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "invalid position: must be a valid football position (e.g. GK, CB, CM, ST)")
	}
	if jerseyNumber <= 0 || jerseyNumber > 99 {
		return nil, derrors.NewErrorf(derrors.ErrorCodeBadRequest, "jersey number must be between 1 and 99")
	}

	now := time.Now()
	return &Player{
		ID:           ulid.GenerateID(),
		TeamID:       teamID,
		Name:         name,
		Height:       height,
		Weight:       weight,
		Position:     position,
		JerseyNumber: jerseyNumber,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (p *Player) Update(name string, height, weight float64, position Position, jerseyNumber int) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "player name is required")
	}
	if len(name) > maxPlayerNameLength {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "player name must not exceed %d characters", maxPlayerNameLength)
	}
	if height < 0 || height > maxHeight {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "height must be between 0 and %.0f cm", maxHeight)
	}
	if weight < 0 || weight > maxWeight {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "weight must be between 0 and %.0f kg", maxWeight)
	}
	if !position.IsValid() {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "invalid position: must be a valid football position (e.g. GK, CB, CM, ST)")
	}
	if jerseyNumber <= 0 || jerseyNumber > 99 {
		return derrors.NewErrorf(derrors.ErrorCodeBadRequest, "jersey number must be between 1 and 99")
	}

	p.Name = name
	p.Height = height
	p.Weight = weight
	p.Position = position
	p.JerseyNumber = jerseyNumber
	p.UpdatedAt = time.Now()
	return nil
}
