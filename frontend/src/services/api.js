const baseURL = import.meta.env.VITE_API_BASE_URL ?? ''

async function request(path, options = {}) {
  const response = await fetch(`${baseURL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed: ${response.status}`)
  }

  return response.json()
}

export const api = {
  health: () => request('/health'),
  query: (question) =>
    request('/api/query', {
      method: 'POST',
      body: JSON.stringify({ question }),
    }),
}
