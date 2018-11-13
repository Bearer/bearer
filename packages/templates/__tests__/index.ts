import templates from '../src/index'
import Authentications from "@bearer/types/lib/authentications";

describe('Templates', () => {
  Object.keys(Authentications).forEach((key) => {
    const intent = templates[Authentications[key]]

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