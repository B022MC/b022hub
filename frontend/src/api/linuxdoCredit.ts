import { apiClient } from './client'

export interface LinuxDoCreditOrder {
  id: number
  provider: string
  out_trade_no: string
  provider_trade_no?: string
  user_id: number
  title: string
  amount: number
  credited_amount: number
  status: string
  paid_at?: string
  created_at: string
  updated_at: string
}

export interface LinuxDoCreditCheckoutResponse {
  order: LinuxDoCreditOrder
  checkout_url: string
  fields: Record<string, string>
}

export async function createLinuxDoCreditOrder(
  amount: number
): Promise<LinuxDoCreditCheckoutResponse> {
  const { data } = await apiClient.post<LinuxDoCreditCheckoutResponse>('/payments/linuxdo/orders', {
    amount
  })
  return data
}

export async function listLinuxDoCreditOrders(limit = 10): Promise<LinuxDoCreditOrder[]> {
  const { data } = await apiClient.get<{ items: LinuxDoCreditOrder[] }>('/payments/linuxdo/orders', {
    params: { limit }
  })
  return data.items || []
}

export async function getLinuxDoCreditOrder(
  outTradeNo: string,
  sync = false
): Promise<LinuxDoCreditOrder> {
  const { data } = await apiClient.get<LinuxDoCreditOrder>(
    `/payments/linuxdo/orders/${encodeURIComponent(outTradeNo)}`,
    {
      params: sync ? { sync: 1 } : undefined
    }
  )
  return data
}
