import kebabCase from 'lodash.kebabcase'
import * as React from 'react'

import BearerLoader from './bearer-loader'

interface IBearerProviderProps {
  clientId: string
  initialContext?: any
  onUpdate?(currentState: any): void
}

export interface IBearerContextValue {
  getState?(): any
  handlePropUpdates?(e: any): void
}

export const BearerContext = React.createContext<IBearerContextValue>({})

export default class BearerProvider extends React.Component<IBearerProviderProps, any> {
  private readonly contextValue: IBearerContextValue

  constructor(props: IBearerProviderProps) {
    super(props)
    this.state = props.initialContext || {}
    this.contextValue = {
      handlePropUpdates: this.handlePropUpdates,
      getState: this.getContextState
    }
  }

  public render() {
    return (
      <React.Fragment>
        <BearerLoader clientId={this.props.clientId} />
        <BearerContext.Provider value={this.contextValue}>{this.props.children}</BearerContext.Provider>
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

  private readonly getContextState = () => {
    return { ...this.state }
  }
}
