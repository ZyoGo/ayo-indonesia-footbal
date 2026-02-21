package postgres

const (
	queryStandings = `
		WITH team_stats AS (
			-- Stats when playing as home team
			SELECT 
				m.home_team_id AS team_id,
				COUNT(*) AS played,
				SUM(CASE WHEN mr.home_score > mr.away_score THEN 1 ELSE 0 END) AS won,
				SUM(CASE WHEN mr.home_score = mr.away_score THEN 1 ELSE 0 END) AS drawn,
				SUM(CASE WHEN mr.home_score < mr.away_score THEN 1 ELSE 0 END) AS lost,
				SUM(mr.home_score) AS gf,
				SUM(mr.away_score) AS ga
			FROM matches m
			JOIN match_results mr ON m.id = mr.match_id
			WHERE m.deleted_at IS NULL AND mr.deleted_at IS NULL
			GROUP BY m.home_team_id

			UNION ALL

			-- Stats when playing as away team
			SELECT 
				m.away_team_id AS team_id,
				COUNT(*) AS played,
				SUM(CASE WHEN mr.away_score > mr.home_score THEN 1 ELSE 0 END) AS won,
				SUM(CASE WHEN mr.away_score = mr.home_score THEN 1 ELSE 0 END) AS drawn,
				SUM(CASE WHEN mr.away_score < mr.home_score THEN 1 ELSE 0 END) AS lost,
				SUM(mr.away_score) AS gf,
				SUM(mr.home_score) AS ga
			FROM matches m
			JOIN match_results mr ON m.id = mr.match_id
			WHERE m.deleted_at IS NULL AND mr.deleted_at IS NULL
			GROUP BY m.away_team_id
		),
		aggregated_stats AS (
			SELECT 
				team_id,
				SUM(played) AS played,
				SUM(won) AS won,
				SUM(drawn) AS drawn,
				SUM(lost) AS lost,
				SUM(gf) AS gf,
				SUM(ga) AS ga,
				SUM(gf) - SUM(ga) AS gd,
				(SUM(won) * 3) + SUM(drawn) AS points
			FROM team_stats
			GROUP BY team_id
		)
		SELECT 
			a.team_id,
			t.name AS team_name,
			a.played,
			a.won,
			a.drawn,
			a.lost,
			a.gf,
			a.ga,
			a.gd,
			a.points
		FROM aggregated_stats a
		JOIN teams t ON t.id = a.team_id AND t.deleted_at IS NULL
		ORDER BY a.points DESC, a.gd DESC, a.gf DESC, t.name ASC
	`

	queryTopScorers = `
		SELECT 
			g.player_id,
			p.name AS player_name,
			t.name AS team_name,
			COUNT(*) AS goals
		FROM goals g
		JOIN players p ON p.id = g.player_id AND p.deleted_at IS NULL
		JOIN teams t ON t.id = g.team_id AND t.deleted_at IS NULL
		WHERE g.deleted_at IS NULL
		GROUP BY g.player_id, p.name, t.name
		ORDER BY goals DESC, p.name ASC
		LIMIT 20
	`
)
