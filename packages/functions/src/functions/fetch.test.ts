import { FetchData, FetchActionExecutionError } from './fetch'
import * as d from '../declaration'

describe('FetchData function', () => {
  describe('.init', () => {
    it('forwards context and params (query + body)', async () => {
      const { func, event } = setup(PassContextFunction)
      const result = await func(event)

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
      const { func, event } = setup(FetchFunction)

      const result = await func(event)

      expect(result).toMatchObject({ data: ['returned-data', 'somethingFromParams', 'iDefinedDataExisting'] })
    })

    describe('when errors occur', () => {
      describe('when error is thrown within the action', () => {
        it('fails gracefully', async () => {
          const { func, event } = setup(HardFailingFunction)

          expect(func(event)).rejects.toEqual(new FetchActionExecutionError('sponge Bob Died'))
        })
      })

      describe('when action returns an error payload', () => {
        it('fails gracefully', async () => {
          const { func, event } = setup(FailingFunction)

          const result = await func(event)

          expect(result).toMatchObject({ error: '😨 No luck today' })
        })
      })
    })
  })
})

class PassContextFunction extends FetchData implements FetchData<any, any, any> {
  async action(event: d.TFetchActionEvent) {
    return { data: event }
  }
}

class FetchFunction extends FetchData implements FetchData<string[], any, d.TAPIKEYAuthContext> {
  async action(event: d.TFetchActionEvent<{ typedParam: string }, d.TAPIKEYAuthContext, { iDefinedData: string }>) {
    return { data: ['returned-data', event.params.typedParam, event.context.iDefinedData] }
  }
}

class FailingFunction extends FetchData implements FetchData<any, string, d.TAPIKEYAuthContext> {
  async action(_event: d.TFetchActionEvent) {
    return { error: '😨 No luck today' }
  }
}

class HardFailingFunction extends FetchData implements FetchData {
  async action(_event: any) {
    return await new Promise(() => {
      throw 'sponge Bob Died'
    })
  }
}

/**
 * setup
 * @param functionClass FetchFunction implementation
 */
function setup(functionClass: any) {
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

  const func = functionClass.init()
  return {
    event,
    func
  }
}
