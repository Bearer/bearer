import { Event, EventEmitter, Listen, RootComponent } from '@bearer/core'

@RootComponent({
  group: 'complex-feature',
  role: 'display'
})
export class FeatDisplayRootComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('config:saved')
  eventHandler() {}

  @Listen('body:click')
  clickHandler() {}
}
