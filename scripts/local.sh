#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

# postgresql
export POSTGRESQL_IMAGE=learn-go-database-postgresql
export POSTGRESQL_IMAGE_TAG=development
export POSTGRESQL_CONTAINER=learn-go-database-postgresql
export POSTGRESQL_HOST=learn-go-database-postgresql.service
export POSTGRESQL_USER=learn-go-database
export POSTGRESQL_PASSWORD=password
export POSTGRESQL_DB=learn-go-database
export POSTGRESQL_TIME_ZONE=Asia/Kuala_Lumpur
docker build -t "$POSTGRESQL_IMAGE:$POSTGRESQL_IMAGE_TAG" -f ./manifest-docker/Dockerfile.postgresql ./manifest-docker

# pgadmin
export PGADMIN_IMAGE=learn-go-database-pgadmin
export PGADMIN_IMAGE_TAG=development
export PGADMIN_CONTAINER=learn-go-database-pgadmin
export PGADMIN_HOST=learn-go-database-pgadmin.service
export PGADMIN_EMAIL=admin@golang.org
export PGADMIN_PASSWORD=123
docker build -t "$PGADMIN_IMAGE:$PGADMIN_IMAGE_TAG" -f ./manifest-docker/Dockerfile.pgadmin ./manifest-docker

docker-compose -f ./manifest/docker-compose.yml up -d --build
