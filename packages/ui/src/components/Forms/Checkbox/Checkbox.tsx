import { Component, Prop, Element, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'bearer-checkbox',
  styleUrl: './Checkbox.scss',
  shadow: true
})
export class BearerCheckbox {
  @Element() el: HTMLElement
  @Prop() label?: string
  @Prop() controlName: string
  @Prop() inline: boolean = false
  @Prop({ mutable: true })
  value: Array<string> = []
  @Prop() buttons: Array<{ label: string; value: string; checked?: boolean }>
  @Event() valueChange: EventEmitter

  inputClicked(event) {
    const index = this.value ? this.value.indexOf(event.path[0].value) : -1
    if (index >= 0) {
      this.value.splice(index, 1)
      this.valueChange.emit(this.value)
    } else {
      this.value.push(event.path[0].value)
      this.valueChange.emit(this.value)
    }
  }

  render() {
    const css = this.inline ? 'form-check form-check-inline' : 'form-check'
    return (
      <div class="form-group">
        {this.label ? <label>{this.label}</label> : ''}
        {this.buttons.map(value => {
          return (
            <div class={css}>
              <input
                class="form-check-input"
                type="checkbox"
                name={this.controlName}
                value={value.value}
                checked={this.value && this.value.indexOf(value.value) >= 0 ? true : false}
                onClick={this.inputClicked.bind(this)}
              />
              <label class="form-check-label">{value.label}</label>
            </div>
          )
        })}
      </div>
    )
  }
}
