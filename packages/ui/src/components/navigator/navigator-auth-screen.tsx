import { Component, Element, Event, EventEmitter, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/with-authentication'
import { FWithAuthenticate, FWithRevoke } from '../Authorized/bearer-authorized'

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
  scenarioId = 'BEARER_SCENARIO_ID'

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

  onRevoked = () => {
    this.scenarioAuthorized = false
  }

  goNext = () => {
    console.debug('[BEARER]', 'go to next screen')
    this.scenarioAuthenticate.emit()
    this.stepCompleted.emit()
    this.scenarioAuthorized = true
  }

  onAuthorizeClick = (authenticate: () => Promise<boolean>) => {
    console.debug('[BEARER]', 'onAuthorized')
    authenticate()
      .then(this.goNext)
      .catch(console.error)
  }

  onRevokeClick = (revoke: () => Promise<boolean>) => {
    revoke()
      .then(this.onRevoked)
      .catch(console.error)
  }

  renderUnauthoried: FWithAuthenticate = ({ authenticate }) => (
    // tslint:disable-next-line:react-this-binding-issue
    <bearer-button kind="primary" onClick={() => this.onAuthorizeClick(authenticate)}>
      {' '}
      Login{' '}
    </bearer-button>
  )

  renderAuthorized: FWithRevoke = ({ revoke }) => (
    // tslint:disable-next-line:react-this-binding-issue
    <bearer-button kind="danger" onClick={() => this.onRevokeClick(revoke)}>
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
