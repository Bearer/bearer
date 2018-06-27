import { SaveState } from './intents'
import { HelloWorld, ScenarioState, SavedData } from '../shared'

export default class SaveStateIntent {
  static get intentName(): string {
    return 'saveState'
  }
  static get intentType() {
    return SaveState
  }

  static action(
    _token,
    _params,
    body: HelloWorld,
    state,
    callback: (ScenarioState) => void
  ): void {
    const {
      title,
      lang
    } = body
    const { helloworld = [] }: ScenarioState = state
    const newHelloWorld: SavedData = { title, lang }

    callback({
      ...state,
      helloworld: [...helloworld, newHelloWorld]
    })
  }
}
