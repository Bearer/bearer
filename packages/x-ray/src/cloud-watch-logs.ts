import CloudWatchLogs from 'aws-sdk/clients/cloudwatchlogs'
import uuid from 'uuid/v1'
import { AWS_LAMBDA_LOG_STREAM_NAME, BEARER_XRAY_LOG_GROUP } from './constants'
import logger from './logger'

const errorDebug = logger.extend('error')

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
    logger(`send event %j`, event)
    const putLogEvent = await cloudWatchLogsClient.putLogEvents(event).promise()
    logger('%j', putLogEvent.$response.data)
    logger('%j', putLogEvent.rejectedLogEventsInfo)
    errorDebug('%j', putLogEvent.$response.error)
  } catch (error) {
    errorDebug(`error ${error.toString()}`)
  }
}

const createLogStream = async (client: CloudWatchLogs, streamName: string) => {
  try {
    logger(`create stream %s`, streamName)
    const stream = await client
      .createLogStream({
        logGroupName: BEARER_XRAY_LOG_GROUP!,
        logStreamName: streamName
      })
      .promise()
    logger('%j', stream.$response.data)
    errorDebug('%j', stream.$response.error)
    return stream
  } catch (error) {
    errorDebug(error.toString())
    return { $response: { data: {} } }
  }
}
