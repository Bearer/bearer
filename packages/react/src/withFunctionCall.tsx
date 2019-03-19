import * as React from 'react'

import { BearerContext } from './bearer-provider'

export interface WithFetchProps<TReturnedData> {
  loading: boolean
  error?: any
  data?: TReturnedData
  fetch: (params: any) => void
}

export interface IState<TReturnedData> {
  loading: boolean
  data?: TReturnedData
  error?: any
}

export const withFunctionCall = function<TReturnedData, TP extends object = {}>(
  integrationId: string,
  functionName: string
) {
  return (WrappedComponent: React.ComponentType<TP & WithFetchProps<TReturnedData>>) =>
    class FetcherComponent extends React.Component<TP, IState<TReturnedData>> {
      static displayName = `withFunctionCall(${functionName})(${getDisplayName(WrappedComponent)})`
      static contextType = BearerContext
      context!: React.ContextType<typeof BearerContext>

      constructor(props: TP) {
        super(props)
        this.state = {
          loading: false
        }
      }

      fetch = (params: any = {}) => {
        this.setState({ ...this.state, error: null, loading: true })
        this.context.bearer
          .functionFetch(integrationId, functionName, params)
          .then(({ data }) => {
            this.setState({ ...this.state, data })
          })
          .catch(({ error }) => {
            this.setState({ ...this.state, error })
          })
          .then(() => {
            this.setState({ ...this.state, loading: false })
          })
      }

      public render() {
        return (
          <WrappedComponent
            {...this.props}
            error={this.state.error}
            data={this.state.data}
            fetch={this.fetch}
            loading={this.state.loading}
          />
        )
      }
    }
}

function getDisplayName(WrappedComponent: any) {
  return WrappedComponent.displayName || WrappedComponent.name || 'Component'
}
