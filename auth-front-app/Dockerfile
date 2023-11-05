FROM node:18.18.0-alpine

WORKDIR /app

COPY public/ /app/public/
COPY src/ /app/src/
COPY package.json /app/package.json
COPY server.js /app/server.js

RUN npm install
RUN npm run build

CMD ["npm", "start"]
