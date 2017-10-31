FROM node:latest

RUN npm install request
RUN npm install ovh

WORKDIR /usr/src/app
VOLUME /usr/src/app/data

COPY ./main.js /usr/src/app/main.js

CMD [ "node" ]
