import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

async function choose(page: Page, label: string, option: string) {
  await page.getByRole('combobox', { name: new RegExp(label) }).click()
  await page.getByRole('option', { name: option, exact: true }).click()
}
async function upload(page: Page, title: string, filename = `${title}.pdf`) {
  await page.getByLabel(`Файл: ${title}`, { exact: true }).setInputFiles({
    name: filename,
    mimeType: 'application/pdf',
    buffer: Buffer.from('%PDF-1.4 demo'),
  })
}

test('навигация содержит только нужные разделы и логотип сохраняет кабинет специалиста', async ({
  page,
}) => {
  await page.goto('/#/catalog')
  await expect(page.getByRole('navigation')).toContainText('Кредиты')
  await expect(page.getByRole('navigation')).not.toContainText(
    /Эквайринг|Расчетный счет|Накопления/,
  )
  await page.goto('/#/specialist/XX-001234')
  await expect(page.getByRole('navigation', { name: 'Основная навигация' })).not.toContainText(
    /Клиенты|Отчеты|Кредиты/,
  )
  await page.getByRole('link', { name: 'Альфа-бизнес — главная' }).click()
  await expect(page).toHaveURL(/#\/specialist$/)
  await expect(page.getByRole('heading', { name: 'Заявки на рассмотрение' })).toBeVisible()
})

test('подбор содержит новые типы регистрации, возраст бизнеса и отрасли', async ({ page }) => {
  await page.goto('/#/selection')
  await page.getByRole('combobox', { name: /тип регистрации/ }).click()
  await expect(page.getByRole('option')).toHaveText([
    'ООО',
    'ИП',
    'АО публичное',
    'АО непубличное',
    'Другое',
  ])
  await page.keyboard.press('Escape')
  await page.getByRole('combobox', { name: /Сколько лет/ }).click()
  await expect(page.getByRole('option')).toHaveText([
    'Менее 9 месяцев',
    'От 9 месяцев до 1 года',
    'От 1 года до 3 лет',
    'Более 3 лет',
  ])
  await page.keyboard.press('Escape')
  await page.getByRole('combobox', { name: 'Сфера деятельности', exact: true }).click()
  for (const name of [
    'Торговля',
    'Общественное питание',
    'IT и информационные технологии',
    'Медицина и здравоохранение',
    'Образование',
    'Недвижимость',
    'Гостиничный бизнес и туризм',
  ])
    await expect(page.getByRole('option', { name, exact: true })).toBeVisible()
  await expect(page.locator('.business-extra')).not.toContainText('ⓘ')
})

test('подробности кредитов включают диапазоны, погашение, особенности и требования из PDF', async ({
  page,
}) => {
  await page.goto('/#/catalog')
  await page
    .getByRole('button', { name: 'Подробнее: Бизнес-кредит без залога', exact: true })
    .click()
  const dialog = page.getByRole('dialog')
  await expect(dialog).toContainText('300 000 — 50 000 000 ₽')
  await expect(dialog).toContainText('От 1 года до 5 лет')
  for (const name of ['Цель кредита', 'Варианты погашения', 'Особенности', 'Требования к заёмщику'])
    await expect(dialog.getByRole('heading', { name, exact: true })).toBeVisible()
  await expect(dialog).toContainText('Дифференцированные платежи')
  await page.keyboard.press('Escape')
  await page.getByRole('button', { name: 'Подробнее: Кредит для селлеров', exact: true }).click()
  await expect(dialog).toContainText('От 1,2% в месяц')
  await expect(dialog).toContainText('не годовая ставка')
})

test('заявка из каталога позволяет выбрать параметры и загрузить паспорт в подписанное поле', async ({
  page,
}) => {
  await page.goto('/#/catalog')
  await page.getByRole('button', { name: 'Подробнее: Business One Credit', exact: true }).click()
  await page
    .getByRole('dialog')
    .getByRole('button', { name: 'Оформить заявку', exact: true })
    .click()
  await expect(page.locator('.document-upload')).toHaveCount(1)
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('600000')
  await page.getByRole('spinbutton', { name: 'Срок кредита, месяцев', exact: true }).fill('3')
  await expect(page.locator('.indicators')).toContainText('600 000 ₽')
  await expect(page.locator('.indicators')).toContainText('3 месяца')
  const chooser = page.waitForEvent('filechooser')
  await page.getByRole('button', { name: 'Прикрепить: Паспорт', exact: false }).click()
  await (
    await chooser
  ).setFiles({ name: 'Паспорт.pdf', mimeType: 'application/pdf', buffer: Buffer.from('%PDF-1.4') })
  await page.reload()
  await expect(page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true })).toHaveValue(
    '600000',
  )
  await expect(page.locator('.document-upload')).toContainText('Паспорт.pdf')
  await page.getByRole('checkbox').check()
  await page.getByRole('button', { name: 'Оформить заявку', exact: true }).click()
  await expect(page).toHaveURL(/#\/accepted\/XX-000124$/)
  await page.goto('/#/applications/XX-000124')
  await expect(page.locator('.indicators')).toContainText('600 000 ₽')
  await expect(page.locator('.document-list')).toContainText('Паспорт.pdf')
})

test('комплект документов зависит от продукта, регистрации и наличия кредитов', async ({
  page,
}) => {
  await page.goto('/#/selection')
  await choose(page, 'тип регистрации', 'АО публичное')
  await choose(page, 'Сколько лет', 'Более 3 лет')
  await page.goto('/#/application/business')
  for (const title of [
    'Анкета предприятия или ИП',
    'Анкета физического лица',
    'Паспорт',
    'Годовая финансовая отчётность',
    'Действующие кредиты',
    'Выписка из реестра акционеров',
  ])
    await expect(page.getByLabel(`Файл: ${title}`, { exact: true })).toHaveCount(1)
  await page.goto('/#/application/overdraft')
  await expect(
    page.getByLabel('Файл: Заявка на предоставление овердрафта', { exact: true }),
  ).toHaveCount(1)
  await page.goto('/#/application/business-one')
  await expect(page.locator('.document-upload')).toHaveCount(1)
  await page.getByRole('checkbox').check()
  await page.getByRole('button', { name: 'Оформить заявку', exact: true }).click()
  await expect(page.getByRole('alert').first()).toContainText('Прикрепите')
  await expect(page).toHaveURL(/#\/application\/business-one$/)
})

test('редактирование сохраняет исходную заявку, отмена отбрасывает изменения', async ({ page }) => {
  await page.goto('/#/applications/XX-000123')
  await page.getByRole('button', { name: 'Изменить заявку' }).click()
  await expect(page.getByRole('button', { name: 'Оформить заявку', exact: true })).toHaveCount(0)
  await page.getByLabel('ФИО *', { exact: true }).fill('Петров Пётр Петрович')
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('8000000')
  await page.getByRole('button', { name: 'Отмена', exact: true }).click()
  await expect(page.locator('.indicators')).toContainText('15 000 000 ₽')
  await page.getByRole('button', { name: 'Изменить заявку' }).click()
  await expect(page.getByLabel('ФИО *', { exact: true })).toHaveValue(
    'Иванов Константин Казимирович',
  )
  await page.getByLabel('ФИО *', { exact: true }).fill('Петров Пётр Петрович')
  await page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true }).fill('8000000')
  await page.getByRole('button', { name: 'Сохранить изменения', exact: true }).click()
  await expect(page).toHaveURL(/#\/applications\/XX-000123$/)
  await page.reload()
  await expect(page.locator('.indicators')).toContainText('8 000 000 ₽')
  await expect(page.getByLabel('ФИО ', { exact: true })).toHaveValue('Петров Пётр Петрович')
  await page.goto('/#/applications')
  await expect(page.locator('.application-row')).toHaveCount(4)
})

test('финансы и проверки показывают сведения о компании вместо расчёта кредита', async ({
  page,
}) => {
  await page.goto('/#/specialist/XX-001234')
  await page.getByRole('button', { name: 'Финансы', exact: true }).click()
  for (const label of [
    'Годовая отчётность',
    'Проходит процедуру банкротства',
    'Отсутствие судимостей',
    'Действующие кредиты',
  ])
    await expect(page.locator('.review-company-facts')).toContainText(label)
  await page.getByRole('button', { name: 'Проверки', exact: true }).click()
  for (const label of [
    'Дата регистрации',
    'Статус компании',
    'Руководитель',
    'Учредители',
    'ОКВЭД',
  ])
    await expect(page.locator('.review-company-facts')).toContainText(label)
})

test('запрос нескольких документов и своего варианта проходит от специалиста к клиенту', async ({
  page,
}) => {
  await page.goto('/#/specialist/XX-001234')
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: 'Запросить данные', exact: true }).click()
  await page.getByRole('button', { name: 'Отправить запрос', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText('хотя бы один')
  await page.getByRole('checkbox', { name: 'Паспорт', exact: true }).check()
  await page.getByRole('checkbox', { name: 'Годовая финансовая отчётность', exact: true }).check()
  await page.getByRole('checkbox', { name: 'Другое', exact: true }).check()
  await page.getByLabel('Какие ещё документы нужны?', { exact: true }).fill('Выписка за 6 месяцев')
  await page.getByRole('button', { name: 'Отправить запрос', exact: true }).click()
  await page.goto('/#/applications/XX-001234')
  await page.reload()
  await expect(page.locator('.decision-requested')).toContainText('Паспорт')
  await expect(page.locator('.decision-requested')).toContainText('Выписка за 6 месяцев')
  await page.getByRole('button', { name: 'Добавить документы', exact: true }).click()
  await upload(page, 'Паспорт', 'Новый паспорт.pdf')
  await upload(page, 'Годовая финансовая отчётность', 'Годовая.pdf')
  await upload(page, 'Выписка за 6 месяцев', 'Выписка.pdf')
  await page.getByRole('button', { name: 'Сохранить изменения', exact: true }).click()
  await expect(page).toHaveURL(/#\/applications\/XX-001234$/)
  await expect(page.locator('.application-detail-heading .status-badge')).toHaveText(
    'На рассмотрении',
  )
  await expect(page.locator('.document-list')).toContainText('Выписка.pdf')
  await page.goto('/#/specialist')
  await expect(page.locator('tbody tr').filter({ hasText: 'XX-001234' })).toHaveCount(1)
})

test('форма ограничивает новые параметры, а редактирование сохраняет прежние условия заявки', async ({
  page,
}) => {
  await page.goto('/#/application/business')
  const amount = page.getByRole('spinbutton', { name: 'Сумма кредита, ₽', exact: true })
  await expect(amount).toHaveAttribute('min', '300000')
  await expect(amount).toHaveAttribute('max', '50000000')
  await amount.fill('100000')
  expect(
    await amount.evaluate((element) => (element as HTMLInputElement).validity.rangeUnderflow),
  ).toBe(true)
  await page.goto('/#/application/overdraft/BD-000122')
  await expect(
    page.getByRole('spinbutton', { name: 'Срок кредита, месяцев', exact: true }),
  ).toHaveValue('24')
  await page.getByRole('button', { name: 'Сохранить изменения', exact: true }).click()
  await expect(page).toHaveURL(/#\/applications\/BD-000122$/)
  await expect(page.locator('.indicators')).toContainText('2 года')
})

test('новые состояния документов, компании и запроса имеют доступные подписи', async ({ page }) => {
  test.setTimeout(90_000)
  async function check(state: string) {
    const report = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze()
    expect(
      report.violations
        .filter((item) => item.id !== 'color-contrast')
        .map((item) => ({ id: item.id, targets: item.nodes.map((node) => node.target) })),
      state,
    ).toEqual([])
  }
  await page.goto('/#/application/business-one')
  await upload(page, 'Паспорт')
  await check('uploaded-document')
  await page.goto('/#/specialist/XX-001234')
  for (const tab of ['Финансы', 'Проверки']) {
    await page.getByRole('button', { name: tab, exact: true }).click()
    await check(tab)
  }
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: 'Запросить данные', exact: true }).click()
  await page.getByRole('checkbox', { name: 'Другое', exact: true }).check()
  await check('request-other-document')
})
