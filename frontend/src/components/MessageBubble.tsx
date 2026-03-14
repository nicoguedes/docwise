import { ChatMessage } from '../types'

interface MessageBubbleProps {
  message: ChatMessage
}

export default function MessageBubble({ message }: MessageBubbleProps) {
  const isUser = message.role === 'user'

  return (
    <div style={{
      display: 'flex',
      justifyContent: isUser ? 'flex-end' : 'flex-start',
      marginBottom: '12px',
    }}>
      <div style={{
        maxWidth: '75%',
        padding: '10px 14px',
        borderRadius: '12px',
        backgroundColor: isUser ? 'var(--accent)' : 'var(--bg-tertiary)',
        color: 'var(--text-primary)',
        fontSize: '14px',
        lineHeight: 1.5,
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word',
      }}>
        {message.content || (
          <span style={{ color: 'var(--text-secondary)' }}>Thinking...</span>
        )}
      </div>
    </div>
  )
}
