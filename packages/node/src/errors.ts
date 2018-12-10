export abstract class CustomError extends Error {
  constructor(message: string) {
    super(message)
    this.name = `Bearer:${this.constructor.name}`
    Error.captureStackTrace(this, this.constructor)
  }
}
