# GooGL-Bye

Google recently announced the discontinuation of their URL shortening service. 

https://developers.googleblog.com/en/google-url-shortener-links-will-no-longer-be-available/

This Go project is designed to help developers prepare for this, by automating the process of finding and expanding goo.gl URLs in GitHub repositories, with plans to automatically raise issues/PRs for affected repositories.

## Features
- Search for repositories/users on GitHub 
- Clone repositories locally
- Parse cloned repositories for goo.gl URLs
- Expand found goo.gl URLs
- Save expanded URLs to a database
- Automatically raise issues in repositories with goo.gl URLs
- Web interface for signing in & scanning private repositories (or all repos in account)**

** Partially implemented. Users can search public repos without logging in.

## Prerequisites
- Go 1.22 or higher
- Git
- GitHub API access token [(docs)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-personal-access-token-classic)
- Database (e.g., PostgreSQL, MySQL)

## Installation

1. Clone this repository:
```sh
git clone https://github.com/jwtly10/googl-bye.git
cd googl-bye
```

2. Create a `.env` file in the root directory with the following environment variables:
```sh
ENV=dev # Environment for configuring logs/paths (dev/prod)
DB_NAME=googl-bye-db
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
PARSER_INTERVAL=30 # interval for the parser job to run in seconds (default 30)
GH_TOKEN= # GitHub API access classic token (needs repo.public_repo scope)
```
(note that by default the docker-compose file handles the database setup, only GH_TOKEN is required)

3. Run with docker: 
```sh
docker-compose up --build
```

### Run without docker in dev:

It may be quicker and more dev-friendly to run the services without using docker.

Given db credentials are set in .env, and db/init.sql has been run, you can run the following commands to start the services:

```sh
# Start frontend nodemon service
cd react 
npm i
npm run watch

# Run backend go app
cd ..
go run cmd/server/main.go
```

The application should now be running on the following URLs:

| Component | URL                   | Description                  |
|-----------|------------------------|------------------------------|
| Frontend  | http://localhost:8080  | React application frontend   |
| API       | http://localhost:8080/v1/api | Backend API endpoints |
| pgAdmin   | http://localhost:5050  | PostgreSQL database management |


