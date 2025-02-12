type Point = {
  x: number
  y: number
  time: number // 时间戳，用于计算速度
}

export class InertiaScroller {
  //是否正在获取竞争
  static getIsRunPro = false

  //获取竞争
  static getIsRun = () => {
    if (this.getIsRunPro) {
      return false
    }
    this.getIsRunPro = true
    setTimeout(() => {
      this.getIsRunPro = false
    }, 50)
    return true
  }

  private container: HTMLElement
  private points: Point[] = []
  private inertiaFrame: number | null = null

  constructor(container: HTMLElement) {
    this.container = container
    const overflow = this.container.style.overflow
    if (StrUtil.isNull(overflow)) {
      this.container.style.overflow = 'scroll'
    }
    this.bindEvents()
  }

  private bindEvents() {
    this.container.addEventListener('mousedown', this.onMouseDown)
  }

  private needInertia = false

  private onMouseDown = (e: MouseEvent) => {
    const isRun = InertiaScroller.getIsRun()
    if (!isRun) return

    this.needInertia = false
    this.clearInertia()
    this.points = []
    this.recordPoint(e.clientX, e.clientY)

    let moveNum = 0
    const onMouseMove = (event: MouseEvent) => {
      this.recordPoint(event.clientX, event.clientY)
      this.container.scrollLeft -= event.movementX
      this.container.scrollTop -= event.movementY

      if (moveNum == 1) this.container.style.pointerEvents = 'none'
      moveNum = 1
    }

    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
      this.needInertia = true
      this.startInertia()
      this.container.style.pointerEvents = 'all'
    }

    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  private recordPoint(x: number, y: number) {
    const now = Date.now()
    this.points.push({ x, y, time: now })
    // 保留最近的几个点用于计算速度
    if (this.points.length > 5) {
      this.points.shift()
    }
  }

  private calculateVelocity(): { vx: number; vy: number } {
    if (this.points.length < 2) return { vx: 0, vy: 0 }

    const lastPoint = this.points[this.points.length - 1]
    const firstPoint = this.points[0]

    const timeDiff = (lastPoint.time - firstPoint.time) / 1000 // 转为秒
    if (timeDiff === 0) return { vx: 0, vy: 0 }

    const vx = (lastPoint.x - firstPoint.x) / timeDiff
    const vy = (lastPoint.y - firstPoint.y) / timeDiff

    return { vx, vy }
  }

  private startInertia() {
    let { vx, vy } = this.calculateVelocity()
    if (Math.abs(vx) < 1 && Math.abs(vy) < 1) return // 速度过低不进行惯性滑动

    const friction = 0.95 // 摩擦力（速度衰减系数）

    const step = () => {
      if (!this.needInertia) return
      vx *= friction
      vy *= friction

      // 更新容器的位置或内容的滚动位置
      this.container.scrollLeft -= vx * 0.016
      this.container.scrollTop -= vy * 0.016

      // 速度衰减
      const newVx = vx * friction
      const newVy = vy * friction

      // 如果速度接近于零，停止惯性滑动
      if (Math.abs(newVx) < 10 && Math.abs(newVy) < 10) {
        this.clearInertia()
        return
      }

      this.inertiaFrame = requestAnimationFrame(step)
    }

    this.inertiaFrame = requestAnimationFrame(step)
  }

  private clearInertia() {
    if (this.inertiaFrame !== null) {
      cancelAnimationFrame(this.inertiaFrame)
      this.inertiaFrame = null
    }
  }
}
