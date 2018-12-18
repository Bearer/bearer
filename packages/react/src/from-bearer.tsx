import * as React from 'react'

import { BearerContext } from './bearer-provider'
const BEARER_EVENT_PROP_PREFIX = 'bearer-'

export default function fromBearer<T>(TagName: string) {
  const propSetEvent = `${TagName}-prop-set`

  return class extends React.Component<T> {
    static displayName = `Bearer(${TagName})`
    static contextType = BearerContext
    context!: React.ContextType<typeof BearerContext>
    readonly eventRef: React.RefObject<HTMLInputElement>

    constructor(props: T) {
      super(props)
      this.eventRef = React.createRef()
    }

    public componentDidMount() {
      if (this.eventRef.current) {
        this.eventRef.current.addEventListener(propSetEvent, this.prophandler)
        // NB: This means once a component is mounted we cannot change its event handlers
        // In the future it would be good to improve this to allow dynamic allocation
        this.handlers.forEach(key => {
          const anyProps = this.props as any
          if (anyProps[key]) {
            this.eventRef.current!.addEventListener(key, anyProps[key])
          }
        })
      }
    }

    public componentWillUnmount() {
      if (this.eventRef.current) {
        this.eventRef.current.removeEventListener(propSetEvent, this.prophandler)
        this.handlers.forEach(key => {
          const anyProps = this.props as any
          if (anyProps[key]) {
            this.eventRef.current!.removeEventListener(key, anyProps[key])
          }
        })
      }
    }

    public render() {
      const combinedProps = {
        ...(this.context.state),
        ...(this.props as any)
      }
      return <TagName {...combinedProps} ref={this.eventRef} />
    }

    readonly prophandler = (e: any) => {
      if (this.context.handlePropUpdates) {
        this.context.handlePropUpdates(e)
      }
    }

    get handlers() {
      return Object.keys(this.props).filter((key: string) => key.startsWith(BEARER_EVENT_PROP_PREFIX))
    }
  }
}
