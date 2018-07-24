import { Component, Prop } from '@bearer/core'
import { createStore, combineReducers } from 'redux'
import { Store } from '@stencil/redux'
import singers from '../navigator/data.json'

const SINGERS = singers
  .concat(singers)
  .concat(singers)
  .concat(singers)
  .map((item, index) => ({ ...item, name: `${index} - ${item.name}` }))
const PER_PAGE = 5

/**
 * Store configuration
 */
const demoReducer = (state = SINGERS, action) => {
  switch (action.type) {
    default:
      return state
  }
}

const rootReducer = (combineReducers as any)({
  demoReducer
})

const configureStore = (preloadedState: any) => createStore(rootReducer, preloadedState)

/**
 * Component
 */
@Component({
  tag: 'app-scrollable'
})
export class AppScrollable {
  @Prop({ context: 'store' })
  store: Store

  renderCollection = collection => (
    <bearer-navigator-collection data={collection} renderFunc={item => <navigator-collection-item item={item} />} />
  )

  fetcher = ({ page }: { page: number }): Promise<{ items: Array<any> }> => {
    return new Promise((resolve, _reject) => {
      const start = (page - 1) * PER_PAGE
      setTimeout(() => {
        resolve({
          items: SINGERS.slice(start, start + PER_PAGE)
        })
      }, 1000)
    })
  }

  componentWillLoad() {
    this.store.setStore(configureStore({}))
    this.store.mapStateToProps(this, state => {
      const {
        demoReducer: { name }
      } = state
      return {
        name
      }
    })
    this.store.mapDispatchToProps(this, {
      addPullRequest: () => (dispatch, _state) => dispatch({ type: 'CLICK', payload: { id: Math.random() * 1000 } })
    })
    console.log(this.store)
  }

  render() {
    return (
      <div>
        <h4>Existing collection</h4>
        <bearer-scrollable
          perPage={PER_PAGE}
          store={this.store}
          reducer="demoReducer"
          fetcher={this.fetcher}
          renderCollection={this.renderCollection}
        />
      </div>
    )
  }
}
