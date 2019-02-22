import kebabCase from 'lodash.kebabcase'
import * as React from 'react'
import bearer from '@bearer/js'

interface IBearerProviderProps {
  clientId: string
  intHost?: string // unused for now
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
    return <BearerContext.Provider value={contextValue}>{this.props.children}</BearerContext.Provider>
  }

  componentDidMount() {
    bearer(this.props.clientId)
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
