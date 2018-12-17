import * as React from 'react'

import { BearerContext } from './bearer-provider'

function fromBearer<T>(TagName: string) {
  const propSetEvent = `${TagName}-prop-set`

  class Klass extends React.Component<T> {
    static displayName = `Bearer(${TagName})`
    static contextType = BearerContext
    context!: React.ContextType<typeof BearerContext>
    private readonly eventRef: React.RefObject<HTMLInputElement>

    constructor(props: T) {
      super(props)
      this.eventRef = React.createRef()
    }

    public componentDidMount() {
      this.eventRef.current!.addEventListener(propSetEvent, this.prophandler)
    }

    public componentWillUnmount() {
      this.eventRef.current!.removeEventListener(propSetEvent, this.prophandler)
    }

    public render() {
      const combinedProps = {
        ...(this.context.getState ? this.context.getState() : {}),
        ...(this.props as any)
      }
      return (
        <div ref={this.eventRef}>
          <TagName {...combinedProps} />
        </div>
      )
    }

    private readonly prophandler = (e: any) => {
      if (this.context.handlePropUpdates) {
        this.context.handlePropUpdates(e)
      }
    }
  }
  return Klass
}

export default fromBearer
