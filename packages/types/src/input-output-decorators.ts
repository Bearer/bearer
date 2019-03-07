/**
 * Input Decorator
 */
export type TInputDecoratorOptions = {
  /**
   * Name of the property for the component
   * default: [propName]RefId
   */
  propertyReferenceIdName: string
  /**
   * Change the event you want to listen from
   * default: [propName]Saved
   */
  eventName: string
  /**
   * Name of the key used to send referenceId to the func (as query parameter)
   */
  functionReferenceIdKeyName: string
  /**
   * Function to retrieve data from
   */
  functionName: string
  /**
   * Function arguments
   */
  functionArguments: string[]
  /**
   * Autoload data when component mounted
   * default: true
   */
  autoLoad: boolean
}

export type TInputDecorator = (options?: Partial<TInputDecoratorOptions>) => (target: any, key: string) => void

/**
 * Output Decorator options
 */
export type TOutputDecoratorOptions = {
  /**
   * Specify an func you would like to use instead of the default one
   */
  functionName: string
  /**
   * Name of the key used to send the referenceId value to the function
   * default : referenceId
   */
  functionReferenceIdKeyName: string
  /**
   * Name of property used to populate the func functionReferenceIdKeyName param
   */
  functionReferenceIdValue?: string
  /**
   * Function arguments
   */
  functionArguments: string[]
  /**
   * Name of the key name used to send the data to the func (in the body)
   * default: [propertyName]
   */
  functionPropertyName: string
  /**
   * Event name triggered when the data has been processed by the function
   * default: [propertyName]Saved
   */
  eventName: string
  /**
   * State or Prop you want to watch. Let you save data when a different property change:
   * @example
   * ```typescript
   *  @State() otherData: string = "whatever"
   *  @Ouput({ propertyWatchedName: 'otherData' }) data: string
   * ```
   * event will be emitted only when otherData will be updated.
   *
   * default: [propertyName]
   *
   */
  propertyWatchedName: string
  /**
   * Name of the key used to send referenceId as property of the event.detail
   * default: referenceId
   */
  referenceKeyName: string

  /**
   * Autoload data when component mounted
   * default: true
   */
  autoLoad: boolean
}

/**
 * Ouput decorator
 * @example
 * ```typescript
 * @Ouput() goat: Goat
 * ```
 */

export type TOutputDecorator = (options?: Partial<TOutputDecoratorOptions>) => (target: any, key: string) => void
