import { devPortalClient } from '../../src/utils/devPortal'

describe('devPortalClient', () => {
  it('export a factory and expect a command as argument', () => {
    expect(devPortalClient).toHaveLength(1)
  })

  it('creates an instance with a request async method', () => {
    const instance = devPortalClient({ constants: { DeveloperPortalAPIUrl: 'http://test.bearer.sh' } } as any)
    expect(instance.request).toHaveLength(1)
  })
})
