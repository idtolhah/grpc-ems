FROM node:14.18.1

WORKDIR /app
COPY comment-cmd/package.json .
RUN npm install

COPY comment-cmd/. .

CMD ["npm", "start"]