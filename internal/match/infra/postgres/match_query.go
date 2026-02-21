package postgres

const (
	queryInsertMatch = `
		INSERT INTO matches (id, home_team_id, away_team_id, match_date, match_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	queryFindMatchByID = `
		SELECT id, home_team_id, away_team_id, match_date, match_time, created_at, updated_at, deleted_at
		FROM matches
		WHERE id = $1 AND deleted_at IS NULL
	`

	queryFindAllMatches = `
		SELECT id, home_team_id, away_team_id, match_date, match_time, created_at, updated_at, deleted_at
		FROM matches
		WHERE deleted_at IS NULL
		ORDER BY match_date DESC, match_time DESC
	`

	queryUpdateMatch = `
		UPDATE matches
		SET home_team_id = $1, away_team_id = $2, match_date = $3, match_time = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	queryDeleteMatch = `UPDATE matches SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	queryInsertMatchResult = `
		INSERT INTO match_results (id, match_id, home_score, away_score)
		VALUES ($1, $2, $3, $4)
	`

	queryInsertGoal = `
		INSERT INTO goals (id, result_id, player_id, team_id, goal_minute)
		VALUES ($1, $2, $3, $4, $5)
	`

	queryFindResultByMatchID = `
		SELECT id, match_id, home_score, away_score, deleted_at
		FROM match_results
		WHERE match_id = $1 AND deleted_at IS NULL
	`

	queryFindGoalsByResultID = `
		SELECT g.id, g.result_id, g.player_id, p.name AS player_name, g.team_id, g.goal_minute, g.deleted_at
		FROM goals g
		JOIN players p ON p.id = g.player_id
		WHERE g.result_id = $1 AND g.deleted_at IS NULL
		ORDER BY g.goal_minute ASC
	`

	queryExistsResultByMatchID = `
		SELECT EXISTS(SELECT 1 FROM match_results WHERE match_id = $1 AND deleted_at IS NULL)
	`

	queryMatchReport = `
		SELECT
			m.id AS match_id,
			TO_CHAR(m.match_date, 'YYYY-MM-DD') AS match_date,
			m.match_time,
			ht.name AS home_team_name,
			at.name AS away_team_name,
			mr.home_score,
			mr.away_score,
			CASE
				WHEN mr.home_score > mr.away_score THEN 'Home Win'
				WHEN mr.away_score > mr.home_score THEN 'Away Win'
				ELSE 'Draw'
			END AS match_status,
			COALESCE(ts.player_name, '') AS top_scorer,
			COALESCE(ts.goal_count, 0) AS top_scorer_goals,
			(SELECT COUNT(*) FROM match_results mr2
				JOIN matches m2 ON m2.id = mr2.match_id
				WHERE (m2.home_team_id = m.home_team_id AND mr2.home_score > mr2.away_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
				   OR (m2.away_team_id = m.home_team_id AND mr2.away_score > mr2.home_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
			) AS home_team_wins,
			(SELECT COUNT(*) FROM match_results mr2
				JOIN matches m2 ON m2.id = mr2.match_id
				WHERE (m2.home_team_id = m.away_team_id AND mr2.home_score > mr2.away_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
				   OR (m2.away_team_id = m.away_team_id AND mr2.away_score > mr2.home_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
			) AS away_team_wins
		FROM matches m
		JOIN teams ht ON ht.id = m.home_team_id AND ht.deleted_at IS NULL
		JOIN teams at ON at.id = m.away_team_id AND at.deleted_at IS NULL
		JOIN match_results mr ON mr.match_id = m.id AND mr.deleted_at IS NULL
		LEFT JOIN LATERAL (
			SELECT p.name AS player_name, COUNT(*) AS goal_count
			FROM goals g
			JOIN players p ON p.id = g.player_id AND p.deleted_at IS NULL
			WHERE g.result_id = mr.id AND g.deleted_at IS NULL
			GROUP BY p.name
			ORDER BY goal_count DESC
			LIMIT 1
		) ts ON TRUE
		WHERE m.id = $1 AND m.deleted_at IS NULL
	`

	queryAllMatchReports = `
		SELECT
			m.id AS match_id,
			TO_CHAR(m.match_date, 'YYYY-MM-DD') AS match_date,
			m.match_time,
			ht.name AS home_team_name,
			at.name AS away_team_name,
			mr.home_score,
			mr.away_score,
			CASE
				WHEN mr.home_score > mr.away_score THEN 'Home Win'
				WHEN mr.away_score > mr.home_score THEN 'Away Win'
				ELSE 'Draw'
			END AS match_status,
			COALESCE(ts.player_name, '') AS top_scorer,
			COALESCE(ts.goal_count, 0) AS top_scorer_goals,
			(SELECT COUNT(*) FROM match_results mr2
				JOIN matches m2 ON m2.id = mr2.match_id
				WHERE (m2.home_team_id = m.home_team_id AND mr2.home_score > mr2.away_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
				   OR (m2.away_team_id = m.home_team_id AND mr2.away_score > mr2.home_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
			) AS home_team_wins,
			(SELECT COUNT(*) FROM match_results mr2
				JOIN matches m2 ON m2.id = mr2.match_id
				WHERE (m2.home_team_id = m.away_team_id AND mr2.home_score > mr2.away_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
				   OR (m2.away_team_id = m.away_team_id AND mr2.away_score > mr2.home_score AND m2.deleted_at IS NULL AND mr2.deleted_at IS NULL)
			) AS away_team_wins
		FROM matches m
		JOIN teams ht ON ht.id = m.home_team_id AND ht.deleted_at IS NULL
		JOIN teams at ON at.id = m.away_team_id AND at.deleted_at IS NULL
		JOIN match_results mr ON mr.match_id = m.id AND mr.deleted_at IS NULL
		LEFT JOIN LATERAL (
			SELECT p.name AS player_name, COUNT(*) AS goal_count
			FROM goals g
			JOIN players p ON p.id = g.player_id AND p.deleted_at IS NULL
			WHERE g.result_id = mr.id AND g.deleted_at IS NULL
			GROUP BY p.name
			ORDER BY goal_count DESC
			LIMIT 1
		) ts ON TRUE
		WHERE m.deleted_at IS NULL
		ORDER BY m.match_date DESC
	`
	querySoftDeleteResultByMatchID = `UPDATE match_results SET deleted_at = NOW() WHERE match_id = $1 AND deleted_at IS NULL`
	querySoftDeleteGoalsByMatchID  = `
		UPDATE goals SET deleted_at = NOW() 
		WHERE result_id IN (SELECT id FROM match_results WHERE match_id = $1) 
		AND deleted_at IS NULL
	`
)
