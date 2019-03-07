import {
  BearerFetch,
  BearerRef,
  BearerState,
  Component,
  Input,
  Function,
  FunctionType,
  Output,
  Event,
  EventEmitter,
  t as translate,
  p as pluralize
} from '@bearer/core'

type Panda = {
  fullPandaName: string
}

@Component({
  tag: 'full-component'
})
export class FullComponent {
  @Function('ListRepositories')
  fetcher: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Function('getPullRequest', FunctionType.FetchData)
  fetchResource: BearerFetch
  @Event() anEvent: EventEmitter<Panda>
  @Output() setup: string

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
        <bearer-navigator-auth-screen />
        <bearer-navigator-screen navigationTitle="First Screen">
          <bearer-scrollable fetcher={this.fetcher} />
        </bearer-navigator-screen>
        <bearer-navigator-screen navigationTitle="First Screen">
          <h1>
            <bearer-i18n key="titles.firstScreen" default="Complex one" />
          </h1>
          <bearer-scrollable fetcher={this.fetcher} />
          <span>{translate('text.paragraphs.firstSpan', 'this text is fine {{value}}', { value: 'Sponge bobd' })}</span>
          <span>{pluralize('text.paragraphs.firstSpan', 0, 'Missing translation {{count}}')}</span>
          this text is not fine
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
