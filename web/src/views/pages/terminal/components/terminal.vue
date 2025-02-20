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
// 添加一个标志位，用于记录 WebSocket 是否已经连接
let isWebSocketConnected = false;
let message: any[] =[]
let alldata: any[] =[]
let newarr :any[] = [] //每次scoket后返回的新数组
let scoketobj = {
  is_one_num:0, //用于判发送的次数
  newtype:false, //用于判断newarr里是否是要删除掉换行前的内容，如果为false则不删，如果为true则删除
  is_number:0, //用于判断n换行出现的次数，如果为2则删除
  colnum:0, //用于记录第二个“换行符”所在的行数
}
// 处理终端输入数据的方法
const terminalOnData = (data: string) => {
  // 确保输入数据中的换行符能正常显示
  const lines = data.split(/\r?\n/);
  lines.forEach((line, index) => {
    if (index > 0) {
      xterm?.write('\r\n');
    }
    xterm?.write(line);
  });
  commandBuffer += data;
};
// 处理终端按键事件的方法
const terminalOnKey = (event: { domEvent: KeyboardEvent }) => {
  const { domEvent } = event;
  const ctrlKey = checkCtrlKeyAllSystem(domEvent); // 是否按了Ctrl键
  ws?.send(commandBuffer);
  return 
  // 回车键
  if (domEvent.key === 'Enter') {
    scoketobj.is_one_num++
    scoketobj.colnum  =0
    scoketobj.is_number =0
    console.log('Enter', commandBuffer);
    newarr = []
    // 发送命令
    ws?.send(commandBuffer);
    // 发送命令后添加换行符
    // xterm?.write('\r\n'); 
    commandBuffer = '';
    // 可以在这里添加输入提示，例如 $ 
    // xterm?.write('$ '); 
  // Tab键
  } else if (domEvent.key === 'Tab') {
      // tab补全...
  // Backspace退格键或Delete删除键
  } else if (domEvent.key === 'Backspace' || domEvent.key === 'Delete') {
      // 在终端中模拟删除效果
      // xterm?.write('\b \b');
      if (commandBuffer.length > 0) {
        commandBuffer = commandBuffer.slice(0, -1);
      }
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
    console.log('其他按键', domEvent.key);
      // xterm?.write(domEvent.key);
      // commandBuffer += domEvent.key;
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
  // 只有当 WebSocket 未连接时，才进行连接操作
  if (!isWebSocketConnected) {
    ws = new WebSocket(wsUrl);
    ws.onopen = socketOnOpen;
    ws.onmessage = socketOnMessage;
    ws.onclose = socketOnClose;
    ws.onerror = socketOnError;
    isWebSocketConnected = true;
  }
};
// WebSocket 连接成功的回调方法
const socketOnOpen = () => {
  xterm?.write('连接成功\r\n');
  // 连接成功后添加输入提示
  xterm?.write('$ '); 
};

// 收到 WebSocket 消息的回调方法
const socketOnMessage = async(event: MessageEvent) => {
    const textContent = ref<string>('');
    textContent.value = event.data as string;
    console.log('收到消息:', textContent.value);
        // console.log(textContent.value,scoketobj.is_number)
        // if (scoketobj.is_one_num > 1) {
        //   if (textContent.value.indexOf('[?2004l') !== -1) {
        //     scoketobj.is_number++
        //   }
        //   if (scoketobj.is_number>=2){
        //     xterm?.write(textContent.value+'');
        //   }
        // }else{
        //   xterm?.write(textContent.value);
        // }
        xterm?.write(textContent.value);
};

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
      convertEol: true, // 确保换行符能正确转换
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

// WebSocket 连接关闭的回调方法
const socketOnClose = () => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
      xterm?.write('\r\n\x1b[91m连接已关闭，请刷新重试！');
  }
  // 当连接关闭时，将标志位重置为 false
  isWebSocketConnected = false;
};

// WebSocket 连接出错的回调方法
const socketOnError = (err: Event) => {
  console.error('WebSocket connection error:', err);
  // 当连接出错时，将标志位重置为 false
  isWebSocketConnected = false;
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