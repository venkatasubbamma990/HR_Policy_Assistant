import { PolicyIcon } from '../icons/PolicyIcons'

function EmptyState({ category }) {
  return (
    <div className="chat-empty">
      <div className="chat-empty__icon">
        <PolicyIcon name={category.icon} />
      </div>
      <h2>Ask about {category.shortLabel}</h2>
      <p>
        Get instant answers from NovaTech HR policy documents — powered by RAG
        with cited sources from {category.label.toLowerCase()}.
      </p>
    </div>
  )
}

export default EmptyState
