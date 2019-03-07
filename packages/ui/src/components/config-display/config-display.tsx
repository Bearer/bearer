import Bearer, { Component, Prop, State } from '@bearer/core'

@Component({
  tag: 'bearer-config-display',
  shadow: true
})
export class BearerConfigDisplay {
  @Prop() integrationId = ''
  @State() isConfig: boolean = false
  @State() configId = ''

  componentDidLoad() {
    Bearer.emitter.addListener(`config_success:${this.integrationId}`, data => {
      this.configId = data.referenceID
      this.isConfig = true
    })
  }

  render() {
    if (this.isConfig) {
      return (
        <p>
          Integration is currently configure with Config ID:&nbsp;
          <bearer-badge kind="info">{this.configId}</bearer-badge>
        </p>
      )
    }
    return (
      <p>
        <p>Integration hasn't been configured yet</p>
      </p>
    )
  }
}
