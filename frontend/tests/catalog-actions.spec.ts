import { expect, test } from '@playwright/test'

test('Подробнее открывает условия каждого продукта; оформление запускается отдельной кнопкой', async ({
  page,
}) => {
  await page.goto('/#/catalog')
  const names = await page.locator('.product-card h3').allTextContents()
  const originalStorage = await page.evaluate(() =>
    localStorage.getItem('alfa-credit-applications-v2'),
  )
  for (const name of names) {
    const trigger = page.getByRole('button', { name: `Подробнее: ${name}`, exact: true })
    await trigger.click()
    const dialog = page.getByRole('dialog', { name, exact: true })
    await expect(dialog).toBeVisible()
    await expect(dialog).toContainText('Обеспечение')
    await expect(page).toHaveURL(/#\/catalog$/)
    await expect(page.locator('.application-form')).toHaveCount(0)
    expect(await page.evaluate(() => localStorage.getItem('alfa-credit-applications-v2'))).toBe(
      originalStorage,
    )
    await page.keyboard.press('Escape')
    await expect(dialog).toHaveCount(0)
    await expect(trigger).toBeFocused()
  }
  await page
    .getByRole('button', { name: 'Подробнее: Бизнес-кредит с залогом', exact: true })
    .click()
  const dialog = page.getByRole('dialog')
  await expect(dialog).toContainText('500 000 — 150 000 000 ₽')
  await expect(dialog).toContainText('С залогом')
  await dialog.getByRole('button', { name: 'Оформить заявку', exact: true }).click()
  await expect(page).toHaveURL(/#\/application\/secured$/)
  await expect(page.locator('.selected-product')).toContainText('Бизнес-кредит с залогом')
})

test('гость читает условия без входа, затем входит для оформления выбранного продукта', async ({
  page,
}) => {
  await page.addInitScript(() => sessionStorage.setItem('alfa-credit-session', 'guest'))
  await page.goto('/#/catalog')
  await page.getByRole('button', { name: 'Без залога', exact: true }).click()
  await page.getByRole('searchbox').fill('Business One')
  const trigger = page.getByRole('button', { name: 'Подробнее: Business One Credit', exact: true })
  await trigger.click()
  let dialog = page.getByRole('dialog')
  await expect(dialog).toContainText('Ставка вне льготного периода в таблице не указана')
  await expect(dialog).toContainText('до 120 дней')
  await dialog.getByRole('button', { name: 'Вернуться в каталог' }).click()
  await expect(page.getByRole('searchbox')).toHaveValue('Business One')
  await expect(page.locator('.product-card')).toHaveCount(1)
  await trigger.click()
  dialog = page.getByRole('dialog')
  await dialog.getByRole('button', { name: 'Оформить заявку', exact: true }).click()
  await expect(page).toHaveURL(/#\/login$/)
  await page.locator('.login-card').getByRole('button', { name: 'Войти', exact: true }).click()
  await expect(page).toHaveURL(/#\/application\/business-one$/)
})

for (const decision of [
  { action: 'Одобрить', status: 'Одобрено', remains: false },
  { action: 'Другие условия', status: 'Другие условия', remains: false },
  { action: 'Запросить данные', status: 'Запрошены данные', remains: true },
]) {
  test(`очередь специалиста после решения «${decision.status}» сохраняется правильно`, async ({
    page,
  }) => {
    await page.goto('/#/specialist/XX-001234')
    await page.getByRole('button', { name: 'Решение', exact: true }).click()
    await page.getByRole('button', { name: new RegExp(decision.action) }).click()
    if (decision.action === 'Запросить данные')
      await page.getByRole('checkbox', { name: 'Паспорт', exact: true }).check()
    await page
      .getByRole('button', {
        name:
          decision.action === 'Другие условия'
            ? 'Сохранить предложение'
            : decision.action === 'Запросить данные'
              ? 'Отправить запрос'
              : 'Подтвердить',
        exact: true,
      })
      .click()
    await page.getByRole('button', { name: /К списку заявок/ }).click()
    const row = page.locator('tbody tr').filter({ hasText: 'XX-001234' })
    await expect(row).toHaveCount(decision.remains ? 1 : 0)
    await expect(page.locator('tbody tr')).toHaveCount(decision.remains ? 10 : 9)
    await page.reload()
    await expect(row).toHaveCount(decision.remains ? 1 : 0)
    await page.goto('/#/applications/XX-001234')
    await expect(page.locator('.application-detail-heading .status-badge')).toHaveText(
      decision.status,
    )
  })
}

test('условия продукта помещаются на телефоне и доступны для клавиатуры', async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 650 })
  await page.goto('/#/catalog')
  await page
    .getByRole('button', { name: 'Подробнее: Бизнес-кредит без залога', exact: true })
    .click()
  const dialog = page.getByRole('dialog')
  await expect(dialog).toBeVisible()
  expect(await dialog.evaluate((element) => element.scrollWidth <= element.clientWidth)).toBe(true)
  for (let index = 0; index < 5; index++) {
    await page.keyboard.press('Tab')
    expect(await dialog.evaluate((element) => element.contains(document.activeElement))).toBe(true)
  }
  await expect(dialog.getByRole('button', { name: 'Оформить заявку', exact: true })).toBeVisible()
})
