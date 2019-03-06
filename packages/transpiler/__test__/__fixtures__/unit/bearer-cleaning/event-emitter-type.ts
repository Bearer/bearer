import { EventEmitter, Event, Prop, RootComponent, Input, Output } from '@bearer/core'
type Goat = {
  color: string
}
type Panda = {
  panda: string
}
@RootComponent({
  name: 'feature-display'
})
export class Component {
  @Input() input: string
  @Output() output: string
  @Event() feed: EventEmitter
  // TODO: Keep the type if it is used with prop or method: ex goat
  // note: keeping it aside until we really need it
  @Event() keptAsIs: EventEmitter<{ goat: Goat }>
  @Prop() goat: Goat
  @Prop() setupId: string

  // transformed to any because isn't imported by Stencil later
  @Event()
  pandaIsAnyfied: EventEmitter<Panda>
}
