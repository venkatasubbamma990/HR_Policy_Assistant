function CitationList({ citations }) {
  if (!citations?.length) {
    return null
  }

  return (
    <div className="citation-list">
      <p className="citation-list__title">Sources</p>
      <ul>
        {citations.map((citation) => (
          <li key={`${citation.index}-${citation.chunk_id}`}>
            <span className="citation-list__index">[{citation.index}]</span>
            <span className="citation-list__source">{citation.source}</span>
            {citation.section_path ? (
              <span className="citation-list__section"> — {citation.section_path}</span>
            ) : null}
          </li>
        ))}
      </ul>
    </div>
  )
}

export default CitationList
