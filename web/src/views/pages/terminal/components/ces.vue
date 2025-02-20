



<template>
  <div>
    <el-tabs v-model="activeName" @click="handleClick" style="margin-left: 1vh; margin-right: 1vh">
      <el-tab-pane name="first" label="SSH">
        <div style="text-align: center">
          <el-form ref="formRef" :model="form" status-icon :rules="rules" label-position="left" label-width="80px" style="margin-left: 50vh; width: 50vh">
            <el-form-item label="Ip" prop="ip">
              <el-input v-model="form.ip"/>
            </el-form-item>
            <el-form-item label="Port" prop="port">
              <el-input v-model="form.port"/>
            </el-form-item>
            <el-form-item label="User" prop="user">
              <el-input v-model="form.user"/>
            </el-form-item>
            <el-form-item label="Password" prop="pwd">
              <el-input v-model="form.pwd" type="password"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitForm">连接</el-button>
              <el-button @click="resetForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>
      <el-tab-pane name="second" label="Terminal">
        <div class="ssh-container" ref="terminalRef"></div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount, watch } from 'vue';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import { debounce } from 'lodash';
import { useRoute } from 'vue-router';
import type { FormInstance, FormRules } from 'element-plus';

// 获取当前路由信息
const route = useRoute();

// 表单引用，类型为 ElementPlus 的表单实例
const formRef = ref<FormInstance>();
// 终端容器引用，类型为 HTML 元素或 null
const terminalRef = ref<HTMLElement | null>(null);

// 封装调整终端大小的消息
const packResize = (cols: number, rows: number): string =>
  JSON.stringify({
    type: 'resize',
    cols,
    rows
  });

// 表单验证规则
const rules: FormRules = {
  ip: [
    {
      validator: (rule: any, value: string, callback: (error?: Error) => void): void => {
        if (value === '') {
          callback(new Error('不能为空'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  port: [
    {
      validator: (rule: any, value: string, callback: (error?: Error) => void): void => {
        if (value === '') {
          callback(new Error('不能为空'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  user: [
    {
      validator: (rule: any, value: string, callback: (error?: Error) => void): void => {
        if (value === '') {
          callback(new Error('不能为空'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  pwd: [
    {
      validator: (rule: any, value: string, callback: (error?: Error) => void): void => {
        if (value === '') {
          callback(new Error('不能为空'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ]
};

// 当前激活的标签页名称
const activeName = ref<string>('second');
// 初始化文本
const initText = '连接中...\r\n';
// 是否为首次连接的标志
const first = ref<boolean>(true);
// 终端实例
const term = ref<Terminal | null>(null);
// 终端自适应插件实例
const fitAddon = ref<FitAddon | null>(null);
// WebSocket 实例
const ws = ref<WebSocket | null>(null);
// WebSocket 连接地址
const socketUrl = 'ws://192.168.31.106:8089/v1/ssh/open?Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImZ0TUlhcHV0IiwiZXhwIjoxNzQwMDE5OTA1fQ.yCGM6dj2RXuV-JF8fxzQD4HFZ9RGh6NL3vUhJZpNNDI';

// 表单数据
const form = ref<{
  user: string;
  pwd: string;
  ip: string;
  port: string;
}>({
  user: '',
  pwd: '',
  ip: route.query.ip as string || '',
  port: route.query.port as string || ''
});

// 终端配置选项
const option = {
  lineHeight: 1.0,
  cursorBlink: true,
  cursorStyle: 'block' as 'block' | 'underline' | 'bar',
  fontSize: 18,
  fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
  theme: {
    background: '#181d28'
  },
  cols: 30
};

// UTF-8 编码转 Base64 编码
const utf8_to_b64 = (rawString: string): string =>
  btoa(unescape(encodeURIComponent(rawString)));

// Base64 编码转 UTF-8 编码
const b64_to_utf8 = (encodeString: string): string =>
  decodeURIComponent(escape(atob(encodeString)));

// 字节数转换为可读格式
const bytesHuman = (bytes: number, precision: number = 1): string => {
  if (!/^([-+])?|(\.\d+)(\d+(\.\d+)?|(\d+\.)|Infinity)$/.test(bytes.toString())) {
    return '-';
  }
  if (bytes === 0) return '0';
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB', 'BB'];
  const num = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = (bytes / Math.pow(1024, Math.floor(num))).toFixed(precision);
  return `${value} ${units[num]}`;
};

// 检查 WebSocket 是否处于打开状态
const isWsOpen = (): boolean => ws.value && ws.value.readyState === WebSocket.OPEN ? true : false;

// 初始化终端
const initTerm = (): void => {
  term.value = new Terminal(option);
  fitAddon.value = new FitAddon();
  term.value.loadAddon(fitAddon.value);
  if (terminalRef.value) {
    term.value.open(terminalRef.value);
  }
  setTimeout(() => {
    fitAddon.value?.fit();
  }, 500);
};

// 处理终端按键输入
const onTerminalKeyPress = (): void => {
  term.value?.onData((data: string) => {
    if (isWsOpen()) {
      ws.value?.send(
        JSON.stringify({
          type: 'stdin',
          data: utf8_to_b64(data)
        })
      );
    }
  });
};

// 调整远程终端大小
const resizeRemoteTerminal = (): void => {
  if (term.value) {
    const { cols, rows } = term.value;
    console.log('列数、行数设置为：', cols, rows);
    if (isWsOpen()) {
      ws.value?.send(packResize(cols, rows));
    }
  }
};

// 窗口大小变化时的处理函数
const onResize = debounce(() => {
  fitAddon.value?.fit();
}, 500);

// 监听终端大小变化和窗口大小变化
const onTerminalResize = (): void => {
  window.addEventListener('resize', onResize);
  term.value?.onResize(resizeRemoteTerminal);
};

// 移除窗口大小变化的监听事件
const removeResizeListener = (): void => {
  window.removeEventListener('resize', onResize);
};

// 初始化 WebSocket 连接
const initSocket = (): void => {
  if (term.value) {
    term.value.write(initText);
  }
  ws.value = new WebSocket(socketUrl, ['webssh']);
  onOpenSocket();
  onCloseSocket();
  onErrorSocket();
  if(term.value){
    // term.value._initialized = true;
  } 
  onMessageSocket();
};

// WebSocket 连接打开时的处理函数
const onOpenSocket = (): void => {
  if (ws.value) {
    ws.value.onopen = () => {
      console.log('websocket 已连接');
      ws.value?.send(
        JSON.stringify({ type: 'addr', data: utf8_to_b64(form.value.ip + ':' + form.value.port) })
      );
      ws.value?.send(
        JSON.stringify({ type: 'login', data: utf8_to_b64(form.value.user) })
      );
      ws.value?.send(
        JSON.stringify({ type: 'password', data: utf8_to_b64(form.value.pwd) })
      );
      term.value?.reset();
      setTimeout(() => {
        resizeRemoteTerminal();
      }, 500);
    };
  }
};

// WebSocket 连接关闭时的处理函数
const onCloseSocket = (): void => {
  if (ws.value) {
    ws.value.onclose = () => {
      console.log('关闭连接');
      term.value?.write('未连接， 刷新后重连...\r\n');
    };
  }
};

// WebSocket 连接出错时的处理函数
const onErrorSocket = (): void => {
  if (ws.value) {
    ws.value.onerror = () => {
      term.value?.write('连接失败，请刷新！');
    };
  }
};

// 处理 WebSocket 接收到的消息
const onMessageSocket = (): void => {
  if (ws.value) {
    ws.value.onmessage = (res: MessageEvent) => {
      console.log(res);
      const msg = JSON.parse(res.data);
      const termInstance = term.value;
      if (first.value) {
        first.value = false;
        termInstance?.reset();
        termInstance?.element && termInstance.focus();
        resizeRemoteTerminal();
      }
      termInstance?.write(b64_to_utf8(msg.data));
    };
  }
};

// 初始化 WebSocket 和终端
const initWs = (): void => {
  initTerm();
  initSocket();
  onTerminalResize();
  onTerminalKeyPress();
};

// 提交表单
const submitForm = (): void => {
  formRef.value?.validate((valid: boolean) => {
    if (valid) {
      activeName.value = 'second';
      initWs();
    } else {
      console.log('error submit!!');
    }
  });
};

// 重置表单
const resetForm = (): void => {
  formRef.value?.resetFields();
};

// 处理标签页点击事件
const handleClick = (tab: { name: string }, event: Event): void => {
  if (tab.name === 'second') {
    // 可添加相应逻辑
  }
};

// 组件卸载前的清理工作
onBeforeUnmount(() => {
  removeResizeListener();
  term.value?.dispose();
});
initWs();
</script>

<style >
body {
  margin: 0;
  padding: 0;
}

.ssh-container {
  overflow: hidden;
  height: 85vh;
  border-radius: 4px;
  background: rgb(24, 29, 40);
  padding: 0px;
  color: rgb(255, 255, 255);

  .xterm-scroll-area::-webkit-scrollbar-thumb {
    background-color: #b7c4d1;
    /* 滚动条的背景颜色 */
  }
}
</style>