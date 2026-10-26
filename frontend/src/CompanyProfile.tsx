import { FinancialFields } from './FinancialFields'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { ApplicantFields } from './Applications'
import { saveCompany } from './api'
import type { ApiCompany } from './api'
import { defaultApplicant } from './types'
import type { Selection } from './types'

export function CompanyProfile({
  company,
  selection,
  onSaved,
}: {
  company?: ApiCompany
  selection: Selection
  onSaved: (company: ApiCompany) => void
}) {
  const [applicant, setApplicant] = useState(
    company
      ? {
          company: company.name,
          inn: company.inn,
          registered: company.registeredAt || '',
          name: company.representativeName || '',
          position: company.representativePosition || '',
          phone: company.phone || '',
          email: company.email || '',
        }
      : defaultApplicant,
  )
  const [finance, setFinance] = useState({
    revenue: (company?.monthlyRevenueMinor || 0) / 100,
    expenses: (company?.monthlyExpensesMinor || 0) / 100,
    debts: (company?.monthlyDebtPaymentsMinor || 0) / 100,
    overdue: company?.hasOverdueDebt || false,
  })
  const [pending, setPending] = useState(false)
  const [error, setError] = useState('')

  async function submit(event: FormEvent) {
    event.preventDefault()
    if (pending) return
    setPending(true)
    setError('')
    try {
      const value = await saveCompany(
        applicant,
        selection,
        finance,
        company?.inn === applicant.inn ? company.id : undefined,
      )
      onSaved(value)
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Не удалось сохранить компанию')
    } finally {
      setPending(false)
    }
  }

  return (
    <main className="page-container">
      <h1>Данные компании</h1>
      <p>Укажите сведения о компании для персонального подбора и заполнения заявки.</p>
      <form className="company-profile-form" onSubmit={submit}>
        <fieldset disabled={pending} style={{ border: 0, padding: 0 }}>
          <div className="detail-applicant-grid">
            <ApplicantFields value={applicant} onChange={setApplicant} />
          </div>
          <FinancialFields value={finance} onChange={setFinance} />
          {error && (
            <p role="alert" className="error-text">
              {error}
            </p>
          )}
          <button className="button-primary" type="submit">
            {pending ? 'Сохраняем…' : 'Сохранить компанию'}
          </button>
        </fieldset>
      </form>
    </main>
  )
}
