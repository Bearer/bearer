import kebabCase from 'lodash.kebabcase'
import * as React from 'react'
import bearer, { BearerInstance } from '@bearer/js'

interface IBearerProviderProps {
  clientId: string
  intHost?: string // unused for now
  initialContext?: any
  onUpdate?(currentState: any): void
}

export interface IBearerContextValue {
  bearer: BearerInstance | undefined
  state: any
  handlePropUpdates?(e: any): void
}

export const BearerContext = React.createContext<IBearerContextValue>({
  bearer: bearer.instance,
  state: {},
  handlePropUpdates: () => {}
})

export default class BearerProvider extends React.Component<IBearerProviderProps, any> {
  constructor(props: IBearerProviderProps) {
    super(props)
    this.state = props.initialContext || {}
  }

  public render() {
    const value = {
      bearer: bearer.instance,
      handlePropUpdates: this.handlePropUpdates,
      state: this.state
    }
    return <BearerContext.Provider value={value}>{this.props.children}</BearerContext.Provider>
  }

  componentDidMount() {
    // TODO: handle prop update to refresh the clientId
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
