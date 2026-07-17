import { useMemo, useState } from 'react'
import ChatPanel from '../components/chat/ChatPanel'
import Layout from '../components/layout/Layout'
import { DEFAULT_CATEGORY, POLICY_CATEGORIES } from '../constants/policyCategories'
import { useApiStatus } from '../hooks/useApiStatus'
import { useChat } from '../hooks/useChat'

function ChatPage() {
  const [activeCategoryId, setActiveCategoryId] = useState(DEFAULT_CATEGORY.id)
  const { messages, isLoading, sendMessage, stopGeneration, clearChat } = useChat()
  const { status: apiStatus } = useApiStatus()

  const activeCategory = useMemo(
    () => POLICY_CATEGORIES.find((c) => c.id === activeCategoryId) ?? DEFAULT_CATEGORY,
    [activeCategoryId],
  )

  return (
    <Layout
      categories={POLICY_CATEGORIES}
      activeCategoryId={activeCategoryId}
      onSelectCategory={setActiveCategoryId}
      apiStatus={apiStatus}
    >
      <ChatPanel
        category={activeCategory}
        messages={messages}
        isLoading={isLoading}
        onSend={sendMessage}
        onStop={stopGeneration}
        onClear={clearChat}
      />
    </Layout>
  )
}

export default ChatPage
