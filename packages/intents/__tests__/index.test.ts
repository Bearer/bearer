import { SaveState, RetrieveState } from '../src/index'

describe('index', () => {
  it('export SaveState', () => {
    expect(SaveState).toBeTruthy()
  })

  it('export RetrieveState', () => {
    expect(RetrieveState).toBeTruthy()
  })
})

describe('SaveState', () => {
  describe('.intent', () => {
    it('does something', () => {
      expect(true).toBeTruthy()
    })
  })
})
