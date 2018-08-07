import { Component, Intent, BearerFetch, IntentType } from '@bearer/core'

@Component({
  tag: 'overrides-decorator-full'
})
class OverridesDecorarotFull {
  @Intent('collectionIntent')
  fetcher: BearerFetch
  @Intent('collectionIntent', IntentType.FetchData)
  memberFetcher: BearerFetch

  constructor() {}

  componentDidLoad() {
    console.log('componentDidLoad')
  }

  componentWillLoad() {
    console.log('componentWillLoad')
  }

  componentDidUnload() {
    console.log('componentDidUnload')
  }
}
