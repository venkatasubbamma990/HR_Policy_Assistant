import { SAMPLE_QUESTIONS } from '../../constants/sampleQuestions'
import ChatInput from './ChatInput'
import EmptyState from './EmptyState'
import MessageList from './MessageList'
import SuggestedQuestions from './SuggestedQuestions'

function ChatPanel({ messages, isLoading, onSend, onClear }) {
  const hasMessages = messages.length > 0

  return (
    <section className="chat-panel" aria-label="HR policy chat">
      <div className="chat-panel__toolbar">
        <p className="chat-panel__status">
          {isLoading ? 'Retrieving policy answer…' : 'Ready'}
        </p>
        {hasMessages ? (
          <button type="button" className="chat-panel__clear" onClick={onClear}>
            Clear chat
          </button>
        ) : null}
      </div>

      <div className="chat-panel__body">
        {!hasMessages ? <EmptyState /> : null}
        {hasMessages ? <MessageList messages={messages} isLoading={isLoading} /> : null}
      </div>

      {!hasMessages ? (
        <SuggestedQuestions
          questions={SAMPLE_QUESTIONS}
          onSelect={onSend}
          disabled={isLoading}
        />
      ) : null}

      <ChatInput
        onSend={onSend}
        disabled={isLoading}
        placeholder="Ask about leave, notice period, WFH, salary, onboarding, or exit…"
      />
    </section>
  )
}

export default ChatPanel
