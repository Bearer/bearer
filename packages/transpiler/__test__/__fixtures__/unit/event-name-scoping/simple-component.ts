import { Component, Event, EventEmitter, Listen } from '@bearer/core'

@Component({
  tag: 'event'
})
export class SimpleComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('configSaved')
  eventHandler() {}

  @Listen('eventFromChildren')
  eventFromChildrenHandler() {}

  // prevent Alice to listen on anything else than her scenario's events
  @Listen('body:click')
  clickHandler() {}
}
