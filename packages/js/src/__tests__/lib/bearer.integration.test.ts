import fetch from 'jest-fetch-mock'

// @ts-ignore
const customGlobal: any = global
customGlobal.fetch = fetch
customGlobal.fetchMock = customGlobal.fetch

import Bearer from '../../lib/bearer'

// must be run within its own file
describe('e2e testing', () => {
  it('inject scripts', async done => {
    expect.assertions(1)

    // @ts-ignore
    fetch.mockResponseOnce(
      JSON.stringify([{ uuid: 'patrick', asset: 'patrick-url' }, { uuid: 'something', asset: 'something-url' }])
    )

    document.body.innerHTML = `
    <bearer-something></bearer-something>
    <bearer-patrick></bearer-patrick>
    <spongebob-something></spongebob-something>
    <notmatching-anything></notmatching-anything>
  `

    new Bearer('provided-client-id', { refreshDebounceDelay: 1 })

    await new Promise((resolve, _reject) =>
      setTimeout(() => {
        expect(document.body).toMatchSnapshot()
        resolve()
        done()
      }, 300)
    )
  })
})
