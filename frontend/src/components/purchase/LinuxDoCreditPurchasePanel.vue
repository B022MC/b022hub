<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { createLinuxDoCreditOrder, getLinuxDoCreditOrder, listLinuxDoCreditOrders, type LinuxDoCreditOrder } from '@/api/linuxdoCredit'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  exchangeRate: number
  initialOutTradeNo?: string
}>()

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const amount = ref<number>(10)
const loading = ref(false)
const syncingOrder = ref(false)
const errorMessage = ref('')
const activeOrder = ref<LinuxDoCreditOrder | null>(null)
const orders = ref<LinuxDoCreditOrder[]>([])
let pollTimer: number | null = null

const normalizedExchangeRate = computed(() => {
  if (!Number.isFinite(props.exchangeRate) || props.exchangeRate <= 0) {
    return 1
  }
  return props.exchangeRate
})

const estimatedCredit = computed(() => {
  return formatAmount(amount.value * normalizedExchangeRate.value)
})

const hasInitialOrder = computed(() => Boolean(props.initialOutTradeNo))

function formatAmount(value: number): string {
  if (!Number.isFinite(value)) {
    return '0'
  }
  return value.toFixed(value % 1 === 0 ? 0 : 2)
}

function statusLabel(status: string): string {
  if (status === 'paid') {
    return t('purchase.linuxdoCredit.status.paid')
  }
  return t('purchase.linuxdoCredit.status.pending')
}

function stopPolling(): void {
  if (pollTimer !== null) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
}

async function loadOrders(): Promise<void> {
  try {
    orders.value = await listLinuxDoCreditOrders(10)
  } catch (error) {
    console.error('Failed to load LinuxDo credit orders:', error)
  }
}

async function syncOrder(outTradeNo: string): Promise<void> {
  syncingOrder.value = true
  try {
    const order = await getLinuxDoCreditOrder(outTradeNo, true)
    activeOrder.value = order
    await loadOrders()

    if (order.status === 'paid') {
      stopPolling()
      await authStore.refreshUser()
      appStore.showSuccess(t('purchase.linuxdoCredit.paymentSuccess'))
    }
  } catch (error: any) {
    errorMessage.value =
      error?.message || t('purchase.linuxdoCredit.syncFailed')
  } finally {
    syncingOrder.value = false
  }
}

function startPolling(outTradeNo: string): void {
  stopPolling()
  let attempts = 0
  pollTimer = window.setInterval(async () => {
    attempts += 1
    await syncOrder(outTradeNo)
    if (attempts >= 20 || activeOrder.value?.status === 'paid') {
      stopPolling()
    }
  }, 3000)
}

function submitCheckoutForm(checkoutUrl: string, fields: Record<string, string>): void {
  const form = document.createElement('form')
  form.method = 'POST'
  form.action = checkoutUrl
  form.style.display = 'none'

  for (const [key, value] of Object.entries(fields)) {
    const input = document.createElement('input')
    input.type = 'hidden'
    input.name = key
    input.value = value
    form.appendChild(input)
  }

  document.body.appendChild(form)
  form.submit()
  form.remove()
}

async function handleCreateOrder(): Promise<void> {
  errorMessage.value = ''
  if (!Number.isFinite(amount.value) || amount.value <= 0) {
    errorMessage.value = t('purchase.linuxdoCredit.invalidAmount')
    return
  }

  loading.value = true
  try {
    const checkout = await createLinuxDoCreditOrder(amount.value)
    activeOrder.value = checkout.order
    await loadOrders()
    appStore.showInfo(t('purchase.linuxdoCredit.redirectingToPay'))
    submitCheckoutForm(checkout.checkout_url, checkout.fields)
  } catch (error: any) {
    errorMessage.value =
      error?.message || t('purchase.linuxdoCredit.createOrderFailed')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadOrders()

  if (props.initialOutTradeNo) {
    await syncOrder(props.initialOutTradeNo)
    if (activeOrder.value?.status !== 'paid') {
      startPolling(props.initialOutTradeNo)
    }
  }
})

onBeforeUnmount(() => {
  stopPolling()
})
</script>

<template>
  <div class="linuxdo-credit-panel">
    <div class="linuxdo-credit-hero">
      <div>
        <p class="linuxdo-credit-eyebrow">{{ t('purchase.linuxdoCredit.eyebrow') }}</p>
        <h3 class="linuxdo-credit-title">{{ t('purchase.linuxdoCredit.title') }}</h3>
        <p class="linuxdo-credit-description">
          {{ t('purchase.linuxdoCredit.description') }}
        </p>
      </div>
      <div class="linuxdo-credit-rate">
        <span>{{ t('purchase.linuxdoCredit.exchangeRate') }}</span>
        <strong>1 : {{ formatAmount(normalizedExchangeRate) }}</strong>
      </div>
    </div>

    <div class="linuxdo-credit-grid">
      <section class="linuxdo-credit-card">
        <label class="linuxdo-credit-label" for="linuxdo-credit-amount">
          {{ t('purchase.linuxdoCredit.amountLabel') }}
        </label>
        <input
          id="linuxdo-credit-amount"
          v-model.number="amount"
          type="number"
          min="1"
          step="1"
          class="input linuxdo-credit-input"
        />

        <div class="linuxdo-credit-summary">
          <span>{{ t('purchase.linuxdoCredit.estimatedCredit') }}</span>
          <strong>{{ estimatedCredit }}</strong>
        </div>

        <p class="linuxdo-credit-hint">
          {{ t('purchase.linuxdoCredit.returnHint') }}
        </p>

        <button class="btn btn-primary w-full" :disabled="loading" @click="handleCreateOrder">
          <Icon v-if="!loading" name="externalLink" size="sm" class="mr-2" :stroke-width="2" />
          {{ loading ? t('purchase.linuxdoCredit.creating') : t('purchase.linuxdoCredit.payNow') }}
        </button>

        <p v-if="errorMessage" class="linuxdo-credit-error">{{ errorMessage }}</p>
      </section>

      <section class="linuxdo-credit-card">
        <div class="linuxdo-credit-status-head">
          <div>
            <p class="linuxdo-credit-label">{{ t('purchase.linuxdoCredit.latestOrder') }}</p>
            <h4 class="linuxdo-credit-order-title">
              {{ activeOrder?.title || t('purchase.linuxdoCredit.noOrderYet') }}
            </h4>
          </div>
          <span
            v-if="activeOrder"
            class="linuxdo-credit-badge"
            :class="activeOrder.status === 'paid' ? 'is-paid' : 'is-pending'"
          >
            <Icon
              :name="activeOrder.status === 'paid' ? 'checkCircle' : 'clock'"
              size="sm"
              class="mr-1"
            />
            {{ statusLabel(activeOrder.status) }}
          </span>
        </div>

        <div v-if="activeOrder" class="linuxdo-credit-order-meta">
          <div>
            <span>{{ t('purchase.linuxdoCredit.orderNumber') }}</span>
            <code>{{ activeOrder.out_trade_no }}</code>
          </div>
          <div>
            <span>{{ t('purchase.linuxdoCredit.orderAmount') }}</span>
            <strong>{{ formatAmount(activeOrder.amount) }}</strong>
          </div>
          <div>
            <span>{{ t('purchase.linuxdoCredit.orderCredit') }}</span>
            <strong>{{ formatAmount(activeOrder.credited_amount) }}</strong>
          </div>
        </div>

        <button
          v-if="activeOrder && activeOrder.status !== 'paid'"
          class="btn btn-secondary w-full"
          :disabled="syncingOrder"
          @click="syncOrder(activeOrder.out_trade_no)"
        >
          {{ syncingOrder ? t('purchase.linuxdoCredit.syncing') : t('purchase.linuxdoCredit.checkStatus') }}
        </button>

        <p v-else-if="hasInitialOrder" class="linuxdo-credit-success">
          {{ t('purchase.linuxdoCredit.returnDetected') }}
        </p>
      </section>
    </div>

    <section class="linuxdo-credit-card">
      <div class="linuxdo-credit-status-head">
        <div>
          <p class="linuxdo-credit-label">{{ t('purchase.linuxdoCredit.recentOrders') }}</p>
          <h4 class="linuxdo-credit-order-title">{{ t('purchase.linuxdoCredit.recentOrdersHint') }}</h4>
        </div>
      </div>

      <div v-if="orders.length === 0" class="linuxdo-credit-empty">
        {{ t('purchase.linuxdoCredit.emptyOrders') }}
      </div>

      <div v-else class="linuxdo-credit-list">
        <div v-for="order in orders" :key="order.out_trade_no" class="linuxdo-credit-list-item">
          <div>
            <p class="linuxdo-credit-list-title">{{ order.title }}</p>
            <code class="linuxdo-credit-list-code">{{ order.out_trade_no }}</code>
          </div>
          <div class="linuxdo-credit-list-side">
            <strong>{{ formatAmount(order.credited_amount) }}</strong>
            <span>{{ statusLabel(order.status) }}</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.linuxdo-credit-panel {
  @apply space-y-6 p-6;
}

.linuxdo-credit-hero {
  @apply grid gap-4 rounded-2xl border border-orange-200 bg-gradient-to-br from-orange-50 via-white to-amber-50 p-6;
  @apply md:grid-cols-[1fr_auto];
}

.linuxdo-credit-eyebrow {
  @apply text-xs font-semibold uppercase tracking-[0.2em] text-orange-600;
}

.linuxdo-credit-title {
  @apply mt-2 text-2xl font-semibold text-gray-900 dark:text-white;
}

.linuxdo-credit-description {
  @apply mt-2 text-sm text-gray-600 dark:text-gray-300;
}

.linuxdo-credit-rate {
  @apply rounded-2xl bg-white/80 px-4 py-3 text-right text-sm text-gray-600 shadow-sm backdrop-blur;
  @apply dark:bg-dark-800/80 dark:text-gray-300;
}

.linuxdo-credit-rate strong {
  @apply mt-1 block text-xl font-semibold text-gray-900 dark:text-white;
}

.linuxdo-credit-grid {
  @apply grid gap-6;
  @apply lg:grid-cols-2;
}

.linuxdo-credit-card {
  @apply rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-dark-700 dark:bg-dark-900;
}

.linuxdo-credit-label {
  @apply mb-2 block text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-gray-400;
}

.linuxdo-credit-input {
  @apply text-lg font-semibold;
}

.linuxdo-credit-summary {
  @apply mt-4 flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 text-sm dark:bg-dark-800;
}

.linuxdo-credit-summary strong {
  @apply text-lg text-gray-900 dark:text-white;
}

.linuxdo-credit-hint {
  @apply mt-3 text-xs text-gray-500 dark:text-gray-400;
}

.linuxdo-credit-error {
  @apply mt-3 text-sm text-red-600 dark:text-red-400;
}

.linuxdo-credit-success {
  @apply text-sm text-emerald-600 dark:text-emerald-400;
}

.linuxdo-credit-status-head {
  @apply mb-4 flex items-start justify-between gap-4;
}

.linuxdo-credit-order-title {
  @apply text-base font-semibold text-gray-900 dark:text-white;
}

.linuxdo-credit-badge {
  @apply inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold;
}

.linuxdo-credit-badge.is-paid {
  @apply bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300;
}

.linuxdo-credit-badge.is-pending {
  @apply bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300;
}

.linuxdo-credit-order-meta {
  @apply space-y-3 rounded-xl bg-gray-50 p-4 text-sm dark:bg-dark-800;
}

.linuxdo-credit-order-meta div {
  @apply flex items-center justify-between gap-4;
}

.linuxdo-credit-order-meta code {
  @apply break-all rounded bg-white px-2 py-1 text-xs dark:bg-dark-900;
}

.linuxdo-credit-empty {
  @apply rounded-xl bg-gray-50 px-4 py-6 text-center text-sm text-gray-500 dark:bg-dark-800 dark:text-gray-400;
}

.linuxdo-credit-list {
  @apply space-y-3;
}

.linuxdo-credit-list-item {
  @apply flex items-center justify-between gap-4 rounded-xl border border-gray-100 px-4 py-3 dark:border-dark-700;
}

.linuxdo-credit-list-title {
  @apply text-sm font-medium text-gray-900 dark:text-white;
}

.linuxdo-credit-list-code {
  @apply mt-1 inline-block break-all text-xs text-gray-500 dark:text-gray-400;
}

.linuxdo-credit-list-side {
  @apply text-right text-sm text-gray-500 dark:text-gray-400;
}

.linuxdo-credit-list-side strong {
  @apply block text-base text-gray-900 dark:text-white;
}
</style>
