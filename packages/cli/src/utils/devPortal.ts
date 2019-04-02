import axios from 'axios'
import BaseCommand from '../base-command'

export const devPortalClient = (command: BaseCommand) => {
  const instance = axios.create({
    baseURL: command.constants.DeveloperPortalAPIUrl
  })
  return {
    request: async <DataReturned>(data: { query: string; variables?: any }) => {
      const token = await command.bearerConfig.getToken()
      if (token) {
        return instance.post<{ data?: DataReturned }>('', data, {
          headers: {
            Authorization: `Bearer ${token.id_token}`
          }
        })
      }
      throw 'NEED_LOGIN'
    }
  }
}
