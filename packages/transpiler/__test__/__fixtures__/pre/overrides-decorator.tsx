import { Component, Function, BearerFetch, FunctionType } from '@bearer/core'

@Component({
  tag: 'overrides-decorator-full'
})
class OverridesDecorarotFull {
  @Function('collectionFunction')
  fetcher: BearerFetch
  @Function('collectionFunction', FunctionType.FetchData)
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
