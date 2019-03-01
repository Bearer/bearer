import * as React from 'react'

const Factory = (integrationId: string) => {
  return {
    Fetcher: () => <span>Fetcher {integrationId}</span>,
    FetcherFactory: (intentName: string) => () => <span>FetcherFactory: {intentName}</span>,
    Connect: () => <span>Connect</span>
  }
}

export default Factory
