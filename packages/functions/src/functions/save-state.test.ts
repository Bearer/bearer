jest.mock('../db-client')

import { DBClient } from '../db-client'
import * as d from '../declaration'
import { SaveState, SaveStateActionExecutionError, SaveStateSavingStateError } from './save-state'

describe.skip('SaveFunction function', () => {
  const update = jest.spyOn(DBClient.prototype, 'upsertData')
  beforeEach(() => {
    update.mockReset()
  })

  describe('.init', () => {
    it('use provided referenceId', async () => {
      const { func, event } = setup(SuccessFunction)

      const result = await func({
        ...event,
        queryStringParameters: {
          ...event.queryStringParameters,
          referenceId: 'a-reference-provided'
        }
      })

      expect(result.meta).toMatchObject({ referenceId: 'a-reference-provided' })
    })

    it('generate a UUID as referenceId', async () => {
      const { func, event } = setup(SuccessFunction)
      const result = await func(event)

      expect(result).toMatchObject({
        meta: expect.objectContaining({
          referenceId: expect.stringMatching(/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{11}/) // uuid
        })
      })
    })

    it('pass existing data to the action', async () => {
      const { func, event } = setup(SuccessFunction, { existingData: 'ok' })

      const result = await func(event)

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
          const { func, event } = setup(SuccessFunction)
          update.mockImplementation(() => Promise.reject({ error: 'from-user-storage' }))

          return expect(func(event)).rejects.toEqual(new SaveStateSavingStateError())
        })
      })

      describe('when error is thrown within the action', () => {
        it('fails gracefully', async () => {
          const { func, event } = setup(HardFailingFunction)

          return expect(func(event)).rejects.toEqual(new SaveStateActionExecutionError())
        })
      })

      describe('when error is returned by the action', () => {
        it('fails gracefully', async () => {
          const { func, event } = setup(FunctionWithErrorReturned)

          const result = await func(event)

          expect(result).toMatchObject({ error: 'sponge Bob Died smoothly' })
        })
      })
    })
  })

  describe('backend restrictions', () => {
    describe('when restriction exists', () => {
      it('succeeds when isBackend is truthy', async () => {
        const { func, event } = setup(RestrictedFunction)
        const result = await func({
          ...event,
          context: { ...event.context, isBackend: true }
        })

        expect(result.data).toBe('success')
      })

      it('fails and returns error if isBackend is falsy', async () => {
        const { func, event } = setup(RestrictedFunction)
        const result = await func(event)

        expect(result.error).toMatchObject({
          code: 'UNAUTHORIZED_FUNCTION_CALL',
          message: "This function can't be called"
        })
      })
    })
  })
})

class RestrictedFunction extends SaveState implements SaveState<{ savedData: string[] }> {
  static backendOnly = true
  async action(event: d.TSaveActionEvent<{ savedData: string[] }>) {
    // proces typing works
    const existingData = event.state.savedData || []
    return { data: 'success', state: { savedData: ['ok', ...existingData] } }
  }
}

class SuccessFunction extends SaveState implements SaveState<{ savedData: string[] }> {
  async action(event: d.TSaveActionEvent<{ savedData: string[] }>) {
    // proces typing works
    const existingData = event.state.savedData || []
    return { data: event, state: { savedData: ['ok', ...existingData] } }
  }
}

class FunctionWithErrorReturned extends SaveState implements SaveState {
  async action(_event: any): Promise<any> {
    return { error: 'sponge Bob Died smoothly' }
  }
}

class HardFailingFunction extends SaveState implements SaveState {
  async action(_event: any): Promise<any> {
    return await new Promise(() => {
      throw 'sponge Bob Died'
    })
  }
}

/**
 * setup
 * @param functionClass
 * @param returnedData
 */
function setup(functionClass: any, returnedData: any = null) {
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

  const func = functionClass.init()

  return {
    event,
    func
  }
}
