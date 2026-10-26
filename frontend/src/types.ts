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
  company: 'ООО «Лапочка»',
  inn: '7701234567',
  registered: '2018-03-15',
  name: 'Иванов Константин Казимирович',
  position: 'Директор',
  phone: '+7 (915) 123-45-67',
  email: 'ivanov@lapochka.ru',
}
export type ApplicationStatus =
  'На рассмотрении' | 'Одобрено' | 'Отклонено' | 'Черновик' | 'Другие условия' | 'Запрошены данные'
export type LoanApplication = {
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
export const defaultDocuments = ['Бухгалтерская_отчетность_2026', 'Паспортные данные']
export const seedApplications: LoanApplication[] = [
  {
    id: 'XX-000123',
    productId: 'business',
    amount: 15_000_000,
    months: 48,
    rate: 19,
    date: '2026-10-21T14:32:00+03:00',
    status: 'На рассмотрении',
    applicant: defaultApplicant,
    documents: defaultDocuments,
  },
  {
    id: 'BD-000122',
    productId: 'overdraft',
    amount: 2_000_000,
    months: 24,
    rate: 13,
    date: '2026-10-05T10:15:00+03:00',
    status: 'Одобрено',
    applicant: defaultApplicant,
    documents: defaultDocuments,
  },
  {
    id: 'CC-000121',
    productId: 'sellers',
    amount: 3_000_000,
    months: 12,
    rate: 1.2,
    date: '2026-09-01T16:20:00+03:00',
    status: 'Отклонено',
    applicant: defaultApplicant,
    documents: defaultDocuments,
  },
  {
    id: 'MM-000120',
    productId: 'mortgage',
    amount: 15_000_000,
    months: 120,
    rate: 23,
    date: '2026-08-23T11:05:00+03:00',
    status: 'Черновик',
    applicant: defaultApplicant,
    documents: defaultDocuments,
  },
]
const statuses: ApplicationStatus[] = [
  'На рассмотрении',
  'Одобрено',
  'Отклонено',
  'Черновик',
  'Другие условия',
  'Запрошены данные',
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
    finite(value.amount, 1, 150_000_000) &&
    finite(value.months, 1, 120) &&
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
    /^[\p{L}\d]+-\d+$/u.test(value.id) &&
    // Earlier demo applications keep their original terms after the PDF catalog update.
    finite(value.amount, 1, 150_000_000) &&
    finite(value.months, 1, 120) &&
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
    (value.amount === undefined || finite(value.amount, 1, 150_000_000)) &&
    (value.months === undefined || (finite(value.months, 1, 120) && Number.isInteger(value.months)))
  )
}
export function readApplications(): LoanApplication[] {
  try {
    const stored: unknown = JSON.parse(
      localStorage.getItem('alfa-credit-applications-v2') ?? 'null',
    )
    if (Array.isArray(stored)) {
      const unique = new Map<string, LoanApplication>()
      for (const item of stored)
        if (isLoanApplication(item) && !unique.has(item.id)) unique.set(item.id, item)
      if (unique.size || !stored.length) return [...unique.values()]
    }
  } catch {
    /* Storage may be unavailable in private browsing. */
  }
  return seedApplications
}
