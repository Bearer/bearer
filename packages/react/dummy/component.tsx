import * as React from 'react'
import Factory from './factory'

const { Fetcher, FetcherFactory, Connect } = Factory('integration-uuid')
const InvoiceList = FetcherFactory('InvoiceList')

export default class ComponentClass extends React.Component {
  success = ({ authId }) => {
    alert(authId)
  }
  render() {
    return (
      <div>
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
            setupId="ok"
            onSuccess={this.success}
            render={({ loading, connect }) => {
              if (loading) {
                return <div> Loading</div>
              }
              return <button onClick={connect}> Connect baby</button>
            }}
          />
        </div>
      </div>
    )
  }
}
