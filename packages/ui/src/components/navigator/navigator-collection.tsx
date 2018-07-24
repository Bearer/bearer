import { Component, Event, EventEmitter, Prop, State, Watch } from '@bearer/core'
import { TMember, TMemberRenderer } from './types'

@Component({
  tag: 'bearer-navigator-collection',
  styleUrl: 'navigator-collection-screen.scss',
  shadow: true
})
export class BearerNavigatorCollection {
  @Prop() data: any
  @Prop() displayMemberProp: string = 'name'
  @State() collection: Array<TMember> = []
  @Prop() renderFunc: TMemberRenderer<TMember>

  @Event() completeScreen: EventEmitter

  select = (member: TMember) => () => this.completeScreen.emit(member)

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

  defaultRender: TMemberRenderer<TMember> = (member: TMember) => member[this.displayMemberProp]

  render() {
    if (this.collection.length === 0) {
      return (
        <slot name="empty">
          <div class="empty">No results</div>
        </slot>
      )
    }
    const renderer: TMemberRenderer<TMember> = this.renderFunc || this.defaultRender
    return (
      <ul class="list-group">
        {this.collection.map((member: TMember) => (
          <li
            onClick={!member._isDisabled && this.select(member)}
            class={`list-group-item ${member._isDisabled && 'disabled'}`}
            role="button"
          >
            {renderer(member)}
          </li>
        ))}
      </ul>
    )
  }
}
