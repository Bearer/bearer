export default {
  SaveState: `
  static action(
    _context: TNONEAuthContext,
    _params: any,
    body: any,
    state: any,
    callback: TSaveStateCallback
  ): void {
    const { item: { name } } = body
    const { items = [] }: any = state
    const newItem: any = { name }

    callback({
      state: {
        ...state,
        items: [...items, newItem]
      },
      data: [...items, newItem]
    })
  }
  `,
  FetchData: `
  static action(context: TNONEAuthContext, params: any, body: any, callback: TFetchDataCallback) {
    //... your code goes here
    callback({ data: []})
  }`
}
