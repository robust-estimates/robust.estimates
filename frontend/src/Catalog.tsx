import { frontendProductId, useCalculation } from './api'
import type { Recommendation } from './api'
import { useMemo, useState } from 'react'
import { categories, formatTerm, loanFigures, money, products, rateLabel } from './products'
import type { Product } from './products'
import type { Selection } from './types'
import { asset, Icon, Modal, SearchField } from './ui'
import { SelectionPanel } from './SelectionPanel'
import { creditDetails } from './creditData'

export function ProductCard({
  product,
  onSelect,
  mode = 'catalog',
  selection,
  disabled = false,
}: {
  product: Product
  onSelect: (product: Product) => void
  mode?: 'catalog' | 'matched' | 'alternative'
  selection?: Selection
  disabled?: boolean
}) {
  const calculation = useCalculation(
    product.id,
    selection?.amount || 0,
    selection?.months || 0,
    !!selection && !disabled,
  )
  const rate = product.rate
  const total = selection
    ? (calculation?.total ?? loanFigures(selection.amount, rate, selection.months).total)
    : 0
  return (
    <article
      className={`product-card ${disabled ? 'product-unavailable' : ''}`}
      data-product-id={product.id}
    >
      <div className="product-heading">
        <div className={`product-icon ${product.color}`}>
          <Icon name={product.icon} />
        </div>
        <div>
          <h3>{product.name}</h3>
          <p>{product.description}</p>
        </div>
      </div>
      <dl className="product-stats">
        <div>
          <dt>{mode === 'catalog' ? 'Ставка от' : 'Ставка'}</dt>
          <dd>{rateLabel(rate)}</dd>
        </div>
        {mode === 'matched' && product.id === 'business' && selection ? (
          <>
            <div>
              <dt>Сумма выплат</dt>
              <dd>{money(total)} ₽</dd>
            </div>
            <div>
              <dt>Переплата</dt>
              <dd>{money(total - selection.amount)} ₽</dd>
            </div>
          </>
        ) : mode === 'catalog' ? (
          <>
            <div>
              <dt>{product.limit ? 'Лимит' : 'Сумма'}</dt>
              <dd>до {product.amount} млн ₽</dd>
            </div>
            <div>
              <dt>{product.termTitle || 'Срок'}</dt>
              <dd>{product.term || `до ${formatTerm(product.months)}`}</dd>
            </div>
          </>
        ) : null}
      </dl>
      <button
        className="details-button"
        type="button"
        disabled={disabled}
        aria-label={`${mode === 'catalog' ? 'Подробнее' : 'Выбрать продукт'}: ${product.name}`}
        onClick={(event) => {
          event.currentTarget.focus()
          onSelect(product)
        }}
      >
        {mode === 'catalog' ? 'Подробнее' : 'Выбрать продукт'} <span aria-hidden="true">›</span>
      </button>
    </article>
  )
}

export function Catalog({
  expanded,
  selection,
  onSelection,
  onExpand,
  onCollapse,
  onComplete,
  onSelect,
}: {
  expanded: boolean
  selection: Selection
  onSelection: (value: Selection) => void
  onExpand: () => void
  onCollapse: () => void
  onComplete: () => void
  onSelect: (product: Product) => void
}) {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState('Все продукты')
  const [details, setDetails] = useState<Product | null>(null)
  const filtered = useMemo(
    () =>
      products.filter(
        (product) =>
          (category === 'Все продукты' || product.categories.includes(category)) &&
          `${product.name} ${product.description}`
            .toLocaleLowerCase('ru')
            .includes(query.trim().toLocaleLowerCase('ru')),
      ),
    [query, category],
  )
  return (
    <main className={`catalog-page page-container ${expanded ? 'is-expanded' : ''}`}>
      {expanded ? (
        <SelectionPanel
          value={selection}
          onChange={onSelection}
          onComplete={onComplete}
          onCollapse={onCollapse}
        />
      ) : (
        <section className="hero" aria-labelledby="hero-title">
          <h1 id="hero-title">Умный подбор кредита</h1>
          <p>
            Подберём оптимальные кредитные решения под ваш бизнес
            <br />
            на основе данных компании и ваших целей
          </p>
          <button className="button-primary hero-cta" type="button" onClick={onExpand}>
            Подобрать кредит <Icon name="arrow-right" />
          </button>
        </section>
      )}
      <section className="catalog" aria-labelledby="catalog-title" id="catalog">
        <div className="catalog-toolbar">
          <div>
            <h2 id="catalog-title">Каталог кредитных продуктов</h2>
            <p>Готовые решения для разных задач вашего бизнеса</p>
          </div>
          <SearchField
            value={query}
            onChange={setQuery}
            placeholder="Поиск по продуктам, например «оборотные средства»"
            label="Поиск по кредитным продуктам"
          />
        </div>
        <div className="category-filters" aria-label="Категории продуктов">
          {categories.map((item) => (
            <button
              key={item}
              className={`category-chip ${item === category ? 'selected' : ''}`}
              type="button"
              aria-pressed={item === category}
              onClick={() => setCategory(item)}
            >
              {item}
            </button>
          ))}
        </div>
        <span className="sr-only" role="status">
          Найдено продуктов: {filtered.length}
        </span>
        <div className="product-grid">
          {filtered.map((product) => (
            <ProductCard key={product.id} product={product} onSelect={setDetails} />
          ))}
        </div>
        {!filtered.length && (
          <div className="empty-state">
            <h3>Ничего не найдено</h3>
            <p>Попробуйте изменить запрос или категорию.</p>
            <button
              className="button-secondary"
              type="button"
              onClick={() => {
                setCategory('Все продукты')
                setQuery('')
              }}
            >
              Сбросить фильтры
            </button>
          </div>
        )}
      </section>
      {details && (
        <Modal
          title={details.name}
          onClose={() => setDetails(null)}
          className="product-details-modal"
        >
          <div className="product-details-intro">
            <div className={`product-icon ${details.color}`}>
              <Icon name={details.icon} />
            </div>
            <p>{details.description}</p>
          </div>
          <dl className="product-details-conditions">
            {[
              ['Ставка', creditDetails[details.id].rateLabel],
              [details.limit ? 'Лимит' : 'Сумма кредита', creditDetails[details.id].amountLabel],
              [
                details.termTitle === 'льготный период' ? 'Льготный период' : 'Срок',
                creditDetails[details.id].termLabel,
              ],
              ['Обеспечение', details.secured ? 'С залогом' : 'Без залога'],
            ].map(([label, value]) => (
              <div key={label}>
                <dt>{label}</dt>
                <dd>{value}</dd>
              </div>
            ))}
          </dl>
          {[
            ['Цель кредита', [creditDetails[details.id].purpose]],
            ['Варианты погашения', creditDetails[details.id].repayment],
            ['Особенности', creditDetails[details.id].features],
            ['Требования к заёмщику', creditDetails[details.id].borrower],
          ].map(([title, lines]) => (
            <section className="product-details-section" key={String(title)}>
              <h3>{title}</h3>
              <ul>
                {(lines as string[]).map((line) => (
                  <li key={line}>{line}</li>
                ))}
              </ul>
            </section>
          ))}
          <div className="product-details-actions">
            <button className="button-secondary" type="button" onClick={() => setDetails(null)}>
              Вернуться в каталог
            </button>
            <button className="button-primary" type="button" onClick={() => onSelect(details)}>
              Оформить заявку
            </button>
          </div>
        </Modal>
      )}
    </main>
  )
}

export function SelectionResults({
  selection,
  recommendations,
  onReset,
  onSelect,
}: {
  selection: Selection
  recommendations: Recommendation[]
  onReset: () => void
  onSelect: (product: Product) => void
}) {
  const groupProducts = (availability: string) =>
    recommendations
      .filter((item) => item.availability === availability)
      .sort((a, b) => a.rank - b.rank)
      .map((item) => products.find((product) => product.id === frontendProductId(item.product.id)))
      .filter((item): item is Product => !!item)
  const matched = groupProducts('BEST_MATCH')
  const alternatives = groupProducts('ALTERNATIVE')
  const unavailable = groupProducts('UNAVAILABLE')
  const groups = [
    {
      title: 'Наиболее подходящие',
      text: 'Эти продукты лучше всего соответствуют вашим целям',
      items: matched,
      kind: 'matched' as const,
    },
    {
      title: 'Доступные альтернативы',
      text: 'Эти продукты вам доступны, но менее точно соответствуют заявленным целям',
      items: alternatives,
      kind: 'alternative' as const,
    },
    {
      title: 'Недоступные продукты',
      text: 'Сейчас эти продукты недоступны. Вы можете изменить параметры или связаться с менеджером для получения консультации',
      items: unavailable,
      kind: 'catalog' as const,
    },
  ]
  return (
    <main className="results-page page-container">
      <div className="results-heading">
        <span>
          <img src={asset('sparkles.png')} width="39" height="39" alt="" />
          <span>
            Умный подбор кредита<strong>Результат</strong>
          </span>
        </span>
        <button type="button" onClick={onReset}>
          Сбросить ×
        </button>
      </div>
      {groups.map((group) => (
        <section
          className={`result-group result-${group.kind}`}
          key={group.title}
          aria-label={group.title}
        >
          <div className="result-group-intro">
            <h2>{group.title}</h2>
            <p>{group.text}</p>
          </div>
          <div className="result-products">
            {group.items.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onSelect={onSelect}
                mode={group.kind}
                selection={selection}
                disabled={group.kind === 'catalog'}
              />
            ))}
            {!group.items.length && (
              <p className="result-empty">
                {group.kind === 'matched'
                  ? 'По указанным параметрам подходящих предложений пока нет.'
                  : 'Нет продуктов в этой категории.'}
              </p>
            )}
          </div>
        </section>
      ))}
      <p className="sr-only" role="status">
        Подбор завершён. Сумма {money(selection.amount)} рублей, срок {formatTerm(selection.months)}
        .
      </p>
    </main>
  )
}
