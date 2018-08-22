import * as copy from 'copy-template-dir'
import * as path from 'path'

import BaseCommand from '../BaseCommand'

export function printFiles(command: BaseCommand, files: Array<string>) {
  const base = process.cwd()
  files.forEach(file => {
    command.log(command.colors.gray(`    create: `) + command.colors.white(file.replace(base + '/', '')))
  })
}

export function copyFiles(
  command: BaseCommand,
  sourceDirectory: string,
  outPutDirectory: string,
  vars: any,
  silent: boolean = false
): Promise<Array<string>> {
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
