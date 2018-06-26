import {
  Component,
  Method,
  State,
  Event,
  EventEmitter,
  Element,
  Prop
} from '@bearer/core'
import Bearer, { Events } from '@bearer/core'

@Component({
  tag: 'bearer-navigator-auth-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorAuthScreen {
  @Element() el: HTMLStencilElement

  @State() sessionInitialized: boolean = false
  @State() scenarioAuthorized: boolean = false

  @Event() scenarioAuthenticate: EventEmitter
  @Event() stepCompleted: EventEmitter

  @Prop() setupId: string = 'BEARER_SCENARIO_ID'

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

  private listener: any

  authenticate = () => {
    Bearer.instance.askAuthorizations({
      scenarioId: this.setupId,
      setupId: this.setupId
    })
  }

  componentDidLoad() {
    Bearer.instance.maybeInitialized.then(() => {
      this.sessionInitialized = true
      this.scenarioAuthorized = Bearer.instance.hasAuthorized(this.setupId)
      this.listener = Bearer.emitter.addListener(
        Events.SCENARIO_AUTHORIZED,
        // TODO: we need to ensure the tokens are not confused
        ({ scenarioId, authorized }) => {
          if (this.setupId === scenarioId) {
            this.scenarioAuthorized = authorized
            this.goNext()
          }
        }
      )
      this.goNext()
    })
  }

  componentDidUnload() {
    if (this.listener) {
      this.listener.remove()
      this.listener = null
    }
  }

  goNext() {
    if (this.scenarioAuthorized) {
      this.scenarioAuthenticate.emit()
      this.stepCompleted.emit()
    }
  }

  revoke = () => {
    Bearer.instance.revokeAuthorization(this.setupId)
  }

  render() {
    return (
      <bearer-navigator-screen
        id="screen"
        navigationTitle="Authentication"
        class="in"
      >
        {this.sessionInitialized &&
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
