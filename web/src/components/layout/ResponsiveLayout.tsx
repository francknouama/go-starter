import { useState } from 'react'
import type { ReactNode } from 'react'
import MobileTabNavigation from './MobileTabNavigation'

interface ResponsiveLayoutProps {
  configurationPanel: ReactNode
  fileExplorerPanel: ReactNode
}

export type TabId = 'configuration' | 'files'

export default function ResponsiveLayout({
  configurationPanel,
  fileExplorerPanel
}: ResponsiveLayoutProps) {
  const [activeTab, setActiveTab] = useState<TabId>('configuration')

  const renderMobileContent = () => {
    switch (activeTab) {
      case 'configuration':
        return configurationPanel
      case 'files':
        return fileExplorerPanel
      default:
        return configurationPanel
    }
  }

  return (
    <>
      {/* Mobile Layout - Single panel with tab navigation */}
      <div className="lg:hidden">
        <main 
          role="main" 
          aria-label="Project configuration interface"
          className="min-h-[calc(100vh-64px-64px)] pb-20"
        >
          <section 
            aria-label={activeTab === 'configuration' ? 'Project configuration form' : 'Generated files preview'}
            className="p-4"
          >
            {renderMobileContent()}
          </section>
        </main>
        
        {/* Mobile Tab Navigation */}
        <nav aria-label="Main navigation">
          <MobileTabNavigation 
            activeTab={activeTab}
            onTabChange={setActiveTab}
          />
        </nav>
      </div>

      {/* Tablet Layout - Stacked panels */}
      <main className="hidden sm:block lg:hidden" role="main" aria-label="Project configuration interface">
        <div className="space-y-6 p-6">
          <section 
            aria-labelledby="config-section-heading"
            className="bg-white/95 backdrop-blur-sm rounded-xl shadow-sm border border-gray-100 p-6"
          >
            <h2 id="config-section-heading" className="sr-only">Project Configuration Section</h2>
            {configurationPanel}
          </section>
          <section 
            aria-labelledby="files-section-heading"
            className="bg-white/95 backdrop-blur-sm rounded-xl shadow-sm border border-gray-100 p-6"
          >
            <h2 id="files-section-heading" className="sr-only">File Preview Section</h2>
            {fileExplorerPanel}
          </section>
        </div>
      </main>

      {/* Desktop Layout - Side by side */}
      <main className="hidden lg:block" role="main" aria-label="Project configuration interface">
        <div className="grid grid-cols-2 gap-6 max-w-7xl mx-auto">
          <section 
            aria-labelledby="config-section-heading-desktop"
            className="bg-white/95 backdrop-blur-sm rounded-xl shadow-sm border border-gray-100 p-6"
          >
            <h2 id="config-section-heading-desktop" className="sr-only">Project Configuration Section</h2>
            {configurationPanel}
          </section>
          <section 
            aria-labelledby="files-section-heading-desktop"
            className="bg-white/95 backdrop-blur-sm rounded-xl shadow-sm border border-gray-100 p-6"
          >
            <h2 id="files-section-heading-desktop" className="sr-only">File Preview Section</h2>
            {fileExplorerPanel}
          </section>
        </div>
      </main>
    </>
  )
}