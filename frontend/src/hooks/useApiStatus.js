import { useCallback, useEffect, useState } from 'react'
import { api } from '../services/api'

export function useApiStatus() {
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
  }, [check])

  return { status, refresh: check }
}
