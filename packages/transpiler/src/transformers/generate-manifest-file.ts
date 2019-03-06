import * as fs from 'fs'
import * as path from 'path'
import * as ts from 'typescript'
import * as vm from 'vm'

import debug from '../logger'
const logger = debug.extend('manifest-creation')
import { CompileSpec, FileTransformerOptions, RootComponent, SpecComponent } from '../types'

export const MANIFEST_FILE = 'bearer-manifest.json'
export const SPEC_FILE = 'spec.ts'

const compilerOptions = { module: ts.ModuleKind.CommonJS }

const compileSpec: (srcDir: string) => CompileSpec = srcDir => {
  const spec = ts.transpileModule(fs.readFileSync(path.join(srcDir, SPEC_FILE), 'utf8'), { compilerOptions }).outputText

  const sandbox: { exports: { default: any } } = {
    exports: { default: null }
  }

  const context = vm.createContext(sandbox)
  const script = new vm.Script(spec)
  script.runInContext(context)

  return sandbox.exports.default
}

const previewRootComponentTags = (components: SpecComponent[], rootComponents: RootComponent[]) =>
  components.map(component => {
    const { initialTagName, label } = component
    const input = component.input
    const output = component.output
    const { finalTagName, group } = rootComponents.find(({ initialTagName: tag }) => tag === initialTagName) || {
      finalTagName: null,
      group: null
    }
    return { finalTagName, group, label, input, output }
  })

const stringifyManifest: (manifest: any, srcDir: string) => string = (manifest, srcDir) => {
  const { components } = compileSpec(srcDir)
  const rootComponents = manifest.components.filter(({ isRoot }) => isRoot)

  const toBeExported = {
    manifest,
    previewRootComponents: previewRootComponentTags(components, rootComponents)
  }

  return JSON.stringify(toBeExported, null, 2)
}

export function transformer(options: FileTransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  generateManifestFile(options)
  return _transformContext => tsSourceFile => tsSourceFile
}

export default function generateManifestFile({ metadata, outDir, srcDir }: FileTransformerOptions): void {
  if (metadata) {
    fs.writeFileSync(path.join(outDir, MANIFEST_FILE), stringifyManifest(metadata, srcDir), 'utf8')
    logger('manifest generated')
  }
}
