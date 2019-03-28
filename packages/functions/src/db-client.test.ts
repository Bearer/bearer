import factory, { DBClient } from './db-client'

describe('db-client', () => {
  describe('factory', () => {
    it('exports a function', () => {
      expect(factory).toHaveLength(2)
    })

    it('returns a instance', () => {
      expect(factory('url', 'asdasd')).toBeInstanceOf(DBClient)
    })
  })

  describe('DBClient', () => {
    const instance = factory('url', 'asdasd')

    describe('#getData', () => {
      describe.each([false, '', undefined])('%s no reference given', falsyValue => {
        const promise = instance.getData(falsyValue)

        it('returns a promise', () => {
          expect(promise).toBeInstanceOf(Promise)
        })

        it('resolves with no referenceId and ReadAllowed', () => {
          return expect(promise).resolves.toMatchObject({
            Item: {
              ReadAllowed: false,
              referenceId: falsyValue
            }
          })
        })
      })

      describe('items calls', () => {
        it('returns data from the service', async () => {
          const get = jest.fn(() => {
            return Promise.resolve({ data: { Item: { something: true } } })
          })

          // @ts-ignore
          instance.client = { get }
          const result = await instance.getData('a-refe')
          expect(get).toBeCalledWith('api/v2/items/a-refe', { params: { signature: 'asdasd' } })
          expect(result).toEqual({ Item: { something: true } })
        })

        it('returns empty data if not found', async () => {
          const get = jest.fn(() => {
            return Promise.reject({ response: { status: 404 } })
          })

          // @ts-ignore
          instance.client = { get }
          const result = await instance.getData('a-refe')

          expect(get).toBeCalledWith('api/v2/items/a-refe', { params: { signature: 'asdasd' } })
          expect(result).toEqual({ Item: { ReadAllowed: false, referenceId: 'a-refe' } })
        })

        it('throws an error if not found', () => {
          const get = jest.fn(() => {
            return Promise.reject({ response: { status: 500 } })
          })

          // @ts-ignore
          instance.client = { get }
          instance.getData('a-refe')

          expect.assertions(2)

          expect(get).toBeCalledWith('api/v2/items/a-refe', { params: { signature: 'asdasd' } })
          return expect(instance.getData('a-refe')).rejects.toThrowError('Error while retrieving data')
        })
      })
    })

    describe('#upsertData', () => {
      describe.each([false, '', undefined])('%s no reference given', falsyValue => {
        const promise = instance.upsertData(falsyValue, { someData: 'ok' })

        it('returns a promise', () => {
          expect(promise).toBeInstanceOf(Promise)
        })

        it('resolves with no referenceId and ReadAllowed', () => {
          return expect(promise).rejects.toThrowError(`Error while updating data: no reference given`)
        })
      })

      describe('items calls', () => {
        it('returns formatted data', async () => {
          const put = jest.fn(() => {
            return Promise.resolve()
          })

          // @ts-ignore
          instance.client = { put }

          const result = await instance.upsertData('a-reference', { someData: 'ok' })

          expect(put).toBeCalledWith(
            'api/v2/items/a-reference',
            { data: { someData: 'ok' }, ReadAllowed: true, referenceId: 'a-reference' },
            { params: { signature: 'asdasd' } }
          )
          expect(result).toEqual({ Item: { data: { someData: 'ok' }, ReadAllowed: true, referenceId: 'a-reference' } })
        })

        it('throws an error if an error occurs', () => {
          const put = jest.fn(() => {
            return Promise.reject('erro thrown')
          })

          // @ts-ignore
          instance.client = { put }
          instance.getData('a-refe')

          return expect(instance.upsertData('a-reference', { someData: 'ok' })).rejects.toThrowError(
            'Error while updating data: erro thrown'
          )
        })
      })
    })
  })
})
