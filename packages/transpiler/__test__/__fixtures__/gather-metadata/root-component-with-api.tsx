import { Event, EventEmitter, Prop, RootComponent } from '@bearer/core'

type TPullRequest = {
  name: string
}

type TPayloadEvent = {
  something: TPullRequest
  somethingElse: string
  anEnum: 'ok' | 'ko'
}
@RootComponent({
  group: 'complex-feature',
  role: 'display'
})
export class FeatDisplayRootComponent {
  @Prop()
  aStringProp: string = 'ok'
  @Prop()
  aNumberProp: number = 5

  @Event()
  nonTypedEvent: EventEmitter

  @Event()
  typedEvent: EventEmitter<{ pullRequest: TPullRequest; aNumber: number; anEnum: 'ok' | 'ko' }>

  @Event()
  typedEventWithType: EventEmitter<TPayloadEvent>
}
