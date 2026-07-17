import Layout from '../components/layout/Layout'

function HomePage() {
  return (
    <Layout>
      <section className="hero-card">
        <h2>Welcome</h2>
        <p>
          This is the frontend boilerplate for the HR Policy Assistant. The Go API
          runs separately on port <code>8081</code>.
        </p>
        <p className="hero-card__hint">
          Tell me what UI and features you want next — chat, citations, sample
          questions, and more.
        </p>
      </section>
    </Layout>
  )
}

export default HomePage
