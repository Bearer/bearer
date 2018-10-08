import { Component, Event, EventEmitter, Listen } from '@bearer/core'

@Component({
  tag: 'event'
})
export class SimpleComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('config:saved')
  eventHandler() {}

  @Listen('body:click')
  clickHandler() {}
}
