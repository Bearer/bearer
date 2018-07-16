import BearerConfig from './BearerConfig'

describe('Init', () => {
  it('sets defaults', () => {
    const config = new BearerConfig()
    expect(config.integrationHost).toEqual('BEARER_INTEGRATION_HOST')
    expect(config.loadingComponent).toBeUndefined()
  })

  it('override defaults', () => {
    const config = new BearerConfig({
      integrationHost: 'spongebob',
      loadingComponent: 'ok'
    })
    expect(config.integrationHost).toEqual('spongebob')
    expect(config.loadingComponent).toEqual('ok')
  })
})
