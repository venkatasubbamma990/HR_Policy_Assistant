function ChatHeader({ category, isLoading, onClear, hasMessages }) {
  return (
    <header className="chat-header">
      <div>
        <h1 className="chat-header__title">{category.label}</h1>
        <p className="chat-header__subtitle">Local RAG · NovaTech Solutions</p>
      </div>
      <div className="chat-header__actions">
        {isLoading ? <span className="chat-header__badge">Thinking…</span> : null}
        {hasMessages ? (
          <button type="button" className="chat-header__clear" onClick={onClear}>
            Clear chat
          </button>
        ) : null}
      </div>
    </header>
  )
}

export default ChatHeader
