import templates from '../src/index'
import Authentications from "@bearer/types/lib/Authentications";

describe('Templates', () => {
  Object.values(Authentications).forEach((key) => {
    const intent = templates[key]

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