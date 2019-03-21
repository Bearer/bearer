import { DBClient } from './db-client'
import uuid from 'uuid/v1'

export class Store {
  /**
   * TODO
   */
  private dbClient: DBClient

  constructor(signature: string) {
    this.dbClient = DBClient.instance(signature)
  }

  find = async <T extends {} = {}>(referenceId: string) => {
    if (!referenceId) {
      return null
    }
    const { Item } = await this.dbClient.getData<T>(referenceId)
    return Item
  }

  save = async <T extends {} = {}>(referenceId: string, state: T) => {
    const { Item } = await this.dbClient.upsertData<T>(referenceId || uuid(), state)
    return Item
  }
}
