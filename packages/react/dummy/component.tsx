import * as React from 'react'

import Factory from './factory'
import BearerProvider from '../src/bearer-provider'

const integrationId = process.env.REACT_INTEGRATION_ID
const setupId = '939407c0-473d-11e9-a595-499c863fdcda'
const clientId = process.env.REACT_CLIENT_ID

const { Connect, withFunctionCall } = Factory(integrationId)

const Fetcher = withFunctionCall<string[]>('MyDataPlease')(props => {
  if (props.loading) {
    return <div>Loading</div>
  }

  if (props.error) {
    return <div>{JSON.stringify(props.error.error)}</div>
  }

  if (props.data) {
    return (
      <div>
        {props.data.map(d => (
          <li key={d}>{d}</li>
        ))}
      </div>
    )
  }
  return (
    <div>
      Coucou
      <br />
      <button
        onClick={() => {
          props.fetch({ page: 1 })
        }}
      >
        Load
      </button>
    </div>
  )
})

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
    console.log('Integration', p.integration, ' authId: ', p.authId)
  }

  render() {
    const Setup = `bearer-${integrationId}-setup-action`
    const SetupDisplay = `bearer-${integrationId}-setup-display`

    return (
      <BearerProvider
        clientId={this.state.clientId}
        domObserver={false}
        // integrationHost="https://int.staging.bearer.sh"
      >
        <div>
          <h2>Connect:</h2>
          <Setup />
          <SetupDisplay />
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
        <div />
      </BearerProvider>
    )
  }
}
