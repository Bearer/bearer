import { BearerRef, Input, RootComponent, State, Prop } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

type Sponge = {
  id: unknown
  name: unknown
}

@RootComponent({
  group: 'no-options',
  role: 'action'
})
class NoOptionsComponent {
  @Input()
  farmer: BearerRef<Farmer>

  @Input()
  aString: BearerRef<string> = 'ok'

  @Input()
  aStringWithoutInitializer: BearerRef<string>

  @Input()
  object: BearerRef<{ title: string }> = { title: 'Guest' }

  @Input()
  objectWithoutInitializer: BearerRef<{ title: string }>

  @Input({
    group: 'other-group',
    propertyReferenceIdName: 'patrick',
    eventName: 'patrickWasKilled',
    intentName: 'killPatrick',
    autoLoad: false
  })
  spongeBob: BearerRef<Sponge>

  @Input({
    propertyReferenceIdName: 'refNotOverrided'
  })
  spongeBobOverrided: BearerRef<Sponge>
  // won't be injected by the transfomer
  @State() refNotOverrided: string
}
