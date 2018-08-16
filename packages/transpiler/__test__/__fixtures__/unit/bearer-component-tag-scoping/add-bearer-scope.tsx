@RootComponent({
  group: 'feature',
  name: 'action'
})
class RootComponent {
  renderFromANotherProp() {
    return <sponge-bob />
  }
  render() {
    return (
      <div>
        <sponge-bob />
        <patrick />
      </div>
    )
  }
}

@Component({
  tag: 'sponge-bob'
})
class SubComponent {
  render() {
    return <div />
  }
}

@Component({
  tag: 'patrick'
})
class Patrick {
  render() {
    return <div />
  }
}
