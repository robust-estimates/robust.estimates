import { CompanyProfile } from './CompanyProfile'
import {
  getSession,
  setSession as saveSession,
  loadProducts,
  loadApplications,
  request,
  saveCompany,
  purpose,
  backendProductId,
  uploadDocuments,
  overview,
} from './api'
import type { ApiApplication, ApiCompany, Recommendation } from './api'
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
import { defaultSelection, isQuote, isSelection } from './types'
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
export default function App() {
  const [route, navigate] = useRoute()
  const [authenticated, setAuthenticated] = useState(() => !!getSession())
  const [accountOpen, setAccountOpen] = useState(false)
  const [company, setCompany] = useState('Моя компания')
  const [displayName, setDisplayName] = useState('')
  const [profile, setProfile] = useState<ApiCompany | null>(null)
  const [ready, setReady] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [recommendations, setRecommendations] = useState<Recommendation[]>([])
  const pendingDraft = useRef<{ route: string; id: string } | null>(null)
  const working = useRef(false)
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
  const [applications, setApplications] = useState<LoanApplication[]>([])
  const [pendingRoute, setPendingRoute] = useSessionValue(
    'alfa-credit-pending-route-v2',
    '/catalog',
    (value): value is string =>
      typeof value === 'string' && value.startsWith('/') && value.length < 300,
  )
  const [toast, setToast] = useState('')
  const accountRef = useRef<HTMLDivElement>(null)
  const loginPage = route === '/login' || !authenticated
  const staff = authenticated && getSession()?.role === 'EMPLOYEE'
  const specialist = staff && route.startsWith('/specialist')
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
    if (!authenticated && route !== '/login') {
      setPendingRoute(route)
      navigate('/login', true)
    }
  }, [authenticated, route, navigate, setPendingRoute])

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
    if (!value) {
      saveSession(null)
      setApplications([])
      setReady(false)
      pendingDraft.current = null
      setProfile(null)
      try {
        for (const key of Object.keys(sessionStorage))
          if (key.startsWith('alfa-credit-')) sessionStorage.removeItem(key)
      } catch {}
    }
    setAuthenticated(value)
  }

  useEffect(() => {
    const expired = () => {
      setSession(false)
      navigate('/login')
      setError('Сессия истекла. Войдите снова.')
    }
    window.addEventListener('credit-session-expired', expired)
    return () => window.removeEventListener('credit-session-expired', expired)
  }, [navigate])

  useEffect(() => {
    if (!authenticated) return
    let stopped = false
    setReady(false)
    setError('')
    async function initialize() {
      try {
        await loadProducts()
        const me = await request<{ displayName: string }>('/me')
        const result = await loadApplications(getSession()?.role === 'EMPLOYEE')
        if (stopped) return
        setDisplayName(me.displayName)
        setApplications(result)
        if (getSession()?.role === 'CLIENT') {
          const list = await request<ApiCompany[]>('/companies')
          if (stopped) return
          if (list[0]) {
            setCompany(list[0].name)
            setProfile(list[0])
          }
        }
        setReady(true)
      } catch (error) {
        if (!stopped)
          setError(error instanceof Error ? error.message : 'Не удалось загрузить данные')
      }
    }
    void initialize()
    return () => {
      stopped = true
    }
  }, [authenticated])

  useEffect(() => {
    if (!authenticated || !ready) return
    let stopped = false
    let active = false
    const timer = window.setInterval(async () => {
      if (working.current || active || document.hidden) return
      active = true
      try {
        const result = await loadApplications(getSession()?.role === 'EMPLOYEE')
        if (!stopped && !working.current) setApplications(result)
      } catch (error) {
        if (!stopped)
          setError(error instanceof Error ? error.message : 'Не удалось обновить заявки')
      } finally {
        active = false
      }
    }, 4000)
    return () => {
      stopped = true
      window.clearInterval(timer)
    }
  }, [authenticated, ready])

  useEffect(() => {
    if (!authenticated || !ready || route !== '/results' || !profile) return
    let stopped = false
    setRecommendations([])
    async function recommend() {
      try {
        const response = await request<Recommendation[]>('/recommendations', {
          method: 'POST',
          body: JSON.stringify({
            companyId: profile!.id,
            amountMinor: Math.round(selection.amount * 100),
            termMonths: selection.months,
            purpose: purpose(selection),
            collateral:
              selection.collateral === 'Да'
                ? 'AVAILABLE'
                : selection.collateral === 'Рассмотрю позже'
                  ? 'LATER'
                  : 'NONE',
          }),
        })
        if (!stopped) setRecommendations(response)
      } catch (error) {
        if (!stopped) setError(error instanceof Error ? error.message : 'Ошибка подбора')
      }
    }
    void recommend()
    return () => {
      stopped = true
    }
  }, [authenticated, ready, route, selection, profile])

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
  async function updateApplication(value: LoanApplication) {
    if (working.current) throw new Error('Дождитесь завершения предыдущего действия')
    working.current = true
    setBusy(true)
    setError('')
    try {
      let current = await request<ApiApplication>(`/applications/${value.id}`)
      if (getSession()?.role === 'EMPLOYEE') {
        const outcomes: Record<string, string> = {
          Одобрено: 'APPROVED',
          Отклонено: 'REJECTED',
          'Другие условия': 'ALTERNATIVE',
          'Запрошены данные': 'NEEDS_INFO',
        }
        const outcome = outcomes[value.status]
        if (!outcome) throw new Error('Недоступное действие')
        const terms = value.alternative
        const body = {
          outcome,
          reason:
            terms?.comment || value.documentRequest?.other || 'Решение кредитного специалиста',
          ...(outcome === 'NEEDS_INFO'
            ? {
                items: value.documentRequest?.items.map((id) => ({
                  title:
                    id === 'other'
                      ? value.documentRequest?.other || 'Дополнительный документ'
                      : id === 'annual-report'
                        ? 'Финансовая отчётность'
                        : id === 'passport'
                          ? 'Паспорт'
                          : 'Сведения по действующим кредитам',
                  type: id === 'annual-report' ? 'FINANCIAL_REPORT' : 'OTHER',
                  required: true,
                })),
              }
            : {}),
          ...(['APPROVED', 'ALTERNATIVE'].includes(outcome)
            ? {
                offers: [
                  {
                    productId: backendProductId(
                      outcome === 'ALTERNATIVE' && terms ? terms.productId : value.productId,
                    ),
                    amountMinor: Math.round(
                      (outcome === 'ALTERNATIVE' && terms ? terms.amount : value.amount) * 100,
                    ),
                    collateral:
                      outcome === 'ALTERNATIVE' && terms
                        ? terms.collateral
                        : products.find((item) => item.id === value.productId)?.secured
                          ? 'REQUIRED'
                          : 'NONE',
                    termMonths: outcome === 'ALTERNATIVE' && terms ? terms.months : value.months,
                    annualRateBps: Math.round(
                      (outcome === 'ALTERNATIVE' && terms ? terms.rate : value.rate) * 100,
                    ),
                  },
                ],
              }
            : {}),
        }
        await request(
          `/bank/applications/${value.id}/decision`,
          { method: 'POST', body: JSON.stringify(body) },
          current.version,
        )
      } else if (value.status === 'Отклонено') {
        await request(
          `/applications/${value.id}/decline`,
          { method: 'POST', body: JSON.stringify({ reason: 'Клиент отказался от предложения' }) },
          current.version,
        )
      } else {
        const old = applications.find((item) => item.id === value.id)
        if (current.status !== 'ACCEPTED') {
          if (!old?.offer) throw new Error('Нет актуального предложения банка. Обновите заявку.')
          current = await request<ApiApplication>(
            `/applications/${value.id}/accept`,
            { method: 'POST', body: JSON.stringify({ offerId: old.offer.id }) },
            current.version,
          )
        }
        if (value.status === 'Выдано')
          await request(`/applications/${value.id}/issuance`, { method: 'POST' }, current.version)
      }
      setApplications(await loadApplications(getSession()?.role === 'EMPLOYEE'))
      setToast('Изменения сохранены')
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Ошибка сохранения')
      throw error
    } finally {
      working.current = false
      setBusy(false)
    }
  }

  async function submitApplication(value: Omit<LoanApplication, 'id' | 'date' | 'status'>) {
    if (working.current) throw new Error('Дождитесь завершения предыдущего действия')
    working.current = true
    setBusy(true)
    setError('')
    try {
      const finance = value.financials
      if (!finance) throw new Error('Заполните финансовые сведения')
      const company = await saveCompany(value.applicant, selection, finance, existing?.companyId)
      setCompany(company.name)
      setProfile(company)
      const need = {
        companyId: company.id,
        amountMinor: Math.round(value.amount * 100),
        termMonths: value.months,
        purpose: value.purposeCode || existing?.purposeCode || 'OTHER',
        description:
          selection.other || selection.purpose || existing?.description || 'Финансирование бизнеса',
      }
      let app: ApiApplication
      const id =
        existing?.id ||
        (pendingDraft.current?.route === route ? pendingDraft.current.id : undefined)
      if (id) {
        app = await request<ApiApplication>(`/applications/${id}`)
        app = await request<ApiApplication>(
          `/applications/${id}`,
          {
            method: 'PATCH',
            body: JSON.stringify({
              need,
              ...(app.status === 'IN_PROGRESS'
                ? { productId: backendProductId(value.productId) }
                : {}),
            }),
          },
          app.version,
        )
      } else {
        app = await request<ApiApplication>('/applications', {
          method: 'POST',
          body: JSON.stringify({
            need,
            productId: backendProductId(value.productId),
            consents: { personalDataProcessing: true, creditHistoryCheck: true },
          }),
        })
        pendingDraft.current = { route, id: app.id }
      }
      app = await uploadDocuments(app.id, value.documentFiles || {}, existing?.information)
      const validation = await request<{
        valid: boolean
        errors: { message: string }[]
        missingFields: string[]
      }>(`/applications/${app.id}/validation`, { method: 'POST' }, app.version)
      if (!validation.valid)
        throw new Error(
          validation.errors.map((item) => item.message).join('; ') ||
            `Заполните: ${validation.missingFields.join(', ')}`,
        )
      app = await request<ApiApplication>(
        `/applications/${app.id}/submit`,
        { method: 'POST' },
        app.version,
      )
      const saved = await overview(app, false)
      setApplications((current) => [saved, ...current.filter((item) => item.id !== saved.id)])
      pendingDraft.current = null
      setEditingId(null)
      navigate(`/accepted/${app.id}`)
    } finally {
      working.current = false
      setBusy(false)
    }
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
                aria-label={`Меню профиля: ${displayName}`}
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
                  <span>{specialist ? displayName : company}</span>
                  <small>{specialist ? 'Специалист' : displayName}</small>
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
                  {!staff && (
                    <button
                      type="button"
                      onClick={() => {
                        setAccountOpen(false)
                        navigate('/company')
                      }}
                    >
                      Мои данные
                    </button>
                  )}
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
                      if (staff) navigate(specialist ? '/catalog' : '/specialist')
                      else {
                        setSession(false)
                        navigate('/login')
                      }
                    }}
                  >
                    {staff
                      ? specialist
                        ? 'Каталог'
                        : 'Кредитный специалист'
                      : 'Войти как сотрудник'}
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
        {error && (
          <div className="integration-error" role="alert">
            {error}
            <button type="button" onClick={() => setError('')}>
              ×
            </button>
          </div>
        )}
        {busy && (
          <p className="integration-status" role="status">
            Сохраняем данные…
          </p>
        )}
        {loginPage ? (
          <Login
            onLogin={(isSpecialist) => {
              setSession(true)
              navigate(isSpecialist ? '/specialist' : pendingRoute)
            }}
          />
        ) : !ready ? (
          <main className="empty-state">
            <p>Загружаем данные кабинета…</p>
            <button type="button" onClick={() => window.location.reload()}>
              Повторить
            </button>
          </main>
        ) : route === '/catalog' || route === '/selection' ? (
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
        ) : route === '/company' || (route === '/results' && !profile) ? (
          <CompanyProfile
            company={profile ?? undefined}
            selection={selection}
            onSaved={(value) => {
              setProfile(value)
              setCompany(value.name)
              if (route === '/company') navigate('/catalog')
            }}
          />
        ) : route === '/results' ? (
          <SelectionResults
            selection={selection}
            recommendations={recommendations}
            onReset={() => {
              setSelection(defaultSelection)
              navigate('/selection')
            }}
            onSelect={(item) => selectProduct(item, true)}
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
            initialCompany={profile ?? undefined}
            initialPurpose={selection.purpose ? purpose(selection) : undefined}
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
            key={`${current.id}-${current.version}`}
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
        ) : staff && route === '/specialist' ? (
          <SpecialistList applications={applications} onOpen={openReview} />
        ) : staff && /^\/specialist\/[^/]+$/.test(route) && review ? (
          <Review
            key={`${review.id}-${applications.find((item) => item.id === review.id)?.version}`}
            item={review}
            application={applications.find((item) => item.id === review.id)!}
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
