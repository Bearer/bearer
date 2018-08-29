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
    ;(this.el.shadowRoot.querySelector('#screen') as any).willAppear()
  }

  @Method()
  willDisappear() {
    console.log('[BEARER]', 'Auth screen willAppear')
    ;(this.el.shadowRoot.querySelector('#screen') as any).willDisappear()
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

  authenticate = () => {
    ;(this.el.shadowRoot.querySelector('#authorizer') as any).authenticate()
  }

  revoke = () => {
    ;(this.el.shadowRoot.querySelector('#authorizer') as any).revoke()
  }

  render() {
    return (
      <bearer-navigator-screen id="screen" navigationTitle="Authentication" class="in">
        <bearer-authorized
          id="authorizer"
          renderUnauthorized={() => (
            <bearer-button color="primary" onClick={this.authenticate}>
              {' '}
              Login{' '}
            </bearer-button>
          )}
          renderAuthorized={() => (
            <bearer-button color="warning" onClick={this.revoke}>
              {' '}
              Logout{' '}
            </bearer-button>
          )}
        />
      </bearer-navigator-screen>
    )
  }
}
