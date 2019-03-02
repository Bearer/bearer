import * as React from 'react'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from './Connect'

const Factory = (integrationId: string) => {
  return {
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory
