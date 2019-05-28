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
    return await instance.post<{ data?: DataReturned; errors?: any }>('', data, {
      headers: {
        Authorization: `Bearer ${access_token}`
      }
    })
  }

  async function access_token() {
    let token = await command.bearerConfig.getToken()

    // If no token is configured, do we have an ENV var
    // available to access the Bearer API?
    if (!token && process.env.BEARER_API_KEY) {
      return process.env.BEARER_API_KEY
    }

    // No token configured and no API key
    // Let's ask the user to log in
    if (!token) {
      command.debug('no token found, trying to log you in')
      await promptToLogin(command)
      token = await command.bearerConfig.getToken()
    }

    if (token) {
      return token.access_token
    }

    return ''
  }

  return {
    request: requestFunction
  }
}
