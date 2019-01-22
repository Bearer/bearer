jest.mock('../db-client')

import { DBClient } from '../db-client'
import * as d from '../declaration'
import { SaveState, SaveStateActionExecutionError, SaveStateSavingStateError } from './save-state'

describe('SaveIntent intent', () => {
  describe('.init', () => {
    it('use provided referenceId', async () => {
      const { intent, event } = setup(SuccessIntent)

      const result = await intent({
        ...event,
        queryStringParameters: {
          ...event.queryStringParameters,
          referenceId: 'a-reference-provided'
        }
      })

      expect(result.meta).toMatchObject({ referenceId: 'a-reference-provided' })
    })

    it('generate a UUID as referenceId', async () => {
      const { intent, event } = setup(SuccessIntent)
      const result = await intent(event)

      expect(result).toMatchObject({
        meta: expect.objectContaining({
          referenceId: expect.stringMatching(/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{11}/) // uuid
        })
      })
    })

    it('pass existing data to the action', async () => {
      const { intent, event } = setup(SuccessIntent, { existingData: 'ok' })

      const result = await intent(event)

      expect(result.data).toMatchObject({
        context: 'something',
        params: {
          firstParams: 'firstValue',
          overriden: 'thisOneWeCare'
        },
        state: { existingData: 'ok' }
      })
    })

    describe('when errors occur', () => {
      describe('when update fails', () => {
        it('fails gracefully', () => {
          const { intent, event } = setup(SuccessIntent)
          const update = jest.spyOn(DBClient.prototype, 'updateData')

          update.mockImplementation(() => Promise.reject({ error: 'from-user-storage' }))

          return expect(intent(event)).rejects.toEqual(new SaveStateSavingStateError())
        })
      })

      describe('when error is thrown within the action', () => {
        it('fails gracefully', async () => {
          const { intent, event } = setup(HardFailingIntent)

          return expect(intent(event)).rejects.toEqual(new SaveStateActionExecutionError())
        })
      })

      describe('when error is returned by the action', () => {
        it('fails gracefully', async () => {
          const { intent, event } = setup(IntentWithErrorReturned)

          const result = await intent(event)

          expect(result).toMatchObject({ error: 'sponge Bob Died smoothly' })
        })
      })
    })
  })
})

class SuccessIntent extends SaveState implements SaveState<{ savedData: string[] }> {
  async action(event: d.TSaveActionEvent<{ savedData: string[] }>) {
    // proces typing works
    const existingData = event.state.savedData || []
    return { data: event, state: { savedData: ['ok', ...existingData] } }
  }
}

class IntentWithErrorReturned extends SaveState implements SaveState {
  async action(_event: any): Promise<any> {
    return { error: 'sponge Bob Died smoothly' }
  }
}

class HardFailingIntent extends SaveState implements SaveState {
  async action(_event: any): Promise<any> {
    return await new Promise(() => {
      throw 'sponge Bob Died'
    })
  }
}

/**
 * setup
 * @param intentClass
 * @param returnedData
 */
function setup(intentClass: any, returnedData: any = null) {
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

  const intent = intentClass.init()

  return {
    event,
    intent
  }
}
