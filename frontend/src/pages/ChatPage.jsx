import ChatPanel from '../components/chat/ChatPanel'
import Layout from '../components/layout/Layout'
import { useChat } from '../hooks/useChat'

function ChatPage() {
  const { messages, isLoading, sendMessage, clearChat } = useChat()

  return (
    <Layout>
      <ChatPanel
        messages={messages}
        isLoading={isLoading}
        onSend={sendMessage}
        onClear={clearChat}
      />
    </Layout>
  )
}

export default ChatPage
