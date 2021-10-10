#!/bin/sh

# Fill the next lines with the database connection string, the discord app clientId, the discord app clientSecret and the website url
export DATABASE_CONNECTION_STRING="217240:Rj4s7E!s*hBf@tcp(mysql-villechenaud-simon.alwaysdata.net:3306)/villechenaud-simon_les-randoms"
export DISCORD_CLIENTID="481156786779324416"
export DISCORD_CLIENTSECRET="Nlzu8wj6z_yb-IZIjeAzQZPZmsVEG8Hu"
export WEBSITE_URL="http://localhost:5000"
#

echo "Ctrl+C to stop running the local server"
echo "See on http://localhost:5000/"
go build -o bin/les-randoms.exe -v .
heroku local web -f Procfile.windows
