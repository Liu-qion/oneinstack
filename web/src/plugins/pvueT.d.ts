interface ScrollConf {
  /**
   * 滚动到顶部
   * @param top 滚动到的位置
   * @param useAni 是否使用动画
   */
  toTop: (top: number, useAni?: boolean) => void
  /**
   * 滚动到左侧
   * @param left 滚动到的位置
   * @param useAni 是否使用动画
   */
  toLeft: (left: number, useAni?: boolean) => void
  /**
   * 滚动到最右侧
   * @param useAni 是否使用动画
   */
  toRight: (useAni?: boolean) => void
  /**
   * 滚动到最下侧
   * @param useAni 是否使用动画
   */
  toBottom: (useAni?: boolean) => void
  /**
   * 将dom完全显示出来
   * @param dom 元素或元素选择器
   * @param useAni 是否使用动画
   */
  showDom: (dom: HTMLElement | string, useAni?: boolean) => void
}
