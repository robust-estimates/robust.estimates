import { useEffect, useState } from 'react'

export function useSessionValue<T>(key: string, initial: T, valid: (value: unknown) => value is T) {
  const [value, setValue] = useState<T>(() => {
    try {
      const stored: unknown = JSON.parse(sessionStorage.getItem(key) ?? 'null')
      if (valid(stored)) return stored
    } catch {}
    return initial
  })
  useEffect(() => {
    try {
      sessionStorage.setItem(key, JSON.stringify(value))
    } catch {}
  }, [key, value])
  return [value, setValue] as const
}

export function clearSessionValue(key: string) {
  try {
    sessionStorage.removeItem(key)
  } catch {}
}
