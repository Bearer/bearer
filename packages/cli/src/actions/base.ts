import baseCommand from '../base-command'

export function createExport<T extends any[] = any[]>(
  klass: new (logger: baseCommand) => { run: (...args: T) => Promise<void> }
) {
  return async (logger: baseCommand, ...args: T) => await new klass(logger).run(...args)
}

export default class BaseAction {
  constructor(readonly logger: baseCommand) {}
}
