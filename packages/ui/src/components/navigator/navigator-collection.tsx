import {
  Component,
  Event,
  EventEmitter,
  Prop,
  State,
  Watch
} from '@bearer/core'

@Component({
  tag: 'bearer-navigator-collection',
  styleUrl: 'navigator-collection-screen.scss',
  shadow: true
})
export class BearerNavigatorCollection {
  @Prop() data: any
  @State() collection: Array<any> = []
  @Prop() renderFunc: (member: any) => void

  @Event() completeScreen: EventEmitter

  select = member => () => this.completeScreen.emit(member)

  @Watch('data')
  dataWatcher(newValue) {
    if (newValue) {
      if (typeof newValue === 'string') {
        this.collection = JSON.parse(newValue)
      } else {
        this.collection = newValue
      }
    } else {
      this.collection = []
    }
  }

  componentDidLoad() {
    this.dataWatcher(this.data)
  }

  render() {
    if (this.collection.length === 0) {
      return (
        <slot name="empty">
          <div class="empty">No results</div>
        </slot>
      )
    }
    return (
      <ul class="list-group">
        {this.collection.map(member => (
          <li
            onClick={!member._isDisabled && this.select(member)}
            class={`list-group-item ${member._isDisabled && 'disabled'}`}
            role="button"
          >
            {this.renderFunc(member)}
          </li>
        ))}
      </ul>
    )
  }
}
