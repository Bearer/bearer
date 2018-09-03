import axios, { AxiosInstance } from 'axios'

import BaseCommand from '../BaseCommand'

export class ScenarioClient {
  private client: AxiosInstance

  constructor(baseUrl: string, authorization?: string, version?: string) {
    const headers = {
      Authorization: authorization,
      ['BEARER-CLI-VERSION']: version
    }
    this.client = axios.create({
      baseURL: baseUrl,
      headers
    })
  }

  async getScenarioArchiveUploadUrl(orgId: string, scenarioId: string): Promise<string> {
    try {
      const response = await this.client.post('scenario-archive-url', { orgId, scenarioId })
      return response.data.url
    } catch (e) {
      throw e
    }
  }
}

export default (command: BaseCommand): ScenarioClient => {
  return new ScenarioClient(
    command.bearerConfig.IntegrationServiceUrl,
    command.bearerConfig.bearerConfig.authorization.AuthenticationResult!.IdToken,
    command.config.version
  )
}
