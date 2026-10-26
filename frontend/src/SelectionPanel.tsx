import { useState } from 'react'
import type { FormEvent } from 'react'
import type { Selection } from './types'
import { asset, Icon, MoneyInput, RangeField, Segmented, SelectField } from './ui'
import { money } from './products'

export const purposes = [
  'Развитие бизнеса',
  'Закупка товара',
  'Оплата поставщикам',
  'Покупка оборудования',
  'Рефинансирование',
  'Покупка недвижимости',
  'Другое',
]
export function SelectionPanel({
  value,
  onChange,
  onComplete,
  onCollapse,
}: {
  value: Selection
  onChange: (value: Selection) => void
  onComplete: () => void
  onCollapse: () => void
}) {
  const [attempted, setAttempted] = useState(false)
  const [help, setHelp] = useState(false)
  function update<K extends keyof Selection>(key: K, next: Selection[K]) {
    onChange({ ...value, [key]: next })
  }
  function submit(event: FormEvent) {
    event.preventDefault()
    setAttempted(true)
    if (!value.registration || !value.age) {
      requestAnimationFrame(() =>
        document.querySelector<HTMLElement>('.selection-panel [aria-invalid="true"]')?.focus(),
      )
      return
    }
    onComplete()
  }
  const helpMin = Math.min(
    50_000_000,
    Math.max(50_000, Math.round((value.turnover * 1.2) / 50_000) * 50_000),
  )
  const helpMax = Math.min(
    50_000_000,
    Math.max(helpMin, Math.round((value.turnover * 3.2) / 50_000) * 50_000),
  )
  return (
    <section className="selection-panel" aria-label="Умный подбор кредита">
      <div className="selection-heading">
        <span>
          <img src={asset('sparkles.png')} width="39" height="39" alt="" />
          Умный подбор кредита
        </span>
        <button type="button" onClick={onCollapse}>
          Скрыть <img src={asset('chevron.png')} alt="" className="collapse-chevron" />
        </button>
      </div>
      <form onSubmit={submit} noValidate>
        <div className="selection-grid">
          <div className="business-information">
            <div className="selection-section-heading">
              <span className="step-number first-step">1</span>
              <h2>Сведения о бизнесе</h2>
              <span className="heading-rule">
                <Icon name="section-rule" />
              </span>
            </div>
            <div className="business-fields">
              <div className="business-column">
                <div>
                  <SelectField
                    label="Для чего вашему бизнесу нужны деньги?"
                    value={value.purpose}
                    options={purposes}
                    onChange={(next) => update('purpose', next)}
                  />
                  {value.purpose === 'Другое' && (
                    <div className="other-purpose">
                      <textarea
                        aria-label="Опишите цель кредита"
                        placeholder="Например: сезонные расходы, участие в тендерах, выплата дивидендов и т.д."
                        maxLength={300}
                        value={value.other}
                        onChange={(event) => update('other', event.target.value)}
                      />
                      <span>{value.other.length}/300</span>
                    </div>
                  )}
                </div>
                <SelectField
                  label="Какой у вас тип регистрации бизнеса?"
                  required
                  value={value.registration}
                  options={['ООО', 'ИП', 'АО публичное', 'АО непубличное', 'Другое']}
                  onChange={(next) => update('registration', next)}
                  error={attempted && !value.registration}
                />
                <SelectField
                  label="Сколько лет работает ваш бизнес?"
                  required
                  value={value.age}
                  options={[
                    'Менее 9 месяцев',
                    'От 9 месяцев до 1 года',
                    'От 1 года до 3 лет',
                    'Более 3 лет',
                  ]}
                  onChange={(next) => update('age', next)}
                  error={attempted && !value.age}
                />
              </div>
              <div className="business-column business-extra">
                <Segmented
                  label="Есть ли залог?"
                  value={value.collateral}
                  options={['Да', 'Нет', 'Рассмотрю позже']}
                  onChange={(next) => update('collateral', next)}
                  hint={
                    <>
                      Наличие залога может увеличить сумму кредита
                      <br />и снизить ставку
                    </>
                  }
                />
                <MoneyInput
                  label="Среднемесячный оборот"
                  value={value.turnover}
                  onChange={(next) => update('turnover', next)}
                />
                <SelectField
                  label="Сфера деятельности"
                  value={value.industry}
                  options={[
                    'Торговля',
                    'Производство',
                    'Строительство',
                    'Услуги',
                    'Транспорт и логистика',
                    'Сельское хозяйство',
                    'Общественное питание',
                    'IT и информационные технологии',
                    'Медицина и здравоохранение',
                    'Образование',
                    'Недвижимость',
                    'Гостиничный бизнес и туризм',
                    'Другое',
                  ]}
                  onChange={(next) => update('industry', next)}
                />
                <Segmented
                  label="Есть ли действующие кредиты?"
                  value={value.loans}
                  options={['Да', 'Нет', 'Не могу ответить']}
                  onChange={(next) => update('loans', next)}
                />
              </div>
            </div>
          </div>
          <div className="amount-information">
            <div className="selection-section-heading">
              <span className="step-number">2</span>
              <h2>Сумма и срок</h2>
            </div>
            <RangeField
              kind="amount"
              value={value.amount}
              onChange={(next) => update('amount', next)}
            />
            <button
              className="button-outline amount-help"
              type="button"
              onClick={() => setHelp(!help)}
              aria-expanded={help}
            >
              ✦ Помогите определить сумму
            </button>
            {help && (
              <div className="amount-suggestion" role="status">
                <strong>
                  {money(helpMin)} — {money(helpMax)} ₽
                </strong>
                <p>Этого может быть достаточно для достижения вашей цели</p>
              </div>
            )}
            <RangeField
              kind="term"
              value={value.months}
              onChange={(next) => update('months', next)}
            />
          </div>
        </div>
        <div className="selection-footer">
          <p>* - поля, обязательные для заполнения</p>
          <button className="button-primary" type="submit">
            Завершить подбор
          </button>
        </div>
      </form>
    </section>
  )
}
