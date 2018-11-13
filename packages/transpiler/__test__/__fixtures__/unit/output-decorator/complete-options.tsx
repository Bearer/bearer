import { BearerRef, Output, RootComponent } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

@RootComponent({
  group: 'no-options',
  role: 'action'
})
class NoOptionsComponent {
  @Output()
  farmer: BearerRef<Farmer>

  @Output({
    eventName: 'milked',
    propertyWatchedName: 'aPanda',
    referenceKeyName: 'aPandaKey'
  })
  farmer: BearerRef<Farmer>
}
