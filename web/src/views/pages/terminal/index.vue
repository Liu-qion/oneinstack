<template>
  <div id="terminal">
    <div class="output" ref="outputRef">
      <pre v-for="(line, index) in outputLines" :key="index">{{ line }}</pre>
    </div>
    <div class="input">
      <span>{{ prompt }}</span>
      <input
        type="text"
        v-model="inputValue"
        @keyup.enter="executeCommand"
        ref="inputRef"
        autofocus
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// 输出区域的引用
const outputRef = ref<HTMLDivElement | null>(null);
// 输入框的引用
const inputRef = ref<HTMLInputElement | null>(null);
// 命令输入值
const inputValue = ref('');
// 终端提示符
const prompt = ref('$ ');
// 输出的行数组
const outputLines = ref<string[]>([]);

// 执行命令的函数
const executeCommand = () => {
  const command = inputValue.value.trim();
  if (command) {
    // 将用户输入的命令添加到输出中
    outputLines.value.push(`${prompt.value}${command}`);

    // 简单的命令处理示例
    let output: string;
    switch (command) {
      case 'help':
        output = '支持的命令: help, hello';
        break;
      case 'hello':
        output = 'Hello, World!';
        break;
      default:
        output = `未知命令: ${command}`;
    }

    // 将命令执行结果添加到输出中
    outputLines.value.push(output);

    // 清空输入框
    inputValue.value = '';

    // 滚动到输出区域的底部
    if (outputRef.value) {
      outputRef.value.scrollTop = outputRef.value.scrollHeight;
    }
  }
};
</script>

<style scoped>
#terminal {
  font-family: 'Courier New', Courier, monospace;
  background-color: #000;
  color: #fff;
  padding: 10px;
  height: 300px;
  overflow-y: auto;
}

.output {
  margin-bottom: 10px;
}

.input {
  display: flex;
  align-items: center;
}

.input span {
  margin-right: 5px;
}

.input input {
  background-color: #000;
  color: #fff;
  border: none;
  outline: none;
  flex: 1;
}
</style>