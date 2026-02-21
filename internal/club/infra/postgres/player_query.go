package postgres

const (
	queryInsertPlayer = `
		INSERT INTO players (id, team_id, name, height, weight, position, jersey_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	queryFindPlayerByID = `
		SELECT id, team_id, name, height, weight, position, jersey_number, created_at, updated_at, deleted_at
		FROM players
		WHERE id = $1 AND deleted_at IS NULL
	`

	queryFindPlayersByTeamID = `
		SELECT id, team_id, name, height, weight, position, jersey_number, created_at, updated_at, deleted_at
		FROM players
		WHERE team_id = $1 AND deleted_at IS NULL
		ORDER BY jersey_number ASC
	`

	queryUpdatePlayer = `
		UPDATE players
		SET name = $1, height = $2, weight = $3, position = $4, jersey_number = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	querySoftDeletePlayer = `UPDATE players SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	queryIsJerseyNumberTaken = `
		SELECT EXISTS(
			SELECT 1 FROM players
			WHERE team_id = $1 AND jersey_number = $2 AND deleted_at IS NULL AND id != $3
		)
	`
)
