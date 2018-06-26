import { Component } from '@bearer/core'
import '@stencil/router'
import Bearer from '@bearer/core'

Bearer.init({ integrationId: 'demo-id' })

@Component({
  tag: 'app-demo',
  shadow: true,
  styleUrl: 'app-demo.scss'
})
export class MyApp {
  render() {
    return (
      <div class="root">
        <header>
          <h1>Demo App</h1>
          <stencil-route-link url="/" exact={true}>
            Popover
          </stencil-route-link>
          <stencil-route-link url="/navigator">Navigator</stencil-route-link>
          <stencil-route-link url="/paginator">Paginator</stencil-route-link>
          <stencil-route-link url="/forms">Forms</stencil-route-link>
          <stencil-route-link url="/scrollable">Scrollable</stencil-route-link>
          <stencil-route-link url="/ui">UI Components</stencil-route-link>
          <stencil-route-link url="/setup">Setup</stencil-route-link>
        </header>

        <main>
          <stencil-router>
            <stencil-route url="/" component="app-home" exact={true} />
            <stencil-route url="/navigator" component="app-navigator" />
            <stencil-route url="/paginator" component="app-pagination" />
            <stencil-route url="/forms" component="app-forms" />
            <stencil-route url="/scrollable" component="app-scrollable" />
            <stencil-route url="/ui" component="app-ui" />
            <stencil-route url="/setup" component="app-setup" />
          </stencil-router>
        </main>
      </div>
    )
  }
}
