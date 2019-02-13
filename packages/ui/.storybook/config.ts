// @ts-ignore
import { configure, setAddon } from '@storybook/react'
import '@storybook/addon-console'
import JSXAddon from 'storybook-addon-jsx'

setAddon(JSXAddon)

// automatically import all files ending in *.stories.tsx
// @ts-ignore
const req = require.context('../stories', true, /.stories.tsx$/)

function loadStories() {
  req.keys().forEach(req)
}

configure(loadStories, module)
