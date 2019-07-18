import exportOfIndex from './index'
import defaultOfClient from './client'

describe('index', () => {
  describe('exports', () => {
    it('client is exported', () => {
      expect(exportOfIndex).toEqual(defaultOfClient)
    })
  })
})
