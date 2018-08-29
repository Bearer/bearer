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

            // styleUrl attribute already provided: skipping
            if (objectLiteral.properties.find(p => (p.name as ts.Identifier).escapedText === 'styleUrl')) {
              return deco
            }

            const styleUrlsProp = objectLiteral.properties.find(
              p => (p.name as ts.Identifier).escapedText === 'styleUrls'
            ) as ts.PropertyAssignment

            if (!styleUrlsProp) {
              return deco
            }
            const findMDStyleUrl = (styleUrlsProp.initializer as ts.ObjectLiteralExpression).properties.find(
              p => (p.name as ts.Identifier).escapedText === 'md'
            ) as ts.PropertyAssignment

            if (!findMDStyleUrl) {
              return deco
            }
            const newFirstArgument = ts.createObjectLiteral([
              ...objectLiteral.properties,
              ts.createPropertyAssignment(
                'styleUrl',
                ts.createLiteral((findMDStyleUrl.initializer as ts.StringLiteral).text)
              )
            ])

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
