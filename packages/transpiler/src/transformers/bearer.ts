import * as ts from 'typescript'
import { Decorators, Component, Module } from '../constants'
// @Prop() BEARER_ID: string;
export function addBearerIdProp(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createProperty(
        [propDecorator()],
        undefined,
        'BEARER_ID',
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
        undefined
      )
    ]
  )
}

// this.BEARER_SCENARIO_ID => replaced during transpilation
export function addBearerScenarioIdAccessor(classNode: ts.ClassDeclaration, scenarioId: string): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createGetAccessor(
        undefined,
        [],
        'SCENARIO_ID',
        undefined,
        undefined,
        ts.createBlock([ts.createReturn(ts.createLiteral(scenarioId))])
      )
    ]
  )
}

// @Prop({ context: 'bearer' }) bearerContext: any
function addBearerContextProp(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createProperty(
        [
          ts.createDecorator(
            ts.createCall(ts.createIdentifier(Decorators.Prop) as ts.Expression, undefined, [
              ts.createObjectLiteral([
                ts.createPropertyAssignment(ts.createLiteral('context'), ts.createLiteral('bearer'))
              ])
            ])
          )
        ],
        undefined,
        Component.bearerContext,
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.AnyKeyword),
        undefined
      )
    ]
  )
}

// @Prop() setupId: string
export function addSetupIdProp(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createProperty(
        [
          ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.Prop) as ts.Expression, undefined, undefined))
        ],
        undefined,
        Component.setupId,
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
        undefined
      )
    ]
  )
}

function methodeNamed(name: string): (node: ts.Node) => boolean {
  return (node: ts.Node): boolean =>
    ts.isMethodDeclaration(node) && (node as ts.MethodDeclaration).name.getText() === name
}

// componentDidLoad(){ this.bearer.setupId = this.setupId }
export function addComponentDidLoad(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  const assignSetupId = ts.createStatement(
    ts.createAssignment(
      ts.createPropertyAccess(ts.createThis(), [Component.bearerContext, Component.setupId].join('.')),
      ts.createPropertyAccess(ts.createThis(), Component.setupId)
    )
  )
  const ifSetupIdPresent = ts.createIf(
    ts.createPropertyAccess(ts.createThis(), Component.setupId),
    ts.createBlock([assignSetupId])
  )
  const predicate = methodeNamed(Component.componentDidLoad)
  const members = classNode.members.filter(n => !predicate(n))

  const componentDidLoad: ts.MethodDeclaration =
    (classNode.members.find(predicate) as ts.MethodDeclaration) ||
    (ts.createMethod(
      /* decorators */ undefined,
      /* modifiers */ undefined,
      /* asteriskToken */ undefined,
      Component.componentDidLoad,
      /* questionToken */ undefined,
      /* typeParameters */ undefined,
      /* parameters */ undefined,
      /* type */ undefined,
      ts.createBlock([])
    ) as ts.MethodDeclaration)

  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...members,
      ts.createMethod(
        componentDidLoad.decorators,
        componentDidLoad.modifiers,
        componentDidLoad.asteriskToken,
        componentDidLoad.name,
        componentDidLoad.questionToken,
        componentDidLoad.typeParameters,
        componentDidLoad.parameters,
        componentDidLoad.type,
        ts.createBlock([ifSetupIdPresent, ...componentDidLoad.body.statements])
      )
    ]
  )
}

function inImportClause(node: ts.ImportClause, libName: string): boolean {
  return (
    ts.forEachChild(node, (node: ts.Node) => {
      if (ts.isNamedImports(node)) {
        return node.elements.reduce((included, element: ts.ImportSpecifier) => {
          return included || element.name.text === libName
        }, false)
      }
    }) || false
  )
}

export function hasImport(node: ts.SourceFile, libName: string): boolean {
  let has = false
  ts.forEachChild(node, node => {
    if (
      ts.isImportDeclaration(node) &&
      coreImport(node) &&
      node.importClause &&
      inImportClause(node.importClause, libName)
    ) {
      has = true
    }
  })

  return has
}

function coreImport(node: ts.ImportDeclaration): boolean {
  return Boolean(node.moduleSpecifier['text'].toString().match(Module.BEARER_CORE_MODULE))
}

function ensureHasImportFromCore(tsSourceFile: ts.SourceFile, importName: string): ts.SourceFile {
  if (hasImport(tsSourceFile, importName)) {
    return tsSourceFile
  }

  const predicate = (statement: ts.Statement): boolean => ts.isImportDeclaration(statement) && coreImport(statement)

  const importDeclaration =
    (tsSourceFile.statements.find(predicate) as ts.ImportDeclaration) ||
    ts.createImportDeclaration(
      undefined,
      undefined,
      ts.createImportClause(undefined, ts.createNamedImports([])),
      ts.createLiteral(Module.BEARER_CORE_MODULE)
    )

  const elements = (importDeclaration.importClause.namedBindings as ts.NamedImports).elements

  const clauseWithNamedImport = ts.updateImportDeclaration(
    importDeclaration,
    importDeclaration.decorators,
    importDeclaration.modifiers,
    ts.updateImportClause(
      importDeclaration.importClause,
      importDeclaration.importClause.name,
      ts.createNamedImports([...elements, ts.createImportSpecifier(undefined, ts.createIdentifier(importName))])
    ),
    importDeclaration.moduleSpecifier
  )

  const statements = [clauseWithNamedImport, ...tsSourceFile.statements.filter(el => !predicate(el))]

  return ts.updateSourceFileNode(
    tsSourceFile,
    statements,
    tsSourceFile.isDeclarationFile,
    tsSourceFile.referencedFiles,
    tsSourceFile.typeReferenceDirectives,
    tsSourceFile.hasNoDefaultLib
  )
}

function ensureHasNotImportFromCore(tsSourceFile: ts.SourceFile, importName: string): ts.SourceFile {
  if (!hasImport(tsSourceFile, importName)) {
    return tsSourceFile
  }

  const predicate = (statement: ts.Statement): boolean => ts.isImportDeclaration(statement) && coreImport(statement)

  const importDeclaration =
    (tsSourceFile.statements.find(predicate) as ts.ImportDeclaration) ||
    ts.createImportDeclaration(
      undefined,
      undefined,
      ts.createImportClause(undefined, ts.createNamedImports([])),
      ts.createLiteral(Module.BEARER_CORE_MODULE)
    )

  const elements = (importDeclaration.importClause.namedBindings as ts.NamedImports).elements.filter(element => {
    if (element.name.text === importName) {
      return
    }
    return element
  })

  const clauseWithNamedImport = ts.updateImportDeclaration(
    importDeclaration,
    importDeclaration.decorators,
    importDeclaration.modifiers,
    ts.updateImportClause(
      importDeclaration.importClause,
      importDeclaration.importClause.name,
      ts.createNamedImports(elements)
    ),
    importDeclaration.moduleSpecifier
  )

  const statements = [clauseWithNamedImport, ...tsSourceFile.statements.filter(el => !predicate(el))]

  return ts.updateSourceFileNode(
    tsSourceFile,
    statements,
    tsSourceFile.isDeclarationFile,
    tsSourceFile.referencedFiles,
    tsSourceFile.typeReferenceDirectives,
    tsSourceFile.hasNoDefaultLib
  )
}

export function ensureBearerContextInjected(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  const has: boolean = ts.forEachChild(
    classNode,
    node => ts.isPropertyDeclaration(node) && node.name['escapedText'] == Component.bearerContext
  )

  return has ? classNode : addBearerContextProp(classNode)
}

export function ensureWatchImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasImportFromCore(tsSourceFile, Decorators.Watch)
}

export function ensurePropImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasImportFromCore(tsSourceFile, Decorators.Prop)
}

export function ensureComponentImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasImportFromCore(tsSourceFile, Decorators.Component)
}

export function ensureRootComponentNotImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasNotImportFromCore(tsSourceFile, Decorators.RootComponent)
}

export function ensureStateImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasImportFromCore(tsSourceFile, Decorators.State)
}

export function ensureElementImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasImportFromCore(tsSourceFile, Decorators.Element)
}

export function propDecorator() {
  return ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.Prop) as ts.Expression, undefined, undefined))
}

export function elementDecorator() {
  return ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Element) as ts.Expression, undefined, undefined)
  )
}

export default {
  addBearerScenarioIdAccessor,
  addBearerIdProp,
  addBearerContextProp,
  addSetupIdProp,
  addComponentDidLoad,
  hasImport,
  coreImport,
  ensureBearerContextInjected,
  ensurePropImported
}
