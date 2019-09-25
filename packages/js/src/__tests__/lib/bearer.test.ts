import Bearer, { findElements, calculateModalPosition } from '../../lib/bearer'
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

  describe('connect', () => {
    const instance = new Bearer('client-id')
    const openSpy = jest.fn()
    window.open = openSpy

    beforeEach(() => {
      openSpy.mockClear()
    })

    test('opens a modal', () => {
      instance.connect('my-integration')

      expect(openSpy).toHaveBeenCalledWith(
        'INTEGRATION_HOST_URL/v2/auth/my-integration?clientId=client-id',
        '',
        // jsdom has no easy to configure jsdom screen size, unit tests are covering calculateModalPosition
        'toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=no, resizable=no, copyhistory=no, width=0, height=0, top=0, left=0'
      )
    })

    describe('calculateModalPosition - screen variations', () => {
      test('center the modal', () => {
        // @ts-ignore
        const values = calculateModalPosition({ screen: { height: 1000, width: 2000 } }, 500, 600)

        expect(values).toEqual({ left: 750, top: 200, calculatedWidth: 500, calculatedHeight: 600 })
      })

      test('resize the modal for a small screen', () => {
        const values = calculateModalPosition({ screen: { height: 400, width: 300 } }, 500, 600)

        expect(values).toEqual({ left: 0, top: 0, calculatedWidth: 300, calculatedHeight: 400 })
      })
    })

    test('returns a promise', () => {
      expect(instance.connect('my-integration')).toBeInstanceOf(Promise)
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
