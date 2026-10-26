import { useEffect, useState } from 'react'
import { products } from './products'
import { creditDetails } from './creditData'
import type { Applicant, LoanApplication, Selection } from './types'

export type Session = { accessToken: string; role: string; expiresIn: number }
export type ApiCompany = {
  id: string
  inn: string
  name: string
  registeredAt?: string
  businessAgeMonths: number
  monthlyRevenueMinor: number
  monthlyExpensesMinor: number
  monthlyDebtPaymentsMinor: number
  representativeName?: string
  representativePosition?: string
  phone?: string
  email?: string
  verificationStatus: string
}
export type ApiDocument = { id: string; type: string; fileName: string; requestItemId?: string }
export type ApiOffer = {
  id: string
  productId: string
  amountMinor: number
  termMonths: number
  annualRateBps: number
  collateral?: string
}
export type ApiApplication = {
  id: string
  number: string
  companyId: string
  productId: string
  version: number
  status: string
  need: { amountMinor: number; termMonths: number; purpose: string; description?: string }
  createdAt: string
  recommendation?: string
}
export type ApiScoring = {
  score: number
  riskLevel: string
  factors: string[]
  modelVersion: string
  recommendation: string
}
export type ApiInformation = {
  id: string
  message: string
  createdAt: string
  items: { id: string; title: string; type: string; required: boolean; status: string }[]
}
export type ApiOverview = {
  application: ApiApplication
  companySnapshot: ApiCompany
  scoring?: ApiScoring
  decision?: { offers: ApiOffer[] }
  documents?: ApiDocument[]
  informationRequests?: ApiInformation[]
}
export type Recommendation = {
  product: { id: string }
  availability: string
  reason: string
  rank: number
}
type ApiProduct = {
  id: string
  name: string
  description: string
  minAmountMinor: number
  maxAmountMinor: number
  minTermMonths: number
  maxTermMonths: number
  annualRateBps: number
  requiredDocumentTypes: string[]
  collateral: string
  paymentType: string
  purposes: string[]
}

const base = (import.meta.env.VITE_API_BASE_URL || '/api/v1').replace(/\/$/, '')
const sessionKey = 'credit-api-session'
let session: Session | null = readSession()
export const selectedFiles = new Map<string, File>()
export const backendProductId = (id: string) => (id === 'business' ? 'unsecured' : id)
export const frontendProductId = (id: string) => (id === 'unsecured' ? 'business' : id)

function readSession(): Session | null {
  try {
    const value = JSON.parse(sessionStorage.getItem(sessionKey) || 'null')
    return value &&
      typeof value.accessToken === 'string' &&
      ['CLIENT', 'EMPLOYEE'].includes(value.role)
      ? value
      : null
  } catch {
    return null
  }
}

export function getSession() {
  return session
}

export function setSession(value: Session | null) {
  session = value
  selectedFiles.clear()
  try {
    if (value) sessionStorage.setItem(sessionKey, JSON.stringify(value))
    else sessionStorage.removeItem(sessionKey)
  } catch {}
}

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message)
  }
}

export async function request<T>(
  path: string,
  options: RequestInit = {},
  version?: number,
): Promise<T> {
  const headers = new Headers(options.headers)
  const requestToken = session?.accessToken
  if (requestToken) headers.set('Authorization', `Bearer ${requestToken}`)
  if (options.body && !(options.body instanceof FormData))
    headers.set('Content-Type', 'application/json')
  if (options.method && !['GET', 'HEAD'].includes(options.method))
    headers.set(
      'Idempotency-Key',
      Array.from(crypto.getRandomValues(new Uint8Array(16)), (value) =>
        value.toString(16).padStart(2, '0'),
      ).join(''),
    )
  if (version !== undefined) headers.set('If-Match', `"${version}"`)

  let response: Response
  try {
    response = await fetch(`${base}${path}`, { ...options, headers })
  } catch (error) {
    if (options.signal?.aborted) throw error
    response = await fetch(`${base}${path}`, { ...options, headers })
  }
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    if (
      response.status === 401 &&
      !path.startsWith('/auth/') &&
      requestToken &&
      session?.accessToken === requestToken
    ) {
      setSession(null)
      window.dispatchEvent(new Event('credit-session-expired'))
    }
    throw new ApiError(response.status, error.message || `Ошибка API: ${response.status}`)
  }
  if (response.status === 204) return undefined as T
  return response.json()
}

export async function authenticate(login: string, password: string, displayName?: string) {
  const result = await request<Session>(
    displayName === undefined ? '/auth/login' : '/auth/register',
    {
      method: 'POST',
      body: JSON.stringify({
        login,
        password,
        ...(displayName === undefined ? {} : { displayName }),
      }),
    },
  )
  setSession(result)
  return result
}

export async function loadProducts() {
  const response = await request<ApiProduct[]>('/products')
  for (const value of response) {
    const id = frontendProductId(value.id)
    const product = products.find((item) => item.id === id)
    if (!product) continue
    Object.assign(product, {
      name: value.name,
      description: value.description,
      amount: value.maxAmountMinor / 100_000_000,
      months: value.maxTermMonths,
      rate: value.annualRateBps / 100,
      term: undefined,
      termTitle: undefined,
      purposes: value.purposes,
      secured: value.collateral !== 'NONE',
    })
    const details = creditDetails[id]
    Object.assign(details, {
      minAmount: value.minAmountMinor / 100,
      minMonths: value.minTermMonths,
      amountLabel: `До ${(value.maxAmountMinor / 100).toLocaleString('ru-RU')} ₽`,
      termLabel: `До ${value.maxTermMonths} месяцев`,
      rateLabel: `От ${value.annualRateBps / 100}% годовых`,
      purpose: value.purposes
        .map(
          (item) =>
            (
              ({
                WORKING_CAPITAL: 'Оборотные средства',
                BUSINESS_EXPANSION: 'Развитие бизнеса',
                EQUIPMENT: 'Оборудование',
                REAL_ESTATE: 'Недвижимость',
                REFINANCING: 'Рефинансирование',
                SELLERS: 'Для селлеров',
                OTHER: 'Другие цели',
              }) as Record<string, string>
            )[item] || item,
        )
        .join(', '),
      repayment: [
        value.paymentType === 'ANNUITY'
          ? 'Аннуитетные платежи'
          : 'Индивидуальный график; расчёт платежа является оценочным',
      ],
      features: [
        'Демонстрационные условия каталога. Итоговые параметры определяются решением по заявке.',
      ],
      borrower: ['Требования проверяются по данным компании и результатам скоринга.'],
      documents: value.requiredDocumentTypes.map((type) => ({
        id: type,
        title: documentTitle(type),
        required: true,
        description: 'PDF, PNG или JPEG; до 10 МБ.',
      })),
    })
  }
}

export function documentTitle(type: string) {
  return (
    (
      {
        FINANCIAL_REPORT: 'Финансовая отчётность',
        BANK_STATEMENT: 'Выписка по счёту',
        OTHER: 'Дополнительный документ',
      } as Record<string, string>
    )[type] || type
  )
}

export function purpose(selection: Selection) {
  return (
    (
      {
        'Пополнение оборотных средств': 'WORKING_CAPITAL',
        'Закупка товара': 'WORKING_CAPITAL',
        'Оплата поставщикам': 'WORKING_CAPITAL',
        'Покупка оборудования': 'EQUIPMENT',
        'Развитие бизнеса': 'BUSINESS_EXPANSION',
        'Покупка недвижимости': 'REAL_ESTATE',
        Рефинансирование: 'REFINANCING',
        'Для селлеров': 'SELLERS',
      } as Record<string, string>
    )[selection.purpose] || 'OTHER'
  )
}

export async function saveCompany(
  applicant: Applicant,
  selection: Selection,
  finance: { revenue: number; expenses: number; debts: number },
  companyId?: string,
) {
  const registered = new Date(applicant.registered)
  const now = new Date()
  const data = {
    inn: applicant.inn,
    name: applicant.company || applicant.inn,
    registeredAt: applicant.registered,
    businessAgeMonths: Math.max(
      0,
      (now.getFullYear() - registered.getFullYear()) * 12 + now.getMonth() - registered.getMonth(),
    ),
    monthlyRevenueMinor: Math.round(finance.revenue * 100),
    monthlyExpensesMinor: Math.round(finance.expenses * 100),
    monthlyDebtPaymentsMinor: Math.round(finance.debts * 100),
    representativeName: applicant.name,
    representativePosition: applicant.position,
    phone: applicant.phone.replace(/[\s()-]/g, ''),
    email: applicant.email,
    marketplaceSales: selection.purpose === 'Для селлеров',
  }
  const companies = await request<ApiCompany[]>('/companies')
  const existing = companyId
    ? companies.find((item) => item.id === companyId)
    : companies.find((item) => item.inn === applicant.inn)
  const company = await request<ApiCompany>(existing ? `/companies/${existing.id}` : '/companies', {
    method: existing ? 'PATCH' : 'POST',
    body: JSON.stringify(data),
  })
  return request<ApiCompany>(`/companies/${company.id}/verification`, { method: 'POST' })
}

export async function overview(value: ApiApplication, staff: boolean): Promise<LoanApplication> {
  let details: ApiOverview
  if (staff) details = await request<ApiOverview>(`/bank/applications/${value.id}`)
  else {
    const [company, documents, information] = await Promise.all([
      request<ApiCompany>(`/companies/${value.companyId}`),
      request<ApiDocument[]>(`/applications/${value.id}/documents`),
      request<ApiInformation[]>(`/applications/${value.id}/information-requests`),
    ])
    const decision = [
      'APPROVED',
      'ALTERNATIVE',
      'REJECTED',
      'ACCEPTED',
      'DECLINED',
      'ISSUED',
    ].includes(value.status)
      ? await request<{ offers: ApiOffer[] }>(`/applications/${value.id}/decision`)
      : undefined
    details = {
      application: value,
      companySnapshot: company,
      documents,
      informationRequests: information,
      decision,
    }
  }
  const app = details.application
  const company = details.companySnapshot
  const offer = details.decision?.offers?.[0]
  const information = details.informationRequests
    ?.filter((item) => item.items.some((part) => part.status !== 'FULFILLED'))
    .at(-1)
  const statuses: Record<string, LoanApplication['status']> = {
    IN_PROGRESS: 'Черновик',
    SCORING: 'На рассмотрении',
    CHECKED: 'На рассмотрении',
    MANUAL_CHECK: 'На рассмотрении',
    NEEDS_INFO: 'Запрошены данные',
    APPROVED: 'Одобрено',
    ALTERNATIVE: 'Другие условия',
    REJECTED: 'Отклонено',
    ACCEPTED: 'Принято',
    ISSUED: 'Выдано',
    DECLINED: 'Отклонено',
    CANCELLED: 'Отменено',
  }
  const accepted = ['ACCEPTED', 'ISSUED'].includes(app.status) && offer
  return {
    id: app.id,
    number: app.number,
    backendStatus: app.status,
    version: app.version,
    companyId: app.companyId,
    purposeCode: app.need.purpose,
    description: app.need.description,
    productId: frontendProductId(accepted ? offer.productId : app.productId),
    amount: (accepted ? offer.amountMinor : app.need.amountMinor) / 100,
    months: accepted ? offer.termMonths : app.need.termMonths,
    rate:
      (offer?.annualRateBps ??
        (products.find((item) => item.id === frontendProductId(app.productId))?.rate ?? 0) * 100) /
      100,
    financials: {
      revenue: company.monthlyRevenueMinor / 100,
      expenses: company.monthlyExpensesMinor / 100,
      debts: company.monthlyDebtPaymentsMinor / 100,
    },
    companyProfile: {
      registrationStatus: company.verificationStatus,
      director: company.representativeName || '',
      founders: [],
      okved: [],
      bankruptcy: null,
      convictions: null,
      activeLoans: company.monthlyDebtPaymentsMinor > 0,
    },
    date: app.createdAt,
    status: statuses[app.status] || 'На рассмотрении',
    scoring: details.scoring,
    applicant: {
      company: company.name,
      inn: company.inn,
      registered: company.registeredAt || '',
      name: company.representativeName || '',
      position: company.representativePosition || '',
      phone: company.phone || '',
      email: company.email || '',
    },
    documents: (details.documents || []).map((item) => item.fileName),
    apiDocuments: details.documents || [],
    documentFiles: Object.fromEntries(
      (details.documents || []).map((item) => [item.requestItemId || item.type, item.fileName]),
    ),
    ...(information && app.status === 'NEEDS_INFO'
      ? {
          documentRequest: {
            items: information.items.map((item) => item.id),
            other: information.message,
            requestedAt: information.createdAt,
          },
          information,
        }
      : {}),
    ...(offer
      ? {
          offer,
          alternative: {
            productId: frontendProductId(offer.productId),
            amount: offer.amountMinor / 100,
            months: offer.termMonths,
            rate: offer.annualRateBps / 100,
            collateral: offer.collateral || '',
            comment: '',
          },
        }
      : {}),
  }
}

export async function loadApplications(staff: boolean) {
  const response = await request<{ items: ApiApplication[] }>(
    `${staff ? '/bank' : ''}/applications?limit=100`,
  )
  return Promise.all(response.items.map((item) => overview(item, staff)))
}

export async function uploadDocuments(
  id: string,
  names: Record<string, string>,
  information?: ApiInformation,
) {
  let app = await request<ApiApplication>(`/applications/${id}`)
  const docs = await request<ApiDocument[]>(`/applications/${id}/documents`)
  for (const document of docs) {
    if (!Object.hasOwn(names, document.requestItemId || document.type)) {
      await request(
        `/applications/${id}/documents/${document.id}`,
        { method: 'DELETE' },
        app.version,
      )
      app = await request<ApiApplication>(`/applications/${id}`)
    }
  }
  for (const [key, name] of Object.entries(names)) {
    const file = selectedFiles.get(key)
    const existing = docs.find((item) => (item.requestItemId || item.type) === key)
    if (!file || file.name !== name) {
      if (existing?.fileName === name) continue
      throw new Error(
        `Выберите файл «${name}» повторно: после обновления страницы браузер не сохраняет его содержимое.`,
      )
    }
    if (existing) {
      await request(
        `/applications/${id}/documents/${existing.id}`,
        { method: 'DELETE' },
        app.version,
      )
      app = await request<ApiApplication>(`/applications/${id}`)
    }
    const item = information?.items.find((item) => item.id === key)
    const body = new FormData()
    body.set('type', item?.type || (/^[A-Z_]+$/.test(key) ? key : 'OTHER'))
    body.set('file', file)
    if (item) body.set('requestItemId', item.id)
    await request(`/applications/${id}/documents`, { method: 'POST', body }, app.version)
    selectedFiles.delete(key)
    app = await request<ApiApplication>(`/applications/${id}`)
  }
  return app
}

export async function downloadDocument(applicationId: string, document: ApiDocument) {
  const response = await fetch(
    `${base}/applications/${applicationId}/documents/${document.id}/download`,
    { headers: { Authorization: `Bearer ${session?.accessToken || ''}` } },
  )
  if (!response.ok) throw new ApiError(response.status, 'Не удалось скачать документ')
  const url = URL.createObjectURL(await response.blob())
  const link = window.document.createElement('a')
  link.href = url
  link.download = document.fileName
  window.document.body.appendChild(link)
  link.click()
  link.remove()
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}

export function useCalculation(productId: string, amount: number, months: number, enabled = true) {
  const key = `${productId}/${amount}/${months}`
  const [result, setResult] = useState<{
    key: string
    monthly: number
    total: number
    overpayment: number
    disclaimer: string
  } | null>(null)
  useEffect(() => {
    if (!enabled || amount <= 0 || months <= 0) return
    const controller = new AbortController()
    const timer = window.setTimeout(() => {
      void request<{
        monthlyPaymentMinor: number
        totalPaymentMinor: number
        overpaymentMinor: number
        disclaimer: string
      }>('/calculations', {
        method: 'POST',
        signal: controller.signal,
        body: JSON.stringify({
          productId: backendProductId(productId),
          amountMinor: Math.round(amount * 100),
          termMonths: months,
        }),
      })
        .then((value) =>
          setResult({
            key,
            monthly: value.monthlyPaymentMinor / 100,
            total: value.totalPaymentMinor / 100,
            overpayment: value.overpaymentMinor / 100,
            disclaimer: value.disclaimer,
          }),
        )
        .catch(() => {})
    }, 250)
    return () => {
      window.clearTimeout(timer)
      controller.abort()
    }
  }, [productId, amount, months, enabled, key])
  return result?.key === key ? result : null
}
