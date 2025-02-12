interface ApiRes {
  /**
   * 错误代码-未处理
   */
  code: any
  /**
   * 错误信息-未处理
   */
  msg: string
  /**
   * 请求数据
   */
  data: any
}
interface ApiData {
  [key: string]: any
  /**
   * 自定义请求错误回调
   * @param msg 错误信息
   * @returns
   */
  toast?: (
    /**
     * 错误代码
     */
    code: any,
    /**
     * 请求响应
     */
    xhr: TXHRRes,
    /**
     * 请求配置
     */
    config: callbackConfig<ApiRes>
  ) => void
  /**
   * 请求完毕后回调
   * @param config
   * @returns
   */
  final?: (
    /** 成功状态 */
    status: boolean,
    /**
     * 请求配置
     */
    config: callbackConfig<ApiRes>,
    /**
     * 请求响应
     */
    xhr: TXHRRes
  ) => void
}

type TXHRRes = {
  /**
   * 请求响应数据
   */
  data: ApiRes
  /**
   * 请求响应类型
   */
  responseType: any
  /**
   * 请求响应代码
   */
  status: number
  /**
   * 请求响应对象
   */
  xhr: XMLHttpRequest
}
