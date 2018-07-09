import * as ts from 'typescript'
export function getSourceCode(node): string {
  const resultFile = ts.createSourceFile(
    'tmp.ts',
    '',
    ts.ScriptTarget.Latest,
    /*setParentNodes*/ false,
    ts.ScriptKind.TS
  )
  const printer = ts.createPrinter({
    newLine: ts.NewLineKind.LineFeed
  })

  return printer.printNode(ts.EmitHint.Unspecified, node, resultFile)
}
