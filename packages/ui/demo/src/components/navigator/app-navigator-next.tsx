import { Component, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'app-navigator-next'
})
export class AppNavigatorNext {
  @Event() stepCompleted: EventEmitter

  render() {
    return (
      <bearer-button kind="primary" onClick={() => this.stepCompleted.emit()}>
        Next Screen
      </bearer-button>
    )
  }
}
