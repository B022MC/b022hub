import http from 'k6/http'
import { check, sleep } from 'k6'

const baseUrl = (__ENV.BASE_URL || 'https://your-host.example.com').replace(/\/$/, '')
const endpoint = __ENV.ENDPOINT || `${baseUrl}/v1/chat/completions`
const apiKey = __ENV.API_KEY || ''
const model = __ENV.MODEL || 'gpt-5.4'
const prompt = __ENV.PROMPT || 'Generate exactly 200 short numbered lines without markdown.'
const maxTokens = Number(__ENV.MAX_TOKENS || 256)
const stream = (__ENV.STREAM || 'false') === 'true'
const stepTargets = (__ENV.CONCURRENCY_STEPS || '5,10,20,40,80')
  .split(',')
  .map((item) => Number(item.trim()))
  .filter((item) => Number.isFinite(item) && item > 0)
const rampDuration = __ENV.RAMP_DURATION || '15s'
const holdDuration = __ENV.HOLD_DURATION || '60s'

if (!apiKey) {
  throw new Error('missing API_KEY env')
}

if (!model) {
  throw new Error('missing MODEL env')
}

const stages = stepTargets.flatMap((target) => [
  { duration: rampDuration, target },
  { duration: holdDuration, target }
])

export const options = {
  scenarios: {
    concurrency_ramp: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages,
      gracefulRampDown: '5s'
    }
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<60000']
  }
}

const payload = JSON.stringify({
  model,
  messages: [
    {
      role: 'user',
      content: prompt
    }
  ],
  max_tokens: maxTokens,
  stream
})

const params = {
  headers: {
    Authorization: `Bearer ${apiKey}`,
    'Content-Type': 'application/json'
  },
  timeout: __ENV.TIMEOUT || '120s'
}

export default function () {
  const response = http.post(endpoint, payload, params)
  check(response, {
    'status is 200 or 429': (res) => res.status === 200 || res.status === 429
  })
  sleep(Number(__ENV.SLEEP_SECONDS || 0))
}
