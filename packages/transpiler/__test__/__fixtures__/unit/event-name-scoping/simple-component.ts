import { Component, Event, EventEmitter, Listen } from '@bearer/core'

@Component({
  tag: 'event'
})
export class SimpleComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('config:saved')
  eventHandler() {}

  @Listen('eventFromChildren')
  eventFromChildrenHandler() {}

  @Listen('body:click')
  clickHandler() {}
}
