# How to

## Use decorators

### Intent

If your component needs fetch capabilities then you'll need to use `Intent` decorator like this

```js
import { Component } from '@bearer/core'
import { Intent, BearerFetch } from '@bearer/core'

@Component({
  tag: 'my-component-tag-name'
})
class MyComponent {
  @Intent('getPullRequest') getPullRequests: BearerFetch

  render() {
    return (
      <component-which-need-collection-fetcher
        fetcher={this.collectionFetcher}
      />
    )
  }
}
```

Advanced usage

```js
import { Component } from '@bearer/core'
import { Intent, BearerFetch, IntentType } from '@bearer/core'

@Component({
  tag: 'ny-component-tag-name'
})
class MyComponent {
  // getCollection is the name you give to the instance property
  @Intent('getPullRequests') getCollection: BearerFetch

  // You can specify the returned type with IntentType.GetResource
  @Intent('getRepository', IntentType.GetResource)
  getAResource: BearerFetch

  collection: any
  data: any

  componentDidLoad() {
    // Intent accept parameters added to the fetch query
    this.getAResource({ id: '42' }).then(({ object }) => {
      // Promise is resolve with { object }: { object: Object}
      this.data = object
    })
    this.getCollection({ page: 2 }).then(({ items }) => {
      // Promise is resolve with { items } : { items: Array<any>}
      this.collection = items
    })
  }

  render() {
    return <div>...</div>
  }
}
```
