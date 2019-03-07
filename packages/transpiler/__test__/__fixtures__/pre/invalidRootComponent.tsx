import { BearerFetch, BearerState, Function, FunctionType, RootComponent } from '@bearer/core'

@RootComponent({
  name: 'attach-pull-request',
  shadow: false
})
export class InvalidRootComponent {
  @Function('ListRepositories')
  fetcher: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Function('getPullRequest', FunctionType.FetchData)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
