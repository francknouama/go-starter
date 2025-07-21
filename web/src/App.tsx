import { useState } from 'react'
import Header from './components/layout/Header'
import ConfigurationPanel from './components/forms/ConfigurationPanel'
import PreviewPanel from './components/preview/PreviewPanel'
import FileExplorerPanel from './components/preview/FileExplorerPanel'

function App() {
  const [disclosureMode, setDisclosureMode] = useState<'basic' | 'advanced'>('basic')

  return (
    <div className="min-h-screen bg-gray-50">
      <Header 
        disclosureMode={disclosureMode}
        onDisclosureModeChange={setDisclosureMode}
      />
      
      <main className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-[calc(100vh-120px)]">
          {/* Configuration Panel */}
          <div className="lg:col-span-1">
            <ConfigurationPanel disclosureMode={disclosureMode} />
          </div>
          
          {/* Preview Panel */}
          <div className="lg:col-span-1">
            <PreviewPanel />
          </div>
          
          {/* File Explorer Panel */}
          <div className="lg:col-span-1">
            <FileExplorerPanel />
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
