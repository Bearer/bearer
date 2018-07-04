import {
  storiesOf
} from '@storybook/html';

storiesOf('Alert', module)
  .add('Example', () => {
    return (
      `<bearer-alert>My alert message</bearer-alert>`
    )
  })
