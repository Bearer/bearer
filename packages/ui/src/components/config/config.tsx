import Bearer, { Component, State, Element, Event, EventEmitter, Prop, StateManager } from '@bearer/core'

import { FieldSet } from '../Forms/Fieldset'
import debug from '../../logger'
const logger = debug('bearer-config')

type TSetupPayload = {
  Item: { referenceId: string }
}

@Component({
  tag: 'bearer-config',
  styleUrl: 'config.scss',
  shadow: true
})
export class BearerConfig {
  @Prop() fields: any[] | string = []
  @Prop() referenceId: string
  @Prop() integrationId: string

  @Element() element: HTMLElement
  @Event() stepCompleted: EventEmitter

  @State() fieldSet: FieldSet
  @State() error: boolean = false
  @State() loading: boolean = false

  handleSubmit = (e: any) => {
    e.preventDefault()
    this.loading = true
    const formSet = this.fieldSet.map(el => ({
      key: el.controlName,
      value: el.value
    }))
    StateManager.storeSetup(formSet.reduce((acc, obj) => ({ ...acc, [obj['key']]: obj['value'] }), {}))
      .then((item: TSetupPayload) => {
        this.loading = false
        logger(this.integrationId)
        Bearer.emitter.emit(`config_success:${this.integrationId}`, {
          referenceID: item.Item.referenceId
        })
      })
      .catch(() => {
        this.error = true
        this.loading = false
        Bearer.emitter.emit(`config_error:${this.integrationId}`, {})
      })
  }

  componentWillLoad() {
    this.fieldSet = new FieldSet(this.fields as any[])
  }

  render() {
    return [
      this.error && <bearer-alert kind="danger">[Error] Unable to store the credentials</bearer-alert>,
      this.loading ? (
        <bearer-loading />
      ) : (
        <bearer-form fields={this.fieldSet} clearOnInput={true} onSubmit={this.handleSubmit} />
      )
    ]
  }
}
