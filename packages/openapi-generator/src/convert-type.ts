import ts from 'typescript'
import merge from 'lodash.merge'

function isPrimitive(flags: number): boolean {
  return (
    (flags & ts.TypeFlags.String) !== 0 || (flags & ts.TypeFlags.Number) !== 0 || (flags & ts.TypeFlags.Boolean) !== 0
  )
}

export function convertType(
  typ: ts.Type | ts.UnionOrIntersectionType,
  checker: ts.TypeChecker,
  definition: any = {}
): any {
  const flags = typ.getFlags()
  const typeNode = checker.typeToTypeNode(typ)!
  if (ts.isArrayTypeNode(typeNode)) {
    let elementType = (checker as any).getElementTypeOfArrayType(typ)
    const elementTypeFlags = elementType.getFlags()
    if (!isPrimitive(elementTypeFlags)) {
      elementType = checker.getApparentType(elementType)
    }
    return { type: 'array', items: convertType(elementType, checker) }
  }
  if (flags & ts.TypeFlags.String) {
    return { type: 'string' }
  }
  if (flags & ts.TypeFlags.Number) {
    return { type: 'number' }
  }
  if (flags & ts.TypeFlags.Boolean) {
    return { type: 'boolean' }
  }
  if (flags & ts.TypeFlags.Intersection || flags & ts.TypeFlags.Union) {
    return (typ as ts.IntersectionType).types.reduce((acc, t) => {
      const converted = convertType(t, checker, acc)
      return merge(acc, converted)
    }, {})
  }

  if (flags & ts.TypeFlags.Object) {
    return {
      type: 'object',
      properties: checker
        .getWidenedType(typ)
        .getProperties()
        .reduce((acc: any, s): any => {
          const nestedType = checker.getTypeOfSymbolAtLocation(s, s.valueDeclaration!)
          acc[checker.symbolToString(s)] = convertType(checker.getWidenedType(nestedType), checker)
          return acc
        }, {})
    }
  }
  if (typ.aliasTypeArguments && typ.aliasTypeArguments.length) {
    return typ.aliasTypeArguments.reduce((acc, t) => {
      acc[checker.symbolToString(t.symbol)] = convertType(t, checker)
      return acc
    }, definition)
  }

  if (typ.aliasSymbol) {
    const aliasedType = checker.getTypeOfSymbolAtLocation(typ.aliasSymbol, typ.aliasSymbol.valueDeclaration!)
    return convertType(aliasedType, checker, definition)
  }
  return definition
}
