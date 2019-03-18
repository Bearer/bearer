import Bearer, { findElements } from '../../lib/bearer'

const clientId = 'a-client-id'

describe('main', () => {
  it('exports a Bearer class', () => {
    expect(Bearer).toBeTruthy()
  })

  describe('constructor', () => {
    describe('readyState ', () => {
      it('register a listerer on dom DOMContentLoaded if loading', () => {
        Object.defineProperty(document, 'readyState', {
          configurable: true,
          get() {
            return 'loading'
          }
        })
        document.addEventListener = jest.fn()

        new Bearer(clientId)

        expect(document.addEventListener).toHaveBeenCalledWith('DOMContentLoaded', expect.any(Function))
      })

      it('does not register listerner if complete', () => {
        Object.defineProperty(document, 'readyState', {
          get() {
            return 'complete'
          }
        })

        document.addEventListener = jest.fn()

        new Bearer(clientId)

        expect(document.addEventListener).not.toHaveBeenCalled()
      })

      it('does not register listerner if interactive', () => {
        Object.defineProperty(document, 'readyState', {
          get() {
            return 'interactive'
          }
        })

        document.addEventListener = jest.fn()

        new Bearer(clientId)

        expect(document.addEventListener).not.toHaveBeenCalled()
      })
    })
  })

  describe('fetch request', () => {
    it.skip('returns payload if success', () => {
      expect(false).toBeTruthy()
    })
  })
  describe('function request', () => {
    it.skip('returns payload if success', () => {
      expect(false).toBeTruthy()
    })
  })
})

describe('findElements', () => {
  beforeAll(() => {
    document.body.innerHTML = `
      <bearer-something></bearer-something>
      <bearer-patrick></bearer-patrick>
      <spongebob-something></spongebob-something>
      <notmatching-anything></notmatching-anything>
    `
  })

  describe('matching elements from default', () => {
    it('returns lowercased matching tab names', () => {
      expect(findElements(document.querySelectorAll('*'))).toEqual(['bearer-something', 'bearer-patrick'])
    })
  })

  describe('matching elements bearer tag by default', () => {
    it('returns lowercased matching tab names', () => {
      expect(findElements(document.querySelectorAll('*'), /^spongebob-/i)).toEqual(['spongebob-something'])
    })
  })
})
