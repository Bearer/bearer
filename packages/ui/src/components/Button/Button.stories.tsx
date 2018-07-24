import { storiesOf } from '@storybook/html'

storiesOf('Button', module).add('kind', () => {
  const kinds = ['primary', 'secondary', 'success', 'danger', 'warning', 'info', 'light', 'dark']

  return (
    '<table><tbody>' +
    kinds
      .map(
        kind =>
          `<tr>
          <td>${kind}</td>
          <td>
            <bearer-button kind="${kind}">Button</bearer-button></td>
          <td>
            <bearer-button kind="${kind}" content="Button content" />
            </td>
        </tr>`
      )
      .join() +
    '</tbody></table>'
  )
})
