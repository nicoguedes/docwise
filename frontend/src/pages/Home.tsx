import DocumentUpload from '../components/DocumentUpload'
import DocumentList from '../components/DocumentList'
import ChatInterface from '../components/ChatInterface'
import { useDocuments } from '../hooks/useDocuments'
import { useChat } from '../hooks/useChat'

interface HomeProps {
  selectedDocId: string | null
  onSelectDoc: (id: string) => void
}

export default function Home({ selectedDocId, onSelectDoc }: HomeProps) {
  const { documents, upload, remove } = useDocuments()
  const { messages, loading, sendMessage, clearMessages } = useChat()

  const handleSelectDoc = (id: string) => {
    if (id !== selectedDocId) {
      clearMessages()
      onSelectDoc(id)
    }
  }

  const selectedDoc = documents.find(d => d.id === selectedDocId)

  return (
    <div style={{
      display: 'flex',
      flex: 1,
      overflow: 'hidden',
    }}>
      {/* Sidebar */}
      <aside style={{
        width: '300px',
        minWidth: '300px',
        borderRight: '1px solid var(--border)',
        padding: '16px',
        display: 'flex',
        flexDirection: 'column',
        gap: '16px',
        overflowY: 'auto',
      }}>
        <DocumentUpload onUpload={async (file) => { await upload(file) }} />
        <DocumentList
          documents={documents}
          selectedId={selectedDocId}
          onSelect={handleSelectDoc}
          onDelete={remove}
        />
      </aside>

      {/* Chat area */}
      <section style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        {selectedDoc && (
          <div style={{
            padding: '10px 20px',
            borderBottom: '1px solid var(--border)',
            fontSize: '14px',
            color: 'var(--text-secondary)',
          }}>
            Chatting with: <strong style={{ color: 'var(--text-primary)' }}>{selectedDoc.filename}</strong>
          </div>
        )}
        <ChatInterface
          messages={messages}
          loading={loading}
          disabled={!selectedDocId || selectedDoc?.status !== 'ready'}
          onSend={(question) => {
            if (selectedDocId) sendMessage(selectedDocId, question)
          }}
        />
      </section>
    </div>
  )
}
