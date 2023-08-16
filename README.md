# BroempSignal

A Discord and Telegram Bot that messages your friends for you!

## Running the Program

### Locally

Create an .env file or set the Envs according to the [.env.example](.env.example) and
Create an sqlite Database
Then run the following Commands

```sh
git clone https://github.com/broemp/broempSignal
go mod install
go run BroempSignal.go
```

### Docker
[Docker Hub](https://hub.docker.com/r/broemp/broemp_signal)

Use the [Docker Compose File](docker-compose.yml)
```yaml
version: '3.8'
services:
  broempSignal:
    image: broemp/broemp_signal:latest
    container_name: broempSignal
    restart: always
    environment:
      # Discord Bot Token https://discord.com/developers/applications
      - DISCORD_TOKEN=<token>
      # Discord Server ID
      - GUILD_ID=<guild_id>
      # Telegram Bot Token https://telegram.me/BotFather
      - TELEGRAM_TOKEN=<token>
    # or use env_file
    #env_file:
    #  - ".env"
    volumes:
      # Make DB persistent 
      - ./data.db:/app/data.db
```
 or run
```
docker run \
  -e DISCORD_TOKEN=<token> \
  -e GUILD_ID=<guild_id> \
  -e TELEGRAM_TOKEN=<token> \
  -v ./data.db:/app/data.db \
  broemp/broemp_signal 
```