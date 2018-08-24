import {
  BearerFetch,
  BearerState,
  Intent,
  IntentType,
  RootComponent
} from '@bearer/core'

@RootComponent({
  group: 'AttachPullRequest',
  role: 'action'
})
export class AttachPullRequestAction {
  @Intent('ListRepositories')
  fetcher: BearerFetch
  @RetrieveStateIntent()
  retrieve: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.FetchData)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
