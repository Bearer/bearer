import { storiesOf } from '@storybook/html'

storiesOf('Button', module).add('kind', () => {
  const kinds = [
    'primary',
    'secondary',
    'success',
    'danger',
    'warning',
    'info',
    'light',
    'dark'
  ]
  return kinds
    .map(
      kind =>
        `<div><bearer-button kind="${kind}">Button: ${kind}</bearer-button></div>`
    )
    .join('<br />')
})
