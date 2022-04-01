#!/bin/sh

. ./setup-secret-keys.sh

export DEBUG_MODE="TRUE"
export PORT="5000"
export WEBSITE_URL="http://localhost:$PORT"

echo "Ctrl+C to stop running the local server"
echo "See on $WEBSITE_URL"
# go build -o bin/les-randoms.exe -v . && heroku local web -f Procfile.windows
go build -o bin/les-randoms.exe -v . && ./bin/les-randoms.exe
# use 'winpty sqlite3 sqlite-database.db' to run sqlite3 command on git bash