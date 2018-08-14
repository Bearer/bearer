declare module 'bearer__bearer-cli'
declare module '@bearer/bearer-cli/dist/bin/index' {
  const _default: (args: any) => void
  export default _default
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
