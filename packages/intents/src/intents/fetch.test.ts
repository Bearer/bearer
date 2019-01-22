import { FetchData, FetchActionExecutionError } from './fetch'
import * as d from '../declaration'

describe('FetchData intent', () => {
  describe('.init', () => {
    it('forwards context and params (query + body)', async () => {
      const { intent, event } = setup(PassContextIntent)
      const result = await intent(event)

      expect(result.data).toMatchObject({
        context: {
          iDefinedData: 'iDefinedDataExisting',
          authAccess: {
            apiKey: 'aKey'
          }
        },
        params: {
          firstParams: 'firstValue',
          overriden: 'thisOneWeCare'
        }
      })
    })

    it('returns a payload', async () => {
      const { intent, event } = setup(FetchIntent)

      const result = await intent(event)

      expect(result).toMatchObject({ data: ['returned-data', 'somethingFromParams', 'iDefinedDataExisting'] })
    })

    describe('when errors occur', () => {
      describe('when error is thrown within the action', () => {
        it('fails gracefully', async () => {
          const { intent, event } = setup(HardFailingIntent)

          expect(intent(event)).rejects.toEqual(new FetchActionExecutionError('sponge Bob Died'))
        })
      })

      describe('when action returns an error payload', () => {
        it('fails gracefully', async () => {
          const { intent, event } = setup(FailingIntent)

          const result = await intent(event)

          expect(result).toMatchObject({ error: '😨 No luck today' })
        })
      })
    })
  })
})

class PassContextIntent extends FetchData implements FetchData<any, any, any> {
  async action(event: d.TFetchActionEvent) {
    return { data: event }
  }
}

class FetchIntent extends FetchData implements FetchData<string[], any, d.TAPIKEYAuthContext> {
  async action(event: d.TFetchActionEvent<{ typedParam: string }, d.TAPIKEYAuthContext, { iDefinedData: string }>) {
    return { data: ['returned-data', event.params.typedParam, event.context.iDefinedData] }
  }
}

class FailingIntent extends FetchData implements FetchData<any, string, d.TAPIKEYAuthContext> {
  async action(_event: d.TFetchActionEvent) {
    return { error: '😨 No luck today' }
  }
}

class HardFailingIntent extends FetchData implements FetchData {
  async action(_event: any) {
    return await new Promise(() => {
      throw 'sponge Bob Died'
    })
  }
}

/**
 * setup
 * @param intencClass FetchIntent implementation
 */
function setup(intencClass: any) {
  const event: d.TLambdaEvent = {
    context: {
      authAccess: { apiKey: 'aKey' },
      iDefinedData: 'iDefinedDataExisting'
    },
    queryStringParameters: {
      firstParams: 'firstValue',
      overriden: 'weDontCare',
      typedParam: 'somethingFromParams'
    },
    body: {
      overriden: 'thisOneWeCare'
    }
  } as any

  const intent = intencClass.init()
  return {
    event,
    intent
  }
}
