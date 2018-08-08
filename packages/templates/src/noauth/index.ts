export default {
  SaveState: `
  static action(
    _context,
    _params,
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
      } 
    })
  }
  `,
  RetrieveState: `
  static action(_context: TnoAuthContext, _params: any, state, callback: TRetrieveStateCallback) {
    callback({ state })
  }
  `,
  FetchData: `
  static action(context: TnoAuthContext, params: any, callback: TFetchDataCallback) {
    //... your code goes here
    callback({ data: []})
  }`,
  PostData: `
  static action(context: TnoAuthContext, params: any, body: any, callback: TPostDataCallback) {
    //... your code goes here
    callback({ data: []})
  }`
}
