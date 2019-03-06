import { Component, RootComponent } from '@bearer/core'

@RootComponent({
  name: 'feature'
})
class RootComponent {
  renderFromANotherProp() {
    return <sponge-bob title="Sponge element" />
  }
  render() {
    return (
      <div>
        <sponge-bob>
          <span>Something</span>
        </sponge-bob>
        <patrick title="Patrick element" />
      </div>
    )
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

@Component({
  tag: 'patrick'
})
class Patrick {
  render() {
    return <div title="Patrick" />
  }
}
