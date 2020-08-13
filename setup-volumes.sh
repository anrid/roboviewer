#!/bin/bash
docker network create robo-net
docker volume create --name robo-dg1-volume
# docker volume create --name local-volume --opt device=:/tmp/docker/local-volume