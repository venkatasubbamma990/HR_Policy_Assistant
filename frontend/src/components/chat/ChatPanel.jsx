import ChatHeader from './ChatHeader'
import ChatInput from './ChatInput'
import EmptyState from './EmptyState'
import MessageList from './MessageList'
import SuggestedQuestions from './SuggestedQuestions'

function ChatPanel({ category, messages, isLoading, onSend, onClear }) {
  const hasMessages = messages.length > 0

  return (
    <section className="chat-panel" aria-label="HR policy chat">
      <ChatHeader
        category={category}
        isLoading={isLoading}
        onClear={onClear}
        hasMessages={hasMessages}
      />

      <div className="chat-panel__body">
        {!hasMessages ? <EmptyState category={category} /> : null}
        {hasMessages ? <MessageList messages={messages} isLoading={isLoading} /> : null}
      </div>

      {!hasMessages ? (
        <SuggestedQuestions
          questions={category.questions}
          onSelect={onSend}
          disabled={isLoading}
        />
      ) : null}

      <ChatInput
        onSend={onSend}
        disabled={isLoading}
        placeholder={`Ask me anything about ${category.label}…`}
      />
    </section>
  )
}

export default ChatPanel
