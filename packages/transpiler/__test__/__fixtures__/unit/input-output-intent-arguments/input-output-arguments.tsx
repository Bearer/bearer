import { BearerRef, Output, RootComponent, Input } from '@bearer/core'

type Farmer = {
  id: string
  name: string
}

type Goat = {
  id: string
  name: string
  milk: number
}

@RootComponent({
  group: 'input-output-arguments',
  role: 'action'
})
class InputOutputArgumentsComponent {
  @Input() farmer: BearerRef<Farmer>
  @Output({ intentArguments: ['farmer'] }) goat: BearerRef<Goat>
}
