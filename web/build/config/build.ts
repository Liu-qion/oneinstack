import { BuildOptions } from 'vite'
import { globalType } from '../env/globalVar'

const buildArr = [
  'tools-javascript',
  'tools-websocket',
  'tools-vue3',
  'swiper',
  'element-plus',
  'vue-router',
  'vue+runtime',
  'vue',
  'qrcode',
  'lottie',
  'open'
]
function extractPackageName(input: string): string | null {
  const packageNameRegex = /^(@?[^@]+(?:\+[^@]+)?)(?:@[\d.]+)?/
  const match = input.match(packageNameRegex)
  if (match) {
    return match[1]
  }
  return null
}

export const getBuild = (env: globalType) => {
  const build: BuildOptions = {
    assetsDir: env.assetsDir,
    minify: 'terser',
    cssCodeSplit: true, // 如果设置为false，整个项目中的所有 CSS 将被提取到一个 CSS 文件中
    reportCompressedSize: false,
    terserOptions: {
      compress: {
        drop_console: true, // 打包时删除console
        drop_debugger: true, // 打包时删除 debugger
        pure_funcs: ['console.log']
      },
      output: {
        comments: false // false-去掉注释内容 true-保留
      }
    },
    rollupOptions: {
      external: (source: string, importer: string | undefined, isResolved: boolean) => {
        return importer?.endsWith('less')
      },
      output: {
        // entryFileNames: `${env.assetsDir}/[name].js`,
        // chunkFileNames: `${env.assetsDir}/[name].js`,
        // assetFileNames: `${env.assetsDir}/[name].[ext]`,
        manualChunks(id) {
          if (id.includes('node_modules') && env.env.BUILDVIEW) {
            const arr = id.toString().split('node_modules/')[1].split('/')
            let npmName: any = arr[0].indexOf('pnpm') != -1 ? arr[1] : arr[0]
            npmName = extractPackageName(npmName)
            for (let i = 0; i < buildArr.length; i++) {
              const key = buildArr[i]
              if (npmName.indexOf(key) != -1) return key
            }
            return '_vendor'
          }
        }
      }
    }
  }
  return build
}
