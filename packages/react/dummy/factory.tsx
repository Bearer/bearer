import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'

const Factory = (integrationId: string) => {
  return {
    Fetcher: () => <span>Fetcher {integrationId}</span>,
    FetcherFactory: (intentName: string) => () => <span>FetcherFactory: {intentName}</span>,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
