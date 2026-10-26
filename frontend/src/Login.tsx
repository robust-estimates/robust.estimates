import { useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { Modal } from './ui'

export function Login({ onLogin }: { onLogin: (specialist?: boolean) => void }) {
  const [mode, setMode] = useState<'login' | 'phone' | 'code'>('login')
  const [login, setLogin] = useState('ivanov123')
  const [password, setPassword] = useState('demo2026')
  const [phone, setPhone] = useState('+7 900 123-45-67')
  const [code, setCode] = useState(['', '', '', ''])
  const [error, setError] = useState('')
  const [information, setInformation] = useState('')
  const digits = useRef<(HTMLInputElement | null)[]>([])
  function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    if (mode === 'phone') {
      if (!/^[78]\d{10}$/.test(phone.replace(/\D/g, ''))) {
        setError('Введите корректный номер телефона')
        return
      }
      setMode('code')
      setCode(['', '', '', ''])
      requestAnimationFrame(() => digits.current[0]?.focus())
      return
    }
    if (mode === 'code') {
      if (code.join('') !== '4267') {
        setError('Неверный код подтверждения')
        return
      }
      onLogin()
      return
    }
    if (!['ivanov123', 'specialist'].includes(login) || password !== 'demo2026') {
      setError('Неправильный логин или пароль')
      return
    }
    onLogin(login === 'specialist')
  }
  return (
    <main className="login-page">
      <form className={`login-card login-${mode}`} onSubmit={submit}>
        <h1>Вход</h1>
        <h2>{mode === 'code' ? 'Подтвердите телефон' : 'С возвращением!'}</h2>
        <p className="login-subtitle">
          {mode === 'code'
            ? `На номер ${phone} отправлен код подтверждения`
            : mode === 'phone'
              ? 'Введите телефон, чтобы продолжить работу'
              : 'Войдите, чтобы продолжить работу'}
        </p>
        {mode === 'login' ? (
          <>
            <label className="field login-field">
              Логин
              <input
                autoComplete="username"
                value={login}
                onChange={(event) => setLogin(event.target.value)}
                required
              />
            </label>
            {error && (
              <p className="login-error" role="alert">
                {error}
              </p>
            )}
            <div className="field login-field password-field">
              <span>
                <label htmlFor="login-password">Пароль</label>
                <button
                  type="button"
                  className="text-button"
                  onClick={(event) => {
                    event.currentTarget.focus()
                    setInformation(
                      'Восстановление пароля будет доступно после подключения авторизации. Для демонстрации используйте логин ivanov123 и пароль demo2026.',
                    )
                  }}
                >
                  Забыли пароль?
                </button>
              </span>
              <input
                id="login-password"
                type="password"
                autoComplete="current-password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>
            <label className="remember-me">
              <input type="checkbox" defaultChecked />
              Запомнить меня
            </label>
          </>
        ) : mode === 'phone' ? (
          <>
            <label className="field phone-field">
              Телефон
              <input
                type="tel"
                autoComplete="tel"
                value={phone}
                onChange={(event) => setPhone(event.target.value)}
                required
              />
            </label>
            {error && (
              <p className="login-error" role="alert">
                {error}
              </p>
            )}
          </>
        ) : (
          <div className="sms-field">
            <p>Код из SMS</p>
            <div
              className="sms-inputs"
              onPaste={(event) => {
                const pasted = event.clipboardData.getData('text').replace(/\D/g, '').slice(0, 4)
                if (pasted.length === 4) {
                  event.preventDefault()
                  setCode(pasted.split(''))
                  digits.current[3]?.focus()
                }
              }}
            >
              {code.map((digit, index) => (
                <input
                  key={index}
                  ref={(element) => {
                    digits.current[index] = element
                  }}
                  inputMode="numeric"
                  autoComplete={index === 0 ? 'one-time-code' : 'off'}
                  maxLength={1}
                  aria-label={`Цифра ${index + 1} кода`}
                  value={digit}
                  required
                  onChange={(event) => {
                    const value = event.target.value.replace(/\D/g, '').slice(-1)
                    setCode(code.map((item, i) => (i === index ? value : item)))
                    if (value) digits.current[index + 1]?.focus()
                  }}
                  onKeyDown={(event) => {
                    if (event.key === 'Backspace' && !digit) digits.current[index - 1]?.focus()
                  }}
                />
              ))}
            </div>
            {error && (
              <p className="login-error" role="alert">
                {error}
              </p>
            )}
            <p className="sms-demo">Демонстрационный код: 4267</p>
          </div>
        )}
        <button className="button-primary login-submit" type="submit">
          {mode === 'code' ? 'Продолжить' : 'Войти'}
        </button>
        {mode === 'code' ? (
          <button
            className="button-secondary change-phone"
            type="button"
            onClick={() => {
              setMode('phone')
              setError('')
            }}
          >
            Изменить номер
          </button>
        ) : (
          <>
            <p className="login-divider">или войдите через</p>
            <div className="login-alternatives">
              <button
                className="button-white"
                type="button"
                onClick={(event) => {
                  event.currentTarget.focus()
                  setInformation(
                    'Вход через Госуслуги будет доступен после подключения сервиса авторизации. Сейчас можно посмотреть фронтенд через вход по логину или демонстрационному SMS-коду.',
                  )
                }}
              >
                Госуслуги
              </button>
              <button
                className="button-white"
                type="button"
                onClick={() => {
                  setMode(mode === 'phone' ? 'login' : 'phone')
                  setError('')
                }}
              >
                {mode === 'phone' ? 'Логин' : 'Телефон'}
              </button>
            </div>
            <button
              className="button-secondary register-button"
              type="button"
              onClick={(event) => {
                event.currentTarget.focus()
                setInformation(
                  'Регистрация будет доступна после подключения банковского сервиса. Демонстрационный вход уже доступен: ivanov123 / demo2026.',
                )
              }}
            >
              Нет аккаунта? &nbsp; Зарегистрироваться
            </button>
          </>
        )}
      </form>
      {information && (
        <Modal title="Вход в кабинет" onClose={() => setInformation('')}>
          <p className="demo-message">{information}</p>
          <button className="button-primary" type="button" onClick={() => setInformation('')}>
            Понятно
          </button>
        </Modal>
      )}
    </main>
  )
}
