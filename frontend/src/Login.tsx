import { useState } from 'react'
import type { FormEvent } from 'react'
import { authenticate } from './api'

export function Login({ onLogin }: { onLogin: (specialist?: boolean) => void }) {
  const [register, setRegister] = useState(false)
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [pending, setPending] = useState(false)

  async function submit(event: FormEvent) {
    event.preventDefault()
    if (pending) return
    setError('')
    setPending(true)
    try {
      const session = await authenticate(login, password, register ? name : undefined)
      onLogin(session.role === 'EMPLOYEE')
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Не удалось войти')
    } finally {
      setPending(false)
    }
  }

  return (
    <main className="login-page">
      <form className="login-card login-login" onSubmit={submit}>
        <h1>{register ? 'Регистрация' : 'Вход'}</h1>
        <h2>{register ? 'Создайте аккаунт' : 'С возвращением!'}</h2>
        <p className="login-subtitle">Кредитный кабинет предпринимателя</p>
        {register && (
          <label className="field login-field">
            Имя
            <input
              value={name}
              onChange={(event) => setName(event.target.value)}
              minLength={2}
              maxLength={100}
              required
              autoComplete="name"
            />
          </label>
        )}
        <label className="field login-field">
          Email
          <input
            type="email"
            value={login}
            onChange={(event) => setLogin(event.target.value)}
            required
            autoComplete="username"
          />
        </label>
        <label className="field login-field">
          Пароль
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
            minLength={register ? 12 : undefined}
            autoComplete={register ? 'new-password' : 'current-password'}
          />
        </label>
        {error && (
          <p className="login-error" role="alert">
            {error}
          </p>
        )}
        <button className="button-primary login-submit" type="submit" disabled={pending}>
          {pending ? 'Подождите…' : register ? 'Зарегистрироваться' : 'Войти'}
        </button>
        <button
          className="button-secondary register-button"
          type="button"
          disabled={pending}
          onClick={() => {
            setRegister(!register)
            setError('')
          }}
        >
          {register ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
        </button>
      </form>
    </main>
  )
}
