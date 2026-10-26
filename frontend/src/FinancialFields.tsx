import { useEffect, useState } from 'react'

type Financials = { revenue: number; expenses: number; debts: number }

function FinancialInput({ value, onChange }: { value: number; onChange: (value: number) => void }) {
  const [text, setText] = useState(String(value))

  useEffect(() => {
    setText((current) => (Number(current) === value ? current : String(value)))
  }, [value])

  return (
    <input
      type="number"
      required
      min="0"
      max="100000000000"
      step="0.01"
      value={text}
      onChange={(event) => {
        const next = event.target.value.replace(/^0+(?=\d)/, '')
        setText(next)
        onChange(Number(next))
      }}
      onBlur={() => {
        if (text !== '') setText(String(Number(text)))
      }}
    />
  )
}

export function FinancialFields({
  value,
  onChange,
  disabled = false,
}: {
  value: Financials
  onChange: (value: Financials) => void
  disabled?: boolean
}) {
  return (
    <fieldset className="financial-fields applicant-card" disabled={disabled}>
      <legend>Финансовые сведения за месяц</legend>
      {(
        [
          ['revenue', 'Выручка, ₽'],
          ['expenses', 'Расходы, ₽'],
          ['debts', 'Платежи по кредитам, ₽'],
        ] as const
      ).map(([key, label]) => (
        <label className="field" key={key}>
          {label}
          <FinancialInput
            value={value[key]}
            onChange={(next) => onChange({ ...value, [key]: next })}
          />
        </label>
      ))}
    </fieldset>
  )
}
