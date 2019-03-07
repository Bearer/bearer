import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'

const Factory = (integrationId: string) => {
  return {
    Fetcher: () => <span>Fetcher {integrationId}</span>,
    FetcherFactory: (functionName: string) => () => <span>FetcherFactory: {functionName}</span>,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
