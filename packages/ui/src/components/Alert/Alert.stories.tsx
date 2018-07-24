import { storiesOf } from '@storybook/html'

storiesOf('Alert', module).add('Kind', () => {
  const kinds = ['primary', 'secondary', 'success', 'danger', 'warning', 'info', 'light', 'dark']
  return kinds.map(kind => `<bearer-alert kind="${kind}">Alert message: ${kind}</bearer-alert>`).join('<br />')
})
