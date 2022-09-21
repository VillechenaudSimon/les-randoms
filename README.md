
# les-randoms

## Dependencies

1. Golang (works on 1.19.1 - should be fine on later versions)   
2. PostgreSQL (works on 14.5) 

## Running Locally

1. Setup a `setup-secret-keys.sh` in order to allow the program to interact with some apis. Example key file at the end.

2. Start `local-run.sh` (works on unix / windows with gitbash)

3. Database will be created and initialised and website up on http://localhost:5000


Note : Some keys are **not** required for launch to succeed, but some parts of the app won't be accessible.

Exemple file for `setup-secret-keys.sh`.
```sh
#!/bin/sh
# Fill the next lines with the database connection string, the discord app clientId, the discord app clientSecret, the website url and the riot api token

# Depends on how your psql db is setup
export DATABASE_URL="host=localhost port=5432 user=postgres password=pwd dbname=blabla sslmode=disable"

# Those 3 keys should be obtained on https://discord.com/developers/applications
export DISCORD_CLIENTID="000000000000000000"
export DISCORD_CLIENTSECRET="000aaa000-000aaa000a_0000aaa0000"
export DISCORD_BOT_TOKEN="AAAAAAAA0000000000aaa000000aaa000X0000X000aaa000000"

# This API key should be obtained on https://developer.riotgames.com/apis
export X_RIOT_TOKEN="AZERT-aa000aa-0000-0000-0000-000aaa000aaa"

# These ones are on https://developer.spotify.com
export SPOTIFY_ID="0a0a0a0a0a0a0a0a0a0a0a0a0a0a"
export SPOTIFY_SECRET="a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a"
```
