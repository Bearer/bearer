import { Component, Element, Event, EventEmitter, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/with-authentication'
import { FWithAuthenticate, FWithRevoke } from '../Authorized/bearer-authorized'
import debug from '../../logger'
const logger = debug('bearer-navigator-auth-screen')

const logError = (error: any) => {
  logger.extend('error')('%j', error)
}

@Component({
  tag: 'bearer-navigator-auth-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorAuthScreen extends AuthenticationListener {
  @Element()
  el: HTMLStencilElement

  @State()
  integrationAuthorized: boolean = null

  @Event()
  integrationAuthenticate: EventEmitter
  @Event()
  stepCompleted: EventEmitter
  @Prop()
  integrationId = 'BEARER_INTEGRATION_ID'

  @Prop()
  authId: string

  @Method()
  willAppear() {
    logger('Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  willDisappear() {
    logger('Auth screen willAppear')
    const screen: HTMLBearerNavigatorScreenElement = this.el.shadowRoot.querySelector('#screen')
    screen.willAppear({})
  }

  @Method()
  getTitle() {
    return 'Authentication'
  }

  onRevoked = () => {
    this.integrationAuthorized = false
  }

  goNext = () => {
    logger('go to next screen')
    this.integrationAuthenticate.emit()
    this.stepCompleted.emit()
    this.integrationAuthorized = true
  }

  onAuthorizeClick = (authenticate: () => Promise<boolean>) => {
    logger('onAuthorized')
    authenticate()
      .then(this.goNext)
      .catch(logError)
  }

  onRevokeClick = (revoke: (authRefId: string) => Promise<boolean>) => {
    revoke(this.authId)
      .then(this.onRevoked)
      .catch(logError)
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
    <bearer-button kind="secondary" onClick={() => this.onRevokeClick(revoke)}>
      {' '}
      Logout{' '}
    </bearer-button>
  )

  render() {
    return (
      <bearer-navigator-screen id="screen" navigationTitle="Authentication" class="in">
        <bearer-authorized
          integrationId={this.integrationId}
          id="authorizer"
          renderUnauthorized={this.renderUnauthoried}
          renderAuthorized={this.renderAuthorized}
        />
      </bearer-navigator-screen>
    )
  }
}
