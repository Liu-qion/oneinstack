import * as path from 'path'
import { FileUtil } from 'tools-mock'
import pinyin from 'pinyin'
import { rename, renameSync, writeFileSync } from 'fs'

//@ts-ignore
const pinyin1 = pinyin.default ?? (pinyin as any)

const dirUrl = path.join(__dirname, '../../../dist')
const publicUrl = path.join(__dirname, '../../../dist/static/page')

const isNeedStr = (fl: string) => fl.endsWith('.js')

function extractIndexFilenames(htmlString: string): string[] {
  // 定义正则表达式来匹配文件名
  const regex = /index-[A-Za-z0-9_\-]+\.js|index-[A-Za-z0-9_\-]+\.css/g

  // 使用正则表达式提取匹配的文件名
  const matches = htmlString.match(regex)

  // 如果有匹配项，则返回匹配的文件名数组，否则返回空数组
  return matches ? matches : []
}

const init = async () => {
  const htmlDir = path.join(dirUrl, 'index.html')
  let indexStr = await FileUtil.readFile(htmlDir)
  let asArr = extractIndexFilenames(indexStr)
  for (let i = 0; i < asArr.length; i++) {
    const _item = asArr[i]
    if (_item.endsWith('.css')) {
      renameSync(path.join(publicUrl, _item), path.join(publicUrl, 'index-main.css'))
      indexStr = indexStr.replace(new RegExp(`${_item}`, 'g'), 'index-main.css')
    }
    if (_item.endsWith('.js')) {
      renameSync(path.join(publicUrl, _item), path.join(publicUrl, 'index-main.js'))
      indexStr = indexStr.replace(new RegExp(`${_item}`, 'g'), 'index-main.js')
    }
  }
  writeFileSync(htmlDir, indexStr)

  let jsPath = asArr.filter((v) => v.endsWith('.js'))[0]
  console.log('hash', jsPath)
  jsPath = jsPath.substring(0, jsPath.length - 3)
  const reg = new RegExp(jsPath, 'g')
  let fileArr = [] as any[]
  await FileUtil.readDirectory(publicUrl, async (fl) => {
    if (isNeedStr(fl)) fileArr.push(fl)
  })

  for (let i = 0; i < fileArr.length; i++) {
    const fl = fileArr[i]
    let str = await FileUtil.readFile(fl)
    str = str.replace(reg, 'index-main')
    writeFileSync(fl, str)
  }
}
const delay = (millisecond: number = 0) => {
  const _this = this
  return new Promise<void>((resolve) => {
    setTimeout(() => resolve(), millisecond)
  })
}
export const nohashindex = async () => {
  await delay(500)
  await init()
  await delay(500)
}
