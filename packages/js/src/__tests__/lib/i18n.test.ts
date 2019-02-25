import I18n from '../../lib/i18n'

describe('I18n', () => {
  it('export something :-)', () => {
    expect(I18n).toBeTruthy()
  })

  it('has a default locale', () => {
    expect(I18n.locale).toEqual('en')
  })

  describe('set a locale', () => {
    beforeAll(() => {
      I18n.locale = 'pl'
    })

    it(' has changed the locale', () => {
      expect(I18n.locale).toEqual('pl')
    })
  })

  describe('translation loading', () => {
    beforeAll(() => {
      I18n.locale = 'en'
    })

    it('has a load method', () => {
      expect(I18n.load).toBeInstanceOf(Function)
    })

    it('has a get method', () => {
      expect(I18n.get).toBeInstanceOf(Function)
    })

    describe('#load & #get', () => {
      it('accepts object dictionnary ', async () => {
        const dictionnary = { I: { Love: 'Bearer' } }

        await I18n.load('myInte', dictionnary)

        expect(I18n.get('myInte', 'I.Love')).toEqual('Bearer')
      })

      it('passing undefined let you load multiple integrations', async () => {
        const dictionnary = { IntegrationOne: { I: { Love: 'Bearer' } }, IntegrationTwo: { I: { Love: 'Sponge Bob' } } }

        await I18n.load(null, dictionnary)

        expect(I18n.get('IntegrationOne', 'I.Love')).toEqual('Bearer')
        expect(I18n.get('IntegrationTwo', 'I.Love')).toEqual('Sponge Bob')
      })

      it('loads the translation for another locale', async () => {
        const dictionnary = { You: { Love: 'Bearer' } }

        await I18n.load('withLocale', dictionnary, { locale: 'ru' })

        expect(I18n.get('IntegrationTwo', 'You.Love')).toBeUndefined()
        expect(I18n.get('withLocale', 'You.Love')).toBeUndefined()
        expect(I18n.get('withLocale', 'You.Love', { locale: 'ru' })).toEqual('Bearer')
      })

      it('accepts promise returning a TranslationValue ', async () => {
        await I18n.load('promise', Promise.resolve({ other: { key: 'Lazy Bob' } }))
        expect(I18n.get('promise', 'other.key')).toEqual('Lazy Bob')
      })

      it('does not replace existing translations', async () => {
        await I18n.load('existing', Promise.resolve({ other: { key: 'Sponge Bob' } }))
        await I18n.load('existing', Promise.resolve({ somethingElse: { key: 'Patrick' } }))

        expect(I18n.get('existing', 'other.key')).toEqual('Sponge Bob')
        expect(I18n.get('existing', 'somethingElse.key')).toEqual('Patrick')
      })
    })
    describe('events', () => {
      beforeEach(() => {
        I18n.locale = 'en'
      })

      it('triggers an event when the locale dictionnary is loaded', async () => {
        const handler = jest.fn()
        document.addEventListener('bearer-locale-changed', handler)

        await I18n.load('existing', { other: { key: 'Sponge Bob' } })

        expect(handler).toHaveBeenCalledWith(
          expect.objectContaining({ detail: expect.objectContaining({ locale: 'en' }) })
        )
      })

      it('triggers an event when the locale change', () => {
        const handler = jest.fn()
        document.addEventListener('bearer-locale-changed', handler)

        I18n.locale = 'fr'

        expect(handler).toHaveBeenCalledWith(
          expect.objectContaining({ detail: expect.objectContaining({ locale: 'fr' }) })
        )
      })
    })
  })
})
