# Dockerfile for containers executing JS code
# for Letto.
#
# This container may be run like this:
#
# $ docker build -t my-nodejs-app .
# $ docker run -it --rm --name my-running-app my-nodejs-app
#
# The image is defined [here](https://store.docker.com/images/node)
# and [here on Github](https://github.com/nodejs/docker-node).
#
# The image assumes that your application has a file named 
# package.json listing its dependencies and defining its start script.
#
# You may run a simple script using the following command, without
# this Dockerfile:
#
# $ docker run -it --rm --name my-running-script -v "$PWD":/usr/src/app \
#   -w /usr/src/app node:4 node your-daemon-or-script.js

FROM node:4-onbuild
