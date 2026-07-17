function SuggestedQuestions({ questions, onSelect, disabled }) {
  if (!questions.length) {
    return null
  }

  return (
    <div className="suggested-questions">
      <p className="suggested-questions__label">Try asking:</p>
      <div className="suggested-questions__list">
        {questions.map((question) => (
          <button
            key={question}
            type="button"
            className="suggested-questions__chip"
            onClick={() => onSelect(question)}
            disabled={disabled}
          >
            {question}
          </button>
        ))}
      </div>
    </div>
  )
}

export default SuggestedQuestions
