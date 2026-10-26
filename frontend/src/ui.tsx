import { useEffect, useId, useRef, useState } from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { formatTerm, money } from './products'

export const asset = (filename: string) => `${import.meta.env.BASE_URL}assets/${filename}`
export function Icon({ name, className = '' }: { name: string; className?: string }) {
  return <img alt="" className={`icon ${className}`} src={asset(`${name}.svg`)} draggable="false" />
}
export function Modal({
  title,
  children,
  onClose,
  className = '',
}: {
  title: string
  children: ReactNode
  onClose: () => void
  className?: string
}) {
  const ref = useRef<HTMLDialogElement>(null)
  const titleId = useId()
  useEffect(() => {
    const previous = document.activeElement as HTMLElement | null
    const dialog = ref.current!
    const overflow = document.body.style.overflow
    dialog.showModal()
    document.body.style.overflow = 'hidden'
    return () => {
      dialog.close()
      document.body.style.overflow = overflow
      previous?.focus()
    }
  }, [])
  return (
    <dialog
      ref={ref}
      className={`modal ${className}`}
      aria-labelledby={titleId}
      onKeyDown={(event) => {
        if (event.key !== 'Tab') return
        const controls = Array.from(
          event.currentTarget.querySelectorAll<HTMLElement>(
            'button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
          ),
        ).filter((element) => element.getClientRects().length > 0)
        const first = controls[0]
        const last = controls.at(-1)
        if (!first || !last) {
          event.preventDefault()
          return
        }
        if (
          event.shiftKey &&
          (document.activeElement === first || document.activeElement === event.currentTarget)
        ) {
          event.preventDefault()
          last.focus()
        } else if (!event.shiftKey && document.activeElement === last) {
          event.preventDefault()
          first.focus()
        }
      }}
      onCancel={(event) => {
        event.preventDefault()
        onClose()
      }}
      onClick={(event) => {
        if (event.target !== event.currentTarget) return
        const b = event.currentTarget.getBoundingClientRect()
        if (
          event.clientX < b.left ||
          event.clientX > b.right ||
          event.clientY < b.top ||
          event.clientY > b.bottom
        )
          onClose()
      }}
    >
      <div className="modal-header">
        <h2 id={titleId}>{title}</h2>
        <button type="button" className="close-button" onClick={onClose} aria-label="Закрыть">
          ×
        </button>
      </div>
      {children}
    </dialog>
  )
}

export function SelectField({
  label,
  value,
  options,
  onChange,
  required = false,
  error = false,
}: {
  label: string
  value: string
  options: string[]
  onChange: (value: string) => void
  required?: boolean
  error?: boolean
}) {
  const [open, setOpen] = useState(false)
  const [active, setActive] = useState(0)
  const id = useId()
  const ref = useRef<HTMLDivElement>(null)
  const button = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    if (open) document.getElementById(`${id}-${active}`)?.scrollIntoView({ block: 'nearest' })
  }, [open, active, id])
  useEffect(() => {
    if (!open) return
    const outside = (event: PointerEvent) => {
      if (!ref.current?.contains(event.target as Node)) setOpen(false)
    }
    document.addEventListener('pointerdown', outside)
    return () => document.removeEventListener('pointerdown', outside)
  }, [open])
  return (
    <div className={`field select-field ${error ? 'field-error' : ''}`} ref={ref}>
      <span id={`${id}-label`}>
        {label}
        {required && '*'}
      </span>
      <button
        type="button"
        ref={button}
        className={`select-button ${!value ? 'placeholder' : ''}`}
        role="combobox"
        aria-labelledby={`${id}-label`}
        aria-expanded={open}
        aria-controls={open ? `${id}-options` : undefined}
        aria-activedescendant={open ? `${id}-${active}` : undefined}
        aria-required={required}
        aria-invalid={error || undefined}
        aria-describedby={error ? `${id}-error` : undefined}
        onClick={(event) => {
          event.currentTarget.focus()
          setActive(Math.max(0, options.indexOf(value)))
          setOpen(!open)
        }}
        onKeyDown={(event) => {
          if (event.key === 'Escape') {
            event.preventDefault()
            setOpen(false)
          }
          if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
            event.preventDefault()
            setOpen(true)
            setActive((index) =>
              open
                ? (index + (event.key === 'ArrowDown' ? 1 : -1) + options.length) % options.length
                : value
                  ? Math.max(0, options.indexOf(value))
                  : event.key === 'ArrowDown'
                    ? 0
                    : options.length - 1,
            )
          }
          if (open && (event.key === 'Home' || event.key === 'End')) {
            event.preventDefault()
            setActive(event.key === 'Home' ? 0 : options.length - 1)
          }
          if ((event.key === 'Enter' || event.key === ' ') && open) {
            event.preventDefault()
            onChange(options[active])
            setOpen(false)
          }
          if (event.key === 'Tab') setOpen(false)
        }}
      >
        {value || 'Выберите нужный вариант'}
        <Icon name="select-chevron" className={open ? 'rotated' : ''} />
      </button>
      {open && (
        <div
          className="select-options"
          role="listbox"
          id={`${id}-options`}
          aria-labelledby={`${id}-label`}
        >
          {options.map((option, index) => (
            <div
              role="option"
              tabIndex={-1}
              id={`${id}-${index}`}
              aria-selected={option === value}
              className={index === active ? 'option-active' : ''}
              key={option}
              onPointerMove={() => setActive(index)}
              onMouseDown={(event) => event.preventDefault()}
              onClick={() => {
                onChange(option)
                setOpen(false)
                button.current?.focus()
              }}
            >
              {option}
            </div>
          ))}
        </div>
      )}
      {error && (
        <small className="error-text" id={`${id}-error`}>
          Выберите вариант
        </small>
      )}
    </div>
  )
}

export function Segmented({
  label,
  value,
  options,
  onChange,
  hint,
}: {
  label: string
  value: string
  options: string[]
  onChange: (value: string) => void
  hint?: ReactNode
}) {
  return (
    <fieldset className="segmented-field">
      <legend>{label}</legend>
      {hint && <p className="field-hint">{hint}</p>}
      <div className="segmented-options">
        {options.map((option) => (
          <button
            type="button"
            aria-pressed={value === option}
            className={value === option ? 'selected' : ''}
            onClick={() => onChange(option)}
            key={option}
          >
            {option}
          </button>
        ))}
      </div>
    </fieldset>
  )
}

export function MoneyInput({
  label,
  value,
  onChange,
  className = '',
  min = 0,
  max = 999_999_999,
}: {
  label: string
  value: number
  onChange: (value: number) => void
  className?: string
  min?: number
  max?: number
}) {
  const [editing, setEditing] = useState<string | null>(null)
  return (
    <label className={`field ${className}`}>
      {label}
      <input
        inputMode="numeric"
        aria-label={label.replaceAll(' ⓘ', '')}
        value={editing ?? `${money(value)} ₽`}
        onChange={(event) => {
          const digits = event.target.value.replace(/\D/g, '')
          setEditing(digits)
          if (digits) onChange(Math.min(max, Math.max(min, Number(digits))))
        }}
        onBlur={() => {
          if (editing === '') onChange(min)
          setEditing(null)
        }}
      />
    </label>
  )
}

export function RangeField({
  kind,
  value,
  onChange,
}: {
  kind: 'amount' | 'term'
  value: number
  onChange: (value: number) => void
}) {
  const amount = kind === 'amount'
  const min = amount ? 50_000 : 6
  const max = amount ? 50_000_000 : 120
  const id = useId()
  const [editing, setEditing] = useState<string | null>(null)
  const progress = `${Math.max(0, Math.min(100, ((value - min) / (max - min)) * 100))}%`
  return (
    <div className="range-field">
      {amount ? (
        <MoneyInput
          className="range-value"
          label="Сумма кредита"
          value={value}
          min={min}
          max={max}
          onChange={onChange}
        />
      ) : (
        <label className="field range-value">
          Срок кредита
          <input
            inputMode="numeric"
            aria-label="Срок кредита"
            value={editing ?? formatTerm(value)}
            onChange={(event) => {
              const input = event.target.value
              const digits = input.replace(/\D/g, '')
              setEditing(input)
              if (digits) {
                const months = Number(digits) * (/(?:год|лет)/i.test(input) ? 12 : 1)
                onChange(Math.min(max, Math.max(min, months)))
              }
            }}
            onBlur={() => {
              if (editing === '') onChange(min)
              setEditing(null)
            }}
          />
        </label>
      )}
      <label className="sr-only" htmlFor={id}>
        {amount ? 'Изменить сумму кредита' : 'Изменить срок кредита'}
      </label>
      <input
        id={id}
        className="range-slider"
        type="range"
        min={min}
        max={max}
        step={amount ? 50_000 : 1}
        value={value}
        onChange={(event) => onChange(Number(event.target.value))}
        aria-valuetext={amount ? `${money(value)} рублей` : formatTerm(value)}
        style={{ '--range-progress': progress } as CSSProperties}
      />
      <div className="range-ends">
        <span>{amount ? '50 000' : '6 месяцев'}</span>
        <span>{amount ? '50 000 000' : '10 лет'}</span>
      </div>
    </div>
  )
}

export function SearchField({
  value,
  onChange,
  placeholder,
  label,
}: {
  value: string
  onChange: (value: string) => void
  placeholder: string
  label: string
}) {
  return (
    <div className="search-field">
      <Icon name="search" />
      <input
        type="search"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        aria-label={label}
      />
      {value && (
        <button type="button" aria-label="Очистить поиск" onClick={() => onChange('')}>
          ×
        </button>
      )}
    </div>
  )
}
