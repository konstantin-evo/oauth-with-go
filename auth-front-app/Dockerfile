FROM node:18.18.0-alpine

WORKDIR /app

COPY public/ /app/public
COPY src/ /app/src
COPY package*.json /app/

RUN npm install --silent

CMD ["npm", "start"]
