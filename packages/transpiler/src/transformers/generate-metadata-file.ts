import * as fs from 'fs'
import * as path from 'path'
import * as ts from 'typescript'

import { FileTransformerOptions } from '../types'

export default function generateMetadataFile(
  { metadata, outDir }: FileTransformerOptions = { outDir }
): ts.TransformerFactory<ts.SourceFile> {
  if (metadata) {
    fs.writeFileSync(path.join(outDir, 'metadata.json'), JSON.stringify(metadata), 'utf8')
  }
  return _transformContext => {
    return tsSourceFile => {
      return tsSourceFile
    }
  }
}
