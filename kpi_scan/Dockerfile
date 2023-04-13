FROM ubuntu

RUN apt-get update && apt-get install -y curl git jq

RUN mkdir /app
ADD run.sh /app/
WORKDIR /app

CMD ["/app/run.sh"]
