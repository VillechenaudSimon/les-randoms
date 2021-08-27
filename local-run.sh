#!/bin/sh

# Fill the next lines with the database connection string, the discord app clientId, the discord app clientSecret and the website url
export DATABASE_CONNECTION_STRING=""
export DISCORD_CLIENTID=""
export DISCORD_CLIENTSECRET=""
export WEBSITE_URL=""
#

echo "Ctrl+C to stop running the local server"
echo "See on http://localhost:5000/"
go build -o bin/les-randoms.exe -v .
heroku local web -f Procfile.windows