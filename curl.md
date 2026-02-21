# API Endpoints & Curl Commands

This document provides a comprehensive list of all API endpoints for the Ayo Indonesia project, ordered logically by module and flow.

**Base URL:** `http://localhost:4000/api/v1`

**Note:** Protected routes require `Authorization: Bearer <token>` header. Get the token from `/auth/login`.

## 1. Authentication

### Register
```bash
curl -X POST http://localhost:4000/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin",
       "password": "password123"
     }'
```

### Login
```bash
curl -X POST http://localhost:4000/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "username": "admin",
       "password": "password123"
     }'
```

---

## 2. Club Management

### Teams

#### Create Team
```bash
curl -X POST http://localhost:4000/api/v1/teams \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
     -d '{
       "name": "Persija Jakarta",
       "logo_url": "https://example.com/persija.png",
       "year_founded": 1928,
       "address": "Jl. Rasuna Said",
       "city": "Jakarta"
     }'
```

#### Get All Teams
```bash
curl -X GET http://localhost:4000/api/v1/teams
```

#### Get Team by ID
```bash
curl -X GET http://localhost:4000/api/v1/teams/{team_id}
```

#### Update Team
```bash
curl -X PUT http://localhost:4000/api/v1/teams/{team_id} \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
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
curl -X DELETE http://localhost:4000/api/v1/teams/{team_id} \
     -H "Authorization: Bearer <token>"
```

---

### Players

#### Create Player
```bash
curl -X POST http://localhost:4000/api/v1/players \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
     -d '{
       "team_id": "{team_id}",
       "name": "Bambang Pamungkas",
       "jersey_number": 20,
       "position": "CF",
       "height": 170.0,
       "weight": 68.5
     }'
```

#### Get Player by ID
```bash
curl -X GET http://localhost:4000/api/v1/players/{player_id}
```

#### Get Players by Team ID
```bash
curl -X GET http://localhost:4000/api/v1/teams/{team_id}/players
```

#### Update Player
```bash
curl -X PUT http://localhost:4000/api/v1/players/{player_id} \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
     -d '{
       "team_id": "{team_id}",
       "name": "Bambang Pamungkas Updated",
       "jersey_number": 20,
       "position": "CF",
       "height": 171.0,
       "weight": 69.0
     }'
```

#### Delete Player
```bash
curl -X DELETE http://localhost:4000/api/v1/players/{player_id} \
     -H "Authorization: Bearer <token>"
```

---

## 2. Match Management

### Create Match
```bash
curl -X POST http://localhost:4000/api/v1/matches \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
     -d '{
       "home_team_id": "{home_team_id}",
       "away_team_id": "{away_team_id}",
       "match_date": "2026-10-15",
       "match_time": "19:00",
       "stadium": "Gelora Bung Karno"
     }'
```

### Get All Matches
```bash
curl -X GET http://localhost:4000/api/v1/matches
```

### Get Match by ID
```bash
curl -X GET http://localhost:4000/api/v1/matches/{match_id}
```

### Report Match Result
```bash
curl -X POST http://localhost:4000/api/v1/matches/{match_id}/result \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <token>" \
     -d '{
       "home_score": 1,
       "away_score": 0,
       "goals": [
         {
           "player_id": "{player_id}",
           "team_id": "{home_team_id}",
           "goal_minute": 45
         }
       ]
     }'
```

### Get Match Report
```bash
curl -X GET http://localhost:4000/api/v1/matches/{match_id}/report
```

### Get All Match Reports
```bash
curl -X GET http://localhost:4000/api/v1/reports/matches
```

---

## 3. Global Reporting

### Get League Standings
```bash
curl -X GET http://localhost:4000/api/v1/reporting/standings
```

### Get Top Scorers
```bash
curl -X GET http://localhost:4000/api/v1/reporting/top-scorers
```

---

## 4. Upload

### Upload File
```bash
curl -X POST http://localhost:4000/api/v1/uploads \
     -H "Authorization: Bearer <token>" \
     -F "file=@/path/to/file.jpg"
```
