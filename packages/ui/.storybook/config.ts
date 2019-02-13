// @ts-ignore
import { configure } from '@storybook/html'
import '@storybook/addon-console'
// automatically import all files ending in *.stories.tsx
// @ts-ignore
const req = require.context('../stories', true, /.stories.tsx$/)

function loadStories() {
  req.keys().forEach(req)
}

configure(loadStories, module)
