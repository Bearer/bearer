declare module 'post-robot' {
  // Warning: This is not actually a Promise, but the interface is the same.
  type ZalgoPromise<T> = Promise<T>

  // For our purposes, Window is cross domain enough. For now at least.
  type CrossDomainWindowType = Window
  type WindowResolverType = CrossDomainWindowType | string | HTMLIFrameElement

  // Client
  // Loosely based on: https://github.com/krakenjs/post-robot/blob/master/src/public/client.js
  type RequestOptionsType = {
    window?: WindowResolverType
    domain?: string | Array<string> | RegExp
    name?: string
    data?: Object
    fireAndForget?: boolean
    timeout?: number
  }
  type ResponseMessageEvent = {
    source: CrossDomainWindowType
    origin: string
    data: any
  }

  function send(
    window: WindowResolverType,
    name: string,
    data?: Object,
    options?: any
  ): ZalgoPromise<ResponseMessageEvent>
  type Sendable = { send: (name: string, data?: Object) => ZalgoPromise<ResponseMessageEvent> }
  function client(options?: RequestOptionsType): Sendable

  // Server
  // Loosely based on: https://github.com/krakenjs/post-robot/blob/master/src/public/server.js
  type ErrorHandlerType = (err: any) => void
  type HandlerType = ({ source: Window, origin: string, data: Object }) => void | any | ZalgoPromise<any>

  type ServerOptionsType = {
    handler?: HandlerType
    errorHandler?: ErrorHandlerType
    window?: CrossDomainWindowType
    name?: string
    domain?: string | RegExp | Array<string>
    once?: boolean
    errorOnClose?: boolean
  }

  type Cancellable = { cancel: () => void }
  function listener(options?: ServerOptionsType): { on: (name: string, handler: HandlerType) => Cancellable }
  function listen(options: ServerOptionsType): Cancellable
  function on(name: string, options: ServerOptionsType | HandlerType, handler?: HandlerType): Cancellable
}
