import { Component } from '@bearer/core'
import singers from '../navigator/data.json'

const SINGERS = singers
  .concat(singers)
  .concat(singers)
  .concat(singers)
  .concat(singers)
  .map((item, index) => ({ ...item, name: `${index} - ${item.name}` }))

const perPage = 5

@Component({
  tag: 'app-pagination'
})
export class AppPagination {
  renderCollection = collection => (
    <bearer-navigator-collection
      data={collection}
      renderFunc={item => <navigator-collection-item item={item} />}
    />
  )

  fetcher = ({ page }: { page: number }): Promise<{ items: Array<any> }> => {
    return new Promise((resolve, _reject) => {
      setTimeout(() => {
        const start = (page - 1) * perPage
        resolve({
          items: SINGERS.slice(start, start + perPage)
        })
      }, 1000)
    })
  }

  emptyFetcher = (): Promise<{ items: Array<any> }> => {
    return new Promise((resolve, _reject) => {
      resolve({
        items: []
      })
    })
  }

  render() {
    return (
      <div>
        <h4>Existing collection</h4>
        <bearer-paginator
          fetcher={this.fetcher}
          renderCollection={this.renderCollection}
          perPage={perPage}
        />
        <hr />
        <h4>No item found</h4>
        <bearer-paginator
          fetcher={this.emptyFetcher}
          renderCollection={this.renderCollection}
          perPage={perPage}
        />
      </div>
    )
  }
}
