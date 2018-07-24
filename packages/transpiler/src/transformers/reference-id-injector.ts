import * as ts from 'typescript'
import { hasPropDecoratedWithName } from './decorator-helpers'
import { ensurePropImported, propDecorator } from './bearer'
import { Decorators } from './constants'

type TransformerOptions = {
  verbose?: true
}

export default function BearerReferenceIdInjector({ verbose }: TransformerOptions = {}): ts.TransformerFactory<
  ts.SourceFile
> {
  return transformContext => {
    return tsSourceFile => {
      if (!hasRetrieveOrSaveStateIntent(tsSourceFile)) {
        return tsSourceFile
      }

      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isClassDeclaration(node)) {
          return injectBearerReferenceIdProp(node)
        }
        return ts.visitEachChild(node, visit, transformContext)
      }
      return ts.visitEachChild(ensurePropImported(tsSourceFile), visit, transformContext)
    }
  }
}

function injectBearerReferenceIdProp(tsClass: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    tsClass,
    tsClass.decorators,
    tsClass.modifiers,
    tsClass.name,
    tsClass.typeParameters,
    tsClass.heritageClauses,
    [
      ...tsClass.members,
      ts.createProperty(
        [propDecorator()],
        undefined,
        'referenceId',
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
        undefined
      )
    ]
  )
}
function hasRetrieveOrSaveStateIntent(tsSourceFile: ts.SourceFile): boolean {
  return ts.forEachChild(
    tsSourceFile,
    (node): boolean => {
      if (ts.isClassDeclaration(node)) {
        return (
          hasPropDecoratedWithName(node, Decorators.RetrieveStateIntent) ||
          hasPropDecoratedWithName(node, Decorators.SaveStateIntent)
        )
      }
      return false
    }
  )
}
