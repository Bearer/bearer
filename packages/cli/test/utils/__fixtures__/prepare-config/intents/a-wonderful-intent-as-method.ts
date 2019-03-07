import { SaveState } from '@bearer/intents'

export default class MyFunction extends SaveState implements SaveState {
  async action(_event: any) {
    return { data: 'something', state: 'something' }
  }
}
