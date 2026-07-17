import CitationList from './CitationList'
import { PolicyIcon } from '../icons/PolicyIcons'

function MessageBubble({ message }) {
  const isUser = message.role === 'user'

  return (
    <div className={`message-row message-row--${message.role}`}>
      {!isUser ? (
        <span className="message-row__avatar message-row__avatar--bot">
          <PolicyIcon name="spark" />
        </span>
      ) : null}

      <div
        className={[
          'message-bubble',
          `message-bubble--${message.role}`,
          message.isError ? 'message-bubble--error' : '',
        ]
          .filter(Boolean)
          .join(' ')}
      >
        <p className="message-bubble__text">{message.content}</p>

        {!isUser && !message.isError ? (
          <>
            {message.meta?.chatModel ? (
              <p className="message-bubble__meta">
                {message.meta.chatModel}
                {message.meta.intent ? ` · ${message.meta.intent}` : ''}
              </p>
            ) : null}
            <CitationList citations={message.citations} />
          </>
        ) : null}
      </div>
    </div>
  )
}

export default MessageBubble
