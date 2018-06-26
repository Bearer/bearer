import { Component, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'app-home',
  styleUrl: 'app-home.scss'
})
export class AppHome {
  @Event() stepCompleted: EventEmitter

  render() {
    return (
      <div class="padded">
        <bearer-popover-navigator direction="right" button="popover">
          <bearer-navigator-screen navigationTitle="Screen one">
            <div style={{ background: '#00b3ff' }}>
              Screen one <app-navigator-next />
            </div>
          </bearer-navigator-screen>
          <bearer-navigator-screen navigationTitle="Final Screen">
            Three <app-navigator-finish />
          </bearer-navigator-screen>
        </bearer-popover-navigator>

        <bearer-popover-navigator direction="right" button="popover">
          <bearer-navigator-auth-screen />
          <bearer-navigator-screen navigationTitle="Final Screen">
            Three <app-navigator-finish />
          </bearer-navigator-screen>
        </bearer-popover-navigator>
      </div>
    )
  }
}
