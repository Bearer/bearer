import * as React from 'react'
type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'
import { withFunctionCall } from '../src/withFunctionCall'

const Factory = (integrationId: string) => {
  return {
    withFunctionCall: function<TReturnedData = any>(functionName: string) {
      return withFunctionCall<TReturnedData>(integrationId, functionName)
    },
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
