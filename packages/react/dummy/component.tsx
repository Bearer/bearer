import * as React from 'react'

import Factory from './factory'
import BearerProvider from '../src/bearer-provider'

const { Fetcher, FetcherFactory, Connect } = Factory(process.env.REACT_INTEGRATION_ID)
const InvoiceList = FetcherFactory('InvoiceList')
const setupId = process.env.REACT_SETUP_ID
const clientId = process.env.REACT_CLIENT_ID

export default class ComponentClass extends React.Component<{}, { clientId: string }> {
  constructor(props) {
    super(props)
    this.state = { clientId }
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
        <div>
          <h2>Connect:</h2>{' '}
          <Connect
            setupId={setupId}
            onSuccess={this.success}
            render={({ loading, connect }) => {
              if (loading) {
                return <div> Loading</div>
              }
              return <button onClick={connect}> Connect</button>
            }}
          />
          <button onClick={this.changeClientId}>Chande clientId</button>
        </div>

        <div>
          <h2>Fetcher:</h2>
          <Fetcher />
        </div>
        <div>
          <h2>InvoiceList:</h2> <InvoiceList />
        </div>
      </BearerProvider>
    )
  }
}
