import { SaveState } from '@bearer/functions'

export default class MyFunction extends SaveState implements SaveState {
  action = async (_event: any) => {
    return { data: 'something', state: 'something' }
  }
}
