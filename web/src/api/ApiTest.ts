export const ApiTest = {
  /** 上传文件 */
  upload: (obj: any) => {
    return http.post(
      {
        url: '/system/upload',
        headers: {
          'Content-Type': 'multipart/form-data'
        },
        onProgress(percent, loaded, total) {
          if (obj?.onProgress) obj.onProgress(percent, loaded, total)
        }
      },
      obj
    )
  },

  /** 获取系统时间 */
  getTime: (obj: any) => {
    return http.get('/system/getTime', obj)
  }
}
