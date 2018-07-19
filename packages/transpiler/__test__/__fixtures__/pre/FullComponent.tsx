import {
  Component,
  Intent,
  BearerFetch,
  IntentType,
  BearerState
} from '@bearer/core'

@Component({
  tag: 'full-component'
})
export class FullComponent {
  @Intent('ListRepositories') fetcher: BearerFetch
  @BearerState() attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.GetResource)
  fetchResource: BearerFetch

  render() {
    return <bearer-scrollable fetcher={this.fetcher} />
  }
}
