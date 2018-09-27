import { Component, Element, Event, EventEmitter, Method, State } from '@bearer/core'

import WithAuthentication, { IAuthenticated, WithAuthenticationMethods } from '../../decorators/withAuthentication'

@WithAuthentication()
@Component({
  tag: 'bearer-navigator-auth-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorAuthScreen extends WithAuthenticationMethods implements IAuthenticated {
  @Element()
  el: HTMLStencilElement

  @State()
  scenarioAuthorized: boolean = null

  @Event()
  scenarioAuthenticate: EventEmitter
  @Event()
  stepCompleted: EventEmitter

  @Method()
  willAppear() {
    console.log('[BEARER]', 'Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  willDisappear() {
    console.log('[BEARER]', 'Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  getTitle() {
    return 'Authentication'
  }

  onAuthorized = () => {
    console.log('[BEARER]', 'onAuthorized')
    this.goNext()
  }

  onRevoked = () => {
    this.scenarioAuthorized = false
  }

  goNext() {
    console.log('[BEARER]', 'go to next screen')
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
          id="authorizer"
          renderUnauthorized={this.renderUnauthoried}
          renderAuthorized={this.renderAuthorized}
        />
      </bearer-navigator-screen>
    )
  }
}
