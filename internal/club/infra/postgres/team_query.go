package postgres

const (
	queryInsertTeam = `
		INSERT INTO teams (id, name, logo_url, year_founded, address, city, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	queryFindTeamByID = `
		SELECT id, name, logo_url, year_founded, address, city, created_at, updated_at, deleted_at
		FROM teams
		WHERE id = $1 AND deleted_at IS NULL
	`

	queryFindAllTeams = `
		SELECT id, name, logo_url, year_founded, address, city, created_at, updated_at, deleted_at
		FROM teams
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	queryUpdateTeam = `
		UPDATE teams
		SET name = $1, logo_url = $2, year_founded = $3, address = $4, city = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	querySoftDeleteTeam = `UPDATE teams SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	queryExistsTeamByName = `
		SELECT EXISTS (
			SELECT 1 FROM teams
			WHERE name = $1 AND id != $2 AND deleted_at IS NULL
		)
	`
)
