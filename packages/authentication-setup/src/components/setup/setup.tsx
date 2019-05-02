import { Component, Prop, Event, EventEmitter, Element } from '@stencil/core'
import { Authentications } from './types'
// import bearer from '@bearer/js'
declare const window: Window & { bearer: any }
@Component({
  tag: 'bearer-save-setup',
  styleUrl: 'setup.css'
})
export class BearerSetup {
  /**
   * Your Bearer clientId
   */
  @Prop() clientId: string

  /**
   * Authentication Type of the integration
   */
  @Prop() type: Authentications | 'unknown' = 'unknown'

  /**
   * Integration identifier
   */
  @Prop() integration: string

  /**
   * Optionally provide your custom setupId
   */
  @Prop() setupId?: string

  @Prop() noError: boolean

  @Event() saved: EventEmitter<{ setupId: string }>

  @Element() el: HTMLElement

  onSave = (event: Event) => {
    event.preventDefault()

    const inputs = Array.from(this.el.querySelector('form').elements).filter(
      e => e.tagName.toLowerCase() === 'input'
    ) as HTMLInputElement[]
    const setup = inputs.reduce(
      (acc, node) => {
        acc[node.name] = node.value
        return acc
      },
      {
        type: this.type
      }
    )
    window
      .bearer(this.clientId)
      .invoke(this.integration, 'bearer-setup-save', {
        setup
      })
      .then(data => {
        this.saved.emit(data)
      })
      .catch(console.error)
  }

  render() {
    const Form = formMapper[this.type.toUpperCase()] || Unknown
    return [
      <form onSubmit={this.onSave}>
        <Form />
        <button>
          <slot name="submit-text">
            <span>Save</span>
          </slot>
        </button>
      </form>,
      !this.noError && (
        <div class="alert alert-warning errors">
          <ul>
            <li class="error missing-client">
              <slot name="errors-missing-client-id">
                <span>missing client id</span>
              </slot>
            </li>
            <li class="error missing-integration">
              <slot name="errors-missing-integration">
                <span>missing integration</span>
              </slot>
            </li>
            <li class="error missing-type">
              <slot name="errors-missing-type-id">
                <span>missing type</span>
              </slot>
            </li>
          </ul>
        </div>
      )
    ]
  }
}

// type FormProps = {
//   onSave: (data: any) => void
// }
function scopedId(name: string, scope = Date.now()) {
  return ['bearer', name, scope].join('-')
}
const OAuth2Form = () => {
  const clientId = scopedId('client-id')
  const clientSecret = scopedId('client-secret')
  return [
    <label htmlFor={clientId}>Client ID</label>,
    <input id={clientId} name="clientID" placeholder="Client Id" required />,
    <label htmlFor={clientSecret}>Client Secret</label>,
    <input id={clientSecret} name="clientSecret" placeholder="Client Secret" required type="password" />
  ]
}
const OAuth1Form = () => {
  const consumerKey = scopedId('consumer-key')
  const consumerSecret = scopedId('consumer-secret')
  return [
    <label htmlFor={consumerKey}>Client ID</label>,
    <input id={consumerKey} name="consumerKey" placeholder="Consumer key" required />,
    <label htmlFor={consumerSecret}>Client ID</label>,
    <input id={consumerSecret} name="consumerSecret" placeholder="Consumer Secret" required type="password" />
  ]
}

const BasicForm = () => {
  const username = scopedId('username')
  const password = scopedId('password')
  return [
    <label htmlFor={username}>Username</label>,
    <input id={username} name="username" placeholder="Username" required />,
    <label htmlFor={password}>Password</label>,
    <input id={password} name="password" placeholder="Password" required type="password" />
  ]
}

const ApikeyForm = () => {
  const key = scopedId('key')
  return [<label htmlFor={key}>API key</label>, <input id={key} name="apiKey" placeholder="Api key" required />]
}

const Unknown = () => (
  <div class="alert alert-warning bearer-setup-unknown" role="alert">
    Unknow authentication type
  </div>
)

const formMapper: Record<Authentications, any> = {
  ['APIKEY']: ApikeyForm,
  ['OAUTH2']: OAuth2Form,
  ['OAUTH1']: OAuth1Form,
  ['BASIC']: BasicForm
}
