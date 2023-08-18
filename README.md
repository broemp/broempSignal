# BroempSignal

A Discord and Telegram Bot that messages your friends for you!

## Running the Program

### Locally

Create an .env file or set the Envs according to the [.env.example](.env.example) and
run a postgres db

```sh
git clone https://github.com/broemp/broempSignal
go mod install
go run cmd/BroempSignal.go
```

### Docker
[Docker Hub](https://hub.docker.com/r/broemp/broemp_signal)

Use the [Docker Compose File](docker-compose.yml)
```yaml
version: '3.8'
services:
  broempSignal:
    image: broemp/broemp-signal:latest
    container_name: broempSignal
    restart: always
    environment:
      #DB Settings 
      - DB_DRIVER=postgres
      - DB_SOURCE=postgresql://{USER}:{PASSWORD}@{IP}:{PORT}/broempSignal?sslmode=disable
      #Server Settings
      - SERVER_ADDRESS=0.0.0.0:8080
```

# API Endpoints

## Users
```
# POST: Create User
/user
## Body
{
	"username": "{name}",
	"discordid": {id}
}

# GET: List all users
/users

# GET: List all data for one user
/users/{id}

# POST: Add telegramID to User
```

## AFK
```
# POST: Create AFK entry
/afk
## Body
{
  "discordid": {id}
}

# GET: List all AFKs order by count
/afk/list

# GET: AFK List for one user
/afk/list/{id}

# GET: AFK count for one user
/afk/count/{id}

```
