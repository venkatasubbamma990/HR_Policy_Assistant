function TypingIndicator() {
  return (
    <div className="message-row message-row--assistant">
      <div className="message-bubble message-bubble--assistant message-bubble--typing">
        <span className="typing-dot" />
        <span className="typing-dot" />
        <span className="typing-dot" />
      </div>
    </div>
  )
}

export default TypingIndicator
