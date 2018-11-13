import { BearerRef, Output, RootComponent, BearerFetch, Intent } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

@RootComponent({
  group: 'no-options',
  role: 'action'
})
class NoOptionsComponent {
  @Intent('doSomething') fetcher: BearerFetch<Farmer>

  @Output()
  farmer: BearerRef<Farmer>

  @Output({
    eventName: 'milked',
    propertyWatchedName: 'aPanda',
    referenceKeyName: 'aPandaKey'
  })
  farmerAndPanda: BearerRef<Farmer>

  componentWillLoad() {
    this.fetcher().then(({ data }) => {
      this.farmerAndPanda = data
    })
  }
}
