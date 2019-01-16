import { FetchData, FetchActionExecutionError } from './fetch'
import * as d from '../declaration'

describe('Intents', () => {
  describe('Fetch', () => {
    it('export a class', () => {
      expect(FetchData).toBeTruthy()
    })

    it('execute promise', done => {
      const event = {
        queryStringParameters: {},
        context: {
          bearerBaseURL: 'none',
          authAccess: {} as any
        }
      }
      const context = {}
      const callback = jest.fn((_error, payload) => {
        expect(payload).toEqual({ data: { ko: 'ko' } })
        done()
      })

      const action = (_context, _params, _body, callback) => {
        callback({ data: { ko: 'ko' } })
      }

      FetchData.intent(action)(event, context, callback)
    })
  })

  describe('Async fetch', () => {
    const defaultAction = () =>
      jest.fn(() => {
        return { data: 'returned-data' }
      })

    function setup(action = defaultAction()) {
      const event: d.TLambdaEvent = {
        context: 'something',
        queryStringParameters: {
          firstParams: 'firstValue',
          overriden: 'weDontCare'
        },
        body: {
          overriden: 'thisOneWeCare'
        }
      } as any

      const intent = FetchData.intentPromise(action as any)
      return {
        action,
        event,
        intent
      }
    }

    it('forward context and params (query + body)', async () => {
      const { intent, event, action } = setup()
      await intent(event)

      expect(action).toHaveBeenCalledWith({
        context: 'something',
        params: {
          firstParams: 'firstValue',
          overriden: 'thisOneWeCare'
        }
      })
    })

    it('return a payload', async () => {
      const { intent, event } = setup()

      const result = await intent(event)

      expect(result).toMatchObject({ data: 'returned-data' })
    })

    describe('error handling', () => {
      it('fails gracefully when action throw an error', async () => {
        const hardFailingAction = jest.fn(() => {
          throw 'sponge Bob Died'
        })
        const { intent, event } = setup(hardFailingAction)
        return expect(intent(event)).rejects.toEqual(new FetchActionExecutionError('sponge Bob Died'))
      })

      it('fails gracefully when action return error payload', async () => {
        const errorFromPayloadAction = jest.fn(() => {
          return { error: '😨 No luck today' } as any
        })
        const { intent, event } = setup(errorFromPayloadAction)

        const result = await intent(event)
        return expect(result).toMatchObject({ error: '😨 No luck today' })
      })
    })
  })
})
