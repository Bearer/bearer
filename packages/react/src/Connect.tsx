import * as React from 'react'
import { BearerInstance } from '@bearer/js'
import { BearerContext } from './bearer-provider'

type TAuthPayload = { authId: string; integration: string }

export interface IConnectProps {
  integration: string
  authId?: string
  setupId?: string
  onSuccess: (data: TAuthPayload) => void
  onError?: (data: { authId?: string; integration: string; error: Error }) => void
  render: (props: { loading: boolean; connect: () => void; error: any }) => JSX.Element
}

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

  connect(bearer: BearerInstance) {
    return () => {
      this.setError(null)
      bearer
        .connect(this.props.integration, this.props.setupId, { authId: this.props.authId })
        .then((data: TAuthPayload) => {
          if (this.props.integration === data.integration) {
            this.props.onSuccess(data)
          }
        })
        .catch((error: any) => {
          if (this.props.onError) {
            this.props.onError({ error, authId: this.props.authId, integration: this.props.integration })
          }
          this.setError(error)
        })
    }
  }
  render() {
    return (
      <BearerContext.Consumer>
        {({ bearer }) => this.props.render({ loading: false, connect: this.connect(bearer), error: this.state.error })}
      </BearerContext.Consumer>
    )
  }
}

export default Connect
