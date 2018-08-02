import Transpiler, { TranpilerOptions } from '../../src/index'

export function TranspilerFactory(options: Partial<TranpilerOptions>): Transpiler {
  const defaults: TranpilerOptions = {
    verbose: false,
    watchFiles: false,
    buildFolder: '../../../.build/'
  }

  return new Transpiler({
    ...defaults,
    ...options
  })
}
