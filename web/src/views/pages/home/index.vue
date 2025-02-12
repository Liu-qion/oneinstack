<script setup lang="ts">
import memo from './components/memo.vue'
import { markRaw, onMounted, reactive } from 'vue'
import sapp from '@/sstore/sapp'
import { Api } from '@/api/Api'
import sutil from '@/sstore/sutil'
import { Scope } from 'tools-vue3'
import { EChartsOption } from 'echarts'
import * as echarts from 'echarts'
import basicChart from '@/components/echarts/basic-chart.vue'
import { ElMessage } from 'element-plus'
import System from '@/utils/System'

type MonitorType = 'network' | 'disk'

interface Options {
  label: string
  value: string | number | object
}

interface ChartData {
  times: string[]
  ascend: { value: number; strValue: string }[]
  descend: { value: number; strValue: string }[]
}

const conf = reactive({
  themeColor: {
    light: ['#F7911C'],
    dark: ['#EAB170']
  },
  category: [
    {
      name: '网站-全部',
      icon: 'home-website',
      value: 0,
      linkFn: () => System.router.push('/website')
    },
    {
      name: '数据-全部',
      icon: 'home-data',
      value: 0,
      linkFn: () => System.router.push('/database')
    },
    {
      name: '安全风险',
      icon: 'home-software',
      value: 0
    },
    {
      name: '备忘录',
      icon: 'home-mome',
      value: '当前内容为空，点击编辑',
      linkFn: () => conf.memo.open()
    }
  ],
  getSysCount: async () => {
    const { data: wbsiteCount } = await Api.getWebsiteCount()
    const { data: databaseCount } = await Api.getDatabaseCount()
    conf.category[0].value = wbsiteCount
    conf.category[1].value = databaseCount
  },
  monitorData: {
    selectedType: 'network' as MonitorType,
    selectedCard: {
      label: 'all',
      value: 0
    },
    network: [
      {
        label: '上行',
        value: '--'
      },
      {
        label: '下行',
        value: '--'
      },
      {
        label: '总发送',
        value: '--'
      },
      {
        label: '总接收',
        value: '--'
      }
    ] as Options[],
    disk: [
      {
        label: '读取',
        value: '--'
      },
      {
        label: '写入',
        value: '--'
      },
      {
        label: '读取',
        value: '--'
      },
      {
        label: '写入',
        value: '--'
      }
    ] as Options[],
    options: [] as Options[],
    handleChangeType: (type: MonitorType) => {
      conf.monitorData.selectedType = type
      conf.monitorData.chartData = {
        times: [],
        ascend: [],
        descend: []
      }
      conf.monitorData.update(true)
    },
    handleChangeCard: () => conf.monitorData.update(),
    chartData: {
      times: [],
      ascend: [],
      descend: []
    } as ChartData,
    chartOptions: null as EChartsOption | null,
    draw: ({ ascend, descend }: { ascend: ChartData['ascend'][0]; descend: ChartData['descend'][0] }) => {
      conf.monitorData.chartData.times.push(new Date().toLocaleTimeString())
      conf.monitorData.chartData.ascend.push(ascend)
      conf.monitorData.chartData.descend.push(descend)
      conf.monitorData.chartOptions = {
        color: ['#FFAD79', '#79D1FF'],
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          },
          formatter: (params: any) => {
            const { name, marker: marker1, data: data1, seriesName: seriesName1 } = params[0]
            const { marker: marker2, data: data2, seriesName: seriesName2 } = params[1]
            return `${name} </br> ${seriesName1} ${marker1} ${data1.strValue}/s </br> ${seriesName2} ${marker2} ${data2.strValue}/s`
          }
        },
        grid: {
          left: 0,
          right: '4%',
          top: '5%',
          bottom: '5% ',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            boundaryGap: false,
            data: conf.monitorData.chartData.times,
            axisLabel: {
              margin: 20,
              color: sapp.theme == 'light' ? '#A2A2A2' : '#A2A2A2'
            },
            axisLine: {
              show: false
            },
            axisTick: {
              show: false
            }
          }
        ],
        yAxis: [
          {
            type: 'value',
            axisLabel: {
              margin: 20,
              color: sapp.theme == 'light' ? '#A2A2A2' : '#A2A2A2'
            },
            splitLine: {
              lineStyle: {
                color: sapp.theme == 'light' ? '#D6D6D699' : '#435B7199',
                type: [5]
              }
            }
          }
        ],
        series: [
          {
            name: conf.monitorData.selectedType == 'network' ? '上行' : '读取',
            type: 'line',
            smooth: false,
            lineStyle: {
              width: 1,
              color: '#FFAD79'
            },
            showSymbol: false,
            areaStyle: {
              opacity: 0.5,
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#FFAD79' },
                { offset: 0.5, color: 'rgba(255,173,121,0.3)' },
                { offset: 1, color: 'rgba(255,173,121,0.1)' }
              ])
            },
            emphasis: {
              focus: 'series'
            },
            data: conf.monitorData.chartData.ascend
          },
          {
            name: conf.monitorData.selectedType == 'network' ? '下行' : '写入',
            type: 'line',
            smooth: false,
            lineStyle: {
              width: 1,
              color: '#79D1FF'
            },
            showSymbol: false,
            areaStyle: {
              opacity: 0.5,
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#79D1FF' },
                { offset: 0.5, color: 'rgba(121,240,255,0.3)' },
                { offset: 1, color: 'rgba(121,240,255,0.1)' }
              ])
            },
            emphasis: {
              focus: 'series'
            },
            data: conf.monitorData.chartData.descend
          }
        ]
      }
    },
    update: async (isCardChange = false) => {
      const { data: res } = await Api.getSysMonitor()
      if (!conf.monitorData.options.length || isCardChange) {
        conf.monitorData.options = (res[conf.monitorData.selectedType] as any[])
          .map((item: any, i: number) => {
            const option = { label: item.Name, value: i }
            if (item.Name == 'all') conf.monitorData.selectedCard = option
            return option
          })
          .sort((a, b) => a.label.localeCompare(b.label))
      }
      let ascend: ChartData['ascend'][0], descend: ChartData['descend'][0]
      switch (conf.monitorData.selectedType) {
        case 'network':
          {
            const { SendRate, RecvRate, BytesSent, BytesRecv } = res.network.find(
              (item: any) => item.Name == conf.monitorData.selectedCard.label
            )
            ascend = {
              value: parseFloat((SendRate / 1024 / 1024).toFixed(2)),
              strValue: sutil.bytesTransform(SendRate).strValue
            }
            descend = {
              value: parseFloat((RecvRate / 1024 / 1024).toFixed(2)),
              strValue: sutil.bytesTransform(RecvRate).strValue
            }
            conf.monitorData.network = [
              { label: '上行', value: `${sutil.bytesTransform(SendRate).strValue}/s` },
              { label: '下行', value: `${sutil.bytesTransform(RecvRate).strValue}/s` },
              { label: '总发送', value: sutil.bytesTransform(BytesSent).strValue },
              { label: '总接收', value: sutil.bytesTransform(BytesRecv).strValue }
            ]
          }
          break
        case 'disk':
          {
            const { ReadSpeed, WriteSpeed, ReadOpsPerSec, WriteOpsPerSec, AvgIoLatency } = res.disk.find(
              (item: any) => item.Name == conf.monitorData.selectedCard.label
            )
            ascend = {
              value: parseFloat((ReadSpeed / 1024 / 1024).toFixed(2)),
              strValue: sutil.bytesTransform(ReadSpeed).strValue
            }
            descend = {
              value: parseFloat((WriteSpeed / 1024 / 1024).toFixed(2)),
              strValue: sutil.bytesTransform(WriteSpeed).strValue
            }
            conf.monitorData.disk = [
              { label: '读取', value: sutil.bytesTransform(ReadSpeed).strValue },
              { label: '写入', value: sutil.bytesTransform(WriteSpeed).strValue },
              { label: '读写次数', value: `${ReadOpsPerSec + WriteOpsPerSec}次/s` },
              { label: '平均延迟', value: `${AvgIoLatency.toFixed(2)}ms` }
            ]
          }
          break
      }
      conf.monitorData.draw({ ascend, descend })
    }
  },
  statusData: {
    chartOptions: null as EChartsOption | null,
    draw: () => {
      const gaugeData = [
        {
          value: conf.statusData.usage.usedPercent,
          name: conf.statusData.selected.label
        }
      ]
      conf.statusData.chartOptions = {
        series: [
          {
            type: 'gauge',
            startAngle: 90,
            endAngle: -270,
            radius: '100%',
            zlevel: 2,
            pointer: {
              show: false
            },
            progress: {
              show: true,
              overlap: false,
              roundCap: true,
              clip: false,
              itemStyle: {
                borderWidth: 0,
                color: conf.themeColor[sapp.theme][0]
              }
            },
            axisLine: {
              lineStyle: {
                width: 20,
                color: [[1, `rgba(${sutil.getCssVariable('--category-item-bg-color')}, 0.88)`]]
              }
            },
            splitLine: {
              show: false
            },
            axisTick: {
              show: false
            },
            axisLabel: {
              show: false,
              distance: 50
            },
            data: gaugeData,
            title: {
              fontSize: 14,
              color: sutil.getCssVariable('--font-color-gray-light'),
              offsetCenter: [0, '30%']
            },
            detail: {
              width: 200,
              height: 200,
              borderWidth: 0,
              offsetCenter: [0, '60%'],
              rich: {
                a: {
                  fontSize: 24,
                  color: sutil.getCssVariable('--font-color-black'),
                  fontWeight: 500
                },
                b: {
                  fontSize: 18,
                  color: sutil.getCssVariable('--font-color-black'),
                  fontWeight: 500,
                  padding: [5, 0, 0, 0]
                }
              },
              formatter: (value) => {
                const int = value.toFixed(2).split('.')[0]
                const flt = value.toFixed(2).split('.')[1]
                return `{a|${int}}{b|.${flt}%}`
              },
              valueAnimation: true
            }
          },
          {
            animationType: 'scale',
            color: `rgba(${sutil.getCssVariable('--category-item-bg-color')}, 0.88)`,
            name: 'Access From',
            type: 'pie',
            radius: '75%',
            label: {
              show: false
            },
            labelLine: {
              show: false
            },
            data: [{ value: 100, name: 'Search Engine' }]
          }
        ]
      }
    },
    selected: {
      value: 1,
      label: '内存'
    },
    options: markRaw([
      {
        value: 1,
        label: '内存'
      },
      {
        value: 2,
        label: '磁盘'
      },
      {
        value: 3,
        label: 'CPU'
      }
    ]),
    usage: {
      total: '--',
      available: '--',
      used: '--',
      usedPercent: 0
    } as { total: string; available: string; used: string; usedPercent: number },
    cpuInfo: '',
    handleStatusChange: () => conf.statusData.update(),
    update: async () => {
      const { data: res } = await Api.getSysinfo()
      switch (conf.statusData.selected.value) {
        case 1:
          {
            const { total, used, available, usedPercent } = res.memory_usage
            conf.statusData.usage.total = sutil.bytesTransform(total).strValue
            conf.statusData.usage.used = sutil.bytesTransform(used).strValue
            conf.statusData.usage.available = sutil.bytesTransform(available).strValue
            conf.statusData.usage.usedPercent = usedPercent
          }
          break
        case 2:
          {
            const { total, free, used, usedPercent } = res.disk_usage[0]
            conf.statusData.usage.total = sutil.bytesTransform(total).strValue
            conf.statusData.usage.used = sutil.bytesTransform(used).strValue
            conf.statusData.usage.available = sutil.bytesTransform(free).strValue
            conf.statusData.usage.usedPercent = usedPercent
          }
          break
        case 3:
          {
            const [usedPercent] = res.cpu_usage
            const { modelName, cores } = res.cpu_info[0]
            conf.statusData.usage.total = '--'
            conf.statusData.usage.used = '--'
            conf.statusData.usage.available = '--'
            conf.statusData.usage.usedPercent = usedPercent
            conf.statusData.cpuInfo = `${modelName} ${cores}核`
          }
          break
      }
      conf.statusData.draw()
    }
  },
  memo: {
    data: {
      id: 1,
      content: ''
    },
    show: false,
    open: async () => {
      await conf.memo.getData()
      conf.memo.show = true
    },
    close: () => (conf.memo.show = false),
    getData: async () => {
      const { data: res } = await Api.getSysRemark()
      conf.memo.data = res
      conf.category[3].value = res.content
    },
    update: async () => {
      await Api.updateSysRemark(conf.memo.data)
      ElMessage.success('保存成功')
      conf.memo.getData()
      conf.memo.show = false
    }
  }
})

const timer = Scope.Timer()
onMounted(() => {
  timer.on(
    () => {
      conf.statusData.update()
      conf.monitorData.update()
    },
    5000,
    true
  )
  conf.monitorData.draw({ ascend: { value: 0, strValue: '0' }, descend: { value: 0, strValue: '0' } })
  conf.getSysCount()
  conf.memo.getData()
})
</script>

<template>
  <div class="home-container">
    <div class="column fit-height fit-width">
      <div class="col column fit-width">
        <div class="col relative fit-width">
          <div class="absolute fit-height fit-width flex column no-wrap" style="gap: 24px">
            <div class="basic-card">
              <div class="basic-card__title">概览</div>
            </div>
            <el-row :gutter="20">
              <el-col v-for="item in conf.category" :lg="6" :md="12" :sm="24">
                <div class="category-item" @click="item.linkFn?.()">
                  <div class="icon">
                    <v-s-icon :name="item.icon" size="52" :color="conf.themeColor[sapp.theme]" />
                  </div>
                  <div class="text column items-center">
                    <span :class="{ 'link': typeof item.value === 'string' }">
                      {{ item.value }}
                    </span>
                    <div class="name">{{ item.name }}</div>
                  </div>
                </div>
              </el-col>
            </el-row>
            <el-row :gutter="24" style="flex: 1">
              <el-col :lg="16" :md="16" :sm="24">
                <div class="basic-card flex column no-wrap fit-height">
                  <div class="basic-card__header">
                    <div class="basic-card__title">监控</div>
                    <div class="miscellaneous">
                      <div class="switch">
                        <span>{{ conf.monitorData.selectedType == 'network' ? '网卡' : '磁盘' }}</span>
                        <el-select
                          v-model="conf.monitorData.selectedCard"
                          placeholder="请选择"
                          style="width: 100%"
                          @change="conf.monitorData.handleChangeCard"
                        >
                          <el-option
                            v-for="item in conf.monitorData.options"
                            :key="item.label"
                            :label="item.label"
                            :value="item"
                          />
                        </el-select>
                      </div>
                      <div class="menu">
                        <div
                          class="item"
                          :class="conf.monitorData.selectedType == 'network' ? 'active' : ''"
                          @click="conf.monitorData.handleChangeType('network')"
                        >
                          流量
                        </div>
                        <div
                          class="item"
                          :class="conf.monitorData.selectedType == 'disk' ? 'active' : ''"
                          @click="conf.monitorData.handleChangeType('disk')"
                        >
                          磁盘IO
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="basic-card__body flex column no-wrap fit-height">
                    <div class="flow">
                      <el-space class="lefts" :size="20" spacer="|">
                        <div
                          v-for="item in conf.monitorData[conf.monitorData.selectedType]"
                          :key="item.label"
                          class="item"
                        >
                          <span class="label">{{ item.label }}：</span>
                          <span class="value">{{ item.value }}</span>
                        </div>
                      </el-space>
                      <div class="rights">
                        <div class="upper">
                          <div class="yuan"></div>
                          <span>{{ conf.monitorData.selectedType == 'network' ? '上行' : '读取' }}</span>
                        </div>
                        <div class="below">
                          <div class="yuan"></div>
                          <span>{{ conf.monitorData.selectedType == 'network' ? '下行' : '写入' }}</span>
                        </div>
                      </div>
                    </div>
                    <basic-chart :option="conf.monitorData.chartOptions as EChartsOption" class="chart-box" />
                  </div>
                </div>
              </el-col>
              <el-col :lg="8" :md="8" :sm="24">
                <div ref="statusCard" class="basic-card flex column no-wrap fit-height">
                  <div class="basic-card__header">
                    <div class="basic-card__title">状态</div>
                    <div class="status-right">
                      <el-select
                        v-model="conf.statusData.selected"
                        placeholder="选择状态"
                        style="width: 100px; margin-right: 10px"
                        @change="conf.statusData.handleStatusChange"
                      >
                        <el-option
                          v-for="item in conf.statusData.options"
                          :key="item.value"
                          :label="item.label"
                          :value="item"
                        />
                      </el-select>
                    </div>
                  </div>
                  <div class="basic-card__body" style="flex: 1">
                    <div class="norule">
                      <basic-chart
                        :option="conf.statusData.chartOptions as EChartsOption"
                        style="width: 256px; height: 256px"
                      />
                      <div class="status-title">
                        <span v-if="conf.statusData.selected.value !== 3">
                          {{ conf.statusData.usage.used }} / {{ conf.statusData.usage.total }}
                        </span>
                        <span v-else>{{ conf.statusData.cpuInfo }}</span>
                      </div>
                      <div class="status-menu" style="margin-bottom: 20px">
                        <div class="b1">
                          总数：
                          <span>{{ conf.statusData.usage.total }}</span>
                        </div>
                        <div class="b2">
                          已用：
                          <span>{{ conf.statusData.usage.used }}</span>
                        </div>
                      </div>
                      <div class="status-menu">
                        <div class="b1">
                          可用：
                          <span>{{ conf.statusData.usage.available }}</span>
                        </div>
                        <div class="b2">
                          使用率：
                          <span>{{ conf.statusData.usage.usedPercent.toFixed(2) }} %</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>
            <div class="basic-card" style="background: transparent" />
          </div>
        </div>
      </div>
    </div>
    <memo :show="conf.memo.show" :memo="conf.memo.data" :close="conf.memo.close" :update="conf.memo.update" />
  </div>
</template>

<style scoped lang="less">
@media screen and (max-width: 1200px) {
  .category-item {
    margin-bottom: 16px;
  }
}

.home-container {
  width: 100%;
  height: 100%;
  padding-bottom: 50px;

  .basic-card {
    width: 100%;
    background: rgb(var(--bg-card-color));
    border-radius: 16px;
    padding: 21px 46px;

    &__title {
      display: flex;
      align-items: center;
      position: relative;
      font-family: PingFang SC;
      font-weight: 500;
      font-size: 16px;
      color: var(--font-color-black);

      &::before {
        content: '';
        background: var(--el-color-primary);
        width: 5px;
        height: 22px;
        margin-right: 18px;
      }
    }

    &__header {
      padding: 17px 0 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: 2px solid rgb(var(--border-color-gray-light));
    }

    &__body {
      padding: 24px 0;
    }
  }

  .category-item {
    padding: 36px 32px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    background: rgb(var(--bg-card-color));
    border-radius: 16px;
    border: 2px solid transparent;
    transition: border-color 0.3s;

    &:hover {
      border-color: rgba(var(--primary-color), 0.88);

      .icon {
        border-color: rgba(var(--primary-color), 0.63);
      }
    }

    .icon {
      width: 102px;
      height: 102px;
      background: rgba(var(--category-item-bg-color), 0.88);
      border-radius: 50%;
      display: flex;
      justify-content: center;
      align-items: center;
      flex: 0 0 102px;
      border: 3px solid transparent;
    }

    .text {
      width: 100%;
      gap: 22px;

      span {
        font-weight: 500;
        font-size: 24px;
        color: var(--font-color-black);
        text-align: center;
        display: inline-block;
        width: 120px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.link {
          font-weight: 400;
          font-size: 16px;
          cursor: pointer;
          text-decoration: underline;
        }
      }

      .name {
        font-size: 16px;
        color: var(--font-color-gray-light);
      }
    }
  }

  .miscellaneous {
    height: 36px;
    display: flex;
    flex-direction: row;
    align-items: center;

    .switch {
      width: 202px;
      height: 100%;
      border-radius: 2px;
      border: 0.4px solid rgb(var(--border-color-gray));
      margin-right: 14px;
      border-radius: 2px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      span {
        flex: 0 0 45px;
        color: var(--font-color-gray-light);
        border-right: 1px solid rgba(var(--border-color-gray), 0.5);
        padding: 0 5px 0 10px;
        margin-right: 10px;
      }
    }

    .menu {
      height: 100%;
      display: flex;
      color: var(--font-color-black);
      border-radius: 2px;
      border: 0.4px solid rgb(var(--border-color-gray));
      overflow: hidden;

      .item {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 5px 10px;
        cursor: pointer;
      }

      .active {
        background: var(--el-color-primary);
        border: 1px solid var(--el-color-primary);
        color: var(--font-color-white);
      }
    }
  }

  .flow {
    width: 100%;
    display: flex;
    justify-content: space-between;

    .lefts {
      color: var(--font-color-gray);

      .value {
        color: var(--el-color-primary);
      }
    }

    .rights {
      display: flex;
      flex-direction: row;
      align-items: center;

      .upper,
      .below {
        display: flex;
        flex-direction: row;
        align-items: center;
        color: var(--font-color-gray-light);
        font-size: 14px;
        margin-right: 10px;

        .yuan {
          width: 10px;
          height: 10px;
          border-radius: 50%;
          background: rgb(var(--primary-color));
          margin-right: 10px;
        }
      }

      .below .yuan {
        background: rgb(var(--blue-color));
        margin-left: 5px;
      }
    }
  }

  .chart-box {
    margin-top: 24px;
    width: 100%;
    height: 100%;
  }

  .status-right {
    width: 97px;
    height: 36px;
    border-radius: 2px;
    border: 0.4px solid rgb(var(--border-color-gray));
  }

  .norule {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    .status-title {
      color: var(--font-color-black);
      margin-top: 26px;
      font-size: 18px;
      margin-bottom: 40px;
    }

    .status-menu {
      display: flex;

      .b1 {
        font-size: 14px;
        margin-right: 30px;
        padding-right: 20px;
        border-right: 1px solid rgba(var(--border-color-gray), 0.5);
        color: var(--font-color-gray);

        span {
          color: var(--el-color-primary);
        }
      }

      .b2 {
        font-size: 14px;
        color: var(--font-color-gray);

        span {
          color: var(--el-color-primary);
        }
      }
    }
  }

  :deep(.el-select) {
    &__placeholder {
      color: var(--font-color-black);
    }

    &__wrapper {
      border: none;
      background-color: transparent;
      box-shadow: none;
    }
  }
}
</style>
