import * as React from 'react'

import { BearerContext } from '../../src/bearer-provider'

export default class DummyContextConsumer extends React.Component<any> {
  static contextType = BearerContext
  context!: React.ContextType<typeof BearerContext>

  constructor(props: any) {
    super(props)
  }

  public render() {
    const data = this.context.getState ? this.context.getState() : {}
    return (
      <React.Fragment>
        {JSON.stringify(data, null, 2)}
        <button onClick={this.prophandler} />
      </React.Fragment>
    )
  }

  private readonly prophandler = (e: any) => {
    if (this.context.handlePropUpdates) {
      this.context.handlePropUpdates(e)
    }
  }
}
