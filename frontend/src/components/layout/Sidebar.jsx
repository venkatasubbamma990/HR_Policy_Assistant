import { PolicyIcon } from '../icons/PolicyIcons'

function Sidebar({ categories, activeCategoryId, onSelectCategory, apiStatus }) {
  const statusLabel =
    apiStatus === 'ready' ? 'API: Ready' : apiStatus === 'offline' ? 'API: Offline' : 'Checking…'

  return (
    <aside className="sidebar">
      <div className="sidebar__brand">
        <span className="sidebar__brand-icon">
          <PolicyIcon name="book" />
        </span>
        <span className="sidebar__brand-text">HR Policy Assistant</span>
      </div>

      <p className="sidebar__section-label">POLICIES</p>
      <nav className="sidebar__nav" aria-label="Policy categories">
        {categories.map((category) => {
          const isActive = category.id === activeCategoryId
          return (
            <button
              key={category.id}
              type="button"
              className={`sidebar__item${isActive ? ' sidebar__item--active' : ''}`}
              onClick={() => onSelectCategory(category.id)}
            >
              <span className="sidebar__item-icon">
                <PolicyIcon name={category.icon} />
              </span>
              <span className="sidebar__item-label">{category.label}</span>
            </button>
          )
        })}
      </nav>

      <div className={`sidebar__status sidebar__status--${apiStatus}`}>
        <span className="sidebar__status-dot" />
        <span>{statusLabel}</span>
      </div>
    </aside>
  )
}

export default Sidebar
