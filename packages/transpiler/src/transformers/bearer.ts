import * as ts from 'typescript'

import { BEARER, Component, Decorators, Module, Types, SETUP, SETUP_ID } from '../constants'
import { getNodeName } from '../helpers/node-helpers'
import { getDecoratorNamed } from '../helpers/decorator-helpers'

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
        Component.scenarioIdAccessor,
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
                ts.createPropertyAssignment(ts.createLiteral('context'), ts.createLiteral(BEARER))
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

function usingSetupProperty(node: ts.ClassElement, source: ts.SourceFile) {
  return (
    (getDecoratorNamed(node, Decorators.Input) || getDecoratorNamed(node, Decorators.Output)) &&
    (node.name && node.name.getText(source) === SETUP)
  )
}

function usingSetupId(node: ts.ClassElement, source: ts.SourceFile) {
  return getDecoratorNamed(node, Decorators.Prop) && (node.name && node.name.getText(source) === SETUP_ID)
}

function setupIdAlreadyIncluded(classNode: ts.ClassDeclaration, source: ts.SourceFile): boolean {
  const scanned = classNode.members.map(member => {
    return usingSetupProperty(member, source) || usingSetupId(member, source)
  })
  for (const i in scanned) {
    if (scanned[i]) return true
  }
  return false
}
// @Prop() setupId: string
export function addSetupIdProp(classNode: ts.ClassDeclaration, source: ts.SourceFile): ts.ClassDeclaration {
  if (setupIdAlreadyIncluded(classNode, source)) {
    return classNode
  }
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
  return (node: ts.Node): boolean => ts.isMethodDeclaration(node) && getNodeName(node) === name
}

const componentDidLoadMethod = methodeNamed(Component.componentDidLoad)

export function createOrUpdateComponentDidLoad(
  classNode: ts.ClassDeclaration,
  updater: (block: ts.Block) => ts.Block
): ts.ClassDeclaration {
  const otherMemebers = classNode.members.filter(n => !componentDidLoadMethod(n))

  const componentDidLoad: ts.MethodDeclaration =
    (classNode.members.find(componentDidLoadMethod) as ts.MethodDeclaration) ||
    (ts.createMethod(
      /* decorators */ undefined,
      /* modifiers */ undefined,
      /* asteriskToken */ undefined,
      Component.componentDidLoad,
      /* questionToken */ undefined,
      /* typeParameters */ undefined,
      /* parameters */ undefined,
      /* type */ undefined,
      ts.createBlock([], true)
    ) as ts.MethodDeclaration)
  const newComponentDidload = ts.updateMethod(
    componentDidLoad,
    componentDidLoad.decorators,
    componentDidLoad.modifiers,
    componentDidLoad.asteriskToken,
    componentDidLoad.name,
    componentDidLoad.questionToken,
    componentDidLoad.typeParameters,
    componentDidLoad.parameters,
    componentDidLoad.type,
    updater(componentDidLoad.body || ts.createBlock([], true))
  )
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [...otherMemebers, newComponentDidload]
  )
}

export function addComponentDidLoad(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  const assignSetupId = ts.createStatement(
    ts.createAssignment(
      ts.createPropertyAccess(ts.createThis(), [Component.bearerContext, Component.setupId].join('.')),
      ts.createPropertyAccess(ts.createThis(), Component.setupId)
    )
  )
  const ifSetupIdPresent = ts.createIf(
    ts.createPropertyAccess(ts.createThis(), Component.setupId),
    ts.createBlock([assignSetupId], true)
  )
  return createOrUpdateComponentDidLoad(classNode, block =>
    ts.updateBlock(block, [ifSetupIdPresent, ...block.statements])
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
  return Boolean((node.moduleSpecifier as ts.StringLiteral).text.toString().match(Module.BEARER_CORE_MODULE))
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

export function ensureBearerContextInjected(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  const has: boolean = ts.forEachChild(
    classNode,
    node => ts.isPropertyDeclaration(node) && (node.name as ts.Identifier).escapedText === Component.bearerContext
  )

  return has ? classNode : addBearerContextProp(classNode)
}

export function ensureImportsFromCore(tsSourceFile: ts.SourceFile, decorators: (Decorators | Types)[]): ts.SourceFile {
  return decorators.reduce((sourceFile, decorator) => ensureHasImportFromCore(sourceFile, decorator), tsSourceFile)
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
  addBearerContextProp,
  addSetupIdProp,
  addComponentDidLoad,
  hasImport,
  coreImport,
  ensureBearerContextInjected
}
