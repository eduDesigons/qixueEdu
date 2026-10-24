<template>
  <CommonPage :show-header="false">
    <div class="error-page">
      <p class="error-page__code">403</p>
      <h1 class="error-page__title">没有权限</h1>
      <p class="error-page__desc">抱歉，您暂无权限访问，请联系管理员开通。</p>
      <div class="error-page__actions">
        <n-button v-if="back" secondary @click="router.replace(back)">
          返回上一页
        </n-button>
        <n-button type="primary" @click="router.replace('/')">
          返回首页
        </n-button>
      </div>
    </div>
  </CommonPage>
</template>

<script setup>
const router = useRouter()
const route = useRoute()

const back = history.state.back

if (history.state.from === 'permission-guard') {
  delete history.state.from
}
else if (route.query.path) {
  router.replace(route.query.path)
}
</script>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100%;
  padding: 48px 20px;
  text-align: center;
}

.error-page__code {
  margin: 0;
  font-family: var(--font-display);
  font-size: 64px;
  font-weight: 600;
  letter-spacing: -0.04em;
  line-height: 1;
  color: rgb(var(--primary-color));
}

.error-page__title {
  margin: 16px 0 0;
  font-family: var(--font-display);
  font-size: var(--title-lg, 32px);
  font-weight: 600;
  letter-spacing: -0.028em;
  color: var(--ctl-text, #1a2332);
}

.error-page__desc {
  margin: 8px 0 24px;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #5b6b7c);
}

.error-page__actions {
  display: flex;
  gap: 10px;
}

.error-page :deep(.n-button) {
  --n-height: var(--ctl-h-xl, 36px) !important;
  --n-border-radius: 9px !important;
  min-width: 112px;
  font-weight: 550 !important;
}
</style>
