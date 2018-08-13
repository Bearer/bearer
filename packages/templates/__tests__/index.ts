import { apiKey, basicAuth, noAuth, oauth2 } from '../src/index'

const intents = {
  apiKey, basicAuth, noAuth, oauth2
}

describe('Templates', () => {
  Object.keys(intents).forEach((key) => {
    const intent = intents[key]

    describe(key, () => {
      it('has SaveState', () => {
        expect(intent.SaveState).toMatchSnapshot()
      })

      it('has RetrieveState', () => {
        expect(intent.RetrieveState).toMatchSnapshot()
      })

      it('has FetchData', () => {
        expect(intent.FetchData).toMatchSnapshot()
      })
    })
  })
});