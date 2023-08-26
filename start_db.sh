#!/bin/bash
docker run \
    --rm --name sandbox \
    -p 5432:5432 \
    -d \
    -e "POSTGRES_DB=pgxsandbox" \
    -e "POSTGRES_HOST_AUTH_METHOD=trust" \
    postgres:15.4-alpine
