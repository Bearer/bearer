import Bearer, { Component, Prop, State } from '@bearer/core'

@Component({
  tag: 'bearer-setup-display',
  shadow: true
})
export class BearerSetupDisplay {
  @Prop()
  scenarioId = ''
  @State()
  isSetup: boolean = false
  @Prop({ mutable: true })
  setupId = ''

  componentDidLoad() {
    Bearer.emitter.addListener(`setup_success:${this.scenarioId}`, data => {
      this.setupId = data.referenceId
      this.isSetup = true
    })
  }

  render() {
    if (this.isSetup || this.setupId) {
      return (
        <p>
          Scenario is currently setup with Setup ID:&nbsp;
          <bearer-badge color="info">{this.setupId}</bearer-badge>
        </p>
      )
    } else {
      return (
        <p>
          <p>Scenario hasn't been setup yet</p>
        </p>
      )
    }
  }
}
