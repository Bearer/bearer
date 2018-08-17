import * as ts from 'typescript'

// function appendToStatements(
//   tsMethod: ts.MethodDeclaration,
//   statements: ReadonlyArray<ts.Statement>
// ): ts.MethodDeclaration {
//   return updateMethodStatements(tsMethod, [], statements)
// }

export function prependToStatements(
  tsMethod: ts.MethodDeclaration,
  statements: ReadonlyArray<ts.Statement>
): ts.MethodDeclaration {
  return updateMethodStatements(tsMethod, statements, [])
}

function updateMethodStatements(
  tsMethod: ts.MethodDeclaration,
  prependStatements: ReadonlyArray<ts.Statement>,
  appendStatements: ReadonlyArray<ts.Statement>
): ts.MethodDeclaration {
  return ts.updateMethod(
    tsMethod,
    tsMethod.decorators,
    tsMethod.modifiers,
    tsMethod.asteriskToken,
    tsMethod.name,
    tsMethod.questionToken,
    tsMethod.typeParameters,
    tsMethod.parameters,
    tsMethod.type,
    ts.createBlock([...prependStatements, ...tsMethod.body.statements, ...appendStatements])
  )
}

export function updateMethodOfClass(
  classNode: ts.ClassDeclaration,
  methodeName: string,
  updater: (method: ts.MethodDeclaration) => ts.MethodDeclaration
): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    classNode.members.map(member => {
      if (ts.isMethodDeclaration(member) && (member.name as ts.Identifier).escapedText === methodeName) {
        return updater(member)
      }
      return member
    })
  )
}
