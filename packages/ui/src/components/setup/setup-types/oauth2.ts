import { FieldType } from '../../Forms/Fieldset'

export const OAuth2SetupType = [
  {
    type: 'text' as FieldType,
    label: 'Client API',
    controlName: 'clientId',
    placeholder: 'Type in your Client ID',
    required: true
  },
  {
    type: 'password' as FieldType,
    label: 'Client Secret',
    controlName: 'clientSecret',
    placeholder: 'Type in your Client Secret',
    required: true
  }
]
