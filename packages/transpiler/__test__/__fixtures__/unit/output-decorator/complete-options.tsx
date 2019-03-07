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
  // @Function('setFarmer') setFarmer: BearerFetch<Farmer>
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
    functionName: 'milkWithAbottle',
    functionReferenceIdKeyName: 'anotherProp',
    functionReferenceIdValue: 'functionReferenceIdValueValue',
    functionArguments: ['farmer', 'spongeBobOverrided', 'authId'],
    functionPropertyName: 'functionPropertyNameOption',
    autoLoad: false
  })
  farmerAndPanda: BearerRef<Farmer> = {}
}
