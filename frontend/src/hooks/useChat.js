import { useCallback, useRef, useState } from 'react'
import { api } from '../services/api'

let messageId = 0

function createId() {
  messageId += 1
  return `msg-${messageId}`
}

function createUserMessage(content) {
  return {
    id: createId(),
    role: 'user',
    content,
    createdAt: Date.now(),
  }
}

function createAssistantMessage(response) {
  return {
    id: createId(),
    role: 'assistant',
    content: response.answer?.text ?? 'No answer returned.',
    citations: response.answer?.citations ?? [],
    meta: {
      intent: response.intent,
      chatModel: response.answer?.chat_model,
      chunksRetrieved: response.retrieval?.chunks?.length ?? 0,
    },
    createdAt: Date.now(),
  }
}

function createStoppedMessage() {
  return {
    id: createId(),
    role: 'assistant',
    content: 'Response stopped.',
    isStopped: true,
    createdAt: Date.now(),
  }
}

function createErrorMessage(error) {
  return {
    id: createId(),
    role: 'assistant',
    content: error.message || 'Something went wrong. Please try again.',
    isError: true,
    createdAt: Date.now(),
  }
}

function isAbortError(error) {
  return error?.name === 'AbortError'
}

export function useChat() {
  const [messages, setMessages] = useState([])
  const [isLoading, setIsLoading] = useState(false)
  const abortRef = useRef(null)

  const stopGeneration = useCallback(() => {
    abortRef.current?.abort()
  }, [])

  const sendMessage = useCallback(async (rawQuestion) => {
    const question = rawQuestion.trim()
    if (!question || isLoading) {
      return
    }

    abortRef.current?.abort()

    const controller = new AbortController()
    abortRef.current = controller

    const userMessage = createUserMessage(question)
    setMessages((current) => [...current, userMessage])
    setIsLoading(true)

    try {
      const response = await api.query(question, { signal: controller.signal })
      const assistantMessage = createAssistantMessage(response)
      setMessages((current) => [...current, assistantMessage])
    } catch (error) {
      if (isAbortError(error)) {
        setMessages((current) => [...current, createStoppedMessage()])
        return
      }
      const errorMessage = createErrorMessage(error)
      setMessages((current) => [...current, errorMessage])
    } finally {
      if (abortRef.current === controller) {
        abortRef.current = null
      }
      setIsLoading(false)
    }
  }, [isLoading])

  const clearChat = useCallback(() => {
    abortRef.current?.abort()
    abortRef.current = null
    setIsLoading(false)
    setMessages([])
  }, [])

  return {
    messages,
    isLoading,
    sendMessage,
    stopGeneration,
    clearChat,
  }
}
