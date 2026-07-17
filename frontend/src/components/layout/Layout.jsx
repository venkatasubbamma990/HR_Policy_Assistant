import Sidebar from './Sidebar'

function Layout({ children, categories, activeCategoryId, onSelectCategory, apiStatus }) {
  return (
    <div className="app-shell">
      <Sidebar
        categories={categories}
        activeCategoryId={activeCategoryId}
        onSelectCategory={onSelectCategory}
        apiStatus={apiStatus}
      />
      <main className="app-main">{children}</main>
    </div>
  )
}

export default Layout
