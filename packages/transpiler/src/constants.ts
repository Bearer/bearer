export enum Component {
  bearerContext = 'bearerContext',
  componentWillLoad = 'componentWillLoad',
  componentDidUnload = 'componentDidUnload',
  componentDidLoad = 'componentDidLoad',
  setupId = 'setupId',
  scenarioId = 'scenarioId'
}

export enum Decorators {
  Component = 'Component',
  Element = 'Element',
  RootComponent = 'RootComponent',
  Prop = 'Prop',
  State = 'State',
  Watch = 'Watch',
  Event = 'Event',
  Listen = 'Listen',
  BearerState = 'BearerState',
  Intent = 'Intent',
  SaveStateIntent = 'SaveStateIntent',
  statePropName = 'statePropName',
  Input = 'Input',
  Output = 'Output'
}

export enum Module {
  BEARER_CORE_MODULE = '@bearer/core'
}

export enum Types {
  HTMLElement = 'HTMLElement',
  EventEmitter = 'EventEmitter',
  BearerFetch = 'BearerFetch',
  IntentType = 'IntentType',
  SaveState = 'SaveState'
}

export enum Properties {
  ReferenceId = 'referenceId',
  Element = 'el',
  eventName = 'eventName'
}

export enum Env {
  BEARER_SCENARIO_ID = 'BEARER_SCENARIO_ID'
}

export const BEARER = 'bearer'
