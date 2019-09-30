import { FetchData, FetchActionExecutionError } from './fetch'
import * as d from '../declaration'

const organizationId = 'test-org-id'
const environmentId = 'test-env-id'
const internalCorrelationId = 'test-intcorr-id'
const userCorrelationId = 'test-usercorr-id'
const buid = 'test-buid'

describe('FetchData function', () => {
  describe('.init', () => {
    it('forwards context and params (query + body)', async () => {
      const { func, event } = setup(PassContextFunction)
      const result = await func(event)

      expect(result.data).toMatchObject({
        context: {
          iDefinedData: 'iDefinedDataExisting',
          auth: {
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

    it('returns a payload with StatusCode', async () => {
      const { func, event } = setup(FetchWithStatusCodeFunction)

      const result = await func(event)

      expect(result).toMatchObject({
        data: ['returned-data', 'somethingFromParams', 'iDefinedDataExisting'],
        statusCode: 201
      })
    })

    it('sets environment variables for ids', async () => {
      const { func, event } = setup(FetchFunction)

      await func(event)

      expect(process.env.clientId).toBe(organizationId)
      expect(process.env.environmentId).toBe(environmentId)
      expect(process.env.internalCorrelationId).toBe(internalCorrelationId)
      expect(process.env.scenarioUuid).toBe(buid)
      expect(process.env.userCorrelationId).toBe(userCorrelationId)
    })

    it('sets the correlation id environment variables to empty strings when the ids are undefined', async () => {
      const { func, event } = setup(FetchFunction, { internalCorrelationId: undefined, userCorrelationId: undefined })

      await func(event)

      expect(process.env.internalCorrelationId).toBe('')
      expect(process.env.userCorrelationId).toBe('')
    })

    it('puts the metadata into an environment variable', async () => {
      const { func, event } = setup(FetchFunction, { metadata: { detailedLogs: true } })

      await func(event)

      expect(process.env.bearerMetadata).toBe('{"detailedLogs":true}')
    })

    it('sets the metatadata environment variable to an empty string', async () => {
      const { func, event } = setup(FetchFunction)

      await func(event)

      expect(process.env.bearerMetadata).toBe('')
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

      describe('when action returns an error payload with a statusCode', () => {
        it('fails gracefully', async () => {
          const { func, event } = setup(FailingWithStatusCodeFunction)

          const result = await func(event)

          expect(result).toMatchObject({ statusCode: 404, error: '😨 No luck today' })
        })
      })
    })
  })

  describe('backend restrictions', () => {
    describe('when restriction exists', () => {
      it('succeeds when isBackend is truthy', async () => {
        const { func, event } = setup(RestrictedFunction)
        const result = await func({ ...event, context: { ...event.context, isBackend: true } })

        expect(result.data).toBe('success')
      })

      it('fails and returns error if isBackend is falsy', async () => {
        const { func, event } = setup(RestrictedFunction)
        const result = await func(event)

        expect(result.error).toMatchObject({
          code: 'UNAUTHORIZED_FUNCTION_CALL',
          message:
            // tslint:disable-next-line:max-line-length
            "This function can't be called from the frontend. If you want to call APIs from the frontend, please refer to this link for more information: https://docs.bearer.sh/integration-clients/javascript#calling-apis"
        })
      })
    })
  })
})

class RestrictedFunction extends FetchData implements FetchData<any, any, any> {
  static serverSideRestricted = true

  async action(_event: d.TFetchActionEvent) {
    return { data: 'success' }
  }
}

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

class FetchWithStatusCodeFunction extends FetchData implements FetchData<string[], any, d.TAPIKEYAuthContext> {
  async action(event: d.TFetchActionEvent<{ typedParam: string }, d.TAPIKEYAuthContext, { iDefinedData: string }>) {
    return { data: ['returned-data', event.params.typedParam, event.context.iDefinedData], statusCode: 201 }
  }
}

class FailingFunction extends FetchData implements FetchData<any, string, d.TAPIKEYAuthContext> {
  async action(_event: d.TFetchActionEvent) {
    return { error: '😨 No luck today' }
  }
}

class FailingWithStatusCodeFunction extends FetchData implements FetchData<any, string, d.TAPIKEYAuthContext> {
  async action(_event: d.TFetchActionEvent) {
    return { error: '😨 No luck today', statusCode: 404 }
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
function setup(functionClass: any, additionalContext = {}) {
  const event: d.TLambdaEvent = {
    context: {
      buid,
      internalCorrelationId,
      userCorrelationId,
      auth: { apiKey: 'aKey' },
      clientId: environmentId,
      iDefinedData: 'iDefinedDataExisting',
      organizationIdentifier: organizationId,
      ...additionalContext
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
