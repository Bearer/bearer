import { Component, State } from '@bearer/core'
import { FieldSet } from '../ui/Forms/FieldSet'

@Component({
  tag: 'app-forms',
  styleUrl: 'app-forms.scss'
})
export class AppForms {
  @State()
  inputs: {
    login: string
    password: string
    gender: string
    country: string
    leisure: Array<string>
  } = {
    login: '',
    password: '',
    gender: 'male',
    country: 'GB',
    leisure: []
  }
  @State()
  radioValues = [
    {
      label: 'Male',
      value: 'male'
    },
    {
      label: 'Female',
      value: 'female'
    },
    {
      label: 'Other',
      value: 'other'
    }
  ]
  @State()
  checkboxValues = [
    {
      label: 'Sailing',
      value: 'sailing'
    },
    {
      label: 'Rowing',
      value: 'rowing'
    },
    {
      label: 'Climbing',
      value: 'climbing'
    },
    {
      label: 'Scuba-diving',
      value: 'scuba-diving'
    }
  ]

  @State()
  countries = [
    { label: 'Austria', value: 'A' },
    { label: 'Belgium', value: 'B' },
    { label: 'France', value: 'F' },
    { label: 'Germany', value: 'D' },
    { label: 'Greece', value: 'GR' },
    { label: 'Iceland', value: 'IS' },
    { label: 'Irland', value: 'IRL' },
    { label: 'Italy', value: 'I' },
    { label: 'Luxembourg', value: 'L' },
    { label: 'Netherland', value: 'NL' },
    { label: 'Poland', value: 'PL' },
    { label: 'Portugal', value: 'P' },
    { label: 'South Africa', value: 'ZA' },
    { label: 'Spain', value: 'E' },
    { label: 'Switzerland', value: 'CH' },
    { label: 'United Kingdom', value: 'GB' }
  ]

  @State() showResult: boolean = false

  @State()
  fields: FieldSet = new FieldSet([
    {
      label: 'Username',
      type: 'text',
      value: '',
      controlName: 'username'
    },
    {
      label: 'Email address',
      type: 'email',
      value: '',
      controlName: 'email',
      hint: 'Do not forget the @ sign'
    },
    {
      label: 'Phone number',
      type: 'tel',
      value: '',
      placeholder: 'Enter your phone',
      controlName: 'phonenumber'
    },
    {
      label: 'Password',
      type: 'password',
      value: '',
      controlName: 'password'
    },
    {
      label: 'Age range',
      type: 'radio',
      controlName: 'range',
      inline: true,
      buttons: [
        {
          label: '< 18',
          value: '0'
        },
        {
          label: '18 - 35',
          value: '1'
        },
        {
          label: '35 - 50',
          value: '2'
        },
        {
          label: '50 - 65',
          value: '3'
        },
        {
          label: '> 65',
          value: '4'
        }
      ]
    },
    {
      label: 'Leisure',
      type: 'checkbox',
      controlName: 'leisure',
      buttons: [
        {
          label: 'Polo',
          value: 'polo'
        },
        {
          label: 'Golf',
          value: 'golf'
        },
        {
          label: 'Tennis',
          value: 'tennis'
        }
      ]
    },
    {
      label: 'Country',
      type: 'select',
      controlName: 'country',
      options: this.countries
    },
    {
      label: 'story',
      type: 'textarea',
      controlName: 'story',
      placeholder: 'Your story'
    }
  ])

  handleValue(field, value) {
    this.inputs[field] = value.detail
  }

  handleSubmit(index: number, e: any) {
    e.preventDefault()
    if (index === 1) {
      this.fields = new FieldSet(e.detail.set)
    }
    if (index === 0) {
      this.showResult = true
    }
  }

  render() {
    return (
      <div>
        <div>
          <bearer-typography kind="h2" id="demo">
            Demo
          </bearer-typography>
          <bearer-typography kind="h3" id="input-component">
            Form with inputs
          </bearer-typography>
          <form onSubmit={e => this.handleSubmit(0, e)}>
            <bearer-input
              type="text"
              controlName="login"
              label="Login"
              value={this.inputs.login}
              onValueChange={value => this.handleValue('login', value)}
            />
            <bearer-input
              type="password"
              controlName="password"
              label="Password"
              value={this.inputs.password}
              onValueChange={value => this.handleValue('password', value)}
            />
            <bearer-textarea
              label="Tell the world your story"
              controlName="story"
              placeholder="A long, long time ago, in a galaxy far, far away..."
              hint="wookies are fun"
            />
            <bearer-radio
              label="Gender"
              controlName="gender"
              onValueChange={value => this.handleValue('gender', value)}
              buttons={this.radioValues}
              value={this.inputs.gender}
            />
            <bearer-select
              label="Country"
              controlName="country"
              onValueChange={value => this.handleValue('country', value)}
              options={this.countries}
              value={this.inputs.country}
            />
            <bearer-checkbox
              label="Leisure"
              controlName="leisure"
              onValueChange={value => this.handleValue('leisure', value)}
              buttons={this.checkboxValues}
              value={this.inputs.leisure}
            />
            <bearer-input type="submit" onSubmit={e => this.handleSubmit(0, e)} />
          </form>
          {this.showResult ? (
            <p>
              Login: {this.inputs.login} - Password: {this.inputs.password}
              <br />
              Radio: {this.inputs.gender}
              <br />
              Country: {this.inputs.country}
              <br />
              Checkbox: {this.inputs.leisure.map(el => el + ' ')}
            </p>
          ) : (
            ''
          )}
        </div>
        <div>
          <h3 id="form-component">Generated form</h3>
          <bearer-form fields={this.fields} onSubmit={e => this.handleSubmit(1, e)} />
          <ul>
            {this.fields.map(el => (
              <li>
                {el.label} - {el.value}
              </li>
            ))}
          </ul>
        </div>
      </div>
    )
  }
}
