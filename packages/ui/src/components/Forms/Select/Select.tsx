import { Component, Prop, Element, Event, EventEmitter } from '@bearer/core'

@Component({
  tag: 'bearer-select',
  styleUrl: './Select.scss',
  shadow: true
})
export class BearerSelect {
  @Element() el: HTMLElement
  @Prop() label?: string
  @Prop() controlName: string
  @Prop({ mutable: true })
  value: string
  @Prop({ mutable: true })
  options: Array<{ label: string; value: string; default?: boolean }> = []
  @Event() valueChange: EventEmitter

  onSelectChange = event => {
    this.valueChange.emit(event.path[0].value)
  }

  componentDidLoad() {
    this.options = [{ label: '--- choose an option ---', value: '' }, ...this.options]
  }

  render() {
    return (
      <div class="form-group">
        {this.label ? <label>{this.label}</label> : ''}
        <select class="form-control" onChange={this.onSelectChange}>
          {this.options.map(value => {
            return (
              <option value={value.value} selected={this.value === value.value ? true : false}>
                {value.label}
              </option>
            )
          })}
        </select>
      </div>
    )
  }
}
