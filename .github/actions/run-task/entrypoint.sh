#!/usr/bin/env sh

set -e

for i in `seq 1 ${CONCURRENT_TASKS}`
do
  aws ecs run-task \
      --cluster ${CLUSTER} \
      --count 1 \
      --tags key=service,value=${TASK_DEFINITION} \
      --network-configuration "awsvpcConfiguration={subnets=['${SUBNET}'],securityGroups=['${SECURITY_GROUP}'],assignPublicIp=ENABLED}" \
      --launch-type FARGATE \
      --region eu-west-1 \
      --task-definition ${TASK_DEFINITION}
done