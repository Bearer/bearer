import kebabCase from 'lodash.kebabcase'
import * as React from 'react'

import BearerLoader from './BearerLoader'

interface IBearerProviderProps {
  clientId: string
  intialContext: any
}

export interface IBearerContextValue {
  getState?(): any
  handlePropUpdates?(e: any): void
}

export const BearerContext = React.createContext<IBearerContextValue>({})

export default class BearerProvider extends React.Component<IBearerProviderProps, any> {
  private eventRef: React.RefObject<HTMLDivElement>

  constructor(props: IBearerProviderProps) {
    super(props)
    this.state = props.intialContext || {}
    this.eventRef = React.createRef()
  }

  render() {
    const contextValue: IBearerContextValue = {
      handlePropUpdates: this.handlePropUpdates,
      getState: this.getContextState
    }
    return (
      <div ref={this.eventRef}>
        <BearerLoader clientId={this.props.clientId} />
        <BearerContext.Provider value={contextValue}>{this.props.children}</BearerContext.Provider>
      </div>
    )
  }

  private handlePropUpdates = (e: any) => {
    e.stopPropagation()
    const payload = e.detail
    const kebabPayload = Object.keys(payload).reduce((acc: any, key: string) => {
      acc[kebabCase(key)] = payload[key]
      return acc
    }, {})
    this.setState({ ...kebabPayload })
  }

  private getContextState = () => {
    return this.state
  }
}
