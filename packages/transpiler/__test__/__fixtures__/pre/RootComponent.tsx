import {
  BearerFetch,
  BearerState,
  Function,
  FunctionType,
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
  @Function('ListRepositories')
  fetcher: BearerFetch
  @BearerState()
  attachedPullRequests: Array<any>
  @BearerState({ statePropName: 'goats' })
  ducks: Array<any>
  @Function('getPullRequest', FunctionType.FetchData)
  fetchResource: BearerFetch
  @Event() inlineDefinition: EventEmitter<{ inlineName: string; inlineNumber: number; other: Panda }>
  @Input()
  farmer: BearerRef<Farmer>
  @Input({
    group: 'SCOPE',
    propertyReferenceIdName: 'goatId',
    eventName: 'goatMilked',
    functionName: 'retrieveGoat',
    autoLoad: false
  })
  goat: BearerRef<{ aString: string; aNumber: number; panda: Panda }>

  @Output()
  feedPanda: BearerRef<Panda>

  @Output({
    functionName: 'manualFunctionName',
    functionReferenceIdKeyName: 'manualfunctionReferenceIdKeyName',
    functionPropertyName: 'manuaLFunctionPropertyName',
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
