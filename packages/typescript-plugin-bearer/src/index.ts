import * as ts_module from 'typescript/lib/tsserverlibrary'

function init(_modules: { typescript: typeof ts_module }) {
  // const ts = modules.typescript

  function create(info: ts.server.PluginCreateInfo) {
    info.project.projectService.logger.info('BEARER typescript plugin initialization')
    info.project.projectService.logger.msg('BEARER typescript plugin initialization')
    // Set up decorator
    const proxy: ts.LanguageService = Object.create(null)
    for (let k of Object.keys(info.languageService) as Array<keyof ts.LanguageService>) {
      const x = info.languageService[k]
      proxy[k] = (...args: Array<{}>) => (x as any).apply(info.languageService, args)
    }

    // proxy.getQuickInfoAtPosition = (fileName: string, position: number) => {
    //   // const node = Helper.getNode(fileName, position);
    //   const prior = info.languageService.getQuickInfoAtPosition(fileName, position)!

    //   if (prior && (prior.kind === 'method' || prior.kind === 'property')) {
    //     prior.displayParts!.push({
    //       kind: 'text',
    //       text: 'bearer bro'
    //     })
    //   }

    //   info.project.projectService.logger.info(`[test] QuickInfo "${JSON.stringify(prior, null, 2)}"`)

    //   return prior
    // }
    // proxy.getCompletionsAtPosition = (fileName, position) => {
    //   // const prior = info.languageService.getCompletionsAtPosition(fileName, position, [])
    //   // prior.entries = prior.entries.filter(e => e.name !== 'caller')
    //   return prior
    // }
    return proxy
  }

  return { create }
}

export = init
