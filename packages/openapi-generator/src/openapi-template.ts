type THeader = {
  openapi: string
  info: {
    title: string
  }
  servers: { url: string }[]
  tags: {}[]
}

type TParam = {
  name: string
  schema: { type: string; format: string }
  in: string
  description: string
  required?: boolean
}

export function specPath({
  buid,
  functionName,
  requestBody,
  response,
  oauth
}: {
  buid: string
  functionName: string
  requestBody: any
  response: any
  oauth: boolean
}) {
  return {
    [`/${buid}/${functionName}`]: {
      post: {
        parameters: [
          {
            in: 'header',
            name: 'authorization',
            schema: { type: 'string' },
            description: 'API Key',
            required: true
          },
          oauthParam(oauth)
        ].filter(e => e !== undefined),
        summary: functionName,
        requestBody: { content: { 'application/json': { schema: requestBody } } },
        responses: {
          '200': { description: functionName, content: { 'application/json': { schema: response } } },
          '401': {
            description: 'Access forbidden',
            content: {
              'application/json': { schema: { type: 'object', properties: { error: { type: 'string' } } } }
            }
          },
          '403': {
            description: 'Unauthorized',
            content: {
              'application/json': { schema: { type: 'object', properties: { error: { type: 'string' } } } }
            }
          }
        }
      }
    }
  }
}

function oauthParam(oauth: boolean): TParam | undefined {
  if (oauth) {
    return {
      name: 'authId',
      schema: { type: 'string', format: 'uuid' },
      in: 'query',
      description: 'User Identifier',
      required: true
    }
  }
  return
}

export function topOfSpec(integrationName: string): THeader {
  return {
    openapi: '3.0.0',
    info: {
      title: integrationName
    },
    servers: [{ url: 'https://int.bearer.sh/api/v4/functions/backend/' }],
    tags: []
  }
}
