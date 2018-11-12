import { BearerRef, Input, RootComponent } from '@bearer/core'

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
  aString: BearerRef<string> = "ok"

  @Input()
  aStringWithoutInitializer: BearerRef<string>

  @Input()
  object: BearerRef<{ title: string }> = { title: 'Guest' }

  @Input()
  objectWithoutInitializer: BearerRef<{ title: string }>


  @Input({
    scope: 'other-scope'
    propName: 'patrick'
    eventName: 'patrickWasKilled'
    intentName: 'killPatrick'
    autoUpdate: false
  })
  spongeBob: BearerRef<Sponge>
}
