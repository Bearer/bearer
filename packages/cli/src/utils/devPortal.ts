import axios from 'axios'
import BaseCommand from '../base-command'
import { promptToLogin } from '../actions/login'

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
  async function requestFunction<DataReturned>(data: { query: string; variables?: any }) {
    let token = await command.bearerConfig.getToken()
    if (!token) {
      command.debug('no token found, trying to log you in')
      await promptToLogin(command)
      token = await command.bearerConfig.getToken()
    }

    return await instance.post<{ data?: DataReturned; errors?: any }>('', data, {
      headers: {
        Authorization: `Bearer ${token!.access_token}`
      }
    })
  }

  return {
    request: requestFunction
  }
}
