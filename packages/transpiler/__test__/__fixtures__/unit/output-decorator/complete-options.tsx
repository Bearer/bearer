import { BearerRef, Output, RootComponent } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

@RootComponent({
  group: 'no-options',
  role: 'action'
})
class NoOptionsComponent {
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
    autoLoad: false
  })
  farmerAndPanda: BearerRef<Farmer> = {}
}
