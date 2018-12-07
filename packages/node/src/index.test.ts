import bearer from './index'

describe('index', () => {
  describe('exports', () => {
    it('client is exported', () => {
      expect(bearer.client).toBeTruthy()
    })

    it('middleware exists', () => {
      expect(bearer.middleware).toBeTruthy()
    })
  })
})
