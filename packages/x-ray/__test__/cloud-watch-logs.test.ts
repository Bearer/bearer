import { sendToCloudwatchGroup } from '../src/cloud-watch-logs'
import CloudwatchLogs from 'aws-sdk/clients/cloudwatchlogs'

jest.mock('../src/constants')
const cloudwatchLogs = new CloudwatchLogs()

describe('sendToCloudwatchGroup', () => {
  beforeAll(() => {
    cloudwatchLogs.putLogEvents = jest.fn(() => {
      return {
        promise: () => {
          return new Promise((resolve, _reject) => {
            resolve({})
          })
        }
      }
    }) as any

    cloudwatchLogs.createLogStream = jest.fn(() => {
      return {
        promise: () => {
          return new Promise((resolve, _reject) => {
            resolve({})
          })
        }
      }
    }) as any
  })

  it('sends logs to cloudwatch log log group', async () => {
    const event = {
      timestamp: new Date().getTime(),
      message: 'My message'
    }

    await sendToCloudwatchGroup(event, cloudwatchLogs)
    expect(true).toBeTruthy()
    expect(cloudwatchLogs.createLogStream).toBeCalledTimes(1)
    expect(cloudwatchLogs.putLogEvents).toBeCalledTimes(1)
    expect(cloudwatchLogs.putLogEvents).toBeCalledWith({
      logEvents: [
        {
          message: '"My message"',
          timestamp: expect.any(Number)
        }
      ],
      logGroupName: 'MyLOgGroup',
      logStreamName: expect.any(String)
    })
  })
})
