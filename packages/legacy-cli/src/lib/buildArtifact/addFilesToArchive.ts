import * as globby from 'globby'
import * as fs from 'graceful-fs'
import * as path from 'path'

export default (archive, packagePath) => {
  const fullPath = path.resolve(packagePath)
  return globby([`${fullPath}/*.js`])
    .then(files => {
      files.forEach(file => {
        archive.append(fs.createReadStream(file), {
          name: file.replace(path.resolve(packagePath), '')
        })
      })
    })
    .catch(console.error)
}
