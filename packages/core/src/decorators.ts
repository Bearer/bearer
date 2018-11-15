import * as d from '@stencil/core/dist/declarations/index'

export * from './decorators/Intent'

export {
  ComponentDidLoad,
  ComponentDidUnload,
  ComponentDidUpdate,
  ComponentWillLoad,
  ComponentWillUpdate,
  Config,
  EventEmitter,
  EventListenerEnable,
  FunctionalComponent
} from '@stencil/core/dist/declarations'

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

/**
 *  BearerState Decorator
 */

export interface IBearerStateDecoratorOptions {
  statePropName?: string
}

export type IBearerStateDecorator<T> = (options?: IBearerStateDecoratorOptions) => T

export const BearerState: IBearerStateDecorator<any> = (_options?: IBearerStateDecoratorOptions) => (
  _target: any,
  _key: string
): void => {}

/**
 *  Component Decorator
 */

export type IBearerComponentDecorator<T> = (options?: d.ComponentOptions) => T

export declare const Component: IBearerComponentDecorator<any>

/**
 * RootComponent Decorator
 */

export type BearerRootComponentRole = 'display' | 'action'

export interface BearerRootComponentOptions extends Omit<d.ComponentOptions, 'tag'> {
  group: string
  role: BearerRootComponentRole
  shadow?: boolean
}

export type IBearerRootComponentDecorator<T> = (options?: BearerRootComponentOptions) => T

export declare const RootComponent: IBearerRootComponentDecorator<any>

/**
 * Input Decorator
 */
type TInputDecoratorOptions = {
  group?: string // target a different group of component to listen to eventName
  propName?: string // specifiy a different attribut name to your component
  eventName?: string // listen to a different event name from group
  intentName?: string // specify an intent to use to retrieve data
  autoLoad?: boolean // auto load data when componentDidLoad
}

type TInputDecorator = (options?: TInputDecoratorOptions) => (target: any, key: string) => void

export declare const Input: TInputDecorator

/**
 * Output Decorator
 */
type TOutputDecoratorOptions = {
  intentName?: string // Intent you want to use to save data: must be a SaveState
  intentPropertyName?: string // name used to send data in the body to the intent
  eventName?: string // event triggered when the data is saved
  propertyWatchedName?: string
  referenceKeyName?: string // name of the reference sent within the event: default referenceId
}

type TOutputDecorator = (options?: TOutputDecoratorOptions) => (target: any, key: string) => void

export declare const Output: TOutputDecorator

export type BearerRef<T> = T

/**
 * Build
 */
export declare const Build: d.UserBuildConditionals
// /**
//  * Component
//  */
// export declare const Component: d.ComponentDecorator;
/**
 * Element
 */
export declare const Element: d.ElementDecorator
/**
 * Event
 */
export declare const Event: d.EventDecorator
/**
 * Listen
 */
export declare const Listen: d.ListenDecorator
/**
 * Method
 */
export declare const Method: d.MethodDecorator
/**
 * Prop
 */
export declare const Prop: d.PropDecorator
/**
 * State
 */
export declare const State: d.StateDecorator
/**
 * Watch
 */
export declare const Watch: d.WatchDecorator
/**
 * Deprecated: Please use @Watch decorator instead
 */
export declare const PropWillChange: d.WatchDecorator
/**
 * Deprecated: Please use @Watch decorator instead
 */
export declare const PropDidChange: d.WatchDecorator
export interface HostElement extends HTMLElement {}
