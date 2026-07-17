import Header from './Header'

function Layout({ children }) {
  return (
    <div className="app-shell">
      <Header />
      <main className="app-main">{children}</main>
    </div>
  )
}

export default Layout
