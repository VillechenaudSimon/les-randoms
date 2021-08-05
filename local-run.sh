#!/bin/sh

# Fill the next line with the database connection string
export DATABASE_CONNECTION_STRING=
#

echo "Ctrl+C to stop running the local server"
echo "See on http://localhost:5000/"
go build -o bin/les-randoms.exe -v .
heroku local web -f Procfile.windows