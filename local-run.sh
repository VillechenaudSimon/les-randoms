#!/bin/sh

. ./setup-secret-keys.sh

export DEBUG_MODE="TRUE"
export WEBSITE_URL="http://localhost:5000"

echo "Ctrl+C to stop running the local server"
echo "See on $WEBSITE_URL"
go build -o bin/les-randoms.exe -v . && heroku local web -f Procfile.windows

# use 'winpty sqlite3 sqlite-database.db' to run sqlite3 command on git bash