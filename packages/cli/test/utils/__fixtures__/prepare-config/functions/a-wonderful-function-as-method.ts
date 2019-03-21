import { FetchData } from '@bearer/functions'

export default class MyFunction extends FetchData implements FetchData {
  async action(_event: any) {
    return { data: 'something', referenceId: 'something' }
  }
}
