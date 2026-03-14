import { useState } from 'react'
import './App.css'
import Layout from './components/Layout'
import Home from './pages/Home'

function App() {
  const [selectedDocId, setSelectedDocId] = useState<string | null>(null)

  return (
    <Layout>
      <Home
        selectedDocId={selectedDocId}
        onSelectDoc={setSelectedDocId}
      />
    </Layout>
  )
}

export default App
