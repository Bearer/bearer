import * as React from 'react'
type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'
import { withFunctionFetch } from '../src/withFunctionFetch'

const Factory = (integrationId: string) => {
  return {
    withFunctionFetch: function<TReturnedData = any>(functionName: string) {
      return withFunctionFetch<TReturnedData>(integrationId, functionName)
    },
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
