version: '3'
services:
  datadog-agent:
    image: datadog/agent:latest
    ports:
      - "8125:8125/udp"
      - "8126:8126/tcp"
    logging:
      driver: awslogs
      options:
        awslogs-group: /ecs/production/dev-portal-tag/statsd
        awslogs-region: $REGION
        awslogs-stream-prefix: ecs
    environment:
      - DD_API_KEY=2ed5212a659acd340565896aa765bbb
      - DD_APM_DD_URL=https://trace.agent.datadoghq.eu
    ulimits:
      nofile: 65000

  dev-portal-tag-production:
    image: $IMAGE_URL
    ports:
      - '8080:8080'
    logging:
      driver: awslogs
      options:
        awslogs-group: /ecs/production/dev-portal-tag
        awslogs-region: eu-west-1
        awslogs-stream-prefix: ecs
    environment:
      - BASE_DOMAIN=$PRODUCTIBASE_DOMAIN
      - BEARER_ASSETS_HOST=$BEARER_ASSETS_HOST
      - WEBHOOKS_ADDRESS=http://example.com/webhooks
