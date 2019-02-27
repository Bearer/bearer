# @bearer/js

The hassle-free way to use bearer's integrations into any web application

## Getting started

Bearer lib can be used instantly in your page or with a package system.

### Directly in your page

```html
<script src="https://cdn.jsdelivr.net/npm/@bearer/js@beta5/lib/bearer.production.min.js"></script>
<script>
  // you have now access to a global variable
  bearer('clientId')
</script>
```

### With a build system

```bash
yarn add @bearer/js
# or
npm install @bearer/js
```

In your app

```jsx
import bearer from '@bearer/js'

class MyApp {
  componentDidMount() {
    bearer('clientId')
  }
}
```

## Usage

### i18n

`@bearer/js` comes with a i18n module that let you deal with internationalization of Bearer's integrations

**bearer.i18n.locale**

Let you change the current locale you are using

```js
bearer.i18n.locale = 'es'
```

**bearer.i18n.load**

Let you load custom translation for a given (or a set of) integration

```js
// with a simple dictionnary
const dictionnary = { titles: { welcome: 'Ola!' } }
bearer.i18n.load('integration-uuid', dictionnary)

// with a promise returning a dictionnary
const promiseReturningADictionnary = Promise.new((resolve, reject) => {
  // async stuff
  resolve({ titles: { welcome: 'Ola!' } })
})
bearer.i18n.load('integration-uuid', promiseReturningADictionnary)

// for a given locale
const dictionnary = { titles: { welcome: 'Guten Morgen' } }
bearer.i18n.load('integration-uuid', dictionnary, { locale: 'de' })

// for a set of integrations
const dictionnary = {
  ['integration-one-uuid']: { title: { welcome: 'Hello my friend' } },
  ['integration-two-uuid']: { message: { goodby: 'Bye Bye' } }
}
bearer.i18n.load(null, dictionnary)
```

### init options

_Soon_
