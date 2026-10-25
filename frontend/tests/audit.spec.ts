import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

async function choose(page: Page, label: string, value: string) {
  await page.getByRole('combobox', { name: new RegExp(label) }).click()
  await page.getByRole('option', { name: value, exact: true }).click()
}
async function selectBusiness(page: Page, amount = '15000000') {
  await page.goto('/#/selection')
  await choose(page, 'тип регистрации', 'ООО')
  await choose(page, 'Сколько лет', 'Более 3 лет')
  const input = page.getByRole('textbox', { name: 'Сумма кредита', exact: true })
  await input.fill(amount)
  await input.blur()
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
}

test('параметры подбора и показатели сохраняются после перезагрузки и совпадают с формой', async ({
  page,
}) => {
  await selectBusiness(page)
  const result = page
    .getByRole('region', { name: 'Наиболее подходящие' })
    .locator('[data-product-id="business"]')
  await expect(result).toContainText('20 916 000 ₽')
  await page.reload()
  await expect(result).toContainText('20 916 000 ₽')
  await result.getByRole('button').click()
  await expect(page.locator('.indicators')).toContainText('15 000 000 ₽')
  await expect(page.locator('.indicators')).toContainText('20 916 000 ₽')
  await page.reload()
  await expect(page.locator('.indicators')).toContainText('15 000 000 ₽')
  await expect(page.locator('.indicators')).toContainText('20 916 000 ₽')
})

test('нестандартная сумма, данные представителя и документы переживают обновление формы', async ({
  page,
}) => {
  await selectBusiness(page, '8250000')
  await page
    .getByRole('button', { name: 'Выбрать продукт: Бизнес-кредит без залога', exact: true })
    .click()
  await page.getByLabel('ФИО *', { exact: true }).fill('Петров Алексей Сергеевич')
  await page.getByLabel('Файл: Паспорт', { exact: true }).setInputFiles({
    name: 'Паспорт.pdf',
    mimeType: 'application/pdf',
    buffer: Buffer.from('%PDF-1.4'),
  })
  await page.reload()
  await expect(page.locator('.indicators')).toContainText('8 250 000 ₽')
  await expect(page.getByLabel('ФИО *', { exact: true })).toHaveValue('Петров Алексей Сергеевич')
  await expect(page.locator('.document-upload.is-uploaded')).toHaveCount(1)
  await expect(page.locator('.document-upload.is-uploaded')).toContainText('Паспорт.pdf')
  await expect(page.getByRole('checkbox')).not.toBeChecked()
})

test('прямая ссылка на короткий продукт ограничивает сумму и срок его условиями', async ({
  page,
}) => {
  await page.goto('/#/application/business-one')
  await expect(page.locator('.indicators')).toContainText('1 000 000 ₽')
  await expect(page.locator('.indicators')).toContainText('4 месяца')
  const repayment = page.locator('.indicators > div').filter({ hasText: 'Переплата' })
  await expect(repayment).toContainText('0 ₽')
})

test('продукт с периодом 120 дней недоступен при подборе на 10 лет', async ({ page }) => {
  await page.goto('/#/selection')
  await choose(page, 'тип регистрации', 'ООО')
  await choose(page, 'Сколько лет', 'Более 3 лет')
  const amount = page.getByRole('textbox', { name: 'Сумма кредита', exact: true })
  await amount.fill('500000')
  await amount.blur()
  await page.getByRole('slider', { name: 'Изменить срок кредита' }).focus()
  await page.keyboard.press('End')
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
  await expect(
    page
      .getByRole('region', { name: 'Недоступные продукты' })
      .locator('[data-product-id="business-one"] button'),
  ).toBeDisabled()
})

test('неизвестные номера заявок и продуктов показывают 404 вместо чужой карточки', async ({
  page,
}) => {
  for (const route of [
    'applications/NO-99999',
    'accepted/NO-99999',
    'application/unknown',
    'application/business/NO-99999',
    'specialist/NO-99999',
    'application/business/XX-000123/extra',
  ]) {
    await page.goto(`/#/${route}`)
    await expect(page.getByRole('heading', { name: 'Страница не найдена' }), route).toBeVisible()
    await expect(page.locator('.application-detail,.application-form,.review-card')).toHaveCount(0)
  }
})

test('ссылка пропуска навигации переводит фокус и не меняет маршрут', async ({ page }) => {
  await page.goto('/#/selection')
  await page.keyboard.press('Tab')
  await expect(page.getByRole('link', { name: 'Перейти к содержимому' })).toBeFocused()
  await page.keyboard.press('Enter')
  await expect(page.locator('#main-content')).toBeFocused()
  await expect(page).toHaveURL(/#\/selection$/)
})

test('выпадающие списки выбираются с клавиатуры, меню фильтра закрывается по Escape', async ({
  page,
}) => {
  await page.goto('/#/selection')
  const purpose = page.getByRole('combobox', { name: /Для чего вашему бизнесу/ })
  await purpose.focus()
  await page.keyboard.press('ArrowDown')
  await page.keyboard.press('Enter')
  await expect(purpose).toHaveText(/Развитие бизнеса/)
  await purpose.press('ArrowDown')
  await purpose.press('End')
  await purpose.press(' ')
  await expect(page.getByRole('textbox', { name: 'Опишите цель кредита' })).toBeVisible()
  await page.goto('/#/specialist')
  const filter = page.getByRole('button', { name: 'Фильтры', exact: true })
  await filter.click()
  await page.keyboard.press('Escape')
  await expect(filter).toHaveAttribute('aria-expanded', 'false')
  await expect(filter).toBeFocused()
})

test('выход скрывает приватные страницы, вход возвращает к исходному адресу', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: /Меню профиля/ }).click()
  await page.getByRole('button', { name: 'Выйти', exact: true }).click()
  await page.goto('/#/applications/XX-000123')
  await expect(page.getByRole('heading', { name: 'Вход', exact: true })).toBeVisible()
  await expect(page.locator('.application-detail')).toHaveCount(0)
  await page.locator('.login-card').getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page).toHaveURL(/#\/applications\/XX-000123$/)
  await expect(page.getByRole('heading', { name: 'Заявка № XX-000123' })).toBeVisible()
})

test('правильный логин с неправильным длинным паролем не проходит демонстрационный вход', async ({
  page,
}) => {
  await page.goto('/#/login')
  await page.locator('input[type="password"]').fill('wrong-password')
  await page.locator('.login-card').getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page.getByRole('alert')).toHaveText('Неправильный логин или пароль')
})

test('повреждённые записи хранилища не обрушивают приложение и не дублируют валидные записи', async ({
  page,
}) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  await page.goto('/')
  await page.evaluate(() => {
    localStorage.setItem(
      'alfa-credit-applications-v2',
      JSON.stringify([
        {
          id: 'BAD-1',
          productId: 'business',
          amount: 500000,
          months: 12,
          rate: 14,
          applicant: {},
          date: 'broken',
          documents: [{}],
        },
      ]),
    )
    sessionStorage.setItem('alfa-credit-selection-v2', '{not json')
    sessionStorage.setItem(
      'alfa-credit-quote-v2',
      JSON.stringify({ productId: 'business', amount: -1, months: 0 }),
    )
  })
  await page.goto('/#/applications')
  await page.reload()
  await expect(page.locator('.application-row')).toHaveCount(4)
  await expect(page.getByRole('heading', { name: 'Мои заявки' })).toBeVisible()
  await page.goto('/#/selection')
  await expect(page.getByRole('slider', { name: 'Изменить сумму кредита' })).toHaveValue('12500000')
  expect(errors).toEqual([])
})

test('итоговое решение убирает заявку из очереди, сохраняя статус и скоринг', async ({ page }) => {
  await page.goto('/#/specialist/SK-007821')
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: /Отказать/ }).click()
  await page.getByRole('button', { name: 'Подтвердить', exact: true }).click()
  await page.getByRole('button', { name: /К списку заявок/ }).click()
  const row = page.locator('tbody tr').filter({ hasText: 'SK-007821' })
  await expect(row).toHaveCount(0)
  await page.reload()
  await expect(row).toHaveCount(0)
  await page.goto('/#/applications/SK-007821')
  await expect(page.locator('.application-detail-heading .status-badge')).toHaveText('Отклонено')
  await page.goto('/#/specialist/SK-007821')
  await expect(page.locator('.score-circle')).toContainText('601')
})

test('смена продукта в предложении банка ограничивает сумму и срок', async ({ page }) => {
  await page.goto('/#/specialist/XX-001234')
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: 'Другие условия', exact: true }).click()
  await choose(page, 'Вид кредита', 'Business One Credit')
  await expect(page.getByRole('textbox', { name: 'Сумма кредита', exact: true })).toHaveValue(
    '1 000 000 ₽',
  )
  await expect(page.getByRole('combobox', { name: 'Срок', exact: true })).toHaveText('4 месяца')
  await page.getByRole('button', { name: 'Сохранить предложение' }).click()
  await page.goto('/#/applications/XX-001234')
  await page.getByRole('button', { name: 'Ознакомиться' }).click()
  await expect(page.getByRole('dialog')).toContainText('1 000 000 ₽')
  await page.goto('/#/specialist/XX-001234')
  await page.reload()
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: 'Другие условия', exact: true }).click()
  await expect(page.getByRole('combobox', { name: 'Вид кредита', exact: true })).toHaveText(
    'Business One Credit',
  )
  await expect(page.getByRole('textbox', { name: 'Сумма кредита', exact: true })).toHaveValue(
    '1 000 000 ₽',
  )
})

test('редактирование параметров существующей заявки обновляет её, сохраняя представителя', async ({
  page,
}) => {
  await page.goto('/#/applications/XX-000123')
  await page.getByRole('button', { name: 'Изменить заявку' }).click()
  await page.getByLabel('ФИО *', { exact: true }).fill('Сидоров Иван Петрович')
  const input = page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true })
  await input.fill('9000000')
  await input.blur()
  await expect(page).toHaveURL(/#\/application\/business\/XX-000123$/)
  await expect(page.getByLabel('ФИО *', { exact: true })).toHaveValue('Сидоров Иван Петрович')
  await expect(page.locator('.indicators')).toContainText('9 000 000 ₽')
  await page.getByRole('button', { name: 'Сохранить изменения', exact: true }).click()
  await page.goto('/#/applications')
  await expect(page.locator('.application-row')).toHaveCount(4)
  await expect(page.locator('.application-row').filter({ hasText: 'XX-000123' })).toContainText(
    '9 000 000 ₽',
  )
})

test('неподдерживаемый документ и неверный телефон не отправляются', async ({ page }) => {
  await page.goto('/#/application/business')
  await page.getByLabel('Файл: Паспорт', { exact: true }).setInputFiles({
    name: 'script.exe',
    mimeType: 'application/octet-stream',
    buffer: Buffer.from('demo'),
  })
  await expect(page.getByRole('alert')).toContainText('Поддерживаются PDF')
  await expect(page.locator('.document-upload.is-uploaded')).toHaveCount(0)
  const phone = page.getByLabel('Мобильный телефон *', { exact: true })
  await phone.fill('abc')
  await page.getByRole('checkbox').check()
  await page.getByRole('button', { name: 'Оформить заявку' }).click()
  await expect(phone).toBeFocused()
  await expect(page).toHaveURL(/#\/application\/business$/)
})

test('семантика и подписи доступны на основных экранах; контраст макета записан отдельно', async ({
  page,
}, testInfo) => {
  test.setTimeout(90_000)
  const contrast: unknown[] = []
  async function scan(route: string) {
    await page.evaluate(() => document.fonts.ready)
    const report = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze()
    contrast.push({
      route,
      violations: report.violations.filter((item) => item.id === 'color-contrast'),
    })
    expect(
      report.violations
        .filter((item) => item.id !== 'color-contrast')
        .map((item) => ({ id: item.id, nodes: item.nodes.map((node) => node.target) })),
      route,
    ).toEqual([])
  }
  for (const route of [
    'catalog',
    'selection',
    'results',
    'application/business',
    'applications',
    'applications/XX-000123',
    'specialist',
    'specialist/XX-001234',
    'login',
  ]) {
    await page.goto(`/#/${route}`)
    await scan(route)
  }
  await page.goto('/#/selection')
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
  await scan('selection/errors')
  await page.getByRole('combobox', { name: /Для чего вашему бизнесу/ }).click()
  await scan('selection/options')
  await page.goto('/#/applications/BD-000122')
  await page.getByRole('button', { name: 'Перейти к подписанию' }).click()
  await scan('application/signing-dialog')
  await page.goto('/#/login')
  await page.locator('input[type="password"]').fill('wrong-password')
  await page.locator('.login-card').getByRole('button', { name: 'Войти', exact: true }).click()
  await scan('login/error')
  await testInfo.attach('design-color-contrast', {
    body: JSON.stringify(contrast, null, 2),
    contentType: 'application/json',
  })
})

test('ручной ввод срока понимает годы и месяцы, ползунок сохраняет их после обновления', async ({
  page,
}) => {
  await page.goto('/#/selection')
  const input = page.getByRole('textbox', { name: 'Срок кредита', exact: true })
  const slider = page.getByRole('slider', { name: 'Изменить срок кредита' })
  await input.fill('5 лет')
  await input.blur()
  await expect(slider).toHaveValue('60')
  await input.fill('37')
  await input.blur()
  await expect(input).toHaveValue('37 месяцев')
  await page.reload()
  await expect(slider).toHaveValue('37')
  await input.fill('')
  await input.blur()
  await expect(slider).toHaveValue('6')
})

test('модальное окно удерживает фокус, закрывается по Escape и возвращает фокус кнопке', async ({
  page,
}) => {
  await page.goto('/#/applications/BD-000122')
  const trigger = page.getByRole('button', { name: 'Перейти к подписанию' })
  await trigger.click()
  const dialog = page.getByRole('dialog')
  await expect(dialog).toBeVisible()
  for (let index = 0; index < 7; index++) {
    await page.keyboard.press('Tab')
    expect(await dialog.evaluate((element) => element.contains(document.activeElement))).toBe(true)
  }
  await page.keyboard.press('Escape')
  await expect(dialog).toHaveCount(0)
  await expect(trigger).toBeFocused()
  await page.goto('/#/login')
  await expect(page.getByLabel('Пароль', { exact: true })).toHaveAttribute('type', 'password')
  for (const name of ['Забыли пароль?', 'Госуслуги', /Нет аккаунта\?/]) {
    const open = page.locator('.login-card').getByRole('button', { name })
    await open.click()
    await expect(dialog).toBeVisible()
    await dialog.getByRole('button', { name: 'Понятно', exact: true }).click()
    await expect(dialog).toHaveCount(0)
    await expect(open).toBeFocused()
  }
})

test('страницы работают с областью просмотра при масштабе 200% и в альбомном окне', async ({
  page,
}) => {
  // Browser zoom at 200% halves the CSS viewport of a 1672 × 1016 window.
  for (const viewport of [
    { width: 836, height: 508 },
    { width: 900, height: 375 },
  ]) {
    await page.setViewportSize(viewport)
    for (const route of [
      'catalog',
      'selection',
      'application/business',
      'applications',
      'specialist',
    ]) {
      await page.goto(`/#/${route}`)
      expect(
        await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1),
        route,
      ).toBe(true)
      await expect(page.locator('main')).toBeVisible()
    }
  }
})

test('основные экраны загружают локальные ресурсы без ошибок JavaScript и сети', async ({
  page,
}) => {
  const failures: string[] = []
  page.on('pageerror', (error) => failures.push(error.message))
  page.on('console', (message) => {
    if (message.type() === 'error') failures.push(message.text())
  })
  page.on('requestfailed', (request) =>
    failures.push(`${request.url()}: ${request.failure()?.errorText}`),
  )
  page.on('response', (response) => {
    if (response.status() >= 400) failures.push(`${response.status()} ${response.url()}`)
  })
  for (const route of [
    'catalog',
    'selection',
    'results',
    'application/business',
    'application/business-one',
    'accepted/XX-000123',
    'applications',
    'applications/XX-000123',
    'applications/BD-000122',
    'applications/CC-000121',
    'specialist',
    'specialist/XX-001234',
    'specialist/SK-007821',
    'specialist/KS-7791',
    'login',
  ]) {
    await page.goto(`/#/${route}`)
    await page.evaluate(() => document.fonts.ready)
    await expect
      .poll(
        () =>
          page
            .locator('img')
            .evaluateAll((images) =>
              images.every((image) => image.complete && image.naturalWidth > 0),
            ),
        { message: route },
      )
      .toBe(true)
  }
  expect(failures).toEqual([])
})

test('сенсорный экран позволяет раскрыть подбор, изменить ползунок и прикрепить документ', async ({
  browser,
  browserName,
  baseURL,
}) => {
  const context = await browser.newContext({
    baseURL,
    viewport: { width: 375, height: 812 },
    hasTouch: true,
    ...(browserName !== 'firefox' ? { isMobile: true } : {}),
  })
  try {
    const page = await context.newPage()
    await page.goto('/')
    await page.getByRole('button', { name: 'Подобрать кредит' }).tap()
    const slider = page.getByRole('slider', { name: 'Изменить сумму кредита' })
    await slider.tap({ position: { x: 50, y: 12 } })
    await expect(slider).not.toHaveValue('12500000')
    await page.getByRole('button', { name: 'Помогите определить сумму', exact: false }).tap()
    await expect(page.locator('.amount-suggestion')).toBeVisible()
    await page.goto('/#/application/business')
    const attachment = page.getByRole('button', { name: 'Прикрепить: Паспорт', exact: false })
    await expect(attachment).toBeVisible()
    const chooser = page.waitForEvent('filechooser')
    await attachment.tap()
    await (
      await chooser
    ).setFiles({ name: 'Отчет.pdf', mimeType: 'application/pdf', buffer: Buffer.from('%PDF-1.4') })
    await expect(page.locator('.document-upload.is-uploaded')).toHaveCount(1)
  } finally {
    await context.close()
  }
})
