import * as path from 'path'
import { FileUtil } from 'tools-mock'
import pinyin from 'pinyin'
import { rename, renameSync, writeFileSync } from 'fs'

//@ts-ignore
const pinyin1 = pinyin.default ?? (pinyin as any)

const dirUrl = path.join(__dirname, '../../../dist')
const publicUrl = path.join(__dirname, '../../../dist/static/page')

const isNeedRename = (fl: string) =>
  fl.endsWith('.jpg') ||
  fl.endsWith('.jpeg') ||
  fl.endsWith('.png') ||
  fl.endsWith('.svg') ||
  fl.endsWith('.webp') ||
  fl.endsWith('.gif')
const isNeedStr = (fl: string) => fl.endsWith('.js') || fl.endsWith('.css')

function convertChineseToInitials(input: string): string {
  //@ts-ignore
  const pinyinArray = pinyin1(input, {
    style: pinyin1.STYLE_FIRST_LETTER,
    heteronym: false,
    segment: true
  })
  return pinyinArray.flat().join('')
}

const renameRun = async () => {
  let fileArr = [] as any[]
  await FileUtil.readDirectory(dirUrl, async (fl: any) => {
    if (isNeedRename(fl)) fileArr.push(fl)
  })
  for (let i = 0; i < fileArr.length; i++) {
    const fl: string = fileArr[i]
    const fl1 = convertChineseToInitials(fl)
    const str = fl.substring(fl.lastIndexOf('/'))
    const str1 = fl1.substring(fl1.lastIndexOf('/'))

    for (let i = 0; i < valArr.length; i++) {
      const item = valArr[i]
      item.str = item.str.replace(new RegExp(str, 'g'), str1)
    }

    renameSync(fl, fl1)
  }

  for (let i = 0; i < valArr.length; i++) {
    const fl = valArr[i]
    writeFileSync(fl.fl, fl.str)
  }
}

const valArr = [] as { str: string; fl: string }[]

const init = async () => {
  let fileArr = [] as any[]
  await FileUtil.readDirectory(publicUrl, async (fl) => {
    if (isNeedStr(fl)) fileArr.push(fl)
  })
  for (let i = 0; i < fileArr.length; i++) {
    const fl = fileArr[i]
    let str = await FileUtil.readFile(fl)
    valArr.push({
      str,
      fl
    })
  }
  renameRun()
}

export const filename = async () => {
  await init()
}
