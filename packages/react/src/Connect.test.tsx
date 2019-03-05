import * as React from 'react'
import { render, fireEvent, cleanup } from 'react-testing-library'

const integration = 'my-dummy-integration'
const connectTo = jest.fn((_one, _two, { authId }) =>
  Promise.resolve({ data: { integration, authId: authId || 'random' } })
)

jest.mock('./bearer-provider', () => {
  const BearerContext = React.createContext<any>({ bearer: { connectTo } })
  return { BearerContext }
})

import Connect from './Connect'

describe('Connect', () => {
  // afterEach(cleanup)
  describe('when success authentication', () => {
    const success = jest.fn()
    const { getByText, container } = renderConnect({
      success,
      integration,
      setup: 'my-dummy-setup',
      authId: 'my-dummy-auth-id'
    })

    beforeAll(() => {
      fireEvent.click(getByText(/Click/i))
    })

    it('uses render prop', () => {
      expect(container).toMatchSnapshot()
    })

    it('calls connectTo with props provided', () => {
      expect(connectTo).toHaveBeenCalledWith(integration, 'my-dummy-setup', { authId: 'my-dummy-auth-id' })
    })

    it('calls success', () => {
      expect(success).toHaveBeenCalledWith({ integration: 'my-dummy-integration', authId: 'my-dummy-auth-id' })
    })
  })

  describe('when error during authentication', () => {
    const success = jest.fn()
    const error = jest.fn()

    const { getByText, container } = renderConnect({
      success,
      text: 'Click Fail',
      onError: error
    })

    beforeAll(() => {
      connectTo.mockReset()
      connectTo.mockImplementation(() => {
        return Promise.reject(false)
      })
      fireEvent.click(getByText(/Click Fail/i))
    })

    it('uses render prop', () => {
      expect(container).toMatchSnapshot()
    })

    it('does not call success', () => {
      expect(success).not.toHaveBeenCalled()
    })

    it('calls connectTo with props provided', () => {
      expect(connectTo).toHaveBeenCalledWith('dummy', 'my-setup', { authId: 'auth-id' })
    })

    it('calls onError', () => {
      expect(error).toHaveBeenCalledWith({ error: false, authId: 'auth-id', integration: 'dummy' })
    })
  })
})

function renderConnect({
  text = 'Click',
  success = jest.fn(),
  integration = 'dummy',
  setup = 'my-setup',
  authId = 'auth-id',
  onError
}: any) {
  return render(
    <Connect
      integration={integration}
      setupId={setup}
      onSuccess={success}
      onError={onError}
      authId={authId}
      render={({ connect }) => <button onClick={connect}>{text}</button>}
    />
  )
}
