export interface CheckableInput {
  label: string
  value: string
  default?: boolean
}

export interface Option {
  label: string
  value: string
  checked?: boolean
}

export type FieldType = 'text' | 'password' | 'email' | 'tel' | 'submit' | 'textarea' | 'radio' | 'checkbox' | 'select'

export interface Field {
  type: FieldType
  label: string
  controlName: string
  value?: string
  valueList?: Array<string>

  hint?: string
  placeholder?: string
  inline?: boolean
  buttons?: Array<CheckableInput>
  options?: Array<Option>
}

export class FieldSet {
  private set: Array<Field>

  constructor(set: Array<Field>) {
    this.set = set
  }

  get(controlName: string) {
    return this.set.find(el => el.controlName === controlName)
  }

  getValue(controlName: string) {
    return this.get(controlName).value
  }

  setValue(controlName: string, value: string) {
    this.set.map(el => {
      if (el.controlName === controlName) {
        el.value = value
        return el
      }
      return el
    })
  }

  map(func: (value: Field, index: number, array: Field[]) => {}) {
    return this.set.map(func)
  }

  reduce(func: (previousValue: Field, currentValue: Field, currentIndex: number, array: Field[]) => Field) {
    return this.set.reduce(func)
  }

  filter(func: (value: Field, index: number, array: Field[]) => {}) {
    return this.set.filter(func)
  }
}
