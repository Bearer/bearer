import { Component, State } from '@bearer/core'

@Component({
  tag: 'app-ui'
})
export class AppUi {
  @State() showAlert: boolean = true

  render() {
    return (
      <div>
        <div>
          <h4>Typography</h4>
          <bearer-typography>Typography</bearer-typography>
          <bearer-typography kind="h1" as="h1">
            Typography
          </bearer-typography>
          <bearer-typography kind="h2" as="pre">
            Typography
          </bearer-typography>
          <bearer-typography kind="h3">Typography</bearer-typography>
          <bearer-typography kind="h4">Typography</bearer-typography>
          <bearer-typography kind="h5">Typography</bearer-typography>
          <bearer-typography kind="h6">Typography</bearer-typography>
        </div>
        <div>
          <h4>Loading</h4>
          <bearer-loading />
        </div>
        <div>
          <h4>Badges</h4>
          <bearer-badge kind="primary">Primary</bearer-badge>
          <bearer-badge kind="secondary">Secondary</bearer-badge>
          <bearer-badge kind="success">Success</bearer-badge>
          <bearer-badge kind="danger">Danger</bearer-badge>
          <bearer-badge kind="warning">Warning</bearer-badge>
          <bearer-badge kind="info">Info</bearer-badge>
          <bearer-badge kind="light">Light</bearer-badge>
          <bearer-badge kind="dark">Dark</bearer-badge>
        </div>
        <div>
          <h4>Button</h4>
          <bearer-button kind="primary">Primary</bearer-button>
          <bearer-button kind="secondary">Secondary</bearer-button>
          <bearer-button kind="success">Success</bearer-button>
          <bearer-button kind="danger">Danger</bearer-button>
          <bearer-button kind="warning">Warning</bearer-button>
          <bearer-button kind="info">Info</bearer-button>
          <bearer-button kind="light">Light</bearer-button>
          <bearer-button kind="dark" size="lg">
            Dark
          </bearer-button>
          <bearer-button kind="dark" size="sm">
            Dark
          </bearer-button>
        </div>

        <div>
          <h4>Alerts</h4>
          {this.showAlert && (
            <bearer-alert
              kind="primary"
              onDismiss={() => {
                this.showAlert = false
              }}
            >
              Primary <hr />
            </bearer-alert>
          )}
          <bearer-alert kind="secondary">
            Secondary <hr />
          </bearer-alert>
          <bearer-alert kind="success">
            Success <hr />
          </bearer-alert>
          <bearer-alert kind="danger">Danger</bearer-alert>
          <bearer-alert kind="warning">Warning</bearer-alert>
          <bearer-alert kind="info">Info</bearer-alert>
          <bearer-alert kind="light">Light</bearer-alert>
          <bearer-alert kind="dark">
            Light <hr /> content
          </bearer-alert>
        </div>
      </div>
    )
  }
}
