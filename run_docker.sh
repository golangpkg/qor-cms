#!/bin/sh

# --restart=always
docker run -itd  --rm --link=mariadb \
    -e MYSQL_HOST=mariadb \
    -p 9000:9000 pythonstock/qor-cms:latest