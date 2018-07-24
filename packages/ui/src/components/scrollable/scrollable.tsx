import { Component, Listen, Prop, State, Method, Element } from '@bearer/core'
import { TCollectionRenderer } from './types'

@Component({
  tag: 'bearer-scrollable',
  styleUrl: 'scrollable.scss'
})
export class BearerScrollable {
  @Prop() renderCollection?: TCollectionRenderer
  @Prop() rendererProps?: JSXElements.BearerNavigatorCollectionAttributes
  @Prop() renderFetching?: () => any
  @Prop() perPage?: number = 5
  @Prop() fetcher: ({ page: number }) => Promise<{ items: Array<any> }>

  @State() hasMore: boolean = true
  @State() page: number = 1
  @State() fetching: boolean = false
  @State() collection: Array<any> = []
  @State() content: HTMLElement
  @Element() element: HTMLElement

  @Listen('BearerScrollableNext')
  fetchNext() {
    if (this.hasMore) {
      this.fetching = true
      this.fetcher({ page: this.page })
        .then((data: { items: Array<any> }) => {
          this.hasMore = data.items.length === this.perPage
          this.collection = [...this.collection, ...data.items]
          this.fetching = false
          this.page = this.page + 1
          return data
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
    return (this.renderCollection || this.renderCollectionDefault)(this.collection)
  }

  renderCollectionDefault: TCollectionRenderer = collection => (
    <bearer-navigator-collection {...this.rendererProps} data={collection} />
  )

  _renderFetching = () => (this.renderFetching ? this.renderFetching() : <bearer-loading />)

  componentDidLoad() {
    if (this.element) {
      this.content = this.element.querySelector('.scrollable-list')
      this.content.addEventListener('scroll', this.onScroll)
      this.fetchNext()
    }
  }

  onScroll = () => {
    if (!this.fetching) {
      if (this.content.scrollTop + this.content.clientHeight >= this.content.scrollHeight) {
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
