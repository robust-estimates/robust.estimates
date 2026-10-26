import { defineConfig } from '@playwright/test'

const external = process.env.FRONTEND_URL
const baseURL = external || 'http://127.0.0.1:5173'

export default defineConfig({
  testDir: './tests/integration',
  fullyParallel: false,
  workers: 1,
  timeout: 60_000,
  expect: { timeout: 15_000 },
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL,
    headless: true,
    viewport: { width: 1440, height: 1000 },
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
    launchOptions: { executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE },
  },
  projects: [{ name: 'chromium', use: { browserName: 'chromium' } }],
  ...(external
    ? {}
    : {
        webServer: {
          command: 'npm run dev -- --port 5173 --strictPort',
          url: baseURL,
          reuseExistingServer: !process.env.CI,
        },
      }),
})
