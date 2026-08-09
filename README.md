# Klimaguessr

**Klimaguessr** is a fun little game, the core principle is guessing a location based on a climate graph. It features both singleplayer and multiplayer, but the multiplayer mode was designed for school, therefore you need host. A 1v1 is going to be implemented in the future.

## Code

Klimaguessr can be cut down into these parts:

- **Frontend**: TailwindCSS and Html/CSS/JS
- **Backend**: Go
- **Database**: PostgreSQL
- **Real-time Communication**: Socket.io

## Local Development

### Dev Environment 
1. **Clone & Enter**:
   ```bash
   git clone https://github.com/cns-studios/klimaguessr.git && cd ILoveConversion
   ```

2. **Environment Setup**:
   ```bash
   cp .env.example .env
   #  Put Debug to True
   ```

3. **Deployment**:
   ```bash
   go mod tidy 
   go run ./backend/
   ```
   
### Docker Deployment 
1. **Clone & Enter**:
   ```bash
   git clone https://github.com/cns-studios/klimaguessr.git && cd ILoveConversion
   ```

2. **Environment Setup**:
   ```bash
   cp .env.example .env
   ```

3. **Deployment**:
   ```bash
   docker-compose up -d --build
   ```

## Backend lookup
What code is in which file.

- **climate.go**: Direct climate data
- **config.go**: Configs for the server
- **db.go**: Everything involving the database
- **geo.go**: The math for Distances
- **handlers_http.go**: Handling all Https and http requests
- **handlers_socket.go**: Handling all socket.io connections
- **logging.go**: Everything about the logs
- **main.go**: The main file 
- **models.go**: The database structure
- **state.go**: Everything about lobby generation
- **sysinfo.go**: For an internal info page
- **utils**: Utils that are nice to have