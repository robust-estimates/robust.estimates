import { expect } from '@playwright/test'
import type { Page } from '@playwright/test'

export const password = 'Integration-Credit-2026!'
export const pdf = Buffer.from('%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\n%%EOF\n')

export async function newClient(page: Page) {
  const login = `browser.${crypto.randomUUID()}@example.com`
  const response = await page.request.post('/api/v1/auth/register', {
    data: { login, displayName: 'Тестовый предприниматель', password },
  })
  expect(response.status()).toBe(201)
  const session = await response.json()
  await page.addInitScript(
    (value) => sessionStorage.setItem('credit-api-session', JSON.stringify(value)),
    session,
  )
  await page.goto('/#/catalog')
  await expect(page.locator('.product-card')).toHaveCount(8)
  return { login, session }
}

export async function fillApplication(page: Page, options: { inn?: string } = {}) {
  await page.goto('/#/application/business')
  await page.getByLabel('Выручка, ₽', { exact: true }).fill('3000000')
  await page.getByLabel('Расходы, ₽', { exact: true }).fill('1500000')
  await page.getByLabel('Платежи по кредитам, ₽', { exact: true }).fill('0')
  await page.getByLabel('Наименование компании', { exact: true }).fill('ООО Браузерная проверка')
  await page.getByLabel('ИНН*', { exact: true }).fill(options.inn ?? '7700000088')
  await page.getByLabel('ФИО *', { exact: true }).fill('Тестовый Предприниматель')
  await page.getByLabel('Должность', { exact: true }).fill('Директор')
  await page.getByLabel('Мобильный телефон *', { exact: true }).fill('+79990000000')
  await page.getByLabel('Email *', { exact: true }).fill('browser@example.com')
  await page.getByLabel('Дата регистрации*', { exact: true }).fill('2020-01-01')
  await page.getByLabel('Дата регистрации*', { exact: true }).press('Tab')
  await page
    .getByRole('checkbox', {
      name: 'Согласие на обработку персональных данных и проверку кредитной истории*',
    })
    .check()
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('1000000')
  await page.getByRole('spinbutton', { name: 'Срок кредита, месяцев', exact: true }).fill('12')
  for (const title of ['Финансовая отчётность', 'Выписка по счёту']) {
    await page
      .getByLabel(`Файл: ${title}`, { exact: true })
      .setInputFiles({ name: `${title}.pdf`, mimeType: 'application/pdf', buffer: pdf })
  }
}

export async function submitApplication(page: Page) {
  await page.getByRole('button', { name: 'Оформить заявку', exact: true }).click()
  await expect(page).toHaveURL(/#\/accepted\/[a-f\d]+$/, { timeout: 30_000 })
  const id = page.url().split('/').at(-1)!
  await page.goto(`/#/applications/${id}`)
  return id
}
