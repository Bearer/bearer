import { Component, Intent, BearerFetch } from '@bearer/core'

@Component({
  tag: 'list-star-wars-characters',
  styleUrl: 'ListStarWarsCharacters.css',
  shadow: true
})
export class HelloWorld {
  @Intent('getStarWarsCharacters') fetcher: BearerFetch

  render() {
    return (
      <div class="root">
        <bearer-typography kind="h4" as="h1">
          Listing Star Wars Characters
        </bearer-typography>
        <bearer-paginator
          fetcher={this.fetcher}
          perPage={10}
          renderCollection={collection => (
            <bearer-navigator-collection
              data={collection}
              renderFunc={item => item.name}
            />
          )}
        />
      </div>
    )
  }
}
