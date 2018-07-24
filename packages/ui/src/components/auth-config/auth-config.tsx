import { Component, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'auth-config',
  styleUrl: './auth-config.scss',
  shadow: true
})
export class AuthConfig {
  @Event() submit: EventEmitter

  handleSubmit() {}

  render() {
    return (
      <form onSubmit={() => this.handleSubmit()}>
        <bearer-input type="text" label="hello" controlName="hello" value="Hello" hint="hello" />
        <bearer-input type="submit" onSubmit={() => this.handleSubmit()} />
      </form>
    )
  }
}
