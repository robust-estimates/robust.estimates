import { downloadDocument } from './api'
import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import type { ApplicationStatus, LoanApplication } from './types'
import { formatTerm, money, products, rateLabel } from './products'
import { DocumentList, StatusBadge } from './Applications'
import { Icon, MoneyInput, SearchField, SelectField } from './ui'
import { creditDetails, requestOptions } from './creditData'

export type ReviewItem = {
  id: string
  company: string
  amount: number
  score: number
  recommendation: string
  inn?: string
  date?: string
  status?: ApplicationStatus
}
export function getReviewItems(applications: LoanApplication[]): ReviewItem[] {
  return applications.map((application) => ({
    id: application.id,
    company: application.applicant.company,
    amount: application.amount,
    score: application.scoring?.score ?? 0,
    inn: application.applicant.inn,
    date: application.date,
    status: application.status,
    recommendation:
      application.scoring?.recommendation === 'APPROVE'
        ? 'Одобрение'
        : application.scoring?.recommendation === 'REJECT'
          ? 'Отказать'
          : 'Ручная проверка',
  }))
}
export function SpecialistList({
  onOpen,
  applications,
}: {
  onOpen: (item: ReviewItem) => void
  applications: LoanApplication[]
}) {
  const [query, setQuery] = useState('')
  const [filter, setFilter] = useState('Все')
  const [filters, setFilters] = useState(false)
  const filterRef = useRef<HTMLDivElement>(null)
  const filterButton = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    if (!filters) return
    const outside = (event: PointerEvent) => {
      if (!filterRef.current?.contains(event.target as Node)) setFilters(false)
    }
    const escape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setFilters(false)
        filterButton.current?.focus()
      }
    }
    document.addEventListener('pointerdown', outside)
    document.addEventListener('keydown', escape)
    return () => {
      document.removeEventListener('pointerdown', outside)
      document.removeEventListener('keydown', escape)
    }
  }, [filters])
  const items = getReviewItems(applications).filter(
    (item) =>
      !item.status || item.status === 'На рассмотрении' || item.status === 'Запрошены данные',
  )
  const filtersList = [
    ...new Set([
      'Все',
      'Одобрение',
      'Ручная проверка',
      'Отказать',
      ...items.map((item) => item.recommendation),
    ]),
  ]
  const rows = items.filter(
    (item) =>
      `${item.id} ${item.company} ${item.inn ?? ''}`
        .toLowerCase()
        .includes(query.trim().toLowerCase()) &&
      (filter === 'Все' || item.recommendation === filter),
  )
  return (
    <main className="specialist-list">
      <h1>Заявки на рассмотрение</h1>
      <div className="specialist-controls">
        <SearchField
          label="Поиск заявок специалиста"
          placeholder="Номер заявки, ИНН или компания"
          value={query}
          onChange={setQuery}
        />
        <div className="filter-anchor" ref={filterRef}>
          <button
            className="button-white filter-button"
            ref={filterButton}
            type="button"
            onClick={() => setFilters(!filters)}
            aria-expanded={filters}
          >
            Фильтры <Icon name="review-filter-chevron" />
          </button>
          {filters && (
            <div className="filter-options">
              {filtersList.map((item) => (
                <button
                  type="button"
                  aria-pressed={filter === item}
                  key={item}
                  onClick={() => {
                    setFilter(item)
                    setFilters(false)
                    filterButton.current?.focus()
                  }}
                >
                  {item}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
      <div className="review-table-wrap">
        <table className="review-table">
          <thead>
            <tr>
              <th>№ заявки</th>
              <th>Компания</th>
              <th>Сумма</th>
              <th>Рекомендация</th>
              <th>
                <span className="sr-only">Открыть</span>
              </th>
            </tr>
          </thead>
          <tbody>
            {rows.map((item) => (
              <tr key={item.id}>
                <td>
                  <strong>{item.id}</strong>
                  <span className="review-row-divider" aria-hidden="true">
                    <Icon name="page-next" />
                  </span>
                </td>
                <td>
                  {item.company}
                  <small>ИНН {item.inn ?? ''}</small>
                </td>
                <td>
                  <strong>{money(item.amount)} ₽</strong>
                </td>
                <td>
                  <span
                    className={`recommendation-pill recommendation-${item.status === 'Отклонено' ? 'red' : item.status === 'Одобрено' ? 'green' : item.score >= 700 ? 'green' : item.score >= 500 ? 'yellow' : 'red'}`}
                  >
                    {item.recommendation}
                  </span>
                </td>
                <td>
                  <button
                    className="review-open"
                    type="button"
                    onClick={() => onOpen(item)}
                    aria-label={`Рассмотреть заявку ${item.id}`}
                  >
                    <Icon name="review-arrow" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {!rows.length && <div className="empty-state">Заявки не найдены</div>}
      </div>
      <div className="pagination">
        <button type="button" disabled aria-label="Предыдущая страница">
          ‹
        </button>
        <button type="button" aria-current="page">
          1
        </button>
        <button type="button" disabled aria-label="Следующая страница">
          ›
        </button>
      </div>
    </main>
  )
}

export function Review({
  item,
  application,
  onBack,
  onSave,
}: {
  item: ReviewItem
  application: LoanApplication
  onBack: () => void
  onSave: (application: LoanApplication) => Promise<void>
}) {
  const [tab, setTab] = useState('Обзор')
  const [showScores, setShowScores] = useState(false)
  const [action, setAction] = useState<'approve' | 'reject' | 'terms' | 'request' | null>(null)
  const [saved, setSaved] = useState('')
  const [pending, setPending] = useState(false)
  const [saveError, setSaveError] = useState('')
  const [productId, setProductId] = useState(
    application.alternative?.productId ?? application.productId,
  )
  const [amount, setAmount] = useState(application.alternative?.amount ?? application.amount)
  const [months, setMonths] = useState(
    Math.min(
      application.alternative?.months ?? application.months,
      products.find(
        (product) => product.id === (application.alternative?.productId ?? application.productId),
      )!.months,
    ),
  )
  const [rate, setRate] = useState(application.alternative?.rate ?? application.rate)
  const [collateral, setCollateral] = useState(
    application?.alternative?.collateral ?? 'Без изменения',
  )
  const [comment, setComment] = useState(application?.alternative?.comment ?? '')
  const [requested, setRequested] = useState<string[]>([])
  const [otherDocument, setOtherDocument] = useState('')
  const [requestError, setRequestError] = useState('')
  const value = application
  const tone = item.score >= 700 ? 'green' : item.score >= 500 ? 'yellow' : 'red'
  const profile = value.companyProfile ?? {
    registrationStatus: 'Нет данных',
    director: value.applicant.name,
    founders: [value.applicant.name],
    okved: ['47.91 — Розничная торговля по почте или через интернет'],
    bankruptcy: false,
    convictions: false,
    activeLoans: true,
  }
  const reportName =
    value.documentFiles?.FINANCIAL_REPORT ??
    value.documents.find((name) => /отч[её]т|бухгалтер|финанс/i.test(name))
  const answer = (value: boolean | null) => (value === null ? 'Нет данных' : value ? 'Да' : 'Нет')
  const offeredProduct = products.find((product) => product.id === productId)!
  const termOptions = [...new Set([1, 3, 4, 6, 12, 15, 18, 24, 36, 48, 60, 120, months])]
    .sort((a, b) => a - b)
    .filter(
      (term) => term <= offeredProduct.months && term >= creditDetails[offeredProduct.id].minMonths,
    )
  function changeProduct(name: string) {
    const next = products.find((product) => product.name === name)!
    setProductId(next.id)
    setAmount((current) =>
      Math.max(creditDetails[next.id].minAmount, Math.min(current, next.amount * 1_000_000)),
    )
    setMonths((current) =>
      Math.max(creditDetails[next.id].minMonths, Math.min(current, next.months)),
    )
  }
  async function save(status: LoanApplication['status']) {
    if (pending) return
    setPending(true)
    setSaveError('')
    try {
      await onSave({
        ...value,
        status,
        ...(status === 'Запрошены данные'
          ? {
              documentRequest: {
                items: requested,
                other: requested.includes('other') ? otherDocument.trim() : '',
                requestedAt: new Date().toISOString(),
              },
            }
          : {}),
        ...(status === 'Другие условия'
          ? {
              alternative: {
                productId,
                amount,
                months,
                rate,
                collateral: collateral === 'Без изменения' ? 'Без залога' : collateral,
                comment,
              },
            }
          : {}),
      })
      setSaved(status)
      setAction(null)
    } catch (error) {
      setSaveError(error instanceof Error ? error.message : 'Не удалось сохранить решение')
    } finally {
      setPending(false)
    }
  }
  function terms(event: FormEvent) {
    event.preventDefault()
    save('Другие условия')
  }
  function sendRequest(event: FormEvent) {
    event.preventDefault()
    if (!requested.length) {
      setRequestError('Выберите хотя бы один документ.')
      return
    }
    if (requested.includes('other') && !otherDocument.trim()) {
      setRequestError('Укажите, какой документ нужно предоставить.')
      return
    }
    setRequestError('')
    save('Запрошены данные')
  }
  return (
    <main className="review-page">
      <button className="back-button" type="button" onClick={onBack}>
        ← &nbsp; К списку заявок
      </button>
      {saveError && (
        <p role="alert" className="error-text">
          {saveError}
        </p>
      )}
      {pending && <p role="status">Сохраняем решение…</p>}
      {value.backendStatus !== 'MANUAL_CHECK' && (
        <p role="status">
          Решение доступно только для заявки на ручной проверке. Текущий статус:{' '}
          {value.backendStatus}
        </p>
      )}
      <section className="review-card">
        <header className="review-card-header">
          <p className="review-breadcrumb">Заявки &nbsp; › &nbsp; {item.id}</p>
          <h1>{item.id}</h1>
          <p className="review-date">
            Дата подачи:{' '}
            {new Intl.DateTimeFormat('ru-RU', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
              timeZone: 'Europe/Moscow',
            }).format(new Date(value.date))}
          </p>
          <div className="review-company">
            <span>
              <Icon name="review-building" />
            </span>
            <div>
              <strong>{value.applicant.company}</strong>
              <p>ИНН {value.applicant.inn}</p>
              <p>Дата регистрации: {value.applicant.registered.split('-').reverse().join('.')}</p>
            </div>
          </div>
        </header>
        <nav className="review-tabs" aria-label="Разделы рассмотрения">
          {['Обзор', 'Финансы', 'Проверки', 'Документы', 'Решение'].map((label) => (
            <button
              className={tab === label ? 'active' : ''}
              type="button"
              key={label}
              onClick={() => {
                setTab(label)
                setAction(null)
              }}
            >
              {label}
            </button>
          ))}
        </nav>
        <div className="review-content">
          {saved && (
            <div className="review-saved" role="status">
              Решение сохранено: <StatusBadge status={saved as LoanApplication['status']} />
            </div>
          )}
          {tab === 'Обзор' && (
            <>
              <section className="scoring">
                <h2>
                  Результаты скоринга{' '}
                  <button
                    type="button"
                    aria-label="Подробнее о скоринге"
                    onClick={() => setShowScores(!showScores)}
                  >
                    ⓘ
                  </button>
                </h2>
                <div className="score-circle">
                  <Icon name={tone === 'green' ? 'score-ring' : `score-${tone}`} />
                  <span>
                    <strong>{item.score}</strong>
                    <small>из 1000</small>
                  </span>
                </div>
                {showScores && (
                  <div className="scoring-popover">
                    {(value.scoring?.factors || []).map((factor) => (
                      <p key={factor}>{factor}</p>
                    ))}
                    <p>Модель: {value.scoring?.modelVersion || 'Результат ещё не получен'}</p>
                  </div>
                )}
              </section>
              <div className={`review-recommendation recommendation-${tone}`}>
                {tone === 'green' && <Icon name="review-check" />}
                <div>
                  <h2>
                    {tone === 'green'
                      ? 'Заявка рекомендуется к одобрению'
                      : tone === 'yellow'
                        ? 'Рекомендуется изучить заявку подробнее'
                        : 'Рекомендуется отклонить заявку'}
                  </h2>
                  <p>
                    {tone === 'green'
                      ? 'Система оценила заявку как низкорисковую. Рекомендуется подтвердить решение.'
                      : tone === 'yellow'
                        ? 'Система не оценила заявку как низкорисковую. Рекомендуется проверить данные и принять решение.'
                        : 'Система оценила заявку как рисковую. Рекомендуется проверить данные и принять решение об отмене или альтернативных условиях.'}
                  </p>
                  {tone === 'green' && value.backendStatus === 'MANUAL_CHECK' && (
                    <button
                      className="button-primary"
                      type="button"
                      onClick={() => {
                        setTab('Решение')
                        setAction('approve')
                      }}
                    >
                      Одобрить
                    </button>
                  )}
                </div>
              </div>
              <section className="review-parameters">
                <h2>Параметры заявки</h2>
                <dl>
                  {[
                    ['Запрашиваемая сумма', `${money(value.amount)} ₽`],
                    [
                      'Залог',
                      products.find((item) => item.id === value.productId)?.secured
                        ? 'Требуется'
                        : 'Не требуется',
                    ],
                    ['Срок', formatTerm(value.months)],
                    ['Ставка', `${rateLabel(value.rate)}`],
                    ['Цель кредита', 'Пополнение оборотных средств'],
                    [
                      'Тип лимита',
                      products.find((item) => item.id === value.productId)?.limit
                        ? 'Возобновляемый'
                        : 'Разовый',
                    ],
                    [
                      'Продукт',
                      products.find((product) => product.id === value.productId)?.name ||
                        'Кредитная линия',
                    ],
                    ['Регион', 'Нет данных'],
                  ].map(([label, text]) => (
                    <div key={label}>
                      <dt>{label}</dt>
                      <dd>{text}</dd>
                    </div>
                  ))}
                </dl>
              </section>
            </>
          )}
          {tab === 'Финансы' && (
            <section className="review-company-facts">
              <h2 className="review-section-title">Финансы компании</h2>
              <dl>
                {[
                  ['Выручка в месяц', `${money(value.financials?.revenue || 0)} ₽`],
                  ['Расходы в месяц', `${money(value.financials?.expenses || 0)} ₽`],
                  ['Платежи по кредитам в месяц', `${money(value.financials?.debts || 0)} ₽`],
                  ['Годовая отчётность', reportName ?? 'Не предоставлена'],
                  ['Проходит процедуру банкротства', answer(profile.bankruptcy)],
                  [
                    'Отсутствие судимостей',
                    profile.convictions === null
                      ? 'Нет данных'
                      : profile.convictions
                        ? 'Не подтверждено'
                        : 'Подтверждено',
                  ],
                  [
                    'Действующие кредиты',
                    profile.activeLoans === null
                      ? 'Нет данных'
                      : profile.activeLoans
                        ? 'Есть'
                        : 'Отсутствуют',
                  ],
                ].map(([label, text]) => (
                  <div key={label}>
                    <dt>{label}</dt>
                    <dd>{text}</dd>
                  </div>
                ))}
              </dl>
            </section>
          )}
          {tab === 'Проверки' && (
            <section className="review-company-facts">
              <h2>Проверка компании</h2>
              <dl>
                {[
                  ['ИНН', value.applicant.inn],
                  ['Дата регистрации', value.applicant.registered.split('-').reverse().join('.')],
                  ['Компания', value.applicant.company],
                  ['Статус компании', profile.registrationStatus],
                  ['Руководитель', profile.director || 'Нет данных'],
                  ['Учредители', profile.founders.join('; ') || 'Нет данных'],
                  ['ОКВЭД', profile.okved.join('; ') || 'Нет данных'],
                ].map(([label, text]) => (
                  <div key={label}>
                    <dt>{label}</dt>
                    <dd>{text}</dd>
                  </div>
                ))}
              </dl>
            </section>
          )}
          {tab === 'Документы' && (
            <>
              <h2 className="review-section-title">Документы</h2>
              <DocumentList
                documents={value.documents}
                onDownload={(index) => {
                  const doc = value.apiDocuments?.[index]
                  if (doc)
                    void downloadDocument(value.id, doc).catch((error) =>
                      setSaveError(error.message),
                    )
                }}
              />
            </>
          )}
          {tab === 'Решение' &&
            value.backendStatus === 'MANUAL_CHECK' &&
            (action === 'request' ? (
              <form className="document-request-form" onSubmit={sendRequest}>
                <h2>Запрос дополнительных данных</h2>
                <fieldset>
                  <legend>Какие документы нужны?</legend>
                  <p>Можно выбрать несколько вариантов.</p>
                  {requestOptions.map((option) => (
                    <label key={option.id}>
                      <input
                        type="checkbox"
                        checked={requested.includes(option.id)}
                        onChange={(event) => {
                          setRequested((current) =>
                            event.target.checked
                              ? [...current, option.id]
                              : current.filter((id) => id !== option.id),
                          )
                          setRequestError('')
                        }}
                      />
                      {option.title}
                    </label>
                  ))}
                </fieldset>
                {requested.includes('other') && (
                  <label className="field">
                    Какие ещё документы нужны?
                    <textarea
                      aria-label="Какие ещё документы нужны?"
                      maxLength={500}
                      required
                      value={otherDocument}
                      onChange={(event) => setOtherDocument(event.target.value)}
                      placeholder="Например: выписка по расчётному счёту за 6 месяцев"
                    />
                  </label>
                )}
                {requestError && (
                  <p className="error-text" role="alert">
                    {requestError}
                  </p>
                )}
                <div className="modal-actions">
                  <button
                    className="button-secondary"
                    type="button"
                    onClick={() => setAction(null)}
                  >
                    Отмена
                  </button>
                  <button className="button-primary" type="submit">
                    Отправить запрос
                  </button>
                </div>
              </form>
            ) : action === 'terms' ? (
              <form className="review-terms" onSubmit={terms}>
                <h2>Альтернативные условия</h2>
                <SelectField
                  label="Вид кредита"
                  value={products.find((product) => product.id === productId)!.name}
                  options={products.map((product) => product.name)}
                  onChange={changeProduct}
                />
                <MoneyInput
                  label="Сумма кредита"
                  value={amount}
                  min={creditDetails[productId].minAmount}
                  max={
                    (products.find((product) => product.id === productId)?.amount ?? 150) *
                    1_000_000
                  }
                  onChange={setAmount}
                />
                <div className="review-term-row">
                  <SelectField
                    label="Срок"
                    value={formatTerm(months)}
                    options={termOptions.map(formatTerm)}
                    onChange={(text) =>
                      setMonths(termOptions.find((number) => formatTerm(number) === text)!)
                    }
                  />
                  <SelectField
                    label={productId === 'sellers' ? 'Месячная переплата' : 'Ставка'}
                    value={rateLabel(rate)}
                    options={[0, 1.2, 13, 14, 14.5, 16.5, 19, 23].map(rateLabel)}
                    onChange={(text) => setRate(Number(text.replace('%', '').replace(',', '.')))}
                  />
                </div>
                <SelectField
                  label="Условия по обеспечению"
                  value={collateral}
                  options={['Без изменения', 'Без залога', 'С залогом']}
                  onChange={setCollateral}
                />
                <label className="field review-comment">
                  Комментарий к решению
                  <textarea
                    placeholder="Укажите причину изменения условий…"
                    maxLength={500}
                    value={comment}
                    onChange={(event) => setComment(event.target.value)}
                  />
                  <span>{comment.length}/500</span>
                </label>
                <div className="modal-actions">
                  <button
                    className="button-secondary"
                    type="button"
                    onClick={() => setAction(null)}
                  >
                    Отмена
                  </button>
                  <button className="button-primary" type="submit">
                    Сохранить предложение
                  </button>
                </div>
              </form>
            ) : (
              <>
                <section className="decision-actions">
                  <h2>Принятие решения</h2>
                  <button
                    className="button-white"
                    type="button"
                    onClick={() => setAction('request')}
                  >
                    Запросить данные
                  </button>
                  <button className="button-white" type="button" onClick={() => setAction('terms')}>
                    Другие условия
                  </button>
                  <button
                    className={`button-white ${action === 'approve' ? 'is-chosen' : ''}`}
                    type="button"
                    onClick={() => setAction('approve')}
                  >
                    ✓ &nbsp; Одобрить
                  </button>
                  <button
                    className={`button-white ${action === 'reject' ? 'is-chosen' : ''}`}
                    type="button"
                    onClick={() => setAction('reject')}
                  >
                    × &nbsp; Отказать
                  </button>
                </section>
                {action && (
                  <section className="decision-confirm">
                    <h2>Подтвердите решение</h2>
                    <div>
                      <button
                        className="button-primary"
                        type="button"
                        onClick={() => save(action === 'approve' ? 'Одобрено' : 'Отклонено')}
                      >
                        Подтвердить
                      </button>
                      <button
                        className="button-secondary"
                        type="button"
                        onClick={() => setAction(null)}
                      >
                        Отмена
                      </button>
                    </div>
                  </section>
                )}
              </>
            ))}
        </div>
      </section>
    </main>
  )
}
