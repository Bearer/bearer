import * as React from 'react'
import { BearerContext } from '../src/bearer-provider'
type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

import Connect, { IConnectProps } from '../src/Connect'

interface WithFetchProps<TReturnedData> {
  functionName: string
  loading: boolean
  error: any
  data: TReturnedData
  fetch: (params: any) => void
}

interface IState<TReturnedData> {
  loading: boolean
  data: TReturnedData
  error: any
}

const Factory = (integrationId: string) => {
  function withFetch<TReturnedData = any>(
    functionName: string
  ): <P extends object>(
    WrappedComponent: React.ComponentType<P & WithFetchProps<TReturnedData>>
  ) => React.ComponentType<P> {
    return <P extends object>(WrappedComponent: React.ComponentType<P & WithFetchProps<TReturnedData>>) => {
      class FetcherComponent extends React.Component<P, IState<TReturnedData>> {
        static displayName = `WithFunctionFetch(${functionName})(${getDisplayName(WrappedComponent)})`
        static contextType = BearerContext
        context!: React.ContextType<typeof BearerContext>

        constructor(props) {
          super(props)
          this.state = {
            loading: false,
            data: null,
            error: null
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
              functionName={functionName}
              loading={this.state.loading}
            />
          )
        }
      }

      return FetcherComponent
    }
  }

  return {
    withFetch,
    Connect: (props: Omit<IConnectProps, 'integration'>) => <Connect {...props} integration={integrationId} />
  }
}

export default Factory

function getDisplayName(WrappedComponent) {
  return WrappedComponent.displayName || WrappedComponent.name || 'Component'
}
