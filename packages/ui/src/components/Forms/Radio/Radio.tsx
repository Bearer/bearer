import { Component, Prop, Element, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'bearer-radio',
  styleUrl: './Radio.scss',
  shadow: true
})
export class BearerRadio {
  @Element() el: HTMLElement
  @Prop() label?: string
  @Prop() controlName: string
  @Prop() inline: boolean = false
  @Prop({ mutable: true })
  value: string
  @Prop() buttons: Array<{ label: string; value: string; checked?: boolean }>
  @Event() valueChange: EventEmitter

  inputClicked(event) {
    this.valueChange.emit(event.path[0].value)
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
                type="radio"
                name={this.controlName}
                value={value.value}
                checked={this.value === value.value ? true : false}
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
