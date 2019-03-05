import * as React from 'react'

import Factory from './factory'
import BearerProvider from '../src/bearer-provider'

const { Fetcher, FetcherFactory, Connect } = Factory('2627b8-slack-sharing')
const InvoiceList = FetcherFactory('InvoiceList')

export default class ComponentClass extends React.Component<{}, { clientId: string }> {
  constructor(props) {
    super(props)
    this.state = { clientId: '296885e93fc38510630a9a4c964d1f404ff59448f3b3ad5b8b' }
  }

  changeClientId = () => {
    this.setState(() => ({
      clientId: 'ok'
    }))
  }
  success = p => {
    console.log(p)
  }

  render() {
    return (
      <BearerProvider clientId={this.state.clientId}>
        Component:
        <div>
          <h2>Fetcher:</h2>
          <Fetcher />
        </div>
        <div>
          <h2>InvoiceList:</h2> <InvoiceList />
        </div>
        <div>
          <h2>Connect:</h2>{' '}
          <Connect
            setupId="ae9a79d0-d93a-11e8-aebb-51df6010fd72"
            onSuccess={this.success}
            render={({ loading, connect }) => {
              if (loading) {
                return <div> Loading</div>
              }
              return <button onClick={connect}> Connect baby</button>
            }}
          />
          <Connect
            setupId="ae9a79d0-d93a-11e8-aebb-51df6010fd72"
            authId="qsdqsdsd"
            onSuccess={this.success}
            render={({ loading, connect }) => {
              if (loading) {
                return <div> Loading</div>
              }
              return <button onClick={connect}> Connect wtih quth</button>
            }}
          />
          <button onClick={this.changeClientId}>Chande clientId</button>
        </div>
      </BearerProvider>
    )
  }
}
