import { useState, useMemo } from 'react'
import { ChevronRightIcon, ChevronDownIcon, DocumentIcon, FolderIcon, FolderOpenIcon } from '@heroicons/react/20/solid'
import { ArrowDownTrayIcon, ClipboardDocumentIcon } from '@heroicons/react/24/outline'
import type { ProjectConfig } from '../../types'

interface FileNode {
  name: string
  type: 'file' | 'directory'
  path: string
  children?: FileNode[]
  content?: string
  size?: number
}

interface PreviewData {
  fileStructure: FileNode[]
  estimatedSize?: string
  fileCount?: number
}

interface FileExplorerPanelProps {
  preview?: PreviewData
  config: ProjectConfig
}

// Helper function to generate file structure from config
function generateFileStructureFromConfig(config: ProjectConfig): FileNode[] {
  const projectName = config.projectName || 'my-go-project'
  
  const baseFiles: FileNode[] = [
    { name: 'main.go', type: 'file', path: '/main.go', size: 1024 },
    { name: 'go.mod', type: 'file', path: '/go.mod', size: 256 },
    { name: 'go.sum', type: 'file', path: '/go.sum', size: 512 },
    { name: 'README.md', type: 'file', path: '/README.md', size: 2048 },
  ]

  if (config.projectType === 'web-api') {
    baseFiles.push({ name: 'Dockerfile', type: 'file', path: '/Dockerfile', size: 512 })
  }

  const cmdChildren: FileNode[] = []
  if (config.projectType === 'web-api') {
    cmdChildren.push({ name: 'server.go', type: 'file', path: '/cmd/server.go', size: 2048 })
  } else if (config.projectType === 'cli') {
    cmdChildren.push(
      { name: 'root.go', type: 'file', path: '/cmd/root.go', size: 1536 },
      { name: 'version.go', type: 'file', path: '/cmd/version.go', size: 512 }
    )
  }

  const internalChildren: FileNode[] = [
    { name: 'config.go', type: 'file', path: '/internal/config.go', size: 1024 },
  ]

  // Add project-type specific files
  if (config.projectType === 'web-api') {
    if (config.framework === 'gin') {
      internalChildren.push(
        { name: 'handler.go', type: 'file', path: '/internal/handler.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'router.go', type: 'file', path: '/internal/router.go', size: 2048 }
      )
    } else if (config.framework === 'echo') {
      internalChildren.push(
        { name: 'handlers.go', type: 'file', path: '/internal/handlers.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'server.go', type: 'file', path: '/internal/server.go', size: 2048 }
      )
    } else if (config.framework === 'fiber') {
      internalChildren.push(
        { name: 'routes.go', type: 'file', path: '/internal/routes.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'app.go', type: 'file', path: '/internal/app.go', size: 2048 }
      )
    }
  }

  // Logger-specific files
  if (config.logger && config.logger !== 'slog') {
    internalChildren.push({
      name: 'logger.go',
      type: 'file',
      path: '/internal/logger.go',
      size: 1024
    })
  }

  // Database files
  if (config.features?.database?.driver) {
    const dbChildren: FileNode[] = [
      { name: 'connection.go', type: 'file', path: '/internal/database/connection.go', size: 1536 },
      { name: 'migrations.go', type: 'file', path: '/internal/database/migrations.go', size: 2048 },
      { name: 'models.go', type: 'file', path: '/internal/database/models.go', size: 2560 }
    ]

    const driver = config.features.database.driver
    dbChildren.push({
      name: `${driver}.go`,
      type: 'file',
      path: `/internal/database/${driver}.go`,
      size: 1024
    })

    const orm = config.features.database.orm
    if (orm !== 'sqlx') {
      dbChildren.push({
        name: `${orm}_models.go`,
        type: 'file',
        path: `/internal/database/${orm}_models.go`,
        size: 1536
      })
    }

    internalChildren.push({
      name: 'database',
      type: 'directory',
      path: '/internal/database',
      children: dbChildren
    })
  }

  // Authentication files
  if (config.features?.authentication?.type) {
    const authChildren: FileNode[] = [
      { name: 'middleware.go', type: 'file', path: '/internal/auth/middleware.go', size: 2048 },
      { name: 'service.go', type: 'file', path: '/internal/auth/service.go', size: 3072 }
    ]

    const authType = config.features.authentication.type
    authChildren.push({
      name: `${authType}.go`,
      type: 'file',
      path: `/internal/auth/${authType}.go`,
      size: 1536
    })

    if (authType === 'oauth2') {
      authChildren.push({
        name: 'oauth_config.go',
        type: 'file',
        path: '/internal/auth/oauth_config.go',
        size: 1024
      })
    }

    if (authType === 'jwt') {
      authChildren.push({
        name: 'token_validator.go',
        type: 'file',
        path: '/internal/auth/token_validator.go',
        size: 1280
      })
    }

    internalChildren.push({
      name: 'auth',
      type: 'directory',
      path: '/internal/auth',
      children: authChildren
    })
  }

  // Architecture-specific adjustments
  if (config.architecture === 'clean') {
    internalChildren.push(
      {
        name: 'domain',
        type: 'directory',
        path: '/internal/domain',
        children: [
          { name: 'entities.go', type: 'file', path: '/internal/domain/entities.go', size: 2048 },
          { name: 'repositories.go', type: 'file', path: '/internal/domain/repositories.go', size: 1536 }
        ]
      },
      {
        name: 'usecase',
        type: 'directory',
        path: '/internal/usecase',
        children: [
          { name: 'interfaces.go', type: 'file', path: '/internal/usecase/interfaces.go', size: 1024 },
          { name: 'user_usecase.go', type: 'file', path: '/internal/usecase/user_usecase.go', size: 2560 }
        ]
      }
    )
  } else if (config.architecture === 'hexagonal') {
    internalChildren.push(
      {
        name: 'ports',
        type: 'directory',
        path: '/internal/ports',
        children: [
          { name: 'primary.go', type: 'file', path: '/internal/ports/primary.go', size: 1536 },
          { name: 'secondary.go', type: 'file', path: '/internal/ports/secondary.go', size: 1536 }
        ]
      },
      {
        name: 'adapters',
        type: 'directory',
        path: '/internal/adapters',
        children: [
          { name: 'http.go', type: 'file', path: '/internal/adapters/http.go', size: 2048 },
          { name: 'persistence.go', type: 'file', path: '/internal/adapters/persistence.go', size: 2048 }
        ]
      }
    )
  } else if (config.architecture === 'ddd') {
    internalChildren.push(
      {
        name: 'domain',
        type: 'directory',
        path: '/internal/domain',
        children: [
          { name: 'aggregates.go', type: 'file', path: '/internal/domain/aggregates.go', size: 3072 },
          { name: 'value_objects.go', type: 'file', path: '/internal/domain/value_objects.go', size: 2048 },
          { name: 'domain_services.go', type: 'file', path: '/internal/domain/domain_services.go', size: 2560 }
        ]
      },
      {
        name: 'application',
        type: 'directory',
        path: '/internal/application',
        children: [
          { name: 'commands.go', type: 'file', path: '/internal/application/commands.go', size: 2048 },
          { name: 'queries.go', type: 'file', path: '/internal/application/queries.go', size: 2048 }
        ]
      }
    )
  }

  // Test files
  const testsChildren: FileNode[] = [
    { name: 'main_test.go', type: 'file', path: '/tests/main_test.go', size: 1024 },
    { name: 'integration_test.go', type: 'file', path: '/tests/integration_test.go', size: 2048 }
  ]

  if (config.features?.database?.driver) {
    testsChildren.push({
      name: 'database_test.go',
      type: 'file',
      path: '/tests/database_test.go',
      size: 1536
    })
  }

  if (config.features?.authentication?.type) {
    testsChildren.push({
      name: 'auth_test.go',
      type: 'file',
      path: '/tests/auth_test.go',
      size: 2048
    })
  }

  return [
    {
      name: projectName,
      type: 'directory',
      path: '/',
      children: [
        ...baseFiles,
        {
          name: 'cmd',
          type: 'directory',
          path: '/cmd',
          children: cmdChildren
        },
        {
          name: 'internal',
          type: 'directory',
          path: '/internal',
          children: internalChildren
        },
        {
          name: 'tests',
          type: 'directory',
          path: '/tests',
          children: testsChildren
        }
      ]
    }
  ]
}

export default function FileExplorerPanel({ preview, config }: FileExplorerPanelProps) {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set(['/', '/cmd', '/internal']))
  const [selectedFile, setSelectedFile] = useState<string | null>(null)

  // Generate dynamic file structure based on config
  const fileStructure: FileNode[] = useMemo(() => {
    // Use preview data if available, otherwise generate from config
    if (preview?.fileStructure) {
      return preview.fileStructure
    }
    
    // Fallback: Generate file structure from config
    return generateFileStructureFromConfig(config)
  }, [preview?.fileStructure, config])

  const toggleFolder = (path: string) => {
    setExpandedFolders(prev => {
      const newSet = new Set(prev)
      if (newSet.has(path)) {
        newSet.delete(path)
      } else {
        newSet.add(path)
      }
      return newSet
    })
  }

  const selectFile = (path: string) => {
    setSelectedFile(path)
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(() => {
      // Could add a toast notification here
    }).catch(err => {
      console.error('Failed to copy to clipboard:', err)
    })
  }

  const renderFileTree = (nodes: FileNode[], depth = 0): React.ReactElement[] => {
    return nodes.map((node) => (
      <div key={node.path}>
        {node.type === 'directory' ? (
          <div>
            <div
              className={`flex items-center gap-1 py-1 px-2 text-sm cursor-pointer hover:bg-gray-100 rounded transition-colors`}
              style={{ paddingLeft: `${depth * 16 + 8}px` }}
              onClick={() => toggleFolder(node.path)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault()
                  toggleFolder(node.path)
                }
              }}
              tabIndex={0}
              role="button"
              aria-expanded={expandedFolders.has(node.path)}
              aria-controls={`folder-content-${node.path.replace(/[^a-zA-Z0-9]/g, '-')}`}
              aria-label={`${expandedFolders.has(node.path) ? 'Collapse' : 'Expand'} folder ${node.name}`}
            >
              {expandedFolders.has(node.path) ? (
                <ChevronDownIcon className="h-4 w-4 text-gray-500" />
              ) : (
                <ChevronRightIcon className="h-4 w-4 text-gray-500" />
              )}
              {expandedFolders.has(node.path) ? (
                <FolderOpenIcon className="h-4 w-4 text-blue-500" />
              ) : (
                <FolderIcon className="h-4 w-4 text-blue-500" />
              )}
              <span className="text-gray-900 font-medium">{node.name}</span>
            </div>
            {expandedFolders.has(node.path) && node.children && (
              <div 
                id={`folder-content-${node.path.replace(/[^a-zA-Z0-9]/g, '-')}`}
                role="group"
                aria-label={`Contents of ${node.name} folder`}
              >
                {renderFileTree(node.children, depth + 1)}
              </div>
            )}
          </div>
        ) : (
          <div
            className={`flex items-center gap-1 py-1 px-2 text-sm cursor-pointer rounded transition-colors ${
              selectedFile === node.path 
                ? 'bg-blue-100 text-blue-900' 
                : 'hover:bg-gray-100 text-gray-700'
            }`}
            style={{ paddingLeft: `${depth * 16 + 24}px` }}
            onClick={() => selectFile(node.path)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                selectFile(node.path)
              }
            }}
            tabIndex={0}
            role="button"
            aria-pressed={selectedFile === node.path}
            aria-label={`Select file ${node.name}${node.size ? ` (${(node.size / 1024).toFixed(1)}KB)` : ''}`}
          >
            <DocumentIcon className="h-4 w-4 text-gray-500" />
            <span>{node.name}</span>
            {node.size && (
              <span className="text-xs text-gray-500 ml-auto">
                {(node.size / 1024).toFixed(1)}KB
              </span>
            )}
          </div>
        )}
      </div>
    ))
  }

  const selectedFileNode = useMemo(() => {
    const findFile = (nodes: FileNode[]): FileNode | null => {
      for (const node of nodes) {
        if (node.path === selectedFile) return node
        if (node.children) {
          const found = findFile(node.children)
          if (found) return found
        }
      }
      return null
    }
    return findFile(fileStructure)
  }, [fileStructure, selectedFile])

  return (
    <div className="h-full flex flex-col">
      <div className="flex-1 overflow-hidden flex">
        {/* File Tree */}
        <nav 
          className="w-1/2 border-r overflow-auto"
          aria-label="Project file tree"
          role="tree"
        >
          <div className="p-4">
            {renderFileTree(fileStructure)}
          </div>
        </nav>

        {/* File Preview */}
        <main 
          className="w-1/2 flex flex-col"
          aria-label="File content preview"
        >
          {selectedFileNode ? (
            <div className="h-full flex flex-col">
              <div className="p-4 border-b bg-gray-50">
                <div className="flex items-center justify-between">
                  <h3 className="font-medium text-gray-900">{selectedFileNode.name}</h3>
                  <button
                    onClick={() => copyToClipboard(selectedFileNode.content || '')}
                    className="p-2 text-gray-500 hover:text-gray-600 rounded"
                    title="Copy to clipboard"
                  >
                    <ClipboardDocumentIcon className="h-4 w-4" />
                  </button>
                </div>
              </div>
              <div className="flex-1 p-4 overflow-auto">
                <pre className="text-sm text-gray-700 whitespace-pre-wrap">
                  {selectedFileNode.content || '// File content will be shown here'}
                </pre>
              </div>
            </div>
          ) : (
            <div className="flex-1 flex items-center justify-center text-gray-500">
              Select a file to preview its content
            </div>
          )}
        </main>
      </div>

      {/* File Count Summary */}
      <footer 
        className="border-t p-4 bg-gray-50"
        role="status"
        aria-label="Project statistics"
      >
        <div className="text-sm text-gray-600">
          <span className="font-medium">
            {fileStructure.reduce((count, node) => {
              const countFiles = (n: FileNode): number => {
                if (n.type === 'file') return 1
                return (n.children || []).reduce((sum, child) => sum + countFiles(child), 0)
              }
              return count + countFiles(node)
            }, 0)} files
          </span>
          {preview?.estimatedSize && (
            <span className="ml-4">
              Estimated size: <span className="font-medium">{preview.estimatedSize}</span>
            </span>
          )}
        </div>
      </footer>
    </div>
  )
}