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
    const label = <strong>Setup-id</strong>
    if (this.isSetup || this.setupId) {
      return (
        <div>
          {label}:&nbsp; {this.setupId}
        </div>
      )
    }
    return <div>{label}:&nbsp; Not set</div>
  }
}
