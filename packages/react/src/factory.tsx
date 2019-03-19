import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from './Connect'
import { withFunctionCall as withFetch } from './withFunctionCall'

const Factory = (integrationId: string) => {
  const withFunctionCall = function<TReturnedData = any>(functionName: string) {
    return withFetch<TReturnedData>(integrationId, functionName)
  }
  return {
    withFunctionCall,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
