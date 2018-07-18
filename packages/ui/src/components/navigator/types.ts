export declare type TMemberRenderer<T> = {
  (member: T): any
}

export declare type TMember = {
  [key: string]: any
  _isDisabled?: boolean
}
