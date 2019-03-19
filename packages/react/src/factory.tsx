import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from './Connect'
import { withFunctionFetch as withFetch } from './withFunctionFetch'

const Factory = (integrationId: string) => {
  const withFunctionFetch = function<TReturnedData = any>(functionName: string) {
    return withFetch<TReturnedData>(integrationId, functionName)
  }
  return {
    withFunctionFetch,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
