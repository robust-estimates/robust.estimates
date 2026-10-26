import { creditDetails } from './creditData'

export type Product = {
  id: string
  name: string
  description: string
  icon: string
  color: string
  categories: string[]
  rate: number
  amount: number
  months: number
  secured: boolean
  limit?: boolean
  term?: string
  termTitle?: string
}
export const categories = [
  'Все продукты',
  'Популярные',
  'Без залога',
  'С залогом',
  'Возобновляемые',
  'Рефинансирование',
]
export const products: Product[] = [
  {
    id: 'business',
    name: 'Бизнес-кредит без залога',
    description: 'На любые бизнес-цели без залога\nи длительного оформления',
    icon: 'coin',
    color: 'pink',
    categories: ['Популярные', 'Без залога'],
    rate: 14,
    amount: 50,
    months: 60,
    secured: false,
  },
  {
    id: 'secured',
    name: 'Бизнес-кредит с залогом',
    description: 'Более крупные суммы и выгодные условия',
    icon: 'home',
    color: 'gray',
    categories: ['С залогом'],
    rate: 14,
    amount: 150,
    months: 120,
    secured: true,
  },
  {
    id: 'credit-line',
    name: 'Кредитная линия',
    description: 'Возобновляемый лимит для\nежедневных задач бизнеса',
    icon: 'refresh',
    color: 'pink',
    categories: ['Популярные', 'Без залога', 'Возобновляемые'],
    rate: 14.5,
    amount: 30,
    months: 18,
    secured: false,
    limit: true,
    term: '18 месяцев',
    termTitle: 'Договор',
  },
  {
    id: 'overdraft',
    name: 'Овердрафт',
    description: 'Краткосрочное финансирование\nдля покрытия кассовых разрывов',
    icon: 'wallet',
    color: 'purple',
    categories: ['Без залога', 'Возобновляемые'],
    rate: 13,
    amount: 20,
    months: 15,
    secured: false,
    limit: true,
    term: '15 месяцев',
    termTitle: 'Договор',
  },
  {
    id: 'mortgage',
    name: 'Бизнес-ипотека',
    description: 'Помещения, офисы и коммерческая\nнедвижимость для вашего бизнеса',
    icon: 'building1',
    color: 'green',
    categories: ['С залогом'],
    rate: 23,
    amount: 150,
    months: 120,
    secured: true,
  },
  {
    id: 'refinance',
    name: 'Рефинансирование',
    description: 'Объедините текущие кредиты\nв один с более выгодными условиями',
    icon: 'percent',
    color: 'blue',
    categories: ['Рефинансирование'],
    rate: 19,
    amount: 30,
    months: 60,
    secured: false,
  },
  {
    id: 'sellers',
    name: 'Кредит для селлеров',
    description: 'Специальные условия для продавцов\nна маркетплейсах',
    icon: 'cart',
    color: 'yellow',
    categories: ['Без залога'],
    rate: 1.2,
    amount: 20,
    months: 12,
    secured: false,
  },
  {
    id: 'business-one',
    name: 'Business One Credit',
    description: 'Комплексное решение для\nбизнеса',
    icon: 'brief',
    color: 'gray',
    categories: ['Без залога'],
    rate: 0,
    amount: 1,
    months: 4,
    secured: false,
    term: 'до 120 дней',
    termTitle: 'льготный период',
  },
]
export const money = (value: number) =>
  new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 0 })
    .format(value)
    .replaceAll('\u00a0', ' ')
export const rateLabel = (rate: number) => `${rate.toLocaleString('ru-RU')}%`
export function formatTerm(months: number) {
  if (months % 12 === 0) {
    const years = months / 12
    return `${years} ${years === 1 ? 'год' : years < 5 ? 'года' : 'лет'}`
  }
  const ending =
    months % 100 >= 11 && months % 100 <= 14
      ? 'месяцев'
      : months % 10 === 1
        ? 'месяц'
        : months % 10 >= 2 && months % 10 <= 4
          ? 'месяца'
          : 'месяцев'
  return `${months} ${ending}`
}
export function monthlyPayment(amount: number, annualRate: number, months: number) {
  if (!amount || months <= 0) return 0
  if (!annualRate) return amount / months
  const rate = annualRate / 1200
  return (amount * rate) / (1 - Math.pow(1 + rate, -months))
}

export function loanFigures(amount: number, rate: number, months: number, monthlyRate = false) {
  if (monthlyRate) {
    const overpayment = Math.round(((amount * rate) / 100) * months)
    return {
      monthly: Math.round((amount + overpayment) / months),
      total: amount + overpayment,
      overpayment,
    }
  }
  // The prepared demo quote is supplied by the mockup; use it consistently on every screen.
  const monthly =
    amount === 15_000_000 && rate === 19 && months === 48
      ? 435_750
      : Math.round(monthlyPayment(amount, rate, months))
  const total = rate === 0 ? amount : Math.max(amount, monthly * months)
  return { monthly, total, overpayment: total - amount }
}

export type Quote = { productId: string; amount: number; months: number; applicationId?: string }
export function quoteForProduct(product: Product, input?: Quote): Quote {
  return {
    productId: product.id,
    amount: Math.min(
      product.amount * 1_000_000,
      Math.max(creditDetails[product.id].minAmount, input?.amount ?? 15_000_000),
    ),
    months: Math.min(
      product.months,
      Math.max(creditDetails[product.id].minMonths, Math.round(input?.months ?? 48)),
    ),
  }
}
