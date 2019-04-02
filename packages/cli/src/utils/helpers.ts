import * as copy from 'copy-template-dir'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../base-command'

export function printFiles(command: BaseCommand, files: string[]) {
  const base = process.cwd()
  const { gray, white } = command.colors
  files.forEach(file => {
    command.log(gray(`    create: `) + white(file.replace(`${base}/`, '')))
  })
}

export function copyFiles(
  command: BaseCommand,
  sourceDirectory: string,
  outPutDirectory: string,
  vars: any,
  silent = false
): Promise<string[]> {
  return new Promise((resolve, reject) => {
    copy(path.join(__dirname, '..', '..', 'templates', sourceDirectory), outPutDirectory, vars, (err, createdFiles) => {
      if (err) {
        reject(err)
      } else {
        if (!command.silent && !silent) {
          printFiles(command, createdFiles)
        }
        resolve(createdFiles)
      }
    })
  })
}

export function ensureFolderExists(path: string, empty = false) {
  if (!fs.existsSync(path)) {
    fs.mkdirpSync(path)
  }
  if (empty) {
    fs.emptyDirSync(path)
  }
}

export function ensureSymlinked(target: string, sourcePath: string): void {
  try {
    fs.symlinkSync(target, sourcePath)
  } catch (e) {
    if (e.code !== 'EEXIST') {
      throw e
    }
  }
}

export function toParams(obj: Record<string, string | number>) {
  return Object.keys(obj)
    .map(key => [key, obj[key]].join('='))
    .join('&')
}
