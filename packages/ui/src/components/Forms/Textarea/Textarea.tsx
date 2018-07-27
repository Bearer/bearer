import {
  Component,
  Prop,
  Element,
  Event,
  EventEmitter
  // Watch
} from '@bearer/core'

@Component({
  tag: 'bearer-textarea',
  styleUrl: './Textarea.scss',
  shadow: true
})
export class BearerTextarea {
  @Element() el: HTMLElement
  @Prop() label?: string
  @Prop() controlName: string
  @Prop() placeholder: string
  @Prop({ mutable: true })
  value: string
  @Prop({ mutable: true })
  hint: string
  @Event() valueChange: EventEmitter

  // @Watch('value')
  // valueChanged() {
  //   const inputEl = this.el.shadowRoot.querySelector('input')
  //   if (inputEl.value !== this.value) {
  //     inputEl.value = this.value
  //   }
  // }

  inputChanged(event: any) {
    let val = event.target && event.target.value
    this.value = val
    this.valueChange.emit(this.value)
  }

  render() {
    return (
      <div class="form-group">
        {this.label ? <label>{this.label}</label> : ''}
        <textarea
          name={this.controlName}
          value={this.value}
          placeholder={this.placeholder}
          class="form-control"
          onInput={this.inputChanged.bind(this)}
        >
          {this.value}
        </textarea>
        {this.hint ? <small class="form-text text-muted">{this.hint}</small> : ''}
      </div>
    )
  }
}
