FROM node:10.8-alpine
RUN apk --update add git bash openssh-client

RUN mkdir /root/.ssh/
RUN touch /root/.ssh/known_hosts

RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

RUN mkdir /app
WORKDIR /app

COPY . .

RUN yarn install --frozen-lockfile
RUN chmod +x /app/scripts/release-package.sh
CMD ["/app/scripts/release-package.sh"]
