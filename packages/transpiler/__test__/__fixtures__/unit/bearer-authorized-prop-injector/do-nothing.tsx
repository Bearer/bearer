import { Component, RootComponent } from '@bearer/core'
@RootComponent({
  group: 'feature',
  role: 'action'
})
class RootComponent {
  renderFromANotherProp() {
    return <sponge-bob title="Sponge element" />
  }
  render() {
    return <div />
  }
}

@Component({
  tag: 'sponge-bob'
})
class SubComponent {
  render() {
    return <div title="spongeBNob" />
  }
}
