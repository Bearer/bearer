import Bearer, { Component, State, Element, Event, EventEmitter, Prop, StateManager } from '@bearer/core'
import { FieldSet } from '../Forms/Fieldset'

import { OAuth2SetupType, EmailSetupType, KeySetupType } from './setup-types'

import debug from '../../logger'
const logger = debug('bearer-setup')

type TSetupPayload = {
  Item: { referenceId: string }
}

@Component({
  tag: 'bearer-setup',
  styleUrl: 'setup.scss',
  shadow: true
})
export class BearerSetup {
  @Prop()
  display: 'inline' | 'block'
  @Prop()
  fields: any[] | string = []
  @Prop()
  referenceId: string
  @Prop()
  integrationId: string

  @Element()
  element: HTMLElement
  @Event()
  setupSuccess: EventEmitter

  @State()
  fieldSet: FieldSet
  @State()
  error: boolean = false
  @State()
  loading: boolean = false

  handleSubmit = (e: any) => {
    e.preventDefault()
    this.loading = true
    const formSet = this.fieldSet.map(el => {
      return { key: el.controlName, value: el.value }
    })
    StateManager.storeSetup(formSet.reduce((acc, obj) => ({ ...acc, [obj['key']]: obj['value'] }), {}))
      .then((item: TSetupPayload) => {
        this.loading = false
        const referenceId = item.Item.referenceId
        logger('setup_success %s', `setup_success:${this.integrationId}`)
        Bearer.emitter.emit(`setup_success:${this.integrationId}`, {
          referenceId
        })
        this.setupSuccess.emit({ referenceId, integrationId: this.integrationId })
      })
      .catch(err => {
        this.error = true
        this.loading = false
        logger('setup_error:% %j', this.integrationId, err)
        Bearer.emitter.emit(`setup_error:${this.integrationId}`, {})
      })
  }

  onValueChange(field: string, event: any) {
    if (event) {
      this.fieldSet.setValue(field, event.target.value)
    }
  }

  componentWillLoad() {
    if (typeof this.fields !== 'string') {
      this.fieldSet = new FieldSet(this.fields as any[])
      return
    }
    switch (this.fields) {
      case 'email':
        this.fieldSet = new FieldSet(EmailSetupType)
        break
      case 'type': // legacy
      case 'key':
        this.fieldSet = new FieldSet(KeySetupType)
        break
      case 'oauth2':
      default:
        this.fieldSet = new FieldSet(OAuth2SetupType)
    }
  }

  renderInputs = input => {
    const inputLabel = input => {
      return <label class="form-label">{input.label}</label>
    }
    const inputField = input => {
      return (
        <input
          class="form-input"
          value={input.value}
          onChange={event => this.onValueChange(input.controlName, event)}
          type={input.type}
          name={input.controlName}
          placeholder={input.placeholder}
        />
      )
    }

    return input.label ? (
      <div class="form-group">
        {inputLabel(input)}
        {inputField(input)}
      </div>
    ) : (
      inputField(input)
    )
  }

  render() {
    return [
      this.error && <bearer-alert kind="danger">[Error] Unable to save the credentials</bearer-alert>,
      <form class={this.display && `form-${this.display}`} onSubmit={this.handleSubmit}>
        {this.fieldSet.map(input => {
          return this.renderInputs(input)
        })}
        <div class="form-submit">
          <bearer-btn type="submit">Save</bearer-btn>
        </div>
      </form>
    ]
  }
}
