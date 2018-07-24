import { Component, Method, State, Event, EventEmitter, Element, Prop } from '@bearer/core'
import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-navigator-auth-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorAuthScreen {
  @Element() el: HTMLStencilElement

  @State() sessionInitialized: boolean = false
  @State() scenarioAuthorized: boolean = null

  @Event() scenarioAuthenticate: EventEmitter
  @Event() stepCompleted: EventEmitter
  @Prop({ context: 'bearer' })
  bearerContext: any

  @Method()
  willAppear() {
    this.el.shadowRoot.querySelector('#screen')['willAppear']()
  }

  @Method()
  willDisappear() {
    this.el.shadowRoot.querySelector('#screen')['willDisappear']()
  }

  @Method()
  getTitle() {
    return 'Authentication'
  }

  get SCENARIO_ID() {
    return 'BEARER_SCENARIO_ID'
  }

  authenticate = () => {
    Bearer.instance.askAuthorizations({
      scenarioId: this.SCENARIO_ID,
      setupId: this.bearerContext.setupId
    })
  }

  private authorizedListener: any = null
  private revokedListener: any = null

  componentDidLoad() {
    Bearer.instance.maybeInitialized.then(() => {
      this.sessionInitialized = true
      Bearer.instance
        .hasAuthorized(this.SCENARIO_ID)
        .then(() => {
          console.log('[BEARER]', 'authorized')
          this.goNext()
        })
        .catch(e => {
          console.log('[BEARER]', 'unauthorized', e)
          this.scenarioAuthorized = false
        })

      this.authorizedListener = Bearer.onAuthorized(this.SCENARIO_ID, () => {
        this.goNext()
      })

      this.revokedListener = Bearer.onRevoked(this.SCENARIO_ID, () => {
        this.scenarioAuthorized = false
      })
    })
  }

  componentDidUnload() {
    if (this.authorizedListener) {
      this.authorizedListener.remove()
      this.authorizedListener = null
    }
    if (this.revokedListener) {
      this.revokedListener.remove()
      this.revokedListener = null
    }
  }

  goNext() {
    this.scenarioAuthenticate.emit()
    this.stepCompleted.emit()
    this.scenarioAuthorized = true
  }

  revoke = () => {
    Bearer.instance.revokeAuthorization(this.SCENARIO_ID)
  }

  render() {
    return (
      <bearer-navigator-screen id="screen" navigationTitle="Authentication" class="in">
        {this.sessionInitialized &&
          this.scenarioAuthorized !== null &&
          (this.scenarioAuthorized ? (
            <bearer-button kind="warning" onClick={this.revoke}>
              Logout
            </bearer-button>
          ) : (
            <bearer-button kind="primary" onClick={this.authenticate}>
              Login
            </bearer-button>
          ))}
      </bearer-navigator-screen>
    )
  }
}
