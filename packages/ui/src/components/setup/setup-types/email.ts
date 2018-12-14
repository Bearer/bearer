import { FieldType } from '../../Forms/Fieldset'

export const EmailSetupType = [
  {
    type: 'email' as FieldType,
    label: 'Email',
    controlName: 'email',
    required: true
  },
  {
    type: 'password' as FieldType,
    label: 'Password',
    controlName: 'password',
    required: true
  }
]
