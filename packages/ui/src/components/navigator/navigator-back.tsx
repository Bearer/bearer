import { Component, Prop, Event, EventEmitter } from '@bearer/core'
import BackArrow from './back-arrow'

@Component({
  tag: 'bearer-navigator-back',
  styleUrl: 'navigator-back.scss'
})
export class BearerNavigatorBack {
  @Prop() disabled: boolean = false
  @Event() navigatorGoBack: EventEmitter

  back = () => {
    this.navigatorGoBack.emit()
  }

  render() {
    return (
      <button onClick={this.back} disabled={this.disabled}>
        <BackArrow />
      </button>
    )
  }
}
