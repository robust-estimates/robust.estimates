import type { Page } from '@playwright/test'

export async function attachBusinessDocuments(page: Page) {
  for (const title of [
    'Анкета предприятия или ИП',
    'Анкета физического лица',
    'Паспорт',
    'Годовая финансовая отчётность',
    'Действующие кредиты',
  ]) {
    await page.getByLabel(`Файл: ${title}`, { exact: true }).setInputFiles({
      name: `${title}.pdf`,
      mimeType: 'application/pdf',
      buffer: Buffer.from('%PDF-1.4 demo'),
    })
  }
}
