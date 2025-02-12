import * as fs from 'fs'
import archiver from 'archiver'
import path from 'path'

/**
 * 将指定目录打包到指定目录的文件
 * @param sourceDir - 目标地址
 * @param outPath - 生成的zip文件地址
 */
export async function tozip(sourceDir: string, outPath: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const output = fs.createWriteStream(outPath)
    const archive = archiver.create('zip', {
      zlib: { level: 9 } // Set the compression level
    })

    output.on('close', () => {
      const sizeInBytes = archive.pointer()
      const sizeInKB = sizeInBytes / 1024
      console.log(`打包成zip完毕，大小： ${sizeInKB.toFixed(2)} KB`)
      resolve()
    })

    archive.on('warning', (err) => {
      if (err.code !== 'ENOENT') {
        reject(err)
      } else {
        console.warn(err)
      }
    })

    archive.on('error', (err) => {
      reject(err)
    })

    archive.pipe(output)
    archive.glob('**/*', {
      cwd: sourceDir,
      ignore: [path.relative(sourceDir, outPath)]
    })

    archive.finalize()
  })
}
