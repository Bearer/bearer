type THeader = {
  openapi: string
  info: {
    title: string
  }
  servers: { url: string }[]
  tags: {}[]
  components: {
    securitySchemes: {
      [key: string]: {
        type: string
        in: string
        description: string
        name: string
      }
    }
  }
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
        parameters: [oauthParam(oauth)].filter(e => e !== undefined),
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
        },
        security: [
          {
            'API Key': []
          }
        ]
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
    tags: [],
    components: {
      securitySchemes: {
        'API Key': {
          type: 'apiKey',
          in: 'header',
          name: 'Authorization',
          description: 'You can find your **Bearer API Key** [here](https://app.bearer.sh/keys).'
        }
      }
    }
  }
}
