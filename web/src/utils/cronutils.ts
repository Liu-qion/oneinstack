// src/utils/cronUtils.ts

export default function formatCron(cronStr: string): string {
  if (!cronStr) return '';
  const cronExpressions = cronStr.split(',');
  const descriptions: string[] = [];

  cronExpressions.forEach((expression) => {
      const parts = expression.split(' ');
      if (parts.length!== 5) {
          descriptions.push(`无效的 cron 表达式: ${expression}`);
          return;
      }

      const [minute, hour, dayOfMonth, month, dayOfWeek] = parts;
      let description = '';

      // 处理月份部分
      if (month === '*') {
          description += '每月';
      } else if (/^\d+$/.test(month)) {
          description += `每年 ${month} 月`;
      } else {
          descriptions.push(`无效的月份部分: ${month}`);
          return;
      }

      // 处理日期部分
      if (dayOfMonth === '*') {
          if (dayOfWeek === '*') {
              description += '每日';
          }
      } else if (/^\d+$/.test(dayOfMonth)) {
          description += ` ${dayOfMonth} 日`;
      } else {
          descriptions.push(`无效的日期部分: ${dayOfMonth}`);
          return;
      }

      // 处理星期部分
      if (dayOfWeek!== '*' && /^\d+$/.test(dayOfWeek)) {
          const dayOfWeekMap = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
          description += ` 每周 ${dayOfWeekMap[parseInt(dayOfWeek)]}`;
      }

      // 处理小时和分钟部分
      if (hour === '*' && minute === '*') {
          description += ' 整时整分';
      } else if (hour === '*') {
          description += ` 每分钟的第 ${minute} 分`;
      } else if (minute === '*') {
          description += ` 每小时的第 ${hour} 时`;
      } else if (/^\d+$/.test(hour) && /^\d+$/.test(minute)) {
          description += ` ${hour}:${String(minute).padStart(2, '0')}分`;
      } else {
          descriptions.push(`无效的时间部分: ${hour}:${minute}`);
          return;
      }

      descriptions.push(description.trim());
  });

  return descriptions.join('<br>');
}