import * as _ts_module from 'typescript/lib/tsserverlibrary'

function init(modules: { typescript: typeof _ts_module }) {
  const ts = modules.typescript

  function create(info: ts.server.PluginCreateInfo) {
    function getIntents(): Array<{ name: string; type: 'FetchData' | 'SaveState' | 'RetrieveState' }> {
      return info.languageServiceHost.getScriptFileNames().reduce(
        (acc, fileName) => {
          const sourceFile = info.languageService.getProgram()!.getSourceFile(fileName)
          if (sourceFile) {
            const intentMeta = ts.forEachChild(sourceFile, tsNode => {
              if (ts.isClassDeclaration(tsNode)) {
                const meta: Record<string, string> = {}
                ts.forEachChild(tsNode, child => {
                  if (isIntentName(child)) {
                    meta.name = ((child as ts.PropertyDeclaration).initializer as ts.StringLiteral).text
                  } else if (isIntentType(child)) {
                    meta.type = ((child as ts.PropertyDeclaration).initializer as ts.Identifier).text || 'type'
                  }
                })
                return meta.name ? meta : null
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

    function getSourceFile(fileName: string): ts.SourceFile | undefined {
      return info.languageService.getProgram()!.getSourceFile(fileName)
    }

    function getNode(fileName: string, position: number): ts.Node | undefined {
      const sourceFile = getSourceFile(fileName)
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

    proxy.getSemanticDiagnostics = (fileName: string): ts.Diagnostic[] => {
      const getDiagnostics = info.languageService.getSemanticDiagnostics(fileName) || []

      getDiagnostics.push({
        category: ts.DiagnosticCategory.Error,
        code: 1,
        file: info.languageService
          .getProgram()!
          .getSourceFile('/Users/tanguyantoine/devel/bearer/tmp/AuthScenarioPlugin/views/feature-action.tsx'),
        start: 220,
        length: 20,
        messageText: 'Unknown intent'
      })
      return getDiagnostics
      // const result = [...errors]
      // diagonosticsList.forEach((diagnostics, i) => {
      //   const node = nodes[i]
      //   const nodeLC = this._helper.getLineAndChar(fileName, node.getStart())
      //   diagnostics.forEach(d => {
      //     const sl = nodeLC.line + d.range.start.line
      //     const sc = d.range.start.line ? d.range.start.character : nodeLC.character + d.range.start.character
      //     const el = nodeLC.line + d.range.end.line
      //     const ec = d.range.end.line ? d.range.end.character : nodeLC.character + d.range.end.character
      //     const start = ts.getPositionOfLineAndCharacter(node.getSourceFile(), sl, sc) + 1
      //     const end = ts.getPositionOfLineAndCharacter(node.getSourceFile(), el, ec) + 1
      //     const h = start === end ? 0 : 1
      //     result.push(translateDiagnostic(d, node.getSourceFile(), start - h, end - start))
      //   })
      // })
      // return result
    }

    proxy.getCompletionsAtPosition = (fileName: string, position) => {
      const node = getNode(fileName, position)
      const prior = info.languageService.getCompletionsAtPosition(fileName, position, {})
      if (prior && node && node.parent) {
        // const entries = retrieveIntentNames()
        const entries: Array<ts.CompletionEntry> = getIntents().map(({ name, type }) => ({
          name: `${name} - ${type}`,
          kind: ts.ScriptElementKind.string,
          sortText: name,
          insertText: name
        }))

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

    // proxy.getDefinitionAtPosition = (fileName: string, position: number): ts.DefinitionInfo[] | undefined => {
    //   const prior = info.languageService.getDefinitionAtPosition(fileName, position) || []
    //   if (prior) {
    //     prior.push({
    //       fileName: '/Users/tanguyantoine/devel/bearer/tmp/AuthScenarioPlugin/intents/FetchIntent.ts',
    //       kind: ts.ScriptElementKind.string,
    //       containerKind: ts.ScriptElementKind.string,
    //       containerName: 'definition name',
    //       name: ' name',
    //       textSpan: { start: 0, length: 100 }
    //     })
    //     return prior
    //   }
    // }

    proxy.getQuickInfoAtPosition = (fileName: string, position: number) => {
      // const node = Helper.getNode(fileName, position);
      const prior = info.languageService.getQuickInfoAtPosition(fileName, position)
      if (prior) {
        const displayParts = prior.displayParts || []
        switch (prior.kind) {
          case 'method': {
            prior.displayParts = [
              ...displayParts,
              { kind: 'punctuation', text: '\n' },
              { kind: 'punctuation', text: '(' },
              { kind: 'text', text: 'method' },
              { kind: 'punctuation', text: ')' },
              { kind: 'space', text: ' ' },
              { kind: 'keyword', text: 'propkeywor' }
            ]
            break
          }
          case 'property': {
            prior.displayParts = [
              ...displayParts,
              { kind: 'punctuation', text: '\n' },
              { kind: 'punctuation', text: '(' },
              { kind: 'text', text: 'property' },
              { kind: 'punctuation', text: ')' },
              { kind: 'space', text: ' ' },
              { kind: 'keyword', text: 'propkeywor' }
            ]
            break
          }
          default: {
            prior.displayParts = [
              ...displayParts,
              { kind: 'punctuation', text: '\n' },
              { kind: 'punctuation', text: '(' },
              { kind: 'text', text: 'default' },
              { kind: 'punctuation', text: ')' },
              { kind: 'space', text: ' ' },
              { kind: 'keyword', text: 'propkeywor' }
            ]
          }
        }
      }
      return prior
    }

    proxy.getCompletionEntryDetails = (
      fileName: string,
      position: number,
      name: string,
      formatOptions: ts.FormatCodeOptions,
      source: string,
      preferences: ts.UserPreferences
    ) => {
      const prior = info.languageService.getCompletionEntryDetails(
        fileName,
        position,
        name,
        formatOptions,
        source,
        preferences
      )
      if (prior) {
        prior.displayParts = prior.displayParts || []
        prior.displayParts.push(
          { kind: 'punctuation', text: '\n' },
          { kind: 'punctuation', text: '(' },
          { kind: 'text', text: 'watched' },
          { kind: 'punctuation', text: ')' },
          { kind: 'space', text: ' ' },
          { kind: 'keyword', text: 'item.prop' }
        )
      }
      return prior
    }

    return proxy
  }

  return { create }
}

export = init
