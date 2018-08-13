FROM node:10.8-alpine
RUN apk --update add git bash
RUN mkdir /app
WORKDIR /app
COPY . .
RUN yarn install
RUN yarn run lerna bootstrap
RUN echo '//registry.npmjs.org/:_authToken=${NPM_TOKEN}' > /app/.npmrc
