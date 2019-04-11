import bearer from './index'

describe('index', () => {
  describe('exports', () => {
    it('client is exported', () => {
      expect(bearer).toBeTruthy()
    })
  })
})
