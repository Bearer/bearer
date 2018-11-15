import { TInputDecorator, TOutputDecorator } from '@bearer/types/lib/input-output-decorators'
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

export declare const Input: TInputDecorator

/**
 * Output Decorator
 */

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
