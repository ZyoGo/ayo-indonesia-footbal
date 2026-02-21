# API Endpoints & Curl Commands

This document provides a comprehensive list of all API endpoints for the Ayo Indonesia project, ordered logically by module and flow.

## 1. Club Management

### Teams

#### Create Team
```bash
curl -X POST http://localhost:8080/api/v1/teams \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Persija Jakarta",
       "logo_url": "https://example.com/persija.png",
       "year_founded": 1928,
       "address": "Jakarta",
       "city": "Jakarta"
     }'
```

#### Get All Teams
```bash
curl -X GET http://localhost:8080/api/v1/teams
```

#### Get Team by ID
```bash
curl -X GET http://localhost:8080/api/v1/teams/{team_id}
```

#### Update Team
```bash
curl -X PUT http://localhost:8080/api/v1/teams/{team_id} \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Persija Jakarta Updated",
       "logo_url": "https://example.com/persija_new.png",
       "year_founded": 1928,
       "address": "Jakarta Southern",
       "city": "Jakarta"
     }'
```

#### Delete Team
```bash
curl -X DELETE http://localhost:8080/api/v1/teams/{team_id}
```

---

### Players

#### Create Player
```bash
curl -X POST http://localhost:8080/api/v1/players \
     -H "Content-Type: application/json" \
     -d '{
       "team_id": "{team_id}",
       "name": "Bambang Pamungkas",
       "height": 170.5,
       "weight": 68.0,
       "position": "FW",
       "jersey_number": 20
     }'
```

#### Get Player by ID
```bash
curl -X GET http://localhost:8080/api/v1/players/{player_id}
```

#### Get Players by Team ID
```bash
curl -X GET http://localhost:8080/api/v1/teams/{team_id}/players
```

#### Update Player
```bash
curl -X PUT http://localhost:8080/api/v1/players/{player_id} \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Bambang Pamungkas Updated",
       "height": 171.0,
       "weight": 69.0,
       "position": "FW",
       "jersey_number": 20
     }'
```

#### Delete Player
```bash
curl -X DELETE http://localhost:8080/api/v1/players/{player_id}
```

---

## 2. Match Management

### Create Match
```bash
curl -X POST http://localhost:8080/api/v1/matches \
     -H "Content-Type: application/json" \
     -d '{
       "home_team_id": "{home_team_id}",
       "away_team_id": "{away_team_id}",
       "match_date": "2026-03-01",
       "match_time": "19:00"
     }'
```

### Get All Matches
```bash
curl -X GET http://localhost:8080/api/v1/matches
```

### Get Match by ID
```bash
curl -X GET http://localhost:8080/api/v1/matches/{match_id}
```

### Report Match Result
```bash
curl -X POST http://localhost:8080/api/v1/matches/{match_id}/result \
     -H "Content-Type: application/json" \
     -d '{
       "home_score": 2,
       "away_score": 1,
       "goals": [
         {
           "player_id": "{player_id_1}",
           "team_id": "{home_team_id}",
           "goal_minute": 23
         },
         {
           "player_id": "{player_id_2}",
           "team_id": "{home_team_id}",
           "goal_minute": 45
         },
         {
           "player_id": "{player_id_3}",
           "team_id": "{away_team_id}",
           "goal_minute": 10
         }
       ]
     }'
```

### Get Match Report
```bash
curl -X GET http://localhost:8080/api/v1/matches/{match_id}/report
```

### Get All Match Reports
```bash
curl -X GET http://localhost:8080/api/v1/reports/matches
```

---

## 3. Global Reporting

### Get League Standings
```bash
curl -X GET http://localhost:8080/api/v1/reporting/standings
```

### Get Top Scorers
```bash
curl -X GET http://localhost:8080/api/v1/reporting/top-scorers
```
