import { Component } from '@bearer/core'
import dataSingers from './data.json'
const singers = dataSingers.map(
  (singer, index) =>
    index !== 2
      ? singer
      : {
          ...singer,
          name: `[DISABLED] ${singer.name}`,
          _isDisabled: true
        }
)

@Component({
  tag: 'app-navigator'
})
export class AppNavigator {
  renderFunc = ({
    prev,
    data
  }: {
    data: any
    next: () => void
    prev: () => void
  }) => (
    <div>
      Rendered through a render prop
      <div>
        <pre>{JSON.stringify(data)}</pre>
        <bearer-button onClick={prev}>Previous</bearer-button>
        <app-navigator-next />
      </div>
    </div>
  )

  render() {
    return (
      <div style={{ width: '500px' }}>
        <b>Navigator with authentication</b>
        <br />
        <bearer-navigator>
          <bearer-navigator-auth-screen setupId="4l1c3-goats-for-fun" />
          <bearer-navigator-screen>Authenticated</bearer-navigator-screen>
        </bearer-navigator>
        <hr />
        <bearer-navigator>
          <bearer-navigator-screen
            navigationTitle="Collection Screen"
            name="singer"
          >
            Click an item to go next screen
            <bearer-navigator-collection
              data={singers}
              renderFunc={item => <navigator-collection-item item={item} />}
            />
          </bearer-navigator-screen>

          <bearer-navigator-screen
            navigationTitle="Render Prop"
            renderFunc={this.renderFunc}
          />

          <bearer-navigator-screen
            navigationTitle="Render Prop 2"
            renderFunc={({
              prev,
              data
            }: {
              data: any
              next: () => void
              prev: () => void
            }) => (
              <div>
                Rendered through a render prop
                <pre>{JSON.stringify(data)}</pre>
                <div>
                  <br />
                  <br />
                  <bearer-button onClick={prev}>Previous</bearer-button>
                  <app-navigator-next />
                </div>
              </div>
            )}
          />

          <bearer-navigator-screen navigationTitle="Final screen">
            Three
            <br />
            <app-navigator-finish />
          </bearer-navigator-screen>
        </bearer-navigator>
      </div>
    )
  }
}
