#!/bin/sh

chrome.exe http://vemuni.ga/ https://github.com/Vemuni/les-randoms http://localhost:5000/ &
code . &

# Use this git bash to git add/commit/push and to ./local-run.sh
git-bash.exe &

# Use this git bash to heroku login and git push heroku main and git logs --tail
git-bash.exe &