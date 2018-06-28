import { createStore } from '@bearer/core/dist/state'
import { State, Action } from '../shared'

export enum ActionTypes {
  STATE_RECEIVED = 'STATE_RECEIVED',
  HELLO_WORLD_SELECTED = 'HELLO_WORLD_SELECTED',
  HELLO_WORLD_RECEIVED = 'HELLO_WORLD_RECEIVED',
  HELLO_WORLD_DETACHED = 'HELLO_WORLD_DETACHED'
}

// reducers/scenario.js
const initialState: State = {
  helloworld: []
}

/* Reducers */
const scenario = (state: State = initialState, { type, payload }: Action) => {
  switch (type) {
    case ActionTypes.HELLO_WORLD_SELECTED: {
      return {
        ...state,
        helloworld: [
          ...state.helloworld,
          payload['helloworld']
        ]
      }
    }

    case ActionTypes.HELLO_WORLD_DETACHED: {
      return {
        ...state,
        helloworld: [
          ...state.helloworld.filter(
            hw => hw.title !== payload['helloworld']['title']
          )
        ]
      }
    }

    case ActionTypes.STATE_RECEIVED: {
      return {
        ...state,
        helloworld: payload['helloworld']
      }
    }

    default: {
      return state
    }
  }
}

// end reducers/scenario.js

/* Store */
const configStore = () => {
  return createStore(
    scenario,
    undefined,
    window['__REDUX_DEVTOOLS_EXTENSION__'] &&
      window['__REDUX_DEVTOOLS_EXTENSION__']({
        instanceId: 'BEARER_SCENARIO_ID'
      })
  )
}

export { configStore }
