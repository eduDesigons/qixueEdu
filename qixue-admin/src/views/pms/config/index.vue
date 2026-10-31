<script setup>
import { ConstructOutline, DocumentTextOutline, HelpCircleOutline, InformationCircleOutline, PeopleOutline, PlayOutline, SettingsOutline, TrophyOutline } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { computed, ref } from 'vue'
import { AppPage, SettingsList, SettingsListRow } from '@/components/index.js'
import EditAboutConfig from './components/EditAboutConfig.vue'
import EditCoinConfig from './components/EditCoinConfig.vue'
import EditEnrollConfig from './components/EditEnrollConfig.vue'
import EditFaqConfig from './components/EditFaqConfig.vue'
import EditMemberRule from './components/EditMemberRule.vue'
import EditToolConfig from './components/EditToolConfig.vue'
import EditPlayerConfig from './components/EditPlayerConfig.vue'
import EditSiteConfig from './components/EditSiteConfig.vue'

defineOptions({ name: 'ManageConfig' })

const activeKey = ref(null)

const groups = [
  {
    label: '站点',
    items: [
      { key: 'sysConfig', name: '网站信息', desc: '名称、Logo 与页脚', icon: SettingsOutline, component: EditSiteConfig },
      { key: 'playerConfig', name: '播放器', desc: '快进、倍速等学习端控制', icon: PlayOutline, component: EditPlayerConfig },
    ],
  },
  {
    label: '运营',
    items: [
      { key: 'coinConfig', name: '积分', desc: '发放数额与 QQ 等级换算', icon: TrophyOutline, component: EditCoinConfig },
      { key: 'memberRule', name: '学员规则', desc: '后台添加学员的时间', icon: PeopleOutline, component: EditMemberRule },
    ],
  },
  {
    label: '帮助',
    items: [
      { key: 'faqConfig', name: '常见问题', desc: '学员端帮助列表', icon: HelpCircleOutline, component: EditFaqConfig },
      { key: 'toolConfig', name: '工具库', desc: '学员端外链工具', icon: ConstructOutline, component: EditToolConfig },
      { key: 'aboutConfig', name: '关于我们', desc: '关于页富文本', icon: InformationCircleOutline, component: EditAboutConfig },
      { key: 'enrollConfig', name: '入学要求', desc: '学员列表里的说明文档', icon: DocumentTextOutline, component: EditEnrollConfig },
    ],
  },
]

const allItems = groups.flatMap(g => g.items)
const activeItem = computed(() => allItems.find(item => item.key === activeKey.value) || null)

function openItem(key) {
  activeKey.value = key
}

function closeEditor() {
  activeKey.value = null
}
</script>

<template>
  <AppPage>
    <div class="settings-page">
      <template v-if="!activeItem">
        <header class="page-mast">
          <h1 class="page-mast__title">设置</h1>
          <p class="page-mast__sub">站点、播放与运营相关配置</p>
        </header>
        <SettingsList v-for="group in groups" :key="group.label" :label="group.label">
          <SettingsListRow
            v-for="item in group.items"
            :key="item.key"
            :title="item.name"
            :desc="item.desc"
            @click="openItem(item.key)"
          >
            <template #icon>
              <NIcon :component="item.icon" :size="18" />
            </template>
          </SettingsListRow>
        </SettingsList>
      </template>

      <section v-else class="settings-editor-page">
        <header class="settings-editor-page__head">
          <button type="button" class="page-back" @click="closeEditor">
            <i class="i-material-symbols:chevron-left page-back__icon" />
            <span class="page-back__label">设置</span>
          </button>
          <h1 class="page-mast__title">
            {{ activeItem.name }}
          </h1>
          <p class="page-mast__sub">
            {{ activeItem.desc }}
          </p>
        </header>
        <component
          :is="activeItem.component"
          :key="activeItem.key"
          embedded
          @close="closeEditor"
        />
      </section>
    </div>
  </AppPage>
</template>

<style scoped>
.settings-page {
  padding: 4px 0 28px;
}

.settings-page .page-mast,
.settings-page .settings-list {
  animation: settings-in 0.32s var(--ease-apple, cubic-bezier(0.32, 0.72, 0, 1));
}

.settings-editor-page {
  max-width: 760px;
  margin: 0 auto;
  animation: settings-in 0.32s var(--ease-apple, cubic-bezier(0.32, 0.72, 0, 1));
}

.settings-editor-page__head {
  margin-bottom: 20px;
}

.settings-editor-page__head .page-back {
  margin: 0 0 10px -8px;
}

@keyframes settings-in {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .settings-page .page-mast,
  .settings-page .settings-list,
  .settings-editor-page {
    animation: none;
  }
}
</style>
