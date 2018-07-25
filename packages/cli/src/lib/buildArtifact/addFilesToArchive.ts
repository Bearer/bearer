import path from 'path'
import globby from 'globby'
import fs from 'graceful-fs'

export default (archive, packagePath) => {
  const fullPath = path.resolve(packagePath)
  return globby([`${fullPath}/dist/*.js`])
    .then(files => {
      files.forEach(file => {
        archive.append(fs.createReadStream(file), {
          name: file.replace(path.resolve(packagePath), '')
        })
      })
    })
    .catch(console.error)
}
