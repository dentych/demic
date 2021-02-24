FROM node:alpine

WORKDIR /app
ARG API_BASE_URL
RUN echo $API_BASE_URL

COPY package.json ./
RUN yarn
COPY . .
RUN yarn build

EXPOSE 3000

CMD yarn start
