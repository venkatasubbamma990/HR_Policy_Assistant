import { useState } from 'react'

function ChatInput({ onSend, disabled, placeholder }) {
  const [value, setValue] = useState('')

  const handleSubmit = (event) => {
    event.preventDefault()
    const question = value.trim()
    if (!question || disabled) {
      return
    }

    onSend(question)
    setValue('')
  }

  return (
    <form className="chat-input" onSubmit={handleSubmit}>
      <label htmlFor="chat-question" className="sr-only">
        Ask a question
      </label>
      <textarea
        id="chat-question"
        className="chat-input__field"
        rows={2}
        value={value}
        placeholder={placeholder}
        onChange={(event) => setValue(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault()
            handleSubmit(event)
          }
        }}
        disabled={disabled}
      />
      <button type="submit" className="chat-input__submit" disabled={disabled || !value.trim()}>
        Send
      </button>
    </form>
  )
}

export default ChatInput
