# GooGL-Bye

Google recently announced the discontinuation of their URL shortening service. 

https://developers.googleblog.com/en/google-url-shortener-links-will-no-longer-be-available/

This Go project is designed to help developers prepare for this, by automating the process of finding and expanding goo.gl URLs in GitHub repositories, with plans to automatically raise issues/PRs for affected repositories.

## Features
- Search for repositories on GitHub (given criteria or specific Url/Author** )
- Clone repositories locally
- Parse cloned repositories for goo.gl URLs
- Expand found goo.gl URLs
- Save expanded URLs to a database
- (Planned) Automatically raise issues in repositories with goo.gl URLs**
- Web interface for signing in & scanning private repositories (or all repos in account)**

** Not yet implemented, but planned

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

2. Run with docker: 
```sh
# TODO
# docker compose up
# Inits a postgres instance, and runs the go application
```
