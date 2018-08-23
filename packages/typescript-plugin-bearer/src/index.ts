import * as _ts_module from 'typescript/lib/tsserverlibrary'

function init(modules: { typescript: typeof _ts_module }) {
  const ts = modules.typescript

  function create(info: ts.server.PluginCreateInfo) {
    function getNode(fileName: string, position: number): ts.Node | undefined {
      const sourceFile = info.languageService.getProgram()!.getSourceFile(fileName)
      if (!sourceFile) {
        return
      }
      function find(node: ts.Node): ts.Node | undefined {
        if (position >= node.getStart() && position < node.getEnd()) {
          return ts.forEachChild(node, find) || node
        }
      }
      return find(sourceFile)
    }

    info.project.projectService.logger.info('BEARER typescript plugin initialization')
    // ensure compiler option has decorator enabled
    info.project.setCompilerOptions({ ...info.project.getCompilerOptions(), experimentalDecorators: true })

    // Set up decorator
    const proxy: ts.LanguageService = Object.create(null)
    for (let k of Object.keys(info.languageService) as Array<keyof ts.LanguageService>) {
      const x = info.languageService[k]
      proxy[k] = (...args: Array<{}>) => (x as any).apply(info.languageService, args)
    }
    info.languageServiceHost.getScriptFileNames()
    // Autocompletion

    function isIntentDecorator(decorator: ts.Decorator): boolean {
      return ts.isCallExpression(decorator.expression) && decorator.expression.expression.getText() === 'Intent'
    }

    function isIntentName(tsNode: ts.Node): boolean {
      return Boolean(
        ts.isPropertyDeclaration(tsNode) &&
          tsNode.modifiers &&
          tsNode.modifiers.find(m => m.kind === ts.SyntaxKind.StaticKeyword) &&
          tsNode.name.getText() === 'intentName'
      )
    }

    function isIntentType(tsNode: ts.Node): boolean {
      return Boolean(
        ts.isPropertyDeclaration(tsNode) &&
          tsNode.modifiers &&
          tsNode.modifiers.find(m => m.kind === ts.SyntaxKind.StaticKeyword) &&
          tsNode.name.getText() === 'intentType'
      )
    }

    function retrieveIntentNames() {
      return info.languageServiceHost.getScriptFileNames().reduce(
        (acc, fileName) => {
          const sourceFile = info.languageService.getProgram()!.getSourceFile(fileName)
          if (sourceFile) {
            const intentMeta = ts.forEachChild(sourceFile, tsNode => {
              if (ts.isClassDeclaration(tsNode)) {
                const metas: Record<string, string> = {}
                ts.forEachChild(tsNode, child => {
                  if (isIntentName(child)) {
                    metas.name = ((child as ts.PropertyDeclaration).initializer as ts.StringLiteral).text
                  } else if (isIntentType(child)) {
                    metas.type = ((child as ts.PropertyDeclaration).initializer as ts.Identifier).text || 'type'
                  }
                })
                return metas.name
                  ? {
                      name: `${metas.name} - ${metas.type}`,
                      kind: ts.ScriptElementKind.string,
                      sortText: metas.name,
                      insertText: metas.name
                    }
                  : null
              }
            })
            if (intentMeta) {
              return [...acc, intentMeta]
            }
          }
          return acc
        },
        [] as any[]
      )
    }

    proxy.getCompletionsAtPosition = (fileName: string, position) => {
      const node = getNode(fileName, position)
      const prior = info.languageService.getCompletionsAtPosition(fileName, position, {})
      if (prior && node && node.parent) {
        const entries = retrieveIntentNames()

        const quotedEntries = entries.map(e => ({ ...e, insertText: `"${e.insertText}"` }))

        // has quote (single or double)
        if (isIntentDecorator(node.parent as ts.Decorator)) {
          prior.entries = [...quotedEntries, ...prior.entries]
        } else if (ts.isDecorator(node.parent.parent) && isIntentDecorator(node.parent.parent)) {
          prior.entries = [...entries, ...prior.entries]
        }
      }
      return prior
    }
    return proxy
  }

  return { create }
}

export = init
