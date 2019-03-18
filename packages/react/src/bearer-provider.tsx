import kebabCase from 'lodash.kebabcase'
import * as React from 'react'
import bearer, { BearerInstance } from '@bearer/js'

// TODO: add JSDoc
interface IBearerProviderProps {
  clientId: string
  domObserver?: boolean
  integrationHost?: string
  initialContext?: any
  onUpdate?(currentState: any): void
}

export interface IBearerContextValue {
  bearer: BearerInstance
  state: any
  handlePropUpdates?(e: any): void
}

export const BearerContext = React.createContext<IBearerContextValue>({
  bearer: bearer.instance!,
  state: {},
  handlePropUpdates: () => {}
})

interface IBearerProviderState {
  bearer: BearerInstance
  integrationState: any
}

export default class BearerProvider extends React.Component<IBearerProviderProps, IBearerProviderState> {
  constructor(props: IBearerProviderProps) {
    super(props)
    this.state = {
      bearer: this.initBearerFromProps(props),
      integrationState: props.initialContext || {}
    }
  }

  public render() {
    const value = {
      bearer: this.state.bearer,
      handlePropUpdates: this.handlePropUpdates,
      state: this.state.integrationState
    }
    return <BearerContext.Provider value={value}>{this.props.children}</BearerContext.Provider>
  }

  componentWillReceiveProps(newProps: IBearerProviderProps, oldProps: IBearerProviderProps) {
    if (newProps.clientId !== oldProps.clientId) {
      this.setState(state => ({ ...state, bearer: this.initBearerFromProps(newProps) }))
    }
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
    this.setState(state => ({
      ...state,
      integrationState: {
        ...state.integrationState,
        ...kebabPayload
      }
    }))
  }

  private initBearerFromProps(props: IBearerProviderProps) {
    const { integrationHost, domObserver } = props
    return bearer(props.clientId, { integrationHost, domObserver })
  }
}
