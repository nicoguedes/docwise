import { Document } from '../types'

interface DocumentListProps {
  documents: Document[]
  selectedId: string | null
  onSelect: (id: string) => void
  onDelete: (id: string) => void
}

const statusColors: Record<string, string> = {
  pending: 'var(--warning)',
  processing: 'var(--warning)',
  ready: 'var(--success)',
  error: 'var(--error)',
}

export default function DocumentList({ documents, selectedId, onSelect, onDelete }: DocumentListProps) {
  if (documents.length === 0) {
    return (
      <p style={{ color: 'var(--text-secondary)', fontSize: '14px', padding: '12px 0' }}>
        No documents yet. Upload a PDF to get started.
      </p>
    )
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
      {documents.map(doc => (
        <div
          key={doc.id}
          onClick={() => doc.status === 'ready' && onSelect(doc.id)}
          style={{
            padding: '10px 12px',
            borderRadius: 'var(--radius)',
            backgroundColor: selectedId === doc.id ? 'var(--bg-tertiary)' : 'transparent',
            cursor: doc.status === 'ready' ? 'pointer' : 'default',
            opacity: doc.status === 'ready' ? 1 : 0.6,
            transition: 'background-color 0.2s',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <div style={{ minWidth: 0, flex: 1 }}>
            <p style={{
              fontSize: '14px',
              fontWeight: 500,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            }}>
              {doc.filename}
            </p>
            <p style={{ fontSize: '12px', color: 'var(--text-secondary)', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <span style={{
                width: '6px',
                height: '6px',
                borderRadius: '50%',
                backgroundColor: statusColors[doc.status],
                display: 'inline-block',
              }} />
              {doc.status}
              {doc.page_count > 0 && ` · ${doc.page_count} pages`}
            </p>
          </div>
          <button
            onClick={e => {
              e.stopPropagation()
              onDelete(doc.id)
            }}
            style={{
              background: 'none',
              color: 'var(--text-secondary)',
              fontSize: '18px',
              padding: '0 4px',
              lineHeight: 1,
            }}
            title="Delete document"
          >
            ×
          </button>
        </div>
      ))}
    </div>
  )
}
