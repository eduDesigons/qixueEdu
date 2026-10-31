<template>
  <CommonPage :show-header="false">
    <div class="error-page">
      <p class="error-page__code">404</p>
      <h1 class="error-page__title">页面不存在</h1>
      <p class="error-page__desc">生活总归带点荒谬</p>
      <n-button type="primary" @click="replace('/')">
        返回首页
      </n-button>
    </div>
  </CommonPage>
</template>

<script setup>
const { replace } = useRouter()
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

.error-page :deep(.n-button) {
  --n-height: var(--ctl-h-xl, 36px) !important;
  --n-border-radius: 9px !important;
  min-width: 112px;
  font-weight: 550 !important;
}
</style>
