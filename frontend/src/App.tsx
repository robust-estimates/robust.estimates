import { useCallback, useEffect, useRef, useState } from 'react'
import { Catalog, SelectionResults } from './Catalog'
import {
  ApplicationAccepted,
  ApplicationDetail,
  ApplicationForm,
  ApplicationsList,
} from './Applications'
import { Login } from './Login'
import { getReviewItems, Review, SpecialistList } from './Specialist'
import type { ReviewItem } from './Specialist'
import { products, quoteForProduct } from './products'
import type { Product } from './products'
import { defaultSelection, isQuote, isSelection, readApplications } from './types'
import type { LoanApplication } from './types'
import { asset } from './ui'
import { useSessionValue } from './session'

function useRoute() {
  const [route, setRoute] = useState(() => window.location.hash.slice(1) || '/catalog')
  useEffect(() => {
    const change = () => {
      setRoute(window.location.hash.slice(1) || '/catalog')
      window.scrollTo(0, 0)
    }
    window.addEventListener('hashchange', change)
    return () => window.removeEventListener('hashchange', change)
  }, [])
  const navigate = useCallback((next: string, replace = false) => {
    if (replace) {
      window.history.replaceState(null, '', `#${next}`)
      setRoute(next)
      window.scrollTo(0, 0)
    } else window.location.hash = next
  }, [])
  return [route, navigate] as const
}
function readSession() {
  try {
    return sessionStorage.getItem('alfa-credit-session') !== 'guest'
  } catch {
    return true
  }
}

export default function App() {
  const [route, navigate] = useRoute()
  const [authenticated, setAuthenticated] = useState(readSession)
  const [accountOpen, setAccountOpen] = useState(false)
  const company = 'ООО «Лапочка»'
  const [selection, setSelection] = useSessionValue(
    'alfa-credit-selection-v2',
    defaultSelection,
    isSelection,
  )
  const [quote, setQuote] = useSessionValue(
    'alfa-credit-quote-v2',
    quoteForProduct(products[0]),
    isQuote,
  )
  const [editingId, setEditingId] = useSessionValue<string | null>(
    'alfa-credit-editing-v2',
    null,
    (value): value is string | null => value === null || typeof value === 'string',
  )
  const [applications, setApplications] = useState(readApplications)
  const [pendingRoute, setPendingRoute] = useSessionValue(
    'alfa-credit-pending-route-v2',
    '/catalog',
    (value): value is string =>
      typeof value === 'string' && value.startsWith('/') && value.length < 300,
  )
  const [toast, setToast] = useState('')
  const accountRef = useRef<HTMLDivElement>(null)
  const requiresLogin = /^\/(?:applications?|accepted|specialist)(?:\/|$)/.test(route)
  const loginPage = route === '/login' || (!authenticated && requiresLogin)
  const specialist = authenticated && route.startsWith('/specialist')
  const listPage = route === '/applications'
  const showAccount = authenticated && !loginPage
  const parts = route.split('/')
  const product = products.find((item) => item.id === parts[2])
  const existing = applications.find((item) => item.id === parts[3])
  const current = applications.find((item) => item.id === parts[2])
  const review = getReviewItems(applications).find((item) => item.id === parts[2])
  const effectiveQuote = product
    ? existing?.productId === product.id
      ? { productId: product.id, amount: existing.amount, months: existing.months }
      : quoteForProduct(
          product,
          quote.productId === product.id && quote.applicationId === parts[3]
            ? quote
            : existing
              ? { productId: product.id, amount: existing.amount, months: existing.months }
              : undefined,
        )
    : quote
  const previousRoute = useRef(route)
  useEffect(() => {
    if (
      selection.registration === 'АО' ||
      ['Менее 6 месяцев', 'От 6 месяцев до 1 года'].includes(selection.age)
    ) {
      setSelection((current) => ({
        ...current,
        registration: current.registration === 'АО' ? '' : current.registration,
        age: ['Менее 6 месяцев', 'От 6 месяцев до 1 года'].includes(current.age) ? '' : current.age,
      }))
    }
  }, [selection.registration, selection.age, setSelection])

  useEffect(() => {
    if (!authenticated && requiresLogin) {
      setPendingRoute(route)
      navigate('/login', true)
    }
  }, [authenticated, requiresLogin, route, navigate, setPendingRoute])

  useEffect(() => {
    if (!accountOpen) return
    const outside = (event: PointerEvent) => {
      if (!accountRef.current?.contains(event.target as Node)) setAccountOpen(false)
    }
    const escape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setAccountOpen(false)
        accountRef.current?.querySelector('button')?.focus()
      }
    }
    document.addEventListener('pointerdown', outside)
    document.addEventListener('keydown', escape)
    return () => {
      document.removeEventListener('pointerdown', outside)
      document.removeEventListener('keydown', escape)
    }
  }, [accountOpen])
  useEffect(() => {
    if (!toast) return
    const id = window.setTimeout(() => setToast(''), 5000)
    return () => window.clearTimeout(id)
  }, [toast])
  useEffect(() => {
    document.title = `${loginPage ? 'Вход' : specialist ? 'Рассмотрение заявок' : route.startsWith('/application') ? 'Заявки' : 'Кредиты для бизнеса'} — Альфа-бизнес`
    setAccountOpen(false)
    if (previousRoute.current !== route)
      document.getElementById('main-content')?.focus({ preventScroll: true })
    previousRoute.current = route
  }, [route, loginPage, specialist])

  function setSession(value: boolean) {
    setAuthenticated(value)
    try {
      sessionStorage.setItem('alfa-credit-session', value ? 'client' : 'guest')
    } catch {
      /* Keep the session in memory when storage is unavailable. */
    }
  }
  function guardedNavigate(next: string) {
    if (!authenticated) {
      setPendingRoute(next)
      navigate('/login')
    } else navigate(next)
  }
  function selectProduct(item: Product, matched = false) {
    const edit =
      matched && editingId && applications.some((application) => application.id === editingId)
        ? editingId
        : undefined
    if (!matched) setEditingId(null)
    setQuote({
      productId: item.id,
      applicationId: edit,
      amount: Math.min(matched ? selection.amount : 15_000_000, item.amount * 1_000_000),
      months: Math.min(matched ? selection.months : 48, item.months),
    })
    guardedNavigate(`/application/${item.id}${edit ? `/${edit}` : ''}`)
  }
  function persist(next: LoanApplication[]) {
    setApplications(next)
    try {
      localStorage.setItem('alfa-credit-applications-v2', JSON.stringify(next))
    } catch {
      setToast('Данные сохранены на время сеанса. Хранилище браузера недоступно.')
    }
  }
  function updateApplication(value: LoanApplication) {
    persist(
      applications.some((item) => item.id === value.id)
        ? applications.map((item) => (item.id === value.id ? value : item))
        : [value, ...applications],
    )
  }
  function submitApplication(value: Omit<LoanApplication, 'id' | 'date' | 'status'>) {
    const numbers = applications.map((item) => Number(item.id.split('-').at(-1)) || 0)
    const id = existing?.id ?? `XX-${String(Math.max(123, ...numbers) + 1).padStart(6, '0')}`
    const next: LoanApplication = {
      ...value,
      id,
      date: existing?.date ?? new Date().toISOString(),
      status: existing?.status === 'Черновик' ? 'Черновик' : 'На рассмотрении',
    }
    updateApplication(next)
    setEditingId(null)
    navigate(existing ? `/applications/${id}` : `/accepted/${id}`)
    if (existing) setToast('Изменения сохранены')
  }
  function openReview(item: ReviewItem) {
    navigate(`/specialist/${item.id}`)
  }

  return (
    <>
      <a
        className="skip-link"
        href="#main-content"
        tabIndex={0}
        onClick={(event) => {
          event.preventDefault()
          const content = document.getElementById('main-content')
          content?.focus({ preventScroll: true })
          content?.scrollIntoView()
        }}
      >
        Перейти к содержимому
      </a>
      <header className="site-header">
        <div className="header-inner">
          <a
            className="brand"
            href={specialist ? '#/specialist' : '#/catalog'}
            aria-label="Альфа-бизнес — главная"
          >
            <img src={asset('alfa-logo.png')} width="33" height="38" alt="" />
            <span>Альфа-бизнес</span>
          </a>
          <nav
            className={`main-nav ${specialist ? 'specialist-nav' : ''}`}
            aria-label="Основная навигация"
          >
            {specialist ? (
              <>
                <a className="nav-link active" href="#/specialist" aria-current="page">
                  Заявки
                </a>
              </>
            ) : (
              <>
                <a
                  className={`nav-link ${!listPage ? 'active' : ''}`}
                  href="#/catalog"
                  aria-current={!listPage ? 'page' : undefined}
                >
                  Кредиты
                </a>
                {showAccount && (
                  <button
                    className={`nav-link ${listPage ? 'active' : ''}`}
                    type="button"
                    onClick={() => guardedNavigate('/applications')}
                    aria-current={listPage ? 'page' : undefined}
                  >
                    Мои заявки
                  </button>
                )}
              </>
            )}
          </nav>
          {showAccount ? (
            <div className="account-anchor" ref={accountRef}>
              <button
                className="account-button"
                type="button"
                aria-label={`Меню профиля: ${specialist ? 'Коробов П. О. Специалист' : `${company} Иванов К. К.`}`}
                aria-expanded={accountOpen}
                onClick={() => setAccountOpen(!accountOpen)}
              >
                <span className="account-avatar">
                  <img
                    src={asset(specialist ? 'specialist-avatar.png' : 'company-avatar.png')}
                    width={specialist ? 38 : 23}
                    height={specialist ? 34 : 22}
                    alt=""
                  />
                </span>
                <span className="account-copy">
                  <span>{specialist ? 'Коробов П. О.' : company}</span>
                  <small>{specialist ? 'Специалист' : 'Иванов К. К.'}</small>
                </span>
                <img
                  className="account-chevron"
                  src={asset('chevron.png')}
                  width="16"
                  height="16"
                  alt=""
                />
              </button>
              {accountOpen && (
                <div className="account-menu">
                  <p>{specialist ? 'Рабочее место специалиста' : company}</p>
                  {!specialist && (
                    <button
                      type="button"
                      onClick={() => {
                        setAccountOpen(false)
                        guardedNavigate('/applications')
                      }}
                    >
                      Мои заявки
                    </button>
                  )}
                  <button
                    type="button"
                    onClick={() => {
                      setAccountOpen(false)
                      navigate(specialist ? '/catalog' : '/specialist')
                    }}
                  >
                    {specialist ? 'Кабинет клиента' : 'Кредитный специалист'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setSession(false)
                      setAccountOpen(false)
                      navigate('/catalog')
                    }}
                  >
                    Выйти
                  </button>
                </div>
              )}
            </div>
          ) : (
            <button
              className="login-button"
              type="button"
              onClick={() => {
                setPendingRoute('/catalog')
                navigate('/login')
              }}
            >
              Войти <img src={asset('login.png')} width="33" height="33" alt="" />
            </button>
          )}
        </div>
      </header>
      <div id="main-content" tabIndex={-1}>
        {route === '/catalog' || route === '/selection' ? (
          <Catalog
            expanded={route === '/selection'}
            selection={selection}
            onSelection={setSelection}
            onExpand={() => {
              setEditingId(null)
              navigate('/selection')
            }}
            onCollapse={() => navigate('/catalog')}
            onComplete={() => navigate('/results')}
            onSelect={(item) => selectProduct(item)}
          />
        ) : route === '/results' ? (
          <SelectionResults
            selection={selection}
            onReset={() => {
              setSelection(defaultSelection)
              navigate('/selection')
            }}
            onSelect={(item) => selectProduct(item, true)}
          />
        ) : loginPage ? (
          <Login
            onLogin={(isSpecialist) => {
              setSession(true)
              navigate(isSpecialist ? '/specialist' : pendingRoute)
            }}
          />
        ) : /^\/application\/[^/]+(?:\/[^/]+)?$/.test(route) &&
          product &&
          (!parts[3] || existing) ? (
          <ApplicationForm
            key={route}
            product={product}
            amount={effectiveQuote.amount}
            months={effectiveQuote.months}
            existing={existing}
            registration={selection.registration}
            hasLoans={selection.loans === 'Да'}
            onCancel={() => navigate(existing ? `/applications/${existing.id}` : '/catalog')}
            onChangeParameters={(amount, months) => {
              setEditingId(existing?.id ?? null)
              setSelection({
                ...selection,
                amount: Math.min(50_000_000, amount),
                months: Math.max(6, months),
              })
              navigate('/selection')
            }}
            onSubmit={submitApplication}
          />
        ) : /^\/accepted\/[^/]+$/.test(route) && current ? (
          <ApplicationAccepted
            application={current}
            onList={() => navigate('/applications')}
            onCatalog={() => navigate('/catalog')}
          />
        ) : route === '/applications' ? (
          <ApplicationsList
            applications={applications}
            onOpen={(item) => {
              if (item.status === 'Черновик') {
                setQuote({
                  productId: item.productId,
                  amount: item.amount,
                  months: item.months,
                  applicationId: item.id,
                })
                navigate(`/application/${item.productId}/${item.id}`)
              } else navigate(`/applications/${item.id}`)
            }}
            onNew={() => navigate('/catalog')}
          />
        ) : /^\/applications\/[^/]+$/.test(route) && current ? (
          <ApplicationDetail
            key={current.id}
            application={current}
            onBack={() => navigate('/applications')}
            onNew={() => navigate('/catalog')}
            onEdit={() => {
              setQuote({
                productId: current.productId,
                amount: current.amount,
                months: current.months,
                applicationId: current.id,
              })
              navigate(`/application/${current.productId}/${current.id}`)
            }}
            onUpdate={updateApplication}
          />
        ) : route === '/specialist' ? (
          <SpecialistList applications={applications} onOpen={openReview} />
        ) : /^\/specialist\/[^/]+$/.test(route) && review ? (
          <Review
            key={review.id}
            item={review}
            application={applications.find((item) => item.id === review.id)}
            onBack={() => navigate('/specialist')}
            onSave={updateApplication}
          />
        ) : (
          <main className="empty-state">
            <h1>Страница не найдена</h1>
            <a className="button-primary" href="#/catalog">
              В каталог кредитов
            </a>
          </main>
        )}
      </div>
      {toast && (
        <div className="toast" role="status">
          {toast}
          <button type="button" aria-label="Закрыть сообщение" onClick={() => setToast('')}>
            ×
          </button>
        </div>
      )}
    </>
  )
}
