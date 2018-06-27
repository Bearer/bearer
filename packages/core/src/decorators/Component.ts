// Class Decorator
export function BearerComponent<T extends { new (...args: any[]): {} }>(
  constructor: T
) {
  console.warn(
    '[BEARER]',
    'Please remove BearerComponent decorator and use tranformers'
  )
  // Do not remove : will be replace in the front.
  constructor.prototype['SCENARIO_ID'] = 'BEARER_SCENARIO_ID'
  return constructor
}
