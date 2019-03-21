export class DBClient {
  static instance() {
    return new DBClient()
  }

  constructor() {}

  async getData(): Promise<any> {}

  async upsertData(): Promise<any> {}
}

export default (): DBClient => new DBClient()
