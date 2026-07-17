import { useCallback, useEffect, useState } from 'react'
import { api } from '../services/api'

export function useApiStatus(pollMs = 30000) {
  const [status, setStatus] = useState('checking')

  const check = useCallback(async () => {
    try {
      await api.health()
      setStatus('ready')
    } catch {
      setStatus('offline')
    }
  }, [])

  useEffect(() => {
    check()
    const timer = setInterval(check, pollMs)
    return () => clearInterval(timer)
  }, [check, pollMs])

  return { status, refresh: check }
}
