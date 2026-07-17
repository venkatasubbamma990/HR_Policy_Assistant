import { useState } from 'react'
import { PolicyIcon } from '../icons/PolicyIcons'

function ChatInput({ onSend, onStop, isLoading, placeholder }) {
  const [value, setValue] = useState('')

  const handleSubmit = (event) => {
    event.preventDefault()
    const question = value.trim()
    if (!question || isLoading) {
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
      <div className="chat-input__wrap">
        <textarea
          id="chat-question"
          className="chat-input__field"
          rows={1}
          value={value}
          placeholder={placeholder}
          onChange={(event) => setValue(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === 'Enter' && !event.shiftKey) {
              event.preventDefault()
              handleSubmit(event)
            }
          }}
        />
        {isLoading ? (
          <button
            type="button"
            className="chat-input__stop"
            onClick={onStop}
            aria-label="Stop generation"
          >
            <PolicyIcon name="stop" />
          </button>
        ) : (
          <button
            type="submit"
            className="chat-input__submit"
            disabled={!value.trim()}
            aria-label="Send message"
          >
            <PolicyIcon name="send" />
          </button>
        )}
      </div>
    </form>
  )
}

export default ChatInput
