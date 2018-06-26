import { Component, State } from '@bearer/core'

@Component({
  tag: 'app-setup',
  styleUrl: 'app-setup.scss'
})
export class AppSetup {
  @State() scenarioId: string = 'alice-attach-pull-request'
  @State()
  fields = [
    {
      type: 'text',
      label: 'Some text',
      controlName: 'ulogin'
    },
    {
      type: 'select',
      label: 'Your region',
      controlName: 'region',
      options: [
        {
          label: 'Africa',
          value: 'africa'
        },
        {
          label: 'America',
          value: 'america'
        },
        {
          label: 'Asia',
          value: 'asia'
        },
        {
          label: 'Europe',
          value: 'europe'
        },
        {
          label: 'Oceania',
          value: 'oceania'
        }
      ]
    },
    {
      type: 'password',
      label: 'Your password',
      controlName: 'passwd'
    }
  ]

  scenarioIdChanged = ({ detail }) => {
    this.scenarioId = detail
  }

  render() {
    const innerListener = `BEARER_SCENARIO_ID:setup_success:${this.scenarioId}`
    return (
      <div class="padded">
        <bearer-typography kind="h4">Setup Component</bearer-typography>
        <bearer-dropdown-button innerListener={innerListener}>
          <span slot="buttonText">Setup Scenario</span>
          <bearer-setup scenario-id={this.scenarioId} fields={this.fields} />
        </bearer-dropdown-button>
        <div class="down">
          <bearer-typography kind="h4">Setup Display</bearer-typography>
          <bearer-setup-display setup-id={this.scenarioId} />
        </div>
      </div>
    )
  }
}
