#!/bin/bash
# Assign argument to dbName variable or default to demo
if [ $# -eq 0 ]; then
  dbName="demo"
else
  dbName=$1
fi

# Create a dump of the MongoDB database
mongodump -h localhost:27017 -d $dbName -o ../backup
