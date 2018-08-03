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

export interface IBearerStateDecorator<T> {
  (options?: IBearerStateDecoratorOptions): T
}

export const BearerState: IBearerStateDecorator<any> = (options?: IBearerStateDecoratorOptions) => (
  target: any,
  key: string
): void => {}

/**
 *  Component Decorator
 */

export interface IBearerComponentDecorator<T> {
  (options?: d.ComponentOptions): T
}

export declare const Component: IBearerComponentDecorator<any>

/**
 * RootComponent Decorator
 */

export enum BearerRootComponentRoleEnum {
  Display = 'display',
  Action = 'action'
}

export interface BearerRootComponentOptions extends Omit<d.ComponentOptions, 'tag'> {
  group: string
  name: 'display' | 'action'
  shadow?: boolean
}

export interface IBearerRootComponentDecorator<T> {
  (options?: BearerRootComponentOptions): T
}

export declare const RootComponent: IBearerRootComponentDecorator<any>

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
