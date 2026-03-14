# Frontend

React 18 + TypeScript + Vite. No UI library — plain CSS with CSS custom properties for theming.

## Structure

- `src/api/client.ts` — fetch-based API client (no axios). `askQuestion()` returns a `ReadableStream` for SSE parsing
- `src/types/index.ts` — TypeScript interfaces (Document, ChatMessage)
- `src/hooks/useDocuments.ts` — document list state, upload, delete, polls every 3s for status updates
- `src/hooks/useChat.ts` — chat state, SSE stream parsing, incremental message rendering
- `src/components/` — UI components:
  - `Layout.tsx` — header + two-column layout
  - `DocumentUpload.tsx` — drag-and-drop PDF upload zone
  - `DocumentList.tsx` — sidebar list with status badges
  - `ChatInterface.tsx` — message area + input form
  - `MessageBubble.tsx` — individual chat message (user right-aligned, assistant left-aligned)
- `src/pages/Home.tsx` — composes all components, manages selected document state

## Dev Server

Vite proxies `/api` requests to `http://localhost:8080` (Go backend).

## Styling

Dark theme using CSS custom properties defined in `index.css`. Key variables:
- `--bg-primary`, `--bg-secondary`, `--bg-tertiary` — background layers
- `--accent` — indigo (#6366f1) for interactive elements
- `--text-primary`, `--text-secondary` — text colors

## Building

```bash
npm run build  # outputs to dist/
```

Docker: node:20-alpine build → nginx:alpine serving static files with `/api` reverse proxy to backend.
