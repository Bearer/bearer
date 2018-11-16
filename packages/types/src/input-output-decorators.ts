/**
 * Input Decorator
 */
export type TInputDecoratorOptions = {
  group: string // target a different group of component to listen to eventName
  propName: string // specifiy a different attribut name to your component
  eventName: string // listen to a different event name from group
  intentReferenceIdKeyName: string // key name used to send referenceId to the intent
  intentName: string // specify an intent to use to retrieve data
  autoLoad: boolean // auto load data when componentDidLoad
}

export type TInputDecorator = (options?: Partial<TInputDecoratorOptions>) => (target: any, key: string) => void

/**
 * Output Decorator
 */
export type TOutputDecoratorOptions = {
  intentName: string // Intent you want to use to save data: must be a SaveState
  intentReferenceIdKeyName: string // key name used to send referenceId to the intent
  intentPropertyName: string // name used to send data in the body to the intent
  eventName: string // event triggered when the data is saved
  propertyWatchedName: string
  referenceKeyName: string // name of the reference sent within the event: default referenceId
}

export type TOutputDecorator = (options?: Partial<TOutputDecoratorOptions>) => (target: any, key: string) => void
