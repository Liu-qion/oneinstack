import { reactive } from 'vue'
export const svf = reactive({
  /**
   * 获取校验对象-默认使用对数字进行只保留数字、正数、小数后两位保留修正操作
   * 
   * @example
   * 
   * //template
   * <input
        @input="vfFun($event,'username')"
        v-model.trim="username"
      />
   * 
   * //js-created
   * 在data中
   * 
      vfFun:(e,name)=>this.vf[name](e)
   * 
   * 在created添加中
			this.vf = svf.getVf(this, {
				username: {
					getNum:true,
					getEn:true,
					getCustom:'_.-',
				},
				password: {
					getNum:true,
					getEn:true,
					getCustom:'_.-',
				}
			})
   */
  getVf: (obj: any, config: { [key: string]: Partial<inputReg> }, fun?: any) => {
    const vrf: any = {}
    if (fun === undefined) fun = svf.inputVal
    for (const key in config) {
      vrf[key] = (e: any) => {
        obj[key] = null
        fun({
          val: e.target.value,
          after: (val: any) => {
            obj[key] = val
          },
          ...config[key]
        })
      }
    }
    return vrf
  },
  /**
     * 处理规则
     *
     * @param trim-去除换行回车空格（默认为true）
     *
     * @param int-只要整数
     *
     * @param float-只要小数
     * @param fixed-保留小数后n位（默认为2）
     *
     * @param min-数字最小-小于更改为最小（默认为undefined）
     * @param max-数字最大-大于更改为最大（默认为undefined）
     * 
     * @param length-长度限制（默认为undefined）
     *
     * @param getHanzi-需要汉字
     * @param getNum-需要数字
     * @param getEn-需要字母
     * @param getCustom-自定义需要内容（如-_.等）
     *
     * @param delHanzi-去除汉字
     * @param delNum-去除数字
     * @param delEn-去除字母
     * @param delCustom-自定义去除内容（如-_.等）
     * 
     * 
     * @example
     * 
     * 
			svf.inputVal({
				val:`123asd你好...!!0d__.6s           “”
				
				dd
				`,
				after:(val)=>{
					console.log('val',val);
				},
				// trim:false,
				// length:6,
				// int:true,
				// float:true,
				// fixed:1,
				// min:1,
				// max:10,
				// getNum:true,
				// getEn:true,
				// getHanzi:true,
				// getCustom:['_.'],
				// delNum:true,
				// delEn:true,
				// delHanzi:true,
				// delCustom:['_.'],
			})
     *
     */
  inputVal: (config: inputReg & { after: (val: any) => void; val: string }) => {
    let { val, length, trim = true, fixed = 2, min = 0, max = 1000000000 } = config
    if (trim) val = val.replace(/\s+|\s+$/g, '')
    const conf: any = {
      int: (val: any) => {
        val = val.replace(/[^\d]/g, '')
        val = Number(val) || 0
        val = min ? (val < min ? min : val) : val
        val = max ? (val > max ? max : val) : val
        return val
      },
      float: (val: any) => {
        val = val
          .replace(/^0(\d)/, '$1')
          .replace(/^\./, '0')
          .replace(/[^\d.]/g, '')
          .replace(/^(\d*\.?\d*).*$/, '$1')
          .replace(new RegExp(`^(\\d+)\\.(\\d{${fixed}}).*`), '$1.$2')
        return val
      },
      hanzi: (val: any) => {
        val = val.replace(/[^\u4e00-\u9fa5]/g, '')
        return val
      },
      Num: '0-9',
      En: 'a-zA-Z',
      Hanzi: '\u4e00-\u9fa5'
    }

    if (val.length) {
      if (length && val.length > length) val = val.substring(0, length)
      const getReg: any[] = []
      const delReg: any[] = []
      const pushReg = (key: any, name: any, arr: any) => {
        if (key === name) {
          if (Array.isArray(config[key])) arr.push(...config[key])
          else arr.push(config[key])
        } else {
          arr.push(conf[key.substring(3)])
        }
      }
      Object.keys(config).forEach((key) => {
        if (key.startsWith('get')) {
          pushReg(key, 'getCustom', getReg)
        } else if (key.startsWith('del')) {
          pushReg(key, 'delCustom', delReg)
        } else {
          if (conf[key]) val = conf[key](val)
        }
      })
      if (getReg.length) {
        val = val.replace(new RegExp(`[^${getReg.join('')}]`, 'g'), '')
      }
      if (delReg.length) {
        val = val.replace(new RegExp(`[${delReg.join('')}]`, 'g'), '')
      }
    }

    config.after(val)

    return val
  }
})

export default svf

type inputReg = {
  [key: string]: any
  /**
   * 去除换行回车空格（默认为true）
   */
  trim: boolean
  /**
   * 只要整数
   */
  int: boolean
  /**
   * 只要小数
   */
  float: boolean
  /**
   * 保留小数后n位（默认为2）
   */
  fixed: number
  /**
   * 数字最小-小于更改为最小（默认为undefined）
   */
  min: number
  /**
   * 数字最大-大于更改为最大（默认为undefined）
   */
  max: number
  /**
   * 长度限制（默认为undefined）
   */
  length: number
  /**
   * 需要汉字
   */
  getHanzi: boolean
  /**
   * 需要数字
   */
  getNum: boolean
  /**
   * 需要字母
   */
  getEn: boolean
  /**
   * 自定义需要内容（如-_.等）
   */
  getCustom: string
  /**
   * 去除汉字
   */
  delHanzi: boolean
  /**
   * 去除数字
   */
  delNum: boolean
  /**
   * 去除字母
   */
  delEn: boolean
  /**
   * 自定义去除内容（如-_.等）
   */
  delCustom: string
}
