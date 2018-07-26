import { Component, Intent, BearerFetch, IntentType, BearerState } from '@bearer/core'

@Component({
  tag: 'full-component'
})
export class FullComponent {
  @Intent('ListRepositories') fetcher: BearerFetch
  @RetrieveStateIntent() retrieve: BearerFetch
  @BearerState() attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.GetResource)
  fetchResource: BearerFetch

  screenRenderer = () => {
    return (
      <bearer-navigator-screen navigationTitle="Last Screen">
        <h1>Hello Partick</h1>
      </bearer-navigator-screen>
    )
  }
  render() {
    return (
      <bearer-navigator>
        <bearer-navigator-screen navigationTitle="First Screen">
          <bearer-scrollable fetcher={this.fetcher} />
        </bearer-navigator-screen>
        <bearer-navigator-screen navigationTitle={({ data }) => data.name}>
          <bearer-scrollable fetcher={this.fetcher} />
        </bearer-navigator-screen>

        <bearer-navigator-screen
          renderFunc={({ data, next, prev }) => (
            <last-screen
              next={next}
              complete={({ complete }) => {
                console.log('complete')
                complete()
              }}
            />
          )}
        />
        {this.screenRenderer()}
      </bearer-navigator>
    )
  }
}
