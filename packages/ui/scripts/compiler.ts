import * as fs from 'fs-extra'
import * as ts from 'typescript'

export default (files: Array<string>) => {
  files.map(
    file => {
      const sourceFile = ts.createSourceFile(
        file,
        fs.readFileSync(file, 'utf8'),
        ts.ScriptTarget.Latest,
        false,
        ts.ScriptKind.TSX
      )

      const printer = ts.createPrinter({
        newLine: ts.NewLineKind.LineFeed
      })
      const result: ts.TransformationResult<ts.SourceFile> = ts.transform(sourceFile, [transformer])

      const transformedSourceFile: ts.SourceFile = result.transformed[0]
      const statements = transformedSourceFile.statements
      fs.writeFileSync(file, printer.printList(ts.ListFormat.MultiLine, statements, sourceFile))
    }
    //   let outPath = tsSourceFile.fileName
    //   .replace(srcDirectory, buildDirectory)
    //   .replace(/js$/, 'ts')
    //   .replace(/jsx$/, 'tsx')
    // fs.ensureFileSync(outPath)
    // fs.writeFileSync(outPath, getSourceCode(tsSourceFile))
  )
}

function transformer(context: ts.TransformationContext) {
  return (tsSourceFile: ts.SourceFile) => {
    function addStyleUrl(objectLiteral: ts.ObjectLiteralExpression): ts.ObjectLiteralExpression {
      // styleUrl attribute already provided: skipping
      if (objectLiteral.properties.find(p => (p.name as ts.Identifier).escapedText === 'styleUrl')) {
        return objectLiteral
      }

      const styleUrlsProp = objectLiteral.properties.find(
        p => (p.name as ts.Identifier).escapedText === 'styleUrls'
      ) as ts.PropertyAssignment

      if (!styleUrlsProp) {
        return objectLiteral
      }
      const findMDStyleUrl = (styleUrlsProp.initializer as ts.ObjectLiteralExpression).properties.find(
        p => (p.name as ts.Identifier).escapedText === 'md'
      ) as ts.PropertyAssignment

      if (!findMDStyleUrl) {
        return objectLiteral
      }
      return ts.createObjectLiteral([
        ...objectLiteral.properties,
        ts.createPropertyAssignment('styleUrl', ts.createLiteral((findMDStyleUrl.initializer as ts.StringLiteral).text))
      ])
    }

    function prefixTagWithBearer(objectLiteral: ts.ObjectLiteralExpression): ts.ObjectLiteralExpression {
      return ts.updateObjectLiteral(
        objectLiteral,
        objectLiteral.properties.map(prop => {
          if (ts.isPropertyAssignment(prop) && (prop.name as ts.Identifier).escapedText === 'tag') {
            return ts.updatePropertyAssignment(
              prop,
              prop.name,
              ts.createLiteral((prop.initializer as ts.StringLiteral).text.replace(/ion-/, 'bearer-'))
            )
          }
          return prop
        })
      )
    }

    function updateDecorator(decorators: ts.NodeArray<ts.Decorator>): ts.NodeArray<ts.Decorator> {
      return ts.createNodeArray(
        [...decorators].map(deco => {
          if (
            ts.isCallExpression(deco.expression) &&
            (deco.expression.expression as ts.Identifier).escapedText === 'Component'
          ) {
            // @Component found
            const [first, ...otherArgs] = deco.expression.arguments
            const objectLiteral: ts.ObjectLiteralExpression = first as ts.ObjectLiteralExpression
            const newFirstArgument = prefixTagWithBearer(addStyleUrl(objectLiteral))

            return ts.updateDecorator(
              deco,
              ts.createCall(deco.expression.expression, deco.expression.typeArguments, [newFirstArgument, ...otherArgs])
            )
          }
          return deco
        })
      )
    }

    function visit(tsNode: ts.Node) {
      if (ts.isClassDeclaration(tsNode)) {
        return ts.updateClassDeclaration(
          tsNode,
          updateDecorator(tsNode.decorators),
          tsNode.modifiers,
          tsNode.name,
          tsNode.typeParameters,
          tsNode.heritageClauses,
          tsNode.members
        )
      }
      return tsNode
    }
    return ts.visitEachChild(tsSourceFile, visit, context)
  }
}
