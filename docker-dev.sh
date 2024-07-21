#!/bin/bash

docker stop googl-bye
docker rm googl-bye
docker rmi googl-bye
docker build --tag googl-bye .
docker run --name googl-bye --env-file .env --detach --mount type=bind,source="$(pwd)",target=/app googl-bye