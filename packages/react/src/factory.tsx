import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from './Connect'
import { withInvoke as withFetch } from './withInvoke'

const Factory = (integrationId: string) => {
  const withInvoke = function<TReturnedData = any>(functionName: string) {
    return withFetch<TReturnedData>(integrationId, functionName)
  }
  return {
    withInvoke,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
