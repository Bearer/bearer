import { SaveState } from '@bearer/functions'

export default class MyFunction extends SaveState implements SaveState {
  async action(_event: any) {
    return { data: 'something', state: 'something' }
  }
}
