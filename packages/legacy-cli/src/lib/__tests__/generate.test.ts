import { getComponentVars, getIntentVars } from '../commands/generateCommand'

describe('Generate command', () => {
  describe('get components variables', () => {
    it('formats variables correctly', () => {
      expect(getComponentVars('test')).toEqual({})
    })
  })

  describe('get intents variables', () => {
    it('formats variables correctly', () => {
      expect(getIntentVars('test')).toEqual({})
    })
  })
})
