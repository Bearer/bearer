import { Component } from '@bearer/core'

@Component({
  tag: 'do-nothing'
})
class DoNothingComponent {
  render() {
    return (
      <div>
        <span>I'm not affected</span>
      </div>
    )
  }
}
