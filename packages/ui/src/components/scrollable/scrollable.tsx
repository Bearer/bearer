import { Component, Listen, Prop, State, Method, Element } from '@bearer/core'
import { Store } from '@stencil/redux'

@Component({
  tag: 'bearer-scrollable',
  styleUrl: 'scrollable.scss'
})
export class BearerScrollable {
  @Prop() renderCollection: (collection: Array<any>) => any
  @Prop() renderFetching: () => any
  @Prop() perPage: number = 5
  @Prop() fetcher: ({ page: number }) => Promise<{ items: Array<any> }>
  @Prop() store: Store
  @Prop() reducer: string

  @State() hasMore: boolean = true
  @State() page: number = 1
  @State() fetching: boolean = false
  @State() collection: Array<any> = []
  @State() content: HTMLElement
  @Element() element: HTMLElement

  @Listen('BearerScrollableNext')
  fetchNext() {
    if (this.store && this.hasMore) {
      const redux = this.store.getState()[this.reducer]
      this.fetching = true
      this.hasMore = redux.length > this.collection.length
      this.collection = [
        ...this.collection,
        ...redux.splice(
          (this.page - 1) * this.perPage,
          this.page * this.perPage
        )
      ]
      this.fetching = false
      this.page = this.page + 1
    }
    if (!this.store && this.hasMore) {
      this.fetching = true
      this.fetcher({ page: this.page })
        .then(({ items }) => {
          this.hasMore = items.length === this.perPage
          this.collection = [...this.collection, ...items]
          this.fetching = false
          this.page = this.page + 1
        })
        .catch(() => {
          this.fetching = false
        })
    }
  }

  _renderCollection = () => {
    if (this.fetching && !this.collection.length) {
      return null
    }
    return this.renderCollection(this.collection)
  }

  _renderFetching = () =>
    this.renderFetching ? this.renderFetching() : <bearer-loading />

  componentDidLoad() {
    if (this.element) {
      this.content = this.element.querySelector('.scrollable-list')
      this.content.addEventListener('scroll', this.onScroll)
      this.fetchNext()
    }
  }

  onScroll = () => {
    if (!this.fetching) {
      if (
        this.content.scrollTop + this.content.clientHeight >=
        this.content.scrollHeight
      ) {
        this.fetchNext()
      }
    }
  }

  @Method()
  reset() {
    this.hasMore = true
    this.collection = []
    this.fetchNext()
  }

  render() {
    return (
      <div class="scrollable-list">
        {this._renderCollection()}
        {this.fetching && this._renderFetching()}
      </div>
    )
  }
}
