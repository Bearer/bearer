import { Component, Listen, Prop, State, Method, TFetchBearerData } from '@bearer/core'

@Component({
  tag: 'bearer-paginator',
  styleUrl: 'pagination.scss',
  shadow: true
})
export class BearerPaginator {
  @Prop()
  renderCollection: (collection: Array<any>) => any
  @Prop()
  renderFetching: () => any
  @Prop()
  perPage: number = 5
  @Prop()
  pageCount: number
  @Prop()
  fetcher: (refineParams: { page: number }) => Promise<TFetchBearerData>

  @State()
  fetching: boolean = false
  @State()
  currentPage: number = 0
  @State()
  maxPages: number = 0
  @State()
  collection: Array<any> = []

  @Listen('BearerPaginationNext')
  nextHandler() {
    const itemMaybeIndex = this.currentPage * this.perPage + 1
    this.currentPage = this.currentPage + 1
    if (!this.collection[itemMaybeIndex]) {
      this.fetching = true
      this.fetcher({ page: this.currentPage })
        .then(({ data }: TFetchBearerData) => {
          this.collection = [...this.collection, ...data]
          this.fetching = false
        })
        .catch(e => {
          this.fetching = false
          console.error(e)
        })
    }
  }

  @Listen('BearerPaginationPrev')
  prevHandler() {
    this.currentPage = Math.max(this.currentPage - 1, 1)
  }

  @Listen('BearerPaginationGoTo')
  goToHandler(e) {
    this.currentPage = e.detail
  }

  _renderCollection = () => {
    const start = (this.currentPage - 1) * this.perPage
    return this.renderCollection(this.collection.slice(start, start + this.perPage))
  }

  _renderFetching = () => (this.renderFetching ? this.renderFetching() : <bearer-loading />)

  componentDidLoad() {
    this.nextHandler()
  }

  @Method()
  reset() {
    this.currentPage = 0
    this.collection = []
    this.nextHandler()
  }

  render() {
    return (
      <div>
        {this.fetching ? this._renderFetching() : this._renderCollection()}
        {!!this.collection.length && [
          <br />,
          <bearer-pagination currentPage={this.currentPage} hasNext={true} displayPages={false} />
        ]}
      </div>
    )
  }
}
