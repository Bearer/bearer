import { Component, Element, Listen, Method, Prop, State, TFetchBearerData } from '@bearer/core'

import { TCollectionRenderer } from './types'

import debug from '../../logger'
const logger = debug('bearer-scrollable')

@Component({
  tag: 'bearer-scrollable',
  styleUrl: 'scrollable.scss'
})
export class BearerScrollable {
  @Prop()
  renderCollection?: TCollectionRenderer
  @Prop()
  rendererProps?: JSXElements.BearerNavigatorCollectionAttributes
  @Prop()
  renderFetching?: () => any
  @Prop()
  perPage?: number = 5
  @Prop()
  fetcher: (params: { page: number }) => Promise<TFetchBearerData>

  @State()
  hasMore: boolean = true
  @State()
  page: number = 1
  @State()
  fetching: boolean = false
  @State()
  collection: any[] = []
  @State()
  content: HTMLElement
  @Element()
  element: HTMLElement

  @Listen('BearerScrollableNext')
  fetchNext() {
    if (this.hasMore) {
      this.fetching = true
      this.fetcher({ page: this.page })
        .then(({ data }: TFetchBearerData) => {
          logger('data receiced from fetcher %j', data)
          this.hasMore = data.length === this.perPage
          this.collection = [...this.collection, ...data]
          this.fetching = false
          this.page = this.page + 1
          return data
        })
        .catch(error => {
          logger('Error while fetching %j', error)
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
