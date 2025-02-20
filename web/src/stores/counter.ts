import { defineStore } from 'pinia';

// 使用 defineStore 定义一个名为 counter 的 store
export const useCounterStore = defineStore('counter', {
  // 状态：定义 store 的数据
  state: () => ({
    count: 77,
    theme:'light',//'light' | 'dark',主题
    menu_status: true,//菜单栏状态 打开/关闭
    is_menu: false,//是否显示菜单
    terminal_width:600,//终端宽度
    terminal_left: 206,//终端距离浏览器左侧的距离
    terminal_height: 650,//终端高度
    menu_width_on: 206,//菜单栏打开宽度
    menu_width_off: 78,//菜单栏关闭宽度
    window_width: 0,//浏览器窗口宽度
    title:'Home9'//网页标题
    
  }),
  // 获取器：类似于计算属性，用于获取状态的派生数据
  getters: {
    doubleCount: (state) => {
        state.window_width = window.innerWidth;
        
        state.terminal_width = window.innerWidth-(state.menu_status?state.menu_width_on:state.menu_width_off)-26;
    }

  },
  // 动作：用于修改状态，支持同步和异步操作
  actions: {
    increment() {
      this.count++;
    },
    async incrementAsync() {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      this.increment();
    },
    upDateTitle(title: any) {
        this.title = title
        console.log('title:',this.title)
    },
    updateDimensions(width: number, height: number) {
        
        this.window_width =  window.innerWidth
        this.terminal_width = window.innerWidth-(this.menu_status?this.menu_width_on:this.menu_width_off)-26;
        this.terminal_left = this.menu_status?this.menu_width_on:this.menu_width_off
      }
  }
});