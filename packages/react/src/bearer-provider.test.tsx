import * as React from 'react'
import * as Renderer from 'react-test-renderer'
import * as ShallowRenderer from 'react-test-renderer/shallow'

import DummyContextConsumer from '../__test__/utils/dummy-context-consumer'

import BearerProvider from './bearer-provider'

describe('BearerProvider', () => {
  const clientId = 'abc123'
  const initialContext = {
    'setup-id': '123'
  }

  it('renders basic provider tag', () => {
    const renderer = ShallowRenderer.createRenderer()
    const tree = renderer.render(<BearerProvider clientId={clientId} />)
    // TODO: move to matchObject at some point
    expect(tree).toMatchSnapshot()
  })

  it('allow custom integration host', () => {
    const renderer = ShallowRenderer.createRenderer()
    const tree = renderer.render(
      <BearerProvider clientId={clientId} integrationHost="https://integrations.bearer.sh/" />
    )
    // TODO: move to matchObject at some point
    expect(tree).toMatchSnapshot()
  })

  it('default state is sent to context components', () => {
    const testRender = Renderer.create(
      <BearerProvider clientId={clientId} initialContext={initialContext}>
        <DummyContextConsumer />
      </BearerProvider>
    )
    // TODO: move to matchObject at some point
    expect(testRender.toJSON()).toMatchSnapshot()
  })

  it('updates are handled', () => {
    const changeHandler = jest.fn(_x => null)
    const testRender = Renderer.create(
      <BearerProvider clientId={clientId} initialContext={initialContext} onUpdate={changeHandler}>
        <DummyContextConsumer />
      </BearerProvider>
    )
    const stopPropagation = jest.fn(_x => null)
    const mockEvent = {
      stopPropagation,
      detail: { someRefrenceId: 'abc' }
    }

    testRender.root.findByType('button').props.onClick(mockEvent)

    expect(stopPropagation.mock.calls.length).toBe(1)
    expect(changeHandler.mock.calls.length).toBe(1)
    expect(changeHandler.mock.calls[0][0].integrationState).toEqual({
      ...initialContext,
      'some-refrence-id': 'abc'
    })
  })
})
