import { apiClient } from './client'
import type { OpsRequestOptions, OpsStatusMatrixParams, OpsStatusMatrixResponse } from './admin/ops'

export async function getStatusMatrix(
  params: OpsStatusMatrixParams,
  options: OpsRequestOptions = {}
): Promise<OpsStatusMatrixResponse> {
  const { data } = await apiClient.get<OpsStatusMatrixResponse>('/ops/status-matrix', {
    params,
    signal: options.signal
  })
  return data
}

export const userOpsAPI = {
  getStatusMatrix
}

export default userOpsAPI
