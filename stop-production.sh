#!/bin/sh

kill $(ps -e | grep "les-randoms" | awk '{print $1}') && echo "Program successfully HARD closed"
