import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { attachBusinessDocuments } from './helpers'

async function choose(page: Page, label: string, option: string) {
  await page.getByRole('combobox', { name: new RegExp(label) }).click()
  await page.getByRole('option', { name: option, exact: true }).click()
}

async function completeSelection(page: Page) {
  await choose(page, 'Для чего вашему бизнесу', 'Закупка товара')
  await choose(page, 'тип регистрации', 'ООО')
  await choose(page, 'Сколько лет', 'Более 3 лет')
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
  await expect(page).toHaveURL(/#\/results$/)
}

test('каталог повторяет правый макет: шапка, белый баннер, условия и геометрия карточек', async ({
  page,
}) => {
  const errors: string[] = []
  page.on('pageerror', (error) => errors.push(error.message))
  await page.goto('/')
  await expect(page.getByRole('link', { name: 'Альфа-бизнес — главная' })).toBeVisible()
  await expect(page.locator('.product-card')).toHaveCount(8)
  await expect(page.locator('.category-chip')).toHaveCount(6)
  await expect(page.getByRole('combobox', { name: 'Сортировка продуктов' })).toHaveCount(0)
  await expect(page.locator('.hero img')).toHaveCount(1)
  await expect(page.locator('.hero')).toHaveCSS('background-color', 'rgb(255, 255, 255)')
  await expect(page.locator('[data-product-id="business"] .product-stats')).toContainText('14%')
  await expect(page.locator('[data-product-id="secured"] .product-stats')).toContainText(
    'до 150 млн ₽',
  )
  await expect(page.locator('[data-product-id="business-one"] .product-stats')).toContainText(
    'до 120 дней',
  )
  const hero = await page.locator('.hero').boundingBox()
  expect(hero).toMatchObject({ x: 69, y: 80, width: 1538, height: 234 })
  const first = await page.locator('.product-card').first().boundingBox()
  expect(first!.y).toBe(470)
  expect(first!.height).toBe(223)
  expect((await page.locator('.product-card').nth(4).boundingBox())!.y).toBe(730)
  await expect
    .poll(() =>
      page
        .locator('img')
        .evaluateAll((images) => images.every((image) => image.complete && image.naturalWidth > 0)),
    )
    .toBe(true)
  expect(errors).toEqual([])
  await page.screenshot({ path: 'tests/desktop.png', fullPage: true })
})

test('поиск и категории работают без выдуманных тегов и сортировки', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Популярные', exact: true }).click()
  await expect(page.locator('.product-card')).toHaveCount(2)
  await page.getByRole('searchbox').fill('неизвестный кредит')
  await expect(page.getByRole('heading', { name: 'Ничего не найдено' })).toBeVisible()
  await page.getByRole('button', { name: 'Сбросить фильтры' }).click()
  await expect(page.locator('.product-card')).toHaveCount(8)
  await page.getByRole('searchbox').fill('рефинансирование')
  await expect(page.locator('.product-card')).toHaveCount(1)
  await page.getByRole('button', { name: 'Очистить поиск' }).click()
  await page.getByRole('button', { name: 'С залогом', exact: true }).click()
  await expect(page.locator('.product-card')).toHaveCount(2)
})

test('подбор раскрывается над каталогом; ползунки синхронизированы с полями', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Подобрать кредит' }).click()
  await expect(page).toHaveURL(/#\/selection$/)
  await expect(page.getByRole('dialog')).toHaveCount(0)
  await expect(page.locator('.selection-panel')).toBeVisible()
  await expect(page.locator('.product-card')).toHaveCount(8)
  const amount = page.getByRole('slider', { name: 'Изменить сумму кредита' })
  const term = page.getByRole('slider', { name: 'Изменить срок кредита' })
  await expect(amount).toHaveValue('12500000')
  await amount.focus()
  await page.keyboard.press('Home')
  await page.keyboard.press('ArrowRight')
  await expect(amount).toHaveValue('100000')
  await expect(page.getByRole('textbox', { name: 'Сумма кредита', exact: true })).toHaveValue(
    '100 000 ₽',
  )
  await term.focus()
  await page.keyboard.press('End')
  await expect(page.getByRole('textbox', { name: 'Срок кредита', exact: true })).toHaveValue(
    '10 лет',
  )
  await page.keyboard.press('Home')
  await expect(term).toHaveValue('6')
  await expect(page.getByRole('textbox', { name: 'Срок кредита', exact: true })).toHaveValue(
    '6 месяцев',
  )
  await page.getByRole('button', { name: 'Скрыть' }).click()
  await page.getByRole('button', { name: 'Подобрать кредит' }).click()
  await expect(amount).toHaveValue('100000')
  await expect(term).toHaveValue('6')
})

test('список целей, поле «Другое», подсказка суммы и обязательные поля соответствуют макету', async ({
  page,
}) => {
  await page.goto('/#/selection')
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
  await expect(page.locator('.select-button[aria-invalid="true"]')).toHaveCount(2)
  await choose(page, 'Для чего вашему бизнесу', 'Другое')
  const description = page.getByRole('textbox', { name: 'Опишите цель кредита' })
  await description.fill('Сезонные расходы')
  await expect(page.locator('.other-purpose')).toContainText('16/300')
  await expect(description).toHaveAttribute('maxlength', '300')
  await page.getByRole('button', { name: 'Помогите определить сумму', exact: false }).click()
  await expect(page.locator('.amount-suggestion')).toContainText('3 000 000 — 8 000 000 ₽')
  await page
    .getByRole('group', { name: 'Есть ли залог?' })
    .getByRole('button', { name: 'Да', exact: true })
    .click()
  await expect(
    page
      .getByRole('group', { name: 'Есть ли залог?' })
      .getByRole('button', { name: 'Да', exact: true }),
  ).toHaveAttribute('aria-pressed', 'true')
  await choose(page, 'тип регистрации', 'ООО')
  await choose(page, 'Сколько лет', 'Более 3 лет')
  await page.getByRole('button', { name: 'Завершить подбор' }).click()
  await expect(page).toHaveURL(/#\/results$/)
})

test('результат содержит три группы, недоступные продукты заблокированы, выбор ведёт к оформлению', async ({
  page,
}) => {
  await page.goto('/#/selection')
  await completeSelection(page)
  await expect(
    page.getByRole('region', { name: 'Наиболее подходящие' }).locator('.product-card'),
  ).toHaveCount(1)
  await expect(
    page.getByRole('region', { name: 'Доступные альтернативы' }).locator('.product-card'),
  ).toHaveCount(0)
  const unavailable = page.getByRole('region', { name: 'Недоступные продукты' })
  await expect(unavailable.locator('.product-card')).toHaveCount(7)
  await expect(unavailable.getByRole('button').first()).toBeDisabled()
  await page
    .getByRole('button', { name: 'Выбрать продукт: Бизнес-кредит без залога', exact: true })
    .click()
  await expect(page).toHaveURL(/#\/application\/business$/)
  await expect(page.getByRole('heading', { name: 'Заявка на бизнес-кредит' })).toBeVisible()
  await expect(page.locator('.indicators')).toContainText('12 500 000 ₽')
})

test('оформление требует документы и согласие; новая заявка сохраняется после перезагрузки', async ({
  page,
}) => {
  await page.goto('/')
  await page
    .getByRole('button', { name: 'Подробнее: Бизнес-кредит без залога', exact: true })
    .click()
  await page
    .getByRole('dialog')
    .getByRole('button', { name: 'Оформить заявку', exact: true })
    .click()
  await page.getByRole('checkbox', { name: 'Согласие на обработку персональных данных*' }).check()
  await page.getByRole('button', { name: 'Оформить заявку' }).click()
  await expect(
    page.getByText('Прикрепите все обязательные документы.', { exact: true }),
  ).toBeVisible()
  await attachBusinessDocuments(page)
  await page.getByRole('button', { name: 'Оформить заявку' }).click()
  await expect(page.getByRole('heading', { name: 'Заявка отправлена!' })).toBeVisible()
  await expect(page.locator('.accepted-number')).toContainText('XX-000124')
  await page.getByRole('button', { name: 'Перейти к моим заявкам' }).click()
  await expect(page.locator('.application-row')).toHaveCount(5)
  await page.reload()
  await expect(page.locator('.application-row').first()).toContainText('XX-000124')
  await page.locator('.application-row').first().click()
  await expect(page.locator('.document-list')).toContainText('Годовая финансовая отчётность.pdf')
})

test('альтернативные условия проходят путь от специалиста до принятия клиентом', async ({
  page,
}) => {
  test.setTimeout(60_000)
  await page.goto('/#/specialist')
  await page.getByRole('button', { name: 'Рассмотреть заявку XX-001234' }).click()
  await page.getByRole('button', { name: 'Решение', exact: true }).click()
  await page.getByRole('button', { name: 'Другие условия', exact: true }).click()
  await choose(page, 'Вид кредита', 'Кредитная линия')
  await page
    .getByRole('textbox', { name: 'Комментарий к решению' })
    .fill('Предлагаем уменьшить сумму кредита.')
  await page.getByRole('button', { name: 'Сохранить предложение' }).click()
  await expect(page.locator('.review-saved')).toContainText('Другие условия')
  await page.getByRole('button', { name: 'Меню профиля' }).click()
  await page.getByRole('button', { name: 'Кабинет клиента', exact: true }).click()
  await page.getByRole('button', { name: 'Мои заявки', exact: true }).click()
  await page.locator('.application-row').filter({ hasText: 'XX-001234' }).click()
  await page.getByRole('button', { name: 'Ознакомиться' }).click()
  const dialog = page.getByRole('dialog')
  await expect(dialog).toContainText('10 000 000 ₽')
  await dialog.getByRole('button', { name: 'Принять предложение' }).click()
  await expect(page.getByRole('dialog')).toHaveCount(0)
  await expect(page.locator('.application-detail-heading .status-badge')).toHaveText('Одобрено')
  await expect(page.locator('.indicators')).toContainText('Кредитная линия')
  await page.reload()
  await expect(page.locator('.indicators')).toContainText('10 000 000 ₽')
})

test('вход показывает ошибку и подтверждает телефон четырьмя цифрами', async ({ page }) => {
  await page.goto('/#/login')
  const form = page.locator('.login-card')
  await form.getByRole('textbox', { name: 'Логин', exact: true }).fill('wrong')
  await form.getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page.getByRole('alert')).toHaveText('Неправильный логин или пароль')
  await form.getByRole('button', { name: 'Телефон', exact: true }).click()
  await form.getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Подтвердите телефон' })).toBeVisible()
  for (const [index, digit] of Array.from('4267').entries())
    await page.getByRole('textbox', { name: `Цифра ${index + 1} кода` }).fill(digit)
  await form.getByRole('button', { name: 'Продолжить', exact: true }).click()
  await expect(page).toHaveURL(/#\/catalog$/)
  await expect(page.getByRole('button', { name: 'Меню профиля' })).toBeVisible()
})

for (const width of [320, 375, 768, 1024, 1366]) {
  test(`страницы не переполняются при ширине ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 900 })
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
      await page.evaluate(() => document.fonts.ready)
      expect(
        await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth),
        route,
      ).toBe(true)
      expect(
        await page
          .locator('img')
          .evaluateAll((images) =>
            images.every((image) => image.complete && image.naturalWidth > 0),
          ),
        route,
      ).toBe(true)
    }
    if (width === 375) {
      await page.goto('/#/selection')
      await page.screenshot({ path: 'tests/mobile.png', fullPage: true })
    }
  })
}

test('оригинальные SVG сохраняют размеры на всех основных страницах', async ({ page }) => {
  for (const route of [
    'catalog',
    'selection',
    'application/business',
    'specialist',
    'specialist/XX-001234',
  ]) {
    await page.goto(`/#/${route}`)
    const icons = await page.locator('img.icon').evaluateAll((images) =>
      images.map((image) => {
        const box = image.getBoundingClientRect()
        return {
          name: image.getAttribute('src')!.split('/').at(-1)!,
          width: box.width,
          height: box.height,
        }
      }),
    )
    for (const icon of icons) {
      if (
        [
          'coin.svg',
          'home.svg',
          'refresh.svg',
          'wallet.svg',
          'building1.svg',
          'percent.svg',
          'cart.svg',
          'brief.svg',
          'doc.svg',
          'review-building.svg',
          'review-check.svg',
        ].includes(icon.name)
      ) {
        expect(icon.width, icon.name).toBe(30)
        expect(icon.height, icon.name).toBe(30)
      }
      if (icon.name === 'score-ring.svg') {
        expect(icon.width).toBe(84)
        expect(icon.height).toBe(84)
      }
      if (icon.name === 'search.svg') {
        expect(icon.width).toBe(22)
        expect(icon.height).toBe(22)
      }
    }
  }
})
