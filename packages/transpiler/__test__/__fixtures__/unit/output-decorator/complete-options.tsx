import { BearerRef, Output, RootComponent, Input } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

@RootComponent({
  name: 'output-optionscomplete'
})
class NoOptionsComponent {
  @Input()
  farmer: BearerRef<Farmer>
  // write this
  @Output()
  farmer: BearerRef<Farmer>
  // // generate this
  // @Event() farmerSaved: EventEmitter<Farmer>
  // @State() farmer: BearerRef<Farmer>
  // @Intent('setFarmer') setFarmer: BearerFetch<Farmer>
  // @Watch('farmer')
  // farmerchangeHandler(newValue: BearerRef<Farmer>) {
  //   if (newValue) {
  //     this.setFarmer().then(({ data, referenceId }) => {
  //       this.farmerSaved.emit({ referenceId, farmer: data }) // farmer
  //     })
  //   } else {
  //     this.farmerSaved.emit({ farmer: newValue })
  //   }
  // }

  @Output({
    eventName: 'milked',
    propertyWatchedName: 'aPanda',
    referenceKeyName: 'aPandaKey',
    intentName: 'milkWithAbottle',
    intentReferenceIdKeyName: 'anotherProp',
    intentReferenceIdValue: 'intentReferenceIdValueValue',
    intentArguments: ['farmer', 'spongeBobOverrided', 'authId'],
    intentPropertyName: 'intentPropertyNameOption',
    autoLoad: false
  })
  farmerAndPanda: BearerRef<Farmer> = {}
}
