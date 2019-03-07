import { Component, RootComponent } from '@bearer/core'
@RootComponent({
  name: 'feature-action'
})
class RootComponent {
  render() {
    return (
      <bearer-navigator>
        <bearer-navigator-auth-screen />
      </bearer-navigator>
    )
  }
}

@Component({
  tag: 'sponge-bob'
})
class SimpleComponent {
  render() {
    return (
      <bearer-navigator>
        <bearer-navigator-auth-screen />
      </bearer-navigator>
    )
  }
}

@Component({
  tag: 'sponge-bob'
})
class WithBearerAuthorizedComponent {
  render() {
    return <bearer-authorized renderAuthorized={() => <div />} renderUnauthorized={() => <span />} />
  }
}

@Component({
  tag: 'sponge-bob'
})
class OverrideValueGiven {
  render() {
    return (
      <bearer-authorized
        renderAuthorized={() => <div />}
        renderUnauthorized={() => <span />}
        integrationId="spongebob"
      />
    )
  }
}
