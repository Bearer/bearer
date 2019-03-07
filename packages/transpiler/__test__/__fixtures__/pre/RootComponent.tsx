import {
  BearerFetch,
  BearerState,
  Intent,
  IntentType,
  RootComponent,
  Input,
  Output,
  BearerRef,
  Event,
  EventEmitter
} from '@bearer/core'

type Farmer = {
  fullName: string
}

type Panda = {
  fullPandaName: string
}

@RootComponent({
  name: 'AttachPullRequest',
  styleUrl: './ok.css'
})
export class AttachPullRequestAction {
  @Intent('ListRepositories')
  fetcher: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Intent('getPullRequest', IntentType.FetchData)
  fetchResource: BearerFetch
  @Event() inlineDefinition: EventEmitter<{ inlineName: string; inlineNumber: number; other: Panda }>
  @Input()
  farmer: BearerRef<Farmer>
  @Input({
    group: 'SCOPE',
    propertyReferenceIdName: 'goatId',
    eventName: 'goatMilked',
    intentName: 'retrieveGoat',
    autoLoad: false
  })
  goat: BearerRef<{ aString: string; aNumber: number; panda: Panda }>

  @Output()
  feedPanda: BearerRef<Panda>

  @Output({
    intentName: 'manualIntentName',
    intentReferenceIdKeyName: 'manualintentReferenceIdKeyName',
    intentPropertyName: 'manuaLIntentPropertyName',
    eventName: 'spongeBobdSaved',
    propertyWatchedName: 'manualpropertyWatchedName',
    referenceKeyName: 'manualreferenceKeyName'
  })
  shortPanda: BearerRef<{ inlineName: string }>

  render() {
    return (
      <bearer-authorized>
        <bearer-scrollable fetcher={this.fetcher} />
        <span>Something that must be kept</span>
      </bearer-authorized>
    )
  }
}
