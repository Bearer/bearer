import { Event, EventEmitter, Listen, RootComponent } from '@bearer/core'

@RootComponent({
  name: 'complex-feature'
})
export class FeatDisplayRootComponent {
  @Event()
  mustBeScopedEvent: EventEmitter

  @Event()
  mustBeScopedEvent: EventEmitter

  @Listen('confiSaved')
  eventHandler() {}

  @Listen('body:eventFromAnotherRootComponentSaved')
  eventFromAnotherRootComponentHandler() {}

  @Listen('body:click')
  clickHandler() {}
}
