import * as React from 'react'

import { BearerContext, IBearerContextValue } from '../../src/bearer-provider'

interface IDummyContextProps {
  initialContext?: any
  onHandlePropUpdates?(e: any): void
}

export default class DummyContext extends React.Component<IDummyContextProps, any> {
  private readonly contextValue: IBearerContextValue

  constructor(props: IDummyContextProps) {
    super(props)
    this.state = props.initialContext || {}
    this.contextValue = {
      handlePropUpdates: props.onHandlePropUpdates || (() => {}),
      state: this.state
    }
  }

  public render() {
    return <BearerContext.Provider value={this.contextValue}>{this.props.children}</BearerContext.Provider>
  }
}
