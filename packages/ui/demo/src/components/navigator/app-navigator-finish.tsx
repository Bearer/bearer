import { Component, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'app-navigator-finish'
})
export class AppNavigatorFinish {
  @Event() scenarioCompleted: EventEmitter

  render() {
    return <button onClick={() => this.scenarioCompleted.emit()}>Finish</button>
  }
}
