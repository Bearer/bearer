import kebabCase from 'lodash.kebabcase'
import * as React from 'react'

import BearerLoader from './bearer-loader'

interface IBearerProviderProps {
  clientId: string
  intHost?: string
  initialContext?: any
  onUpdate?(currentState: any): void
}

export interface IBearerContextValue {
  state: any
  handlePropUpdates?(e: any): void
}

export const BearerContext = React.createContext<IBearerContextValue>({
  state: {},
  handlePropUpdates: () => {}
})

export default class BearerProvider extends React.Component<IBearerProviderProps, any> {
  private readonly contextValue: IBearerContextValue

  constructor(props: IBearerProviderProps) {
    super(props)
    this.state = props.initialContext || {}
    this.contextValue = {
      handlePropUpdates: this.handlePropUpdates,
      state: this.state
    }
  }

  public render() {
    const contextValue = {
      handlePropUpdates: this.handlePropUpdates,
      state: this.state
    }
    return (
      <React.Fragment>
        <BearerLoader clientId={this.props.clientId} intHost={this.props.intHost} />
        <BearerContext.Provider value={contextValue}>{this.props.children}</BearerContext.Provider>
      </React.Fragment>
    )
  }

  public componentDidUpdate(_prevProps: IBearerProviderProps, prevState: any) {
    if (this.props.onUpdate && this.state !== prevState) {
      this.props.onUpdate(this.state)
    }
  }

  private readonly handlePropUpdates = (e: CustomEvent) => {
    e.stopPropagation()
    const payload = e.detail
    const kebabPayload = Object.keys(payload).reduce((acc: any, key: string) => {
      acc[kebabCase(key)] = payload[key]
      return acc
    }, {})
    this.setState({ ...kebabPayload })
  }
}
