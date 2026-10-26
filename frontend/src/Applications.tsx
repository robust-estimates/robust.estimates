import { useRef, useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'
import { formatTerm, loanFigures, money, products, quoteForProduct, rateLabel } from './products'
import type { Product } from './products'
import { defaultApplicant, isApplicationDraft } from './types'
import type { Applicant, ApplicationStatus, LoanApplication } from './types'
import { asset, Icon, Modal, SearchField } from './ui'
import { clearSessionValue, useSessionValue } from './session'
import { creditDetails, requestOptions, requirementsFor } from './creditData'
import type { DocumentRequirement } from './creditData'
import { DocumentUploads, legacyDocumentFiles } from './DocumentUploads'
export function Indicators({
  product,
  amount,
  rate,
  months,
  showProduct = false,
}: {
  product: Product
  amount: number
  rate: number
  months: number
  showProduct?: boolean
}) {
  const figures = loanFigures(amount, rate, months, product.id === 'sellers')
  const rows = [
    ...(showProduct ? [['Вид кредита', product.name]] : []),
    ['Сумма кредита', `${money(amount)} ₽`],
    ['Срок', formatTerm(months)],
    ['Ставка', `${rateLabel(rate)}${product.id === 'sellers' ? ' в месяц' : ''}`],
    ['Сумма выплат', `${money(figures.total)} ₽`],
    ['Переплата', `${money(figures.overpayment)} ₽`],
    ['Ежемесячный платеж', `${money(figures.monthly)} ₽`],
  ]
  return (
    <dl className="indicators">
      {rows.map(([label, value]) => (
        <div key={label}>
          <dt>{label}</dt>
          <dd>{value}</dd>
        </div>
      ))}
    </dl>
  )
}
export function ApplicantFields({
  value,
  onChange,
  readOnly = false,
  required = true,
}: {
  value: Applicant
  onChange?: (value: Applicant) => void
  readOnly?: boolean
  required?: boolean
}) {
  function field(
    key: keyof Applicant,
    label: string,
    options: { type?: string; required?: boolean; compact?: boolean } = {},
  ) {
    return (
      <label className={`field applicant-field ${options.compact ? 'compact-field' : ''}`}>
        {label}
        {options.required && required && '*'}
        <input
          type={options.type || 'text'}
          value={value[key]}
          onChange={(event) => onChange?.({ ...value, [key]: event.target.value })}
          readOnly={readOnly}
          required={options.required && required}
          pattern={
            key === 'inn'
              ? '(?:[0-9]{10}|[0-9]{12})'
              : key === 'phone'
                ? '(?:\\+7|8)(?:[\\s\\(\\)\\-]*[0-9]){10}'
                : key === 'name'
                  ? '.*\\S.*'
                  : key === 'email'
                    ? '[^\\s@]+@[^\\s@]+\\.[^\\s@]+'
                    : undefined
          }
          max={
            key === 'registered'
              ? new Intl.DateTimeFormat('en-CA', {
                  timeZone: 'Europe/Moscow',
                  year: 'numeric',
                  month: '2-digit',
                  day: '2-digit',
                }).format(new Date())
              : undefined
          }
          maxLength={key === 'inn' ? 12 : 255}
          autoComplete={
            {
              company: 'organization',
              inn: 'off',
              registered: 'off',
              name: 'name',
              position: 'organization-title',
              phone: 'tel',
              email: 'email',
            }[key]
          }
        />
      </label>
    )
  }
  return (
    <>
      <fieldset className="applicant-card">
        <legend className="sr-only">Информация о компании</legend>
        <h3>Информация о компании</h3>
        {field('company', 'Наименование компании')}
        {field('inn', 'ИНН', { required: true, compact: true })}
        {field('registered', 'Дата регистрации', { type: 'date', required: true, compact: true })}
      </fieldset>
      <fieldset className="applicant-card">
        <legend className="sr-only">Данные представителя</legend>
        <h3>Данные представителя</h3>
        {field('name', 'ФИО ', { required: true })}
        {field('position', 'Должность')}
        <div className="contact-fields">
          {field('phone', 'Мобильный телефон ', { type: 'tel', required: true })}
          {field('email', 'Email ', { type: 'email', required: true })}
        </div>
      </fieldset>
    </>
  )
}
export function DocumentList({
  documents,
  onRemove,
  onAdd,
}: {
  documents: string[]
  onRemove?: (index: number) => void
  onAdd?: (files: string[]) => void
}) {
  const input = useRef<HTMLInputElement>(null)
  const [message, setMessage] = useState('')
  function files(event: ChangeEvent<HTMLInputElement>) {
    const selected = Array.from(event.target.files ?? [])
    event.target.value = ''
    if (selected.some((file) => !/\.(pdf|docx?|xlsx?|png|jpe?g)$/i.test(file.name))) {
      setMessage('Поддерживаются PDF, Word, Excel, PNG и JPEG.')
      return
    }
    if (selected.some((file) => !file.size)) {
      setMessage('Нельзя прикрепить пустой файл.')
      return
    }
    if (selected.some((file) => file.size > 20 * 1024 * 1024)) {
      setMessage('Размер файла не должен превышать 20 МБ.')
      return
    }
    setMessage('')
    onAdd?.(selected.map((file) => file.name))
  }
  return (
    <>
      <div className="document-list">
        {documents.map((name, index) => (
          <div className="document-row" key={`${name}-${index}`}>
            <span className="document-icon">
              <Icon name="doc" />
            </span>
            <span>{name}</span>
            {onRemove && (
              <button
                type="button"
                onClick={() => onRemove(index)}
                aria-label={`Удалить документ ${name}`}
              >
                ×
              </button>
            )}
          </div>
        ))}
        {!documents.length && <p className="documents-empty">Прикрепите необходимые документы</p>}
      </div>
      {onAdd && (
        <>
          <input
            className="sr-only"
            ref={input}
            type="file"
            aria-label="Прикрепить документы"
            multiple
            accept=".pdf,.doc,.docx,.xls,.xlsx,.png,.jpg,.jpeg"
            onChange={files}
            tabIndex={-1}
          />
          <button
            type="button"
            className="text-button attach-button"
            onClick={() => input.current?.click()}
          >
            + Прикрепить документ
          </button>
        </>
      )}
      {message && (
        <p className="error-text" role="alert">
          {message}
        </p>
      )}
    </>
  )
}

export function ApplicationForm({
  product,
  amount,
  months,
  existing,
  registration,
  hasLoans,
  onCancel,
  onChangeParameters,
  onSubmit,
}: {
  product: Product
  amount: number
  months: number
  existing?: LoanApplication
  registration: string
  hasLoans: boolean
  onCancel: () => void
  onChangeParameters: (amount: number, months: number) => void
  onSubmit: (value: Omit<LoanApplication, 'id' | 'date' | 'status'>) => void
}) {
  const editing = !!existing
  const request = existing?.status === 'Запрошены данные' ? existing.documentRequest : undefined
  const draftKey = `alfa-credit-form-${existing?.id ?? product.id}${request ? `-${request.requestedAt}` : ''}`
  const initialFiles = { ...(existing?.documentFiles ?? legacyDocumentFiles(existing?.documents)) }
  if (request && !request.fulfilledAt) for (const id of request.items) delete initialFiles[id]
  const [draft, setDraft] = useSessionValue(
    draftKey,
    {
      applicant: existing?.applicant ?? defaultApplicant,
      documents: existing?.documents ?? [],
      documentFiles: initialFiles,
      amount,
      months,
    },
    isApplicationDraft,
  )
  const { applicant } = draft
  const files = draft.documentFiles ?? legacyDocumentFiles(draft.documents)
  const previousTerms = existing?.productId === product.id ? existing : undefined
  const minAmount = Math.min(creditDetails[product.id].minAmount, previousTerms?.amount ?? Infinity)
  const maxAmount = Math.max(product.amount * 1_000_000, previousTerms?.amount ?? 0)
  const minMonths = Math.min(creditDetails[product.id].minMonths, previousTerms?.months ?? Infinity)
  const maxMonths = Math.max(product.months, previousTerms?.months ?? 0)
  const parameters = previousTerms
    ? {
        productId: product.id,
        amount: Math.max(minAmount, Math.min(maxAmount, draft.amount ?? amount)),
        months: Math.max(minMonths, Math.min(maxMonths, draft.months ?? months)),
      }
    : quoteForProduct(product, {
        productId: product.id,
        amount: draft.amount ?? amount,
        months: draft.months ?? months,
      })
  const documentRequirements: DocumentRequirement[] = requirementsFor(
    product.id,
    registration,
    hasLoans,
  ).map((item) => ({ ...item, required: editing ? false : item.required }))
  for (const id of request?.items ?? []) {
    const found = documentRequirements.find((item) => item.id === id)
    if (found) found.required = true
    else
      documentRequirements.push({
        id,
        title:
          id === 'other'
            ? request!.other
            : (requestOptions.find((option) => option.id === id)?.title ?? id),
        required: true,
        description: 'Запрошено кредитным специалистом.',
      })
  }
  const [consent, setConsent] = useState(false)
  const [documentError, setDocumentError] = useState(false)
  const rate =
    existing?.productId === product.id
      ? existing.rate
      : product.id === 'business'
        ? 19
        : product.rate
  function submit(event: FormEvent) {
    event.preventDefault()
    const incomplete = documentRequirements.some((item) => item.required && !files[item.id])
    setDocumentError(incomplete)
    if (incomplete) {
      requestAnimationFrame(() =>
        document
          .querySelector<HTMLElement>('.document-upload-button[aria-invalid="true"]')
          ?.focus(),
      )
      return
    }
    const previousFiles = existing?.documentFiles ?? legacyDocumentFiles(existing?.documents)
    const documents = [
      ...new Set([
        ...(existing?.documents.filter((name) => !Object.values(previousFiles).includes(name)) ??
          []),
        ...Object.values(files),
      ]),
    ]
    onSubmit({
      productId: product.id,
      amount: parameters.amount,
      months: parameters.months,
      rate,
      applicant,
      documents,
      documentFiles: files,
      companyProfile: existing?.companyProfile ?? {
        registrationStatus: 'Не проверено',
        director: applicant.name,
        founders: [],
        okved: [],
        bankruptcy: null,
        convictions: null,
        activeLoans: hasLoans,
      },
      ...(existing?.documentRequest
        ? {
            documentRequest: {
              ...existing.documentRequest,
              ...(request ? { fulfilledAt: new Date().toISOString() } : {}),
            },
          }
        : {}),
    })
    clearSessionValue(draftKey)
  }
  return (
    <main className="application-form-page">
      <h1>{editing ? 'Изменение заявки' : 'Заявка на бизнес-кредит'}</h1>
      <p className="application-subtitle">
        {editing ? (
          `Заявка № ${existing.id}. Измените данные и сохраните изменения.`
        ) : (
          <>
            Заполните данные и отправьте заявку. Наши специалисты
            <br />
            рассмотрят её и свяжутся с вами
          </>
        )}
      </p>
      <form className="application-form" onSubmit={submit}>
        <div className="application-applicant">
          <ApplicantFields
            value={applicant}
            onChange={(value) => setDraft((current) => ({ ...current, applicant: value }))}
          />
          {!editing && (
            <label className="consent">
              <input
                type="checkbox"
                checked={consent}
                onChange={(event) => setConsent(event.target.checked)}
                required
              />
              Согласие на обработку персональных данных*
            </label>
          )}
          <p className="required-note">* - поля, обязательные для заполнения</p>
        </div>
        <div className="application-summary">
          <h2>Выбранный продукт</h2>
          <div className="selected-product">
            <div className={`product-icon ${product.color}`}>
              <Icon name={product.icon} />
            </div>
            <strong>{product.name}</strong>
          </div>
          <fieldset className="application-parameters">
            <legend>Сумма и срок кредита</legend>
            <label className="field">
              Сумма кредита, ₽
              <input
                type="number"
                inputMode="numeric"
                aria-label="Сумма кредита, ₽"
                min={minAmount}
                max={maxAmount}
                step="1"
                required
                value={draft.amount ?? amount}
                onChange={(event) =>
                  setDraft((current) => ({ ...current, amount: Number(event.target.value) }))
                }
              />
              <small>{creditDetails[product.id].amountLabel}</small>
            </label>
            <label className="field">
              Срок кредита, месяцев
              <input
                type="number"
                inputMode="numeric"
                aria-label="Срок кредита, месяцев"
                min={minMonths}
                max={maxMonths}
                step="1"
                required
                value={draft.months ?? months}
                onChange={(event) =>
                  setDraft((current) => ({ ...current, months: Number(event.target.value) }))
                }
              />
              <small>
                От {minMonths} до {maxMonths} месяцев
              </small>
            </label>
          </fieldset>
          <h2>Показатели</h2>
          <Indicators
            product={product}
            amount={parameters.amount}
            months={parameters.months}
            rate={rate}
          />
          <h2>Документы*</h2>
          {request && (
            <p className="requested-documents-note">
              Специалист запросил дополнительные документы. Прикрепите новые файлы в отмеченные
              поля.
            </p>
          )}
          <DocumentUploads
            requirements={documentRequirements}
            files={files}
            attempted={documentError}
            onChange={(documentFiles) => {
              setDraft((current) => ({ ...current, documentFiles }))
              setDocumentError(false)
            }}
          />
          {documentError && (
            <p role="alert" className="error-text">
              Прикрепите все обязательные документы.
            </p>
          )}
          <div className="application-actions">
            <button
              className="button-secondary"
              type="button"
              onClick={() => {
                if (editing) {
                  clearSessionValue(draftKey)
                  onCancel()
                } else onChangeParameters(parameters.amount, parameters.months)
              }}
            >
              {editing ? 'Отмена' : 'Изменить данные'}
            </button>
            <button className="button-primary" type="submit">
              {editing ? 'Сохранить изменения' : 'Оформить заявку'}
            </button>
          </div>
        </div>
      </form>
    </main>
  )
}

export function ApplicationAccepted({
  application,
  onList,
  onCatalog,
}: {
  application: LoanApplication
  onList: () => void
  onCatalog: () => void
}) {
  return (
    <main className="accepted-page">
      <img src={asset('success.png')} width="121" height="121" alt="" />
      <h1>Заявка отправлена!</h1>
      <div className="accepted-number">
        <p>Номер заявки</p>
        <strong>№ {application.id}</strong>
      </div>
      <div className="accepted-info">
        <h2>Что дальше?</h2>
        <p>
          Мы проведем предварительную проверку
          <br />и свяжемся с вами в течение 1 рабочего дня
        </p>
        <h2>Просмотр заявки</h2>
        <p>Статус заявки будет доступен в разделе «Мои заявки»</p>
      </div>
      <button className="button-primary" type="button" onClick={onList}>
        Перейти к моим заявкам
      </button>
      <button className="button-secondary" type="button" onClick={onCatalog}>
        Вернуться в каталог
      </button>
    </main>
  )
}

export function StatusBadge({ status }: { status: ApplicationStatus }) {
  return (
    <span
      className={`status-badge status-${status === 'На рассмотрении' || status === 'Запрошены данные' ? 'pending' : status === 'Отклонено' ? 'rejected' : status === 'Черновик' ? 'draft' : 'approved'}`}
    >
      {status}
    </span>
  )
}
function formatDate(date: string) {
  const value = new Date(date)
  const day = new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'Europe/Moscow',
  })
    .format(value)
    .replace(' г.', '')
  const time = new Intl.DateTimeFormat('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'Europe/Moscow',
  }).format(value)
  return `${day}, ${time}`
}
export function ApplicationsList({
  applications,
  onOpen,
  onNew,
}: {
  applications: LoanApplication[]
  onOpen: (application: LoanApplication) => void
  onNew: () => void
}) {
  const [query, setQuery] = useState('')
  const [status, setStatus] = useState('Все заявки')
  const filtered = applications.filter(
    (application) =>
      (status === 'Все заявки' || application.status === status) &&
      `${application.id} ${products.find((product) => product.id === application.productId)?.name}`
        .toLowerCase()
        .includes(query.trim().toLowerCase()),
  )
  return (
    <main className="applications-page">
      <div className="applications-heading">
        <h1>Мои заявки</h1>
        <button className="button-primary" type="button" onClick={onNew}>
          <span aria-hidden="true">＋</span> Новая заявка
        </button>
      </div>
      <div className="applications-controls">
        <SearchField
          value={query}
          onChange={setQuery}
          placeholder="Поиск по номеру заявки или продукту…"
          label="Поиск заявок"
        />
        <label className="status-select">
          <span className="sr-only">Статус заявки</span>
          <select value={status} onChange={(event) => setStatus(event.target.value)}>
            {[
              'Все заявки',
              'На рассмотрении',
              'Одобрено',
              'Другие условия',
              'Отклонено',
              'Черновик',
              'Запрошены данные',
            ].map((item) => (
              <option key={item}>{item}</option>
            ))}
          </select>
          <img src={asset('chevron.png')} alt="" />
        </label>
      </div>
      <section className="applications-table" aria-label="Список заявок">
        {filtered.map((application) => (
          <button
            className="application-row"
            type="button"
            key={application.id}
            onClick={() => onOpen(application)}
          >
            <span className="application-row-title">
              <strong>№ {application.id}</strong>
              <span>{products.find((product) => product.id === application.productId)?.name}</span>
              <small>{formatDate(application.date)}</small>
            </span>
            <span className="application-row-value">
              <span>Сумма</span>
              {money(application.amount)} ₽
            </span>
            <span className="application-row-value">
              <span>Срок</span>
              {formatTerm(application.months)}
            </span>
            <StatusBadge status={application.status} />
            <span className="row-arrow" aria-hidden="true">
              ›
            </span>
          </button>
        ))}
        {!filtered.length && (
          <div className="empty-state">
            <h2>Заявки не найдены</h2>
            <p>Измените поиск или создайте новую заявку.</p>
          </div>
        )}
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
      </section>
    </main>
  )
}

export function ApplicationDetail({
  application,
  onBack,
  onEdit,
  onNew,
  onUpdate,
}: {
  application: LoanApplication
  onBack: () => void
  onEdit: () => void
  onNew: () => void
  onUpdate: (application: LoanApplication) => void
}) {
  const [terms, setTerms] = useState(false)
  const [signed, setSigned] = useState(false)
  const product = products.find((item) => item.id === application.productId) ?? products[0]
  const complete = !['На рассмотрении', 'Черновик', 'Запрошены данные'].includes(application.status)
  const alternative = application.alternative ?? {
    productId: 'credit-line',
    amount: 10_000_000,
    months: 48,
    rate: 19,
    collateral: 'Без залога',
    comment: '',
  }
  return (
    <main className="application-detail-page">
      <button className="back-button" type="button" onClick={onBack}>
        ← &nbsp; К списку заявок
      </button>
      <section className="application-detail">
        <div className="application-detail-main">
          <div className="application-detail-heading">
            <h1>Заявка № {application.id}</h1>
            <StatusBadge status={application.status} />
          </div>
          <p className="application-detail-date">от {formatDate(application.date)}</p>
          <ol className="application-timeline">
            {[
              'Заявка создана',
              'Данные получены',
              'Проверка компании',
              'Рассмотрение банком',
              'Решение банка',
            ].map((step, index) => (
              <li
                key={step}
                className={index < 3 || complete ? 'done' : index === 3 ? 'current' : ''}
              >
                <span>{index < 3 || complete ? '✓' : index + 1}</span>
                <p>{step}</p>
              </li>
            ))}
          </ol>
          <div className="detail-applicant-grid">
            <ApplicantFields value={application.applicant} readOnly required={false} />
          </div>
        </div>
        <aside className="application-detail-aside">
          <h2>Документы</h2>
          <DocumentList documents={application.documents} />
          <h2>Показатели</h2>
          <Indicators
            product={product}
            amount={application.amount}
            rate={application.rate}
            months={application.months}
            showProduct
          />
          {application.status === 'Одобрено' ? (
            <div className="decision-callout decision-approved">
              <h3>✓ &nbsp; Кредит одобрен</h3>
              <button
                className="button-primary"
                type="button"
                onClick={(event) => {
                  event.currentTarget.focus()
                  setSigned(true)
                }}
              >
                Перейти к подписанию
              </button>
            </div>
          ) : application.status === 'Другие условия' ? (
            <div className="decision-callout decision-alternative">
              <h3>Предложены другие условия</h3>
              <button
                className="button-primary"
                type="button"
                onClick={(event) => {
                  event.currentTarget.focus()
                  setTerms(true)
                }}
              >
                Ознакомиться
              </button>
            </div>
          ) : application.status === 'Отклонено' ? (
            <div className="decision-callout decision-rejected">
              <h3>✕ Отказ</h3>
              <p>
                Банк не может одобрить вашу заявку
                <br />
                на выбранных условиях
              </p>
              <button className="button-secondary" type="button" onClick={onNew}>
                Подать новую заявку
              </button>
            </div>
          ) : application.status === 'Запрошены данные' && application.documentRequest ? (
            <div className="decision-callout decision-requested">
              <h3>Нужны дополнительные документы</h3>
              <ul>
                {application.documentRequest.items.map((id) => (
                  <li key={id}>
                    {id === 'other'
                      ? application.documentRequest!.other
                      : requestOptions.find((option) => option.id === id)?.title}
                  </li>
                ))}
              </ul>
              <button className="button-primary" type="button" onClick={onEdit}>
                Добавить документы
              </button>
            </div>
          ) : (
            <button className="button-secondary detail-edit" type="button" onClick={onEdit}>
              Изменить заявку
            </button>
          )}
        </aside>
      </section>
      {terms && (
        <Modal title="Обновленные условия" onClose={() => setTerms(false)} className="terms-modal">
          <p className="terms-intro">
            Банк не может одобрить эту заявку, но готов
            <br />
            предложить альтернативные условия
          </p>
          <div className="terms-comparison">
            {[
              {
                title: 'Ваша заявка',
                product,
                amount: application.amount,
                months: application.months,
                rate: application.rate,
                collateral: product.secured ? 'С залогом' : 'Без залога',
              },
              {
                title: 'Предложение банка',
                product: products.find((item) => item.id === alternative.productId) ?? products[2],
                ...alternative,
              },
            ].map((column) => (
              <div key={column.title}>
                <h3>{column.title}</h3>
                <dl>
                  {[
                    [
                      'Вид кредита',
                      typeof column.product === 'object' ? column.product.name : product.name,
                    ],
                    ['Сумма кредита', `${money(column.amount)} ₽`],
                    ['Срок', formatTerm(column.months)],
                    [
                      'Ставка',
                      `${rateLabel(column.rate)}${column.product.id === 'sellers' ? ' в месяц' : ''}`,
                    ],
                    ['Обеспечение', column.collateral],
                  ].map(([label, value]) => (
                    <div key={label}>
                      <dt>{label}</dt>
                      <dd>{value}</dd>
                    </div>
                  ))}
                </dl>
              </div>
            ))}
          </div>
          <div className="modal-actions">
            <button
              className="button-primary"
              type="button"
              onClick={() => {
                onUpdate({
                  ...application,
                  productId: alternative.productId,
                  amount: alternative.amount,
                  months: alternative.months,
                  rate: alternative.rate,
                  status: 'Одобрено',
                })
                setTerms(false)
              }}
            >
              Принять предложение
            </button>
            <button
              className="button-secondary"
              type="button"
              onClick={() => {
                onUpdate({ ...application, status: 'Отклонено' })
                setTerms(false)
              }}
            >
              Отказаться
            </button>
          </div>
        </Modal>
      )}
      {signed && (
        <Modal title="Подписание документов" onClose={() => setSigned(false)}>
          <p className="demo-message">
            Это демонстрация фронтенда. Подписание и выдача кредита будут доступны после подключения
            банковского сервиса.
          </p>
          <button className="button-primary" type="button" onClick={() => setSigned(false)}>
            Понятно
          </button>
        </Modal>
      )}
    </main>
  )
}
