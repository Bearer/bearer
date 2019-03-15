export enum Component {
  bearerContext = 'bearerContext',
  componentWillLoad = 'componentWillLoad',
  componentDidUnload = 'componentDidUnload',
  componentDidLoad = 'componentDidLoad',
  setupId = 'setupId',
  integrationId = 'integrationId',
  integrationIdAccessor = 'INTEGRATION_ID'
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
  Function = 'Function',
  BackendFunction = '_BackendFunction',
  SaveStateFunction = 'SaveStateFunction',
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
  FunctionType = 'FunctionType',
  SaveState = 'SaveState'
}

export enum Properties {
  ReferenceId = 'referenceId',
  Element = 'el',
  eventName = 'eventName'
}

export enum Env {
  BEARER_INTEGRATION_ID = 'BEARER_INTEGRATION_ID'
}

export const BEARER = 'bearer'

export const SETUP_ID = 'setupId'
export const SETUP = 'setup'
