import * as ts from 'typescript'
import decoratorTransformers from './decorator-transformer'

export function Transformers() {
  return [decoratorTransformers()]
}

export default Transformers
