import vue from '@vitejs/plugin-vue'
import { readFileSync, rmSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { visualizer } from 'rollup-plugin-visualizer'
import { UIViteAutoImport } from 'ui-vite/src/autoimport'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { PluginOption } from 'vite'
import { viteComType, viteDef, viteVar } from 'vite-var'
import { globalType } from '../env/globalVar'
import { tozip } from './plugins/tozip'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export const getPlugins = (env: globalType) => {
  const isBuild = env.env.pro === 'build'

  const plugin: PluginOption[] = [
    viteVar(env as any),
    viteDef(env.env.pro),
    viteComType({
      prefix: 'x'
    }),
    UIViteAutoImport({
      isBuild: true
    }),
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.startsWith('lottie-') || tag.startsWith('marquee')
        }
      }
    }),
    AutoImport({
      dts: './build/auto/auto-imports.d.ts',
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      dts: './build/auto/components.d.ts',
      resolvers: [ElementPlusResolver()]
    }),
  ]

  if (isBuild) {
    plugin.push(
      visualizer({
        open: env.env.BUILDVIEW,
        gzipSize: true,
        brotliSize: true
      }) as any,
      {
        name: 'end',
        async closeBundle() {
          const htmlDir = path.join(__dirname, '../../dist/index.html')
          let _str = readFileSync(htmlDir).toString()
          _str = _str.replace(new RegExp(`"./static/`, 'g'), `"${env.static}/`)
          writeFileSync(htmlDir, _str)
          const versionDir = path.join(__dirname, '../../dist/version.json ')
          const appName = `app-${env.version}.zip`
          writeFileSync(versionDir, JSON.stringify({ version: env.version, url: appName }))
          await tozip(path.join(__dirname, '../../dist'), path.join(__dirname, `../../version/${appName}`))
          rmSync(versionDir)
        }
      }
    )
  }

  return plugin
}
