import * as React from 'react'
import * as ShallowRenderer from 'react-test-renderer/shallow'

import BearerLoader from './bearer-loader'

describe('BearerLoader', () => {
  const clientId = 'abc123'

  it('renders basic loader tag', () => {
    const renderer = ShallowRenderer.createRenderer()
    const tree = renderer.render(<BearerLoader clientId={clientId} />)
    expect(tree).toMatchSnapshot()
  })

  it('renders basic forward intHost', () => {
    const renderer = ShallowRenderer.createRenderer()
    const tree = renderer.render(<BearerLoader clientId={clientId} intHost="https://int.bearer.sh/v2/" />)
    expect(tree).toMatchSnapshot()
  })
})
