import type { ApiDocument, ApiOffer, ApiScoring, ApiInformation } from './api'
import { products } from './products'
import type { Quote } from './products'

export type Company = { name: string; inn: string; registered: string }
export type Selection = {
  purpose: string
  other: string
  registration: string
  age: string
  collateral: string
  turnover: number
  industry: string
  loans: string
  amount: number
  months: number
}
export const defaultSelection: Selection = {
  purpose: '',
  other: '',
  registration: '',
  age: '',
  collateral: 'Нет',
  turnover: 2_500_000,
  industry: '',
  loans: 'Да',
  amount: 12_500_000,
  months: 48,
}
export type Applicant = {
  company: string
  inn: string
  registered: string
  name: string
  position: string
  phone: string
  email: string
}
export const defaultApplicant: Applicant = {
  company: '',
  inn: '',
  registered: '',
  name: '',
  position: '',
  phone: '',
  email: '',
}
export type ApplicationStatus =
  | 'На рассмотрении'
  | 'Одобрено'
  | 'Отклонено'
  | 'Черновик'
  | 'Другие условия'
  | 'Запрошены данные'
  | 'Принято'
  | 'Выдано'
  | 'Отменено'
export type LoanApplication = {
  number?: string
  backendStatus?: string
  version?: number
  companyId?: string
  purposeCode?: string
  description?: string
  scoring?: ApiScoring
  apiDocuments?: ApiDocument[]
  offer?: ApiOffer
  information?: ApiInformation
  financials?: { revenue: number; expenses: number; debts: number }
  id: string
  productId: string
  amount: number
  months: number
  rate: number
  date: string
  status: ApplicationStatus
  applicant: Applicant
  documents: string[]
  documentFiles?: Record<string, string>
  documentRequest?: { items: string[]; other: string; requestedAt: string; fulfilledAt?: string }
  companyProfile?: {
    registrationStatus: string
    director: string
    founders: string[]
    okved: string[]
    bankruptcy: boolean | null
    convictions: boolean | null
    activeLoans: boolean | null
  }
  alternative?: {
    productId: string
    amount: number
    months: number
    rate: number
    collateral: string
    comment: string
  }
}
const statuses: ApplicationStatus[] = [
  'На рассмотрении',
  'Одобрено',
  'Отклонено',
  'Черновик',
  'Другие условия',
  'Запрошены данные',
  'Принято',
  'Выдано',
  'Отменено',
]
function record(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}
function text(value: unknown, max = 1000): value is string {
  return typeof value === 'string' && value.length <= max
}
function finite(value: unknown, min: number, max: number): value is number {
  return typeof value === 'number' && Number.isFinite(value) && value >= min && value <= max
}
function applicantFields(value: unknown): value is Applicant {
  return (
    record(value) &&
    ['company', 'inn', 'registered', 'name', 'position', 'phone', 'email'].every((key) =>
      text(value[key]),
    )
  )
}
export function isDocumentNames(value: unknown): value is string[] {
  return Array.isArray(value) && value.every((name) => text(name, 255) && name.trim().length > 0)
}
export function isApplicant(value: unknown): value is Applicant {
  return (
    applicantFields(value) &&
    /^\d{10}(?:\d{2})?$/.test(value.inn) &&
    value.name.trim().length > 0 &&
    /^\d{4}-\d{2}-\d{2}$/.test(value.registered) &&
    Number.isFinite(Date.parse(value.registered)) &&
    /^(?:\+7|8)\d{10}$/.test(value.phone.replace(/[\s()-]/g, '')) &&
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.email)
  )
}
function isTerms(value: unknown): value is NonNullable<LoanApplication['alternative']> {
  if (!record(value)) return false
  const product = products.find((item) => item.id === value.productId)
  return (
    !!product &&
    finite(value.amount, 1, 500_000_000) &&
    finite(value.months, 1, 180) &&
    Number.isInteger(value.months) &&
    finite(value.rate, 0, 100) &&
    text(value.collateral, 100) &&
    text(value.comment, 500)
  )
}
export function isLoanApplication(value: unknown): value is LoanApplication {
  if (!record(value)) return false
  const product = products.find((item) => item.id === value.productId)
  return (
    !!product &&
    text(value.id, 80) &&
    /^[\p{L}\d-]+$/u.test(value.id) &&
    finite(value.amount, 1, 500_000_000) &&
    finite(value.months, 1, 180) &&
    Number.isInteger(value.months) &&
    finite(value.rate, 0, 100) &&
    text(value.date, 50) &&
    Number.isFinite(Date.parse(value.date)) &&
    statuses.includes(value.status as ApplicationStatus) &&
    isApplicant(value.applicant) &&
    isDocumentNames(value.documents) &&
    (value.documentFiles === undefined || isDocumentFiles(value.documentFiles)) &&
    (value.documentRequest === undefined || isDocumentRequest(value.documentRequest)) &&
    (value.companyProfile === undefined || isCompanyProfile(value.companyProfile)) &&
    (value.alternative === undefined || isTerms(value.alternative))
  )
}
export function isSelection(value: unknown): value is Selection {
  return (
    record(value) &&
    ['purpose', 'registration', 'age', 'collateral', 'industry', 'loans'].every((key) =>
      text(value[key], 100),
    ) &&
    text(value.other, 300) &&
    finite(value.turnover, 0, 999_999_999) &&
    finite(value.amount, 50_000, 50_000_000) &&
    finite(value.months, 6, 120) &&
    Number.isInteger(value.months)
  )
}
export function isQuote(value: unknown): value is Quote {
  if (!record(value)) return false
  const product = products.find((item) => item.id === value.productId)
  return (
    !!product &&
    finite(value.amount, 1, product.amount * 1_000_000) &&
    finite(value.months, 1, product.months) &&
    Number.isInteger(value.months) &&
    (value.applicationId === undefined || text(value.applicationId, 80))
  )
}
export type ApplicationDraft = {
  applicant: Applicant
  documents: string[]
  documentFiles?: Record<string, string>
  amount?: number
  months?: number
}
export function isDocumentFiles(value: unknown): value is Record<string, string> {
  return (
    record(value) &&
    Object.keys(value).length <= 50 &&
    Object.entries(value).every(
      ([key, filename]) => text(key, 100) && text(filename, 255) && filename.trim().length > 0,
    )
  )
}
function isDocumentRequest(
  value: unknown,
): value is NonNullable<LoanApplication['documentRequest']> {
  return (
    record(value) &&
    Array.isArray(value.items) &&
    value.items.every((item) =>
      ['passport', 'annual-report', 'current-loans', 'other'].includes(String(item)),
    ) &&
    text(value.other, 500) &&
    text(value.requestedAt, 50) &&
    Number.isFinite(Date.parse(value.requestedAt)) &&
    (value.fulfilledAt === undefined ||
      (text(value.fulfilledAt, 50) && Number.isFinite(Date.parse(value.fulfilledAt))))
  )
}
function isCompanyProfile(value: unknown): value is NonNullable<LoanApplication['companyProfile']> {
  return (
    record(value) &&
    text(value.registrationStatus, 255) &&
    text(value.director, 255) &&
    isDocumentNames(value.founders) &&
    isDocumentNames(value.okved) &&
    ['bankruptcy', 'convictions', 'activeLoans'].every(
      (key) => value[key] === null || typeof value[key] === 'boolean',
    )
  )
}
export function isApplicationDraft(value: unknown): value is ApplicationDraft {
  return (
    record(value) &&
    applicantFields(value.applicant) &&
    isDocumentNames(value.documents) &&
    (value.documentFiles === undefined || isDocumentFiles(value.documentFiles)) &&
    (value.amount === undefined || finite(value.amount, 1, 500_000_000)) &&
    (value.months === undefined || (finite(value.months, 1, 180) && Number.isInteger(value.months)))
  )
}
