import { expect, test } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { fillApplication, newClient, password, pdf, submitApplication } from './helpers'

test('регистрация и повторный вход используют backend', async ({ page }) => {
  const email = `ui.${crypto.randomUUID()}@example.com`
  await page.goto('/#/catalog')
  await page.getByRole('button', { name: /Нет аккаунта/ }).click()
  await page.getByLabel('Имя', { exact: true }).fill('Проверка регистрации')
  await page.getByLabel('Email', { exact: true }).fill(email)
  await page.getByLabel('Пароль', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Зарегистрироваться', exact: true }).click()
  await expect(page.locator('.product-card')).toHaveCount(8)
  await page.getByRole('button', { name: /Меню профиля/ }).click()
  await page.getByRole('button', { name: 'Выйти', exact: true }).click()
  await page.getByLabel('Email', { exact: true }).fill(email)
  await page.getByLabel('Пароль', { exact: true }).fill(password)
  await page.getByRole('button', { name: 'Войти', exact: true }).last().click()
  await expect(page.locator('.product-card')).toHaveCount(8)
})

test('неверный пароль не создаёт локальную сессию', async ({ page }) => {
  await page.goto('/#/login')
  await page.getByLabel('Email', { exact: true }).fill('absent@example.com')
  await page.getByLabel('Пароль', { exact: true }).fill('Wrong-Credit-2026!')
  await page.locator('.login-card').getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText('Неверный логин или пароль')
  await expect(page.locator('.product-card')).toHaveCount(0)
})

test('заявка проходит UI → S3 → scoring → выдача и сохраняется после обновления', async ({
  page,
}) => {
  const { session } = await newClient(page)
  const browserErrors: string[] = []
  page.on('pageerror', (error) => browserErrors.push(error.message))
  await fillApplication(page)
  let lostResponse = false
  const commandKeys: string[] = []
  await page.route('**/api/v1/applications', async (route) => {
    if (route.request().method() !== 'POST') return route.continue()
    commandKeys.push(route.request().headers()['idempotency-key'])
    if (!lostResponse) {
      lostResponse = true
      await route.fetch()
      await route.abort('failed')
    } else await route.continue()
  })
  const id = await submitApplication(page)
  expect(commandKeys).toHaveLength(2)
  expect(commandKeys[0]).toBe(commandKeys[1])
  await expect(page.locator('.status-badge')).toHaveText('Одобрено', { timeout: 30_000 })
  await page.reload()
  await expect(page.locator('.status-badge')).toHaveText('Одобрено')
  const downloading = page.waitForEvent('download')
  await page.getByRole('button', { name: 'Скачать', exact: true }).first().click()
  const download = await downloading
  const stream = await download.createReadStream()
  const chunks: Buffer[] = []
  for await (const chunk of stream!) chunks.push(Buffer.from(chunk))
  expect(Buffer.concat(chunks)).toEqual(pdf)
  await page.getByRole('button', { name: 'Перейти к подписанию' }).click()
  await page.getByRole('button', { name: 'Подтвердить и получить кредит' }).click()
  await expect(page.locator('.status-badge')).toHaveText('Выдано')
  const response = await page.request.get(`/api/v1/applications/${id}`, {
    headers: { Authorization: `Bearer ${session.accessToken}` },
  })
  expect((await response.json()).status).toBe('ISSUED')
  expect(browserErrors).toEqual([])
})

test('неподдерживаемые файлы блокируются до отправки', async ({ page }) => {
  await newClient(page)
  await page.goto('/#/application/business')
  await page.getByLabel('Файл: Финансовая отчётность', { exact: true }).setInputFiles({
    name: 'data.docx',
    mimeType: 'application/octet-stream',
    buffer: Buffer.from('invalid'),
  })
  await expect(page.getByRole('alert')).toContainText('Поддерживаются PDF, PNG и JPEG')
})

test('каталог берёт условия из API и доступен с клавиатуры', async ({ page }) => {
  const { session } = await newClient(page)
  const response = await page.request.get('/api/v1/products/unsecured', {
    headers: { Authorization: `Bearer ${session.accessToken}` },
  })
  const product = await response.json()
  const card = page.locator('[data-product-id="business"]')
  await expect(card).toContainText(`${product.annualRateBps / 100}`.replace('.', ','))
  await card.getByRole('button').click()
  await expect(page.getByRole('dialog')).toContainText('Демонстрационные условия каталога')
  await page.keyboard.press('Escape')
  await expect(page.getByRole('dialog')).toHaveCount(0)
  const audit = await new AxeBuilder({ page }).withTags(['wcag2a', 'wcag2aa']).analyze()
  expect(audit.violations.filter((item) => item.impact === 'critical')).toEqual([])
})

test('клиент не получает банковскую очередь через URL', async ({ page }) => {
  await newClient(page)
  await page.goto('/#/specialist')
  await expect(page.getByRole('heading', { name: 'Страница не найдена' })).toBeVisible()
  await expect(page.locator('.specialist-list')).toHaveCount(0)
})

test('ручная проверка: сотрудник запрашивает документ, клиент загружает его, сотрудник одобряет', async ({
  page,
  browser,
  baseURL,
}) => {
  const { session } = await newClient(page)
  await fillApplication(page, { overdue: true })
  const id = await submitApplication(page)
  const headers = { Authorization: `Bearer ${session.accessToken}` }
  await expect
    .poll(
      async () =>
        (await (await page.request.get(`/api/v1/applications/${id}`, { headers })).json()).status,
    )
    .toBe('MANUAL_CHECK')
  const login = await page.request.post('/api/v1/auth/login', {
    data: {
      login: 'employee@demo.local',
      password: process.env.DEMO_PASSWORD || 'Demo-Credit-2026!',
    },
  })
  expect(login.status()).toBe(200)
  const employee = await login.json()
  const context = await browser.newContext({ baseURL })
  await context.addInitScript(
    (value) => sessionStorage.setItem('credit-api-session', JSON.stringify(value)),
    employee,
  )
  const staff = await context.newPage()
  await staff.goto(`/#/specialist/${id}`)
  await expect(staff.locator('.score-circle strong')).not.toHaveText('0')
  await staff.getByRole('button', { name: 'Подробнее о скоринге' }).click()
  await expect(staff.locator('.scoring-popover')).toContainText(
    'Нет кредитной истории собственника',
  )
  await staff.getByRole('button', { name: 'Решение', exact: true }).click()
  await staff.getByRole('button', { name: 'Запросить данные', exact: true }).click()
  await staff.getByRole('checkbox', { name: /Годовая финансовая отчётность/ }).check()
  await staff.getByRole('button', { name: 'Отправить запрос', exact: true }).click()
  await expect
    .poll(
      async () =>
        (await (await page.request.get(`/api/v1/applications/${id}`, { headers })).json()).status,
    )
    .toBe('NEEDS_INFO')
  await page.reload()
  await page.getByRole('button', { name: 'Добавить документы', exact: true }).click()
  await page
    .getByLabel('Файл: Финансовая отчётность', { exact: true })
    .last()
    .setInputFiles({ name: 'additional.pdf', mimeType: 'application/pdf', buffer: pdf })
  await page.getByRole('button', { name: 'Сохранить изменения', exact: true }).click()
  await expect(page).toHaveURL(/#\/accepted\//)
  await expect
    .poll(
      async () =>
        (await (await page.request.get(`/api/v1/applications/${id}`, { headers })).json()).status,
    )
    .toBe('MANUAL_CHECK')
  await staff.reload()
  await staff.getByRole('button', { name: 'Решение', exact: true }).click()
  await staff.getByRole('button', { name: /Одобрить/ }).click()
  await staff.getByRole('button', { name: 'Подтвердить', exact: true }).click()
  await expect
    .poll(
      async () =>
        (await (await page.request.get(`/api/v1/applications/${id}`, { headers })).json()).status,
    )
    .toBe('APPROVED')
  await page.goto(`/#/applications/${id}`)
  await expect(page.locator('.status-badge')).toHaveText('Одобрено')
  await context.close()
})

test('альтернативное предложение принимается по offerId банка', async ({ page }) => {
  await newClient(page)
  await fillApplication(page)
  await page.getByLabel('Выручка, ₽', { exact: true }).fill('300000')
  await page.getByLabel('Расходы, ₽', { exact: true }).fill('150000')
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('5000000')
  await submitApplication(page)
  await expect(page.locator('.status-badge')).toHaveText('Другие условия', { timeout: 30_000 })
  await page.getByRole('button', { name: 'Ознакомиться', exact: true }).click()
  await expect(page.getByRole('dialog')).toContainText('1 800 000')
  await page.getByRole('button', { name: 'Принять предложение', exact: true }).click()
  await expect(page.locator('.status-badge')).toHaveText('Принято')
  await page.reload()
  await expect(page.locator('.status-badge')).toHaveText('Принято')
})

test('первый персональный подбор сохраняет компанию и получает рекомендации API', async ({
  page,
}) => {
  const { session } = await newClient(page)
  await page.goto('/#/results')
  await expect(page.getByRole('heading', { name: 'Данные компании', exact: true })).toBeVisible()
  await page.getByLabel('Наименование компании', { exact: true }).fill('ООО Персональный подбор')
  await page.getByLabel('ИНН*', { exact: true }).fill('7700000077')
  await page.getByLabel('ФИО *', { exact: true }).fill('Тестовый Предприниматель')
  await page.getByLabel('Должность', { exact: true }).fill('Директор')
  await page.getByLabel('Мобильный телефон *', { exact: true }).fill('+79990000000')
  await page.getByLabel('Email *', { exact: true }).fill('selection@example.com')
  await page.getByLabel('Дата регистрации*', { exact: true }).fill('2020-01-01')
  await page.getByLabel('Дата регистрации*', { exact: true }).press('Tab')
  await page.getByLabel('Выручка, ₽', { exact: true }).fill('3000000')
  await page.getByLabel('Расходы, ₽', { exact: true }).fill('1500000')
  await page.getByRole('button', { name: 'Сохранить компанию', exact: true }).click()
  await expect(
    page.getByRole('region', { name: 'Наиболее подходящие' }).locator('.product-card').first(),
  ).toBeVisible()
  const response = await page.request.get('/api/v1/companies', {
    headers: { Authorization: `Bearer ${session.accessToken}` },
  })
  const companies = await response.json()
  expect(companies).toHaveLength(1)
  expect(companies[0].verificationStatus).toBe('VERIFIED')
  expect(companies[0].monthlyRevenueMinor).toBe(300000000)
})

test('финансовые поля очищаются, цель Business One согласована с API, подписи помещаются', async ({
  page,
}) => {
  await newClient(page)
  for (const width of [1100, 390]) {
    await page.setViewportSize({ width, height: 900 })
    for (const chip of await page.locator('.category-chip').all()) {
      expect(await chip.evaluate((element) => element.scrollWidth <= element.clientWidth + 1)).toBe(
        true,
      )
    }
  }
  await page.setViewportSize({ width: 1280, height: 1000 })
  await fillApplication(page)
  const expenses = page.getByLabel('Расходы, ₽', { exact: true })
  await expenses.fill('')
  await expect(expenses).toHaveValue('')
  await expenses.pressSequentially('052000')
  await expect(expenses).toHaveValue('52000')
  await expenses.fill('1500000')
  await page.goto('/#/application/business-one')
  await expect(page.getByLabel('Цель кредита', { exact: true })).toHaveValue('BUSINESS_EXPANSION')
  await page.getByLabel('Наименование компании', { exact: true }).fill('ООО Проверка цели')
  await page.getByLabel('ИНН*', { exact: true }).fill('7700000088')
  await page.getByLabel('ФИО *', { exact: true }).fill('Предприниматель')
  await page.getByLabel('Мобильный телефон *', { exact: true }).fill('+79990000000')
  await page.getByLabel('Email *', { exact: true }).fill('browser@example.com')
  await page.getByLabel('Дата регистрации*', { exact: true }).fill('2020-01-01')
  await page.getByLabel('Дата регистрации*', { exact: true }).press('Tab')
  await page.getByLabel('Выручка, ₽', { exact: true }).fill('6700000')
  await page.getByLabel('Расходы, ₽', { exact: true }).fill('2052000')
  await page.getByLabel('Платежи по кредитам, ₽', { exact: true }).fill('30000')
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('3150000')
  await page.getByRole('spinbutton', { name: 'Срок кредита, месяцев', exact: true }).fill('48')
  await page
    .getByRole('checkbox', {
      name: 'Согласие на обработку персональных данных и проверку кредитной истории*',
    })
    .check()
  for (const title of ['Финансовая отчётность', 'Выписка по счёту']) {
    await page
      .getByLabel(`Файл: ${title}`, { exact: true })
      .setInputFiles({ name: `${title}.pdf`, mimeType: 'application/pdf', buffer: pdf })
  }
  const consent = await page.locator('.consent').boundingBox()
  const summary = await page.locator('.application-summary').boundingBox()
  expect(consent!.x + consent!.width).toBeLessThanOrEqual(summary!.x)
  const id = await submitApplication(page)
  const response = await page.request.get(`/api/v1/applications/${id}`, {
    headers: {
      Authorization: `Bearer ${(await page.evaluate(() => JSON.parse(sessionStorage.getItem('credit-api-session')!))).accessToken}`,
    },
  })
  expect(response.ok()).toBe(true)
  expect((await response.json()).need.purpose).toBe('BUSINESS_EXPANSION')
})
