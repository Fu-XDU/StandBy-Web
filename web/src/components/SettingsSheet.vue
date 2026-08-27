<template>
  <div class="settings-sheet" :class="{ visible: sheetVisible }">
    <div class="sheet-header">
      <button class="close-btn" @click="$emit('close')" aria-label="关闭">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
             stroke-linecap="round">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
      <span class="sheet-title">设置</span>
    </div>

    <div class="sheet-body">
      <div class="setting-row">
        <span class="setting-label">定时夜间模式</span>
        <el-switch
            v-model="autoNightMode"
            style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
        />
      </div>
      <div v-if="autoNightMode" class="setting-sub">
        <span class="sub-label">时间段</span>
        <el-time-picker v-model="nightModeRange[0]" size="small" arrow-control placeholder="开始"
                        class="time-picker" :clearable="false"/>
        <span class="time-sep">—</span>
        <el-time-picker v-model="nightModeRange[1]" size="small" arrow-control placeholder="结束"
                        class="time-picker" :clearable="false"/>
      </div>

      <div class="divider"/>

      <div class="setting-row">
        <span class="setting-label">定时关闭显示</span>
        <el-switch
            v-model="autoInvisible"
            style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
        />
      </div>
      <div v-if="autoInvisible" class="setting-sub">
        <span class="sub-label">时间段</span>
        <el-time-picker v-model="invisibleRange[0]" size="small" placeholder="开始"
                        class="time-picker" :clearable="false"/>
        <span class="time-sep">—</span>
        <el-time-picker v-model="invisibleRange[1]" size="small" placeholder="结束"
                        class="time-picker" :clearable="false"/>
      </div>

      <div class="divider"/>

      <div class="setting-row">
        <span class="setting-label">全天关闭显示</span>
        <el-switch
            v-model="invisibleDayEnable"
            style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
        />
      </div>
      <div v-if="invisibleDayEnable" class="setting-sub weekday-row">
        <el-checkbox-button
            v-for="(_, index) in invisibleDay"
            :key="index"
            v-model="invisibleDay[index]"
            :label="'周' + weekLabels[index]"
            size="small"
        />
      </div>

      <div class="divider"/>

      <div class="setting-row">
        <span class="setting-label">股票行情展示（数字时钟）</span>
        <span class="stock-count-hint">{{ (numericalStocks || []).length }}/3</span>
      </div>
      <div class="setting-sub stock-setting-sub">
        <div v-if="stockLoadError" class="stock-error-box">
          <span class="stock-error-text">{{ stockLoadError }}</span>
          <button class="retry-btn" type="button" @click="loadStockOptions">重试</button>
        </div>
        <el-select
          v-model="numericalStocks"
          multiple
          :multiple-limit="3"
          placeholder="选择要展示的股票/行情 (最多3项)"
          size="small"
          class="stock-select"
          :loading="isLoadingStocks"
          clearable
        >
          <el-option
            v-for="item in stockOptionsWithSelected"
            :key="item.value"
            :label="item.fullTitle || item.name"
            :value="item.value"
          >
            <div class="stock-option-row">
              <span class="option-name">{{ item.name }}</span>
              <span class="option-price">{{ item.price }}</span>
            </div>
          </el-option>
        </el-select>
        <p class="stock-hint">
          选择 0 项时不发起网络请求以节电；选择 1-3 项将在数字时钟日期下方展示。
        </p>
      </div>

      <div class="divider"/>

      <div class="setting-row">
        <span class="setting-label">远程控制</span>
        <el-switch
            v-model="remoteControlEnabled"
            style="--el-switch-on-color: #ff4949; --el-switch-off-color: #4C4D4F"
        />
      </div>

      <div v-if="remoteControlEnabled" class="device-section">
        <div class="setting-row device-row">
          <span class="setting-label">设备 ID</span>
          <div class="device-id-actions">
            <code class="device-id">{{ deviceId }}</code>
            <button class="reset-btn" type="button" @click="confirmResetDeviceId">重置</button>
          </div>
        </div>
        <p class="device-hint">
          在其他设备打开下方链接修改样式或设置，可在此设备同步生效：
        </p>
        <button
            class="device-link"
            type="button"
            @click="copyRemoteUrl"
        >{{ remoteControlUrl }}</button>
        <p v-if="copyHint" class="copy-hint">{{ copyHint }}</p>
      </div>

      <footer class="footer">
        <a :href="`${repoUrl}/commit/${commitHash}`" target="_blank" rel="noopener noreferrer">
          Commit: {{ commitHash }}
        </a>
        <p>Build time: {{ buildTime }}</p>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {triggerRemoteSync, useFloatConfig} from '@/composables/useFloatConfig'
import {extractStockItems, fetchGlanceMenu, type StockItem} from '@/composables/useGlanceStocks'
import {
  buildRemoteControlUrl,
  registerDevice,
  resetDeviceId,
  useDeviceId,
} from '@/composables/useRemoteDevice'

defineEmits<{ close: [] }>()

const {
  autoNightMode,
  nightModeRange,
  autoInvisible,
  invisibleRange,
  invisibleDayEnable,
  invisibleDay,
  remoteControlEnabled,
  numericalStocks,
} = useFloatConfig()
const deviceId = useDeviceId()
const remoteControlUrl = computed(() => buildRemoteControlUrl(deviceId.value))
const copyHint = ref('')
const weekLabels = ['日', '一', '二', '三', '四', '五', '六']

const repoUrl = __REPO_URL__
const commitHash = __COMMIT_HASH__
const buildTime = __BUILD_TIME__

const sheetVisible = ref(false)

// 股票行情选项加载
const availableStockOptions = ref<StockItem[]>([])
const isLoadingStocks = ref(false)
const stockLoadError = ref<string | null>(null)

const stockOptionsWithSelected = computed(() => {
  const map = new Map<string, StockItem>()
  for (const opt of availableStockOptions.value) {
    map.set(opt.value, opt)
  }
  // 如果已有选择但尚未加载到（或接口故障），补充占位项
  for (const key of numericalStocks.value || []) {
    if (!map.has(key)) {
      map.set(key, {
        value: key,
        name: key.replace(/^stocks:|^fx:/, ''),
        price: '--',
        fullTitle: key,
      })
    }
  }
  return Array.from(map.values())
})

const loadStockOptions = async () => {
  isLoadingStocks.value = true
  stockLoadError.value = null
  try {
    const data = await fetchGlanceMenu()
    availableStockOptions.value = extractStockItems(data.menu)
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : '网络连接失败'
    stockLoadError.value = `行情接口不可用 (${msg})`
    ElMessage.warning('无法连接股票数据接口: https://dev.flxdu.cn/glance/api/menu')
  } finally {
    isLoadingStocks.value = false
  }
}

const copyRemoteUrl = async () => {
  copyHint.value = ''
  try {
    await navigator.clipboard.writeText(remoteControlUrl.value)
    copyHint.value = '已复制到剪贴板'
  } catch {
    copyHint.value = '复制失败，请手动选择复制'
  }
}

const confirmResetDeviceId = async () => {
  try {
    await ElMessageBox.confirm(
        '将生成新的设备 ID，旧 ID 的远程控制链接将失效，是否继续？',
        '重置设备 ID',
        {
          confirmButtonText: '重置',
          cancelButtonText: '取消',
          type: 'warning',
        },
    )
  } catch {
    return
  }

  const newId = resetDeviceId()
  try {
    await registerDevice(newId)
    triggerRemoteSync()
    ElMessage.success('设备 ID 已重置')
  } catch {
    ElMessage.warning('设备 ID 已重置，但注册到服务器失败')
  }
}

onMounted(() => {
  void loadStockOptions()
  requestAnimationFrame(() => {
    sheetVisible.value = true
  })
})
</script>

<style scoped>
.settings-sheet {
  position: fixed;
  inset: 0;
  z-index: 20;
  background: rgba(18, 18, 18, 0.97);
  display: flex;
  flex-direction: column;
  transform: translateY(100%);
  transition: transform 0.4s cubic-bezier(0.32, 0.72, 0, 1);
  overflow: hidden;
}

.settings-sheet.visible {
  transform: translateY(0);
}

.sheet-header {
  display: flex;
  align-items: center;
  padding: 16px 20px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  margin-right: 12px;
}

.sheet-title {
  font-size: 17px;
  font-weight: 600;
  color: #fff;
}

.sheet-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
}

.setting-label {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.85);
}

.setting-sub {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0 10px 4px;
  flex-wrap: wrap;
}

.weekday-row {
  gap: 6px;
}

.sub-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.45);
  white-space: nowrap;
}

.time-picker {
  width: 130px;
}

.time-sep {
  color: rgba(255, 255, 255, 0.4);
  font-size: 14px;
}

.divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
  margin: 4px 0;
}

.stock-count-hint {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.45);
}

.stock-setting-sub {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
  padding: 4px 0 10px;
}

.stock-select {
  width: 100%;
}

.stock-error-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 73, 73, 0.15);
  border: 1px solid rgba(255, 73, 73, 0.3);
  border-radius: 6px;
  padding: 6px 10px;
}

.stock-error-text {
  font-size: 12px;
  color: #ff8a8a;
}

.retry-btn {
  font-size: 12px;
  color: #fff;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  padding: 2px 8px;
  cursor: pointer;
}

.retry-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.stock-hint {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
  color: rgba(255, 255, 255, 0.45);
}

.stock-option-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  gap: 12px;
}

.option-name {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.option-price {
  color: rgba(255, 255, 255, 0.55);
  font-size: 12px;
  white-space: nowrap;
}

.device-row {
  flex-wrap: wrap;
  gap: 6px;
  align-items: flex-start;
}

.device-id-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
  max-width: 62vw;
}

.reset-btn {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.65);
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 6px;
  padding: 4px 10px;
  cursor: pointer;
}

.reset-btn:hover {
  background: rgba(255, 255, 255, 0.12);
}

.device-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 6px 0 4px;
}

.device-hint {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.5);
}

.device-link {
  font-size: 12px;
  line-height: 1.5;
  color: #8ab4f8;
  word-break: break-all;
  user-select: all;
  text-decoration: none;
  text-align: left;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
}

.device-link:hover {
  text-decoration: underline;
}

.copy-hint {
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}

.device-id {
  font-size: 12px;
  color: #8ab4f8;
  word-break: break-all;
  user-select: all;
  text-align: right;
}

.footer {
  margin-top: auto;
  padding-top: 24px;
  text-align: center;
  font-size: 12px;
  color: #555;
}

.footer a {
  color: #555;
  text-decoration: none;
}
</style>
