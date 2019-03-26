import * as React from 'react'
type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'
import { withInvoke } from '../src/withInvoke'

const Factory = (integrationId: string) => {
  return {
    withInvoke: function<TReturnedData = any>(functionName: string) {
      return withInvoke<TReturnedData>(integrationId, functionName)
    },
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
