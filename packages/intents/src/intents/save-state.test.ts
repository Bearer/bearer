jest.mock('../db-client')

import { DBClient } from '../db-client'
import * as d from '../declaration'
import { SaveState, SaveStateActionExecutionError, SaveStateSavingStateError } from './save-state'

describe('Intents => SaveIntent', () => {
  describe('.intent', () => {
    const defaultAction = () =>
      jest.fn(() => {
        return { data: 'returned-data' }
      })

    function setup(action = defaultAction(), returnedData: any = null) {
      const getData = jest.spyOn(DBClient.prototype, 'getData')

      getData.mockImplementation(() => {
        return returnedData ? { Item: returnedData } : null
      })

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

      const intent = SaveState.intent(action as any)

      return {
        action,
        event,
        intent
      }
    }

    it('use provided referenceId', async () => {
      const { intent, event, action } = setup()

      const result = await intent({
        ...event,
        queryStringParameters: {
          ...event.queryStringParameters,
          referenceId: 'a-reference-provided'
        }
      })

      expect(action).toHaveBeenCalledWith({
        context: 'something',
        params: {
          referenceId: 'a-reference-provided',
          firstParams: 'firstValue',
          overriden: 'thisOneWeCare'
        },
        state: {}
      })

      expect(result).toMatchObject({ data: 'returned-data', meta: { referenceId: 'a-reference-provided' } })
    })

    it('generate a UUID as referenceId', async () => {
      const { intent, event } = setup()
      const result = await intent(event)

      expect(result).toMatchObject({
        meta: expect.objectContaining({
          referenceId: expect.stringMatching(/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{11}/) // uuid
        })
      })
    })

    it('pass existing data to the action', async () => {
      const { intent, event, action } = setup(defaultAction(), { existingData: 'ok' })

      await intent(event)

      expect(action).toHaveBeenCalledWith({
        context: 'something',
        params: {
          firstParams: 'firstValue',
          overriden: 'thisOneWeCare'
        },
        state: { existingData: 'ok' }
      })
    })

    it('resolve an error if action does', async () => {
      const errorFromPayloadAction = jest.fn(() => {
        return { error: 'sponge Bob Died' } as any
      })

      const { intent, event } = setup(errorFromPayloadAction)
      const result = await intent(event)

      expect(result).toMatchObject({ error: 'sponge Bob Died' })
    })

    describe('errors', () => {
      it('fails gracefully when update fails', () => {
        const { intent, event } = setup()
        const update = jest.spyOn(DBClient.prototype, 'updateData')

        update.mockImplementation(() => Promise.reject({ error: 'from-user-storage' }))

        return expect(intent(event)).rejects.toEqual(new SaveStateSavingStateError())
      })

      it('fails gracefully when action raises an error', async () => {
        const hardFailingAction = jest.fn(() => {
          throw 'sponge Bob Died error'
        })
        const { intent, event } = setup(hardFailingAction)
        return expect(intent(event)).rejects.toEqual(new SaveStateActionExecutionError())
      })
    })
  })
})
