import { Event, EventEmitter, Listen, RootComponent } from '@bearer/core'

@RootComponent({
  name: 'complex-feature'
})
export class FeatDisplayRootComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('config:saved')
  eventHandler() {}

  @Listen('body:eventFromAnotherRootComponent:saved')
  eventFromAnotherRootComponentHandler() {}

  @Listen('body:click')
  clickHandler() {}
}
