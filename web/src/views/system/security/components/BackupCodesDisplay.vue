<template>
  <div class="backup-codes-display">
    <n-alert type="info" title="备用恢复码" style="margin-bottom: 24px;">
      这些是您的新备用恢复码。请将它们保存在安全的地方，每个代码只能使用一次。
    </n-alert>
    
    <div class="codes-container">
      <div class="codes-grid">
        <div v-for="(code, index) in backupCodes" :key="index" class="backup-code">
          <span class="code-number">{{ String(index + 1).padStart(2, '0') }}.</span>
          <span class="code-value">{{ code }}</span>
        </div>
      </div>
    </div>
    
    <div class="actions">
      <n-space>
        <n-button @click="copyAllCodes">
          <template #icon>
            <n-icon><CopyOutline /></n-icon>
          </template>
          复制所有代码
        </n-button>
        
        <n-button @click="downloadCodes">
          <template #icon>
            <n-icon><DownloadOutline /></n-icon>
          </template>
          下载为文件
        </n-button>
        
        <n-button @click="printCodes">
          <template #icon>
            <n-icon><PrintOutline /></n-icon>
          </template>
          打印
        </n-button>
      </n-space>
    </div>
    
    <div class="security-tips">
      <n-alert type="warning" title="安全提醒">
        <ul>
          <li>请将这些代码保存在安全的地方，如密码管理器或安全的物理位置</li>
          <li>不要将这些代码存储在不安全的地方，如电子邮件或云笔记</li>
          <li>每个代码只能使用一次，使用后将自动失效</li>
          <li>如果您担心代码泄露，可以随时重新生成新的备用码</li>
          <li>建议定期更新备用码以确保安全性</li>
        </ul>
      </n-alert>
    </div>
    
    <div class="close-action">
      <n-button type="primary" @click="$emit('close')">
        我已安全保存
      </n-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import { CopyOutline, DownloadOutline, PrintOutline } from '@vicons/ionicons5'

// 属性定义
interface Props {
  backupCodes: string[]
}

const props = defineProps<Props>()

// 事件定义
const emit = defineEmits<{
  close: []
}>()

// 响应式数据
const message = useMessage()

// 方法
const copyAllCodes = async () => {
  try {
    const codesText = props.backupCodes.map((code, index) => 
      `${String(index + 1).padStart(2, '0')}. ${code}`
    ).join('\n')
    
    await navigator.clipboard.writeText(codesText)
    message.success('所有备用码已复制到剪贴板')
  } catch (error) {
    message.error('复制失败，请手动复制')
  }
}

const downloadCodes = () => {
  const timestamp = new Date().toISOString().slice(0, 19).replace(/[T:]/g, '-')
  const content = [
    '双因素认证备用恢复码',
    `生成时间: ${new Date().toLocaleString()}`,
    '',
    '重要提醒:',
    '- 每个代码只能使用一次',
    '- 请将这些代码保存在安全的地方',
    '- 不要与他人分享这些代码',
    '- 如果丢失，您可以重新生成新的备用码',
    '',
    '备用恢复码:',
    ...props.backupCodes.map((code, index) => 
      `${String(index + 1).padStart(2, '0')}. ${code}`
    )
  ].join('\n')
  
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `2FA-backup-codes-${timestamp}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  message.success('备用码文件已下载')
}

const printCodes = () => {
  const printContent = [
    '<html><head><title>双因素认证备用恢复码</title>',
    '<style>',
    'body { font-family: Arial, sans-serif; margin: 40px; }',
    'h1 { color: #333; border-bottom: 2px solid #333; padding-bottom: 10px; }',
    '.timestamp { color: #666; font-size: 14px; margin-bottom: 20px; }',
    '.warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; margin: 20px 0; border-radius: 5px; }',
    '.codes { margin: 20px 0; }',
    '.code-item { font-family: monospace; font-size: 16px; margin: 8px 0; padding: 8px; background: #f8f9fa; border-radius: 3px; }',
    'ul { margin: 10px 0; padding-left: 20px; }',
    'li { margin: 5px 0; }',
    '</style></head><body>',
    '<h1>双因素认证备用恢复码</h1>',
    `<div class="timestamp">生成时间: ${new Date().toLocaleString()}</div>`,
    '<div class="warning">',
    '<strong>重要提醒:</strong>',
    '<ul>',
    '<li>每个代码只能使用一次</li>',
    '<li>请将这些代码保存在安全的地方</li>',
    '<li>不要与他人分享这些代码</li>',
    '<li>如果丢失，您可以重新生成新的备用码</li>',
    '</ul>',
    '</div>',
    '<div class="codes">',
    '<h3>备用恢复码:</h3>',
    ...props.backupCodes.map((code, index) => 
      `<div class="code-item">${String(index + 1).padStart(2, '0')}. ${code}</div>`
    ),
    '</div>',
    '</body></html>'
  ].join('')
  
  const printWindow = window.open('', '_blank')
  if (printWindow) {
    printWindow.document.write(printContent)
    printWindow.document.close()
    printWindow.focus()
    printWindow.print()
    printWindow.close()
    message.success('打印窗口已打开')
  } else {
    message.error('无法打开打印窗口，请检查浏览器设置')
  }
}
</script>

<style lang="less" scoped>
.backup-codes-display {
  .codes-container {
    margin: 24px 0;
    
    .codes-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 8px;
      max-height: 300px;
      overflow-y: auto;
      padding: 16px;
      background-color: #f8f9fa;
      border-radius: 8px;
      border: 1px solid #e9ecef;
      
      .backup-code {
        display: flex;
        align-items: center;
        padding: 8px 12px;
        background-color: white;
        border-radius: 4px;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 14px;
        
        .code-number {
          color: #666;
          margin-right: 12px;
          min-width: 30px;
        }
        
        .code-value {
          color: #333;
          font-weight: 500;
          letter-spacing: 1px;
        }
      }
    }
  }
  
  .actions {
    margin: 24px 0;
    text-align: center;
  }
  
  .security-tips {
    margin: 24px 0;
    
    :deep(.n-alert__content) {
      ul {
        margin: 8px 0 0 0;
        padding-left: 20px;
        
        li {
          margin: 6px 0;
          line-height: 1.5;
        }
      }
    }
  }
  
  .close-action {
    margin-top: 32px;
    text-align: center;
    padding-top: 24px;
    border-top: 1px solid #f0f0f0;
  }
}
</style>