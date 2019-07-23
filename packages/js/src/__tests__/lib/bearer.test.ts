import Bearer, { findElements } from '../../lib/bearer'
import { IntegrationClient } from '../../lib/integrationClient'

const clientId = 'a-client-id'

describe('bearer', () => {
  it('exports a Bearer class', () => {
    expect(Bearer).toBeTruthy()
  })

  describe('Events', () => {
    it('has singleton event emitter', () => {
      // @ts-ignore
      expect(window['bearer-listeners']).toBeUndefined()

      expect(Bearer.authorizedListener).toEqual(Bearer.authorizedListener)
      // @ts-ignore
      expect(window['bearer-listeners']).not.toBeUndefined()
    })
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

  describe('integration', () => {
    it('returns an instance of an  integration client', () => {
      expect(new Bearer('client-id').integration('integration-id')).toBeInstanceOf(IntegrationClient)
    })
  })

  describe('invoke function', () => {
    it('delegates to Integration Client', () => {
      const instance = new Bearer('client-id')
      const integrationSpy = jest.spyOn(instance, 'integration')
      const invokeSpy = jest.fn()
      integrationSpy.mockReturnValue({ invoke: invokeSpy } as any)

      instance.invoke('integration-name', 'function-name', { someData: 'data value' })

      expect(integrationSpy).toHaveBeenCalledWith('integration-name')
      expect(invokeSpy).toHaveBeenCalledWith('function-name', { someData: 'data value' })
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
