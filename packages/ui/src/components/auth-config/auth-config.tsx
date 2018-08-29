import { Component, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'auth-config',
  styleUrl: './auth-config.scss',
  shadow: true
})
export class AuthConfig {
  @Event()
  submit: EventEmitter

  handleSubmit = () => {}

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <bearer-input type="text" value="Hello" />
        <bearer-button type="submit" onSubmit={this.handleSubmit} />
      </form>
    )
  }
}
