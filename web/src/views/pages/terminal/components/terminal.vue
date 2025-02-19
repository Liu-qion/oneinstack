<template>
  <div class="terminalPage">
      <div class="terminal" ref="terminalDiv"></div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { Terminal } from 'xterm';
import 'xterm/css/xterm.css';
import { FitAddon } from 'xterm-addon-fit';

// 定义终端容器的引用
const terminalDiv = ref<HTMLElement | null>(null);


const baseWsUrl = 'ws://192.168.31.106:8089/v1/ssh/open';
  const authorization = localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo') || '')?.token : '';
  const wsUrl = `${baseWsUrl}?Authorization=${authorization}`;
// 定义状态变量
let xterm: Terminal | null = null;
let fitAddon: FitAddon | null = null;
let ws: WebSocket | null = null;
let commandBuffer = '';

// 初始化终端的方法
const initTerminal = () => {
  xterm = new Terminal({
      rows: 30, // 可视窗口的行数
      cols: 80, // 可视窗口的列数
      cursorBlink: true, // 光标是否闪烁
      // 终端主题
      theme: {
          foreground: '#ffffff',
          background: '#000000',
          cursor: '#ffffff',
      },
      fontFamily: 'Consolas, Courier, monospace',
      fontSize: 14,
  });

  if (terminalDiv.value) {
      xterm.open(terminalDiv.value);
  }

  fitAddon = new FitAddon();
  xterm.loadAddon(fitAddon);
  fitAddon.fit();

  const handleResize = () => {
      try {
          fitAddon?.fit();
      } catch (e) {
          console.error((e as Error).message);
      }
  };

  window.addEventListener('resize', handleResize);

  onBeforeUnmount(() => {
      window.removeEventListener('resize', handleResize);
  });

  xterm.write('Connecting \r\n');
  xterm.onData((data) => terminalOnData(data));
  xterm.onKey((event) => terminalOnKey(event));
};

// 处理终端输入数据的方法
const terminalOnData = (data: string) => {
  xterm?.write(data);
  commandBuffer += data;
};

// 处理终端按键事件的方法
const terminalOnKey = (event: { domEvent: KeyboardEvent }) => {
  const { domEvent } = event;
  const ctrlKey = checkCtrlKeyAllSystem(domEvent); // 是否按了Ctrl键

  // 回车键
  if (domEvent.key === 'Enter') {
    console.log('Enter',commandBuffer);
      // 发送命令
      ws?.send(commandBuffer);
      commandBuffer = '';
  // Tab键
  } else if (domEvent.key === 'Tab') {
      // tab补全...
  // Backspace退格键或Delete删除键
  } else if (domEvent.key === 'Backspace' || domEvent.key === 'Delete') {
      // 在终端中模拟删除效果
      xterm?.write('\b \b');
  // 上箭头
  } else if (domEvent.key === 'ArrowUp') {
      // 历史指令...
  // 下箭头
  } else if (domEvent.key === 'ArrowDown') {
      // 历史指令...
  // Ctrl+C 快捷键
  } else if (ctrlKey && domEvent.key === 'c') {
      // ...
  // Ctrl+Shift+A 快捷键
  } else if (ctrlKey && domEvent.shiftKey && domEvent.key === 'a') {
      // ...
  // 其他按键与快捷键等处理...
  } else {
      xterm?.write(domEvent.key);
  }
};

// 判断是否按了Ctrl键 - 兼容Mac和Win不同系统
const checkCtrlKeyAllSystem = (domEvent: KeyboardEvent): boolean => {
  const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
  const isPressCtrlKey = isMac ? domEvent.metaKey : domEvent.ctrlKey;
  return isPressCtrlKey;
};

// 初始化 WebSocket 的方法
const initSocket = () => {
  ws = new WebSocket('ws://124.222.224.186:8800');
  ws.onopen = socketOnOpen;
  ws.onmessage = socketOnMessage;
  ws.onclose = socketOnClose;
  ws.onerror = socketOnError;
};

// WebSocket 连接成功的回调方法
const socketOnOpen = () => {
  xterm?.write('连接成功');
};

// 收到 WebSocket 消息的回调方法
const socketOnMessage = (event: MessageEvent) => {
  console.log('收到消息', event.data);
  xterm?.write(event.data);
};

// WebSocket 连接关闭的回调方法
const socketOnClose = () => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
      xterm?.write('\r\n\x1b[91mConnection closed');
  }
};

// WebSocket 连接出错的回调方法
const socketOnError = (err: Event) => {
  console.error('WebSocket connection error:', err);
};

// 组件挂载时初始化终端和 WebSocket
onMounted(() => {
  initTerminal();
  initSocket();
  // new WebSocket(wsUrl);
});

// 组件销毁前关闭 WebSocket 连接
onBeforeUnmount(() => {
  ws?.close();
});
</script>

<style scoped>
/* 可以在这里添加样式 */
</style>

