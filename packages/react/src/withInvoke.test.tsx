import React, { createContext } from 'react'
import * as Renderer from 'react-test-renderer'
import * as ShallowRenderer from 'react-test-renderer/shallow'
import { BearerContext } from './bearer-provider'
import { withInvoke } from './withInvoke'

interface TProps {
  loading: boolean
  error?: any
  data?: { title: string }
  invoke: (params: any) => void
  sommething: string
}

class DummyWithFetchDataComponent extends React.Component<TProps> {
  render() {
    if (this.props.loading) {
      return <div>Loading</div>
    }

    if (this.props.error) {
      return <div>Error: {this.props.error}</div>
    }

    if (this.props.data) {
      return <div>Some data passed: {this.props.data}</div>
    }

    return (
      <div>
        <button onClick={this.props.invoke}>Click to invoke</button>
      </div>
    )
  }
}

describe('withInvoke', () => {
  it('exports a function', () => {
    expect(withInvoke).toBeInstanceOf(Function)
    expect(withInvoke).toHaveLength(2)
  })

  describe('Wrapped Component', () => {
    const renderer = ShallowRenderer.createRenderer()
    const WithFetch = withInvoke<{ title: string }, { sommething: string }>('integrationName', 'GimmeData')(
      DummyWithFetchDataComponent
    )

    it('renders the wrapped component adn forward props', () => {
      const tree = renderer.render(<WithFetch sommething="ok" />)

      expect(tree).toMatchSnapshot()
    })

    it('sets the display name', () => {
      expect(WithFetch.displayName).toEqual('withBearerInvoke(GimmeData)(DummyWithFetchDataComponent)')
    })
  })

  describe('invoke behaviours', () => {
    const WithFetch = withInvoke<{ title: string }, any>('integrationName', 'GimmeData')(DummyWithFetchDataComponent)

    function render(bearer: any) {
      return Renderer.create(
        <BearerContext.Provider value={{ bearer, state: {} }}>
          <WithFetch />
        </BearerContext.Provider>
      )
    }

    it('forwards loading on invoke', () => {
      const bearer = {
        functionFetch: jest.fn(() => new Promise(() => {}))
      } as any

      const rendered = render(bearer)

      rendered.root.findByType('button').props.onClick()
      expect(rendered.toJSON()).toMatchSnapshot()
    })

    it('forwards data when success', async () => {
      const bearer = {
        functionFetch: jest.fn(() => Promise.resolve({ data: 'Sponge bob is alive' }))
      } as any

      const rendered = render(bearer)

      rendered.root.findByType('button').props.onClick()
      // we delay expectations ;-)
      await bearer.functionFetch()

      expect(rendered.toJSON()).toMatchSnapshot()
    })

    it('forwards data when failure', async () => {
      const bearer = {
        functionFetch: jest.fn(() => Promise.reject({ error: 'data error' }))
      } as any

      const rendered = render(bearer)

      rendered.root.findByType('button').props.onClick()
      // we delay expectations ;-)
      try {
        await bearer.functionFetch()
      } catch (_e) {}

      expect(rendered.toJSON()).toMatchSnapshot()
    })

    describe('callbacks', () => {
      // TODO
    })
  })
})
