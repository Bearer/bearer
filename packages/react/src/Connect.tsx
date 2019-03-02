import * as React from 'react'
import { BearerContext } from './bearer-provider'

export interface IConnectProps {
  integration: string
  authId?: string
  setupId: string
  onSuccess: (data: { authId: string; integration: string }) => void
  onError?: (data: { authId?: string; integration: string; error: Error }) => void
  render: (props: { loading: boolean; connect: () => void; error: any }) => JSX.Element
}

// TODO: add tests
class Connect extends React.Component<IConnectProps, { error?: any }> {
  static contextType = BearerContext
  context!: React.ContextType<typeof BearerContext>

  constructor(props: IConnectProps) {
    super(props)
    this.state = {}
    this.connect = this.connect.bind(this)
    this.setError = this.setError.bind(this)
  }

  setError(error: any) {
    this.setState(state => ({ ...state, error }))
  }

  connect() {
    this.setError(null)
    this.context
      .bearer!.connectTo(this.props.integration, this.props.setupId, { authId: this.props.authId })
      .then(({ data }) => {
        if (this.props.integration === data.integration) {
          this.props.onSuccess(data)
        }
      })
      .catch(error => {
        if (this.props.onError) {
          this.props.onError({ error, authId: this.props.authId, integration: this.props.integration })
        }
        this.setError(error)
      })
  }
  render() {
    // TODO: check component is injected within a BearerProvider
    return this.props.render({ loading: false, connect: this.connect, error: this.state.error })
  }
}

export default Connect
