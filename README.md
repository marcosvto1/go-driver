GODRIVER

# Show coverage test and show output file
> go test -coverprofile=coverage.out ./internal/... -v

# Visualize navigator coverage tests output in html

go tool cover -html=coverage.out

# Create postgress container and map volume to temp folder container

docker run -it --rm -v $(pwd):/tmp postgres bash

(CONTAINER)>psql -U postgres
(CONTAINER)(psql)>create database godriver;
(CONTAINER)(psql)>create user godriver with encrypted password '12345678';
(CONTAINER)(psql)>grant all privileges on database godriver to godriver;

# Copy files for container running by docker compose

(HOST)> cd scripts/database
(HOST)> docker cp files.sql dbpostgres:/tmp
(HOST)> docker cp folders.sql dbpostgres:/tmp
(HOST)> docker cp users.sql dbpostgres:/tmp
(HOST)> docker exec -it dbpostgres bash
(CONTAINER)> cd tmp

# Create tables

(CONTAINER)> psql -U postgres godriver < users.sql
(CONTAINER)> psql -U postgres godriver < folders.sql
(CONTAINER)> psql -U postgres godriver < files.sql
