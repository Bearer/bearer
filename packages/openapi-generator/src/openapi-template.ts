type THeader = {
  openapi: string
  info: {
    description: string
    version: string
    title: string
    contact: { email: string }
    license: { name: string }
  }
  servers: { url: string }[]
  tags: {
    name: string
    description: string
    externalDocs: { description: string; url: string }
  }[]
}

type TParam = {
  name: string
  schema: { type: string; format: string }
  in: string
  description: string
  required?: boolean
}

export function specPath({
  integrationUuid,
  intentName,
  requestBody,
  response,
  oauth2
}: {
  integrationUuid: string
  intentName: string
  requestBody: any
  response: any
  oauth2: boolean
}) {
  return {
    [`/${integrationUuid}/${intentName}`]: {
      post: {
        parameters: [
          {
            in: 'header',
            name: 'authorization',
            schema: { type: 'string' },
            description: 'API Key',
            required: true
          },
          oauthParam(oauth2)
        ].filter(e => e !== undefined),
        summary: intentName,
        requestBody: { content: { 'application/json': { schema: requestBody } } },
        responses: {
          '200': { description: intentName, content: { 'application/json': { schema: response } } },
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

function oauthParam(oauth2: boolean): TParam | undefined {
  if (oauth2) {
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
      description: `openapi definition file for ${integrationName}`,
      version: '0.0.1',
      title: integrationName,
      contact: { email: 'bearer@bearer.sh' },
      license: { name: 'MIT' }
    },
    servers: [{ url: 'https://int.bearer.sh/api/v2/intents/backend/' }],
    tags: [
      {
        name: 'integration',
        description: `List of endpoints providing backend to backend integration with ${integrationName}`,
        externalDocs: { description: 'Find out more', url: 'https://www.bearer.sh' }
      }
    ]
  }
}
