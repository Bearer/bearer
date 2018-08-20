import * as copy from 'copy-template-dir'
import * as path from 'path'

import BaseCommand from '../BaseCommand'

export function copyFiles(
  command: BaseCommand,
  sourceDirectory: string,
  outPutDirectory: string,
  vars: any
): Promise<boolean | string> {
  const base = process.cwd()
  function printFiles(files: Array<string>) {
    files.forEach(file => {
      command.log(command.colors.gray(`    create: `) + command.colors.white(file.replace(base + '/', '')))
    })
  }

  return new Promise((resolve, reject) => {
    copy(path.join(__dirname, '..', '..', 'templates', sourceDirectory), outPutDirectory, vars, (err, createdFiles) => {
      if (err) {
        reject(err)
      } else {
        printFiles(createdFiles)
        resolve(true)
      }
    })
  })
}
