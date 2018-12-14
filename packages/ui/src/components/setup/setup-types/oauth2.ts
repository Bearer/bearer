import { FieldType } from '../../Forms/Fieldset'

export const OAuth2SetupType = [
  {
    type: 'text' as FieldType,
    label: 'Client ID',
    controlName: 'clientId',
    required: true
  },
  {
    type: 'password' as FieldType,
    label: 'Client Secret',
    controlName: 'clientSecret',
    required: true
  }
]
