#!/bin/bash

echo "Starting Web Crawler Interface..."
echo "Building project..."
make build

echo "Starting web server on port 8080..."
./bin/web

echo "Web server stopped." 