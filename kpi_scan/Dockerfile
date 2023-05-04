FROM ubuntu

RUN apt-get update && apt-get install -y curl git jq ca-certificates
RUN update-ca-certificates
RUN apt-get install apt-transport-https

RUN echo "deb [trusted=yes] https://apt.fury.io/bearer/ /" > /etc/apt/sources.list.d/fury.list
RUN apt-get update && apt-get install -y bearer

RUN mkdir /app
ADD ./kpi_scan/run.sh /app/
WORKDIR /app

CMD ["/app/run.sh"]
