export class DBClient {
  static instance() {
    return new DBClient()
  }

  constructor() {}

  async getData(): Promise<any> {}

  async updateData(): Promise<any> {}

  async saveData(): Promise<any> {}
}

export default (): DBClient => new DBClient()
