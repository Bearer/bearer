import { Component, Element, Event, EventEmitter, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/withAuthentication'

@Component({
  tag: 'bearer-navigator-auth-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorAuthScreen extends AuthenticationListener {
  @Element()
  el: HTMLStencilElement

  @State()
  scenarioAuthorized: boolean = null

  @Event()
  scenarioAuthenticate: EventEmitter
  @Event()
  stepCompleted: EventEmitter
  @Prop()
  scenarioId?: string = 'BEARER_SCENARIO_ID'

  @Method()
  willAppear() {
    console.debug('[BEARER]', 'Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  willDisappear() {
    console.debug('[BEARER]', 'Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  getTitle() {
    return 'Authentication'
  }

  onAuthorized = () => {
    console.debug('[BEARER]', 'onAuthorized')
    this.goNext()
  }

  onRevoked = () => {
    this.scenarioAuthorized = false
  }

  goNext() {
    console.debug('[BEARER]', 'go to next screen')
    this.scenarioAuthenticate.emit()
    this.stepCompleted.emit()
    this.scenarioAuthorized = true
  }

  renderUnauthoried = ({ authenticate }) => (
    <bearer-button kind="primary" onClick={authenticate}>
      {' '}
      Login{' '}
    </bearer-button>
  )

  renderAuthorized = ({ revoke }) => (
    <bearer-button kind="warning" onClick={revoke}>
      {' '}
      Logout{' '}
    </bearer-button>
  )

  render() {
    return (
      <bearer-navigator-screen id="screen" navigationTitle="Authentication" class="in">
        <bearer-authorized
          scenarioId={this.scenarioId}
          id="authorizer"
          renderUnauthorized={this.renderUnauthoried}
          renderAuthorized={this.renderAuthorized}
        />
      </bearer-navigator-screen>
    )
  }
}
