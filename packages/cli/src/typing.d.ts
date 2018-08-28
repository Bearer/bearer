// TODO: remove this file once we migrated everything from @bearer/beare-cli
declare module 'bearer__bearer-cli'
declare module '@bearer/bearer-cli/lib/bin/index' {
  const _default: (args: any) => void
  export default _default
}

declare module '@bearer/bearer-cli/lib/src/lib/serviceClient' {
  const serviceClient: (url: string) => any

  export = serviceClient
}

type CopycallBack = (err: any, createdFiles: Array<string>) => void

declare module 'copy-template-dir' {
  const copyTemplateDir: (inDir: string, outDir: string, vars: Record<string, string>, callBack: CopycallBack) => void

  export = copyTemplateDir
}

declare module 'rc' {
  const rc: (name: string, defaults?: any, argv?: {} | null, parse?: ((content: string) => any) | null) => any

  export = rc
}
