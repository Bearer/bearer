import CloudWatchLogs from 'aws-sdk/clients/cloudwatchlogs'
import uuid from 'uuid/v1'
import { AWS_LAMBDA_LOG_STREAM_NAME, BEARER_XRAY_LOG_GROUP } from './constants'
import logger from './logger'

const cloudWatchLogsClientInstance = new CloudWatchLogs()

export const sendToCloudwatchGroup = async (payload: any, cloudWatchLogsClient = cloudWatchLogsClientInstance) => {
  try {
    const streamName = AWS_LAMBDA_LOG_STREAM_NAME!.concat('-').concat(uuid())
    const event = {
      logGroupName: BEARER_XRAY_LOG_GROUP,
      logStreamName: streamName,
      logEvents: [
        {
          timestamp: payload.timestamp,
          message: JSON.stringify(payload.message)
        }
      ]
    } as CloudWatchLogs.Types.PutLogEventsRequest
    await createLogStream(cloudWatchLogsClient, streamName)
    await cloudWatchLogsClient.putLogEvents(event).promise()
  } catch (error) {
    console.log(`error ${error.message}`)
  }
}

const createLogStream = async (client: CloudWatchLogs, streamName: string) => {
  try {
    return await client
      .createLogStream({
        logGroupName: BEARER_XRAY_LOG_GROUP!,
        logStreamName: streamName
      })
      .promise()
  } catch (error) {
    logger('%j', error)
    return undefined
  }
}
