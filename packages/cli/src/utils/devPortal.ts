import axios from 'axios'
import BaseCommand from '../base-command'

export const devPortalClient = (command: BaseCommand) => {
  const instance = axios.create({
    baseURL: command.constants.DeveloperPortalAPIUrl
  })
  instance.interceptors.response.use(
    r => r,
    error => {
      if (error.response && error.response.status) {
        command.error('Unauthorized action, please run bearer login first')
      }
      return Promise.reject(error)
    }
  )

  return {
    request: async <DataReturned>(data: { query: string; variables?: any }) => {
      const token = await command.bearerConfig.getToken()
      if (token) {
        return instance.post<{ data?: DataReturned }>('', data, {
          headers: {
            Authorization: `Bearer ${token.access_token}`
          }
        })
      }
      throw new UnauthorizedError()
    }
  }
}

export class UnauthorizedError extends Error {
  constructor() {
    super('Unauthorized request')
  }
}
