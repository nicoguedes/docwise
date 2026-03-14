import { useState, useRef, useEffect, FormEvent } from 'react'
import { ChatMessage } from '../types'
import MessageBubble from './MessageBubble'

interface ChatInterfaceProps {
  messages: ChatMessage[]
  loading: boolean
  disabled: boolean
  onSend: (question: string) => void
}

export default function ChatInterface({ messages, loading, disabled, onSend }: ChatInterfaceProps) {
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const trimmed = input.trim()
    if (!trimmed || loading || disabled) return
    onSend(trimmed)
    setInput('')
  }

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100%',
    }}>
      {/* Messages area */}
      <div style={{
        flex: 1,
        overflowY: 'auto',
        padding: '20px',
      }}>
        {messages.length === 0 && (
          <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            height: '100%',
            color: 'var(--text-secondary)',
            fontSize: '14px',
            textAlign: 'center',
          }}>
            {disabled
              ? 'Select a document to start chatting'
              : 'Ask a question about your document'}
          </div>
        )}
        {messages.map((msg, i) => (
          <MessageBubble key={i} message={msg} />
        ))}
        <div ref={messagesEndRef} />
      </div>

      {/* Input area */}
      <form
        onSubmit={handleSubmit}
        style={{
          padding: '16px 20px',
          borderTop: '1px solid var(--border)',
          display: 'flex',
          gap: '8px',
        }}
      >
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder={disabled ? 'Select a document first...' : 'Ask a question...'}
          disabled={disabled || loading}
          style={{
            flex: 1,
            padding: '10px 14px',
            borderRadius: 'var(--radius)',
            border: '1px solid var(--border)',
            backgroundColor: 'var(--bg-secondary)',
            color: 'var(--text-primary)',
            outline: 'none',
          }}
        />
        <button
          type="submit"
          disabled={disabled || loading || !input.trim()}
          style={{
            padding: '10px 20px',
            borderRadius: 'var(--radius)',
            backgroundColor: disabled || loading || !input.trim() ? 'var(--bg-tertiary)' : 'var(--accent)',
            color: 'var(--text-primary)',
            fontWeight: 500,
            transition: 'background-color 0.2s',
          }}
        >
          Send
        </button>
      </form>
    </div>
  )
}
