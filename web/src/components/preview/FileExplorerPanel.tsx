import { useState } from 'react'
import { ChevronRightIcon, ChevronDownIcon, DocumentIcon, FolderIcon, FolderOpenIcon } from '@heroicons/react/20/solid'
import { ArrowDownTrayIcon, ClipboardDocumentIcon } from '@heroicons/react/24/outline'

interface FileNode {
  name: string
  type: 'file' | 'folder'
  path: string
  children?: FileNode[]
  content?: string
  size?: number
  language?: string
}

export default function FileExplorerPanel() {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set(['/', '/cmd', '/internal']))
  const [selectedFile, setSelectedFile] = useState<string | null>(null)

  // Mock file structure
  const fileStructure: FileNode[] = [
    {
      name: 'cmd',
      type: 'folder',
      path: '/cmd',
      children: [
        {
          name: 'server',
          type: 'folder',
          path: '/cmd/server',
          children: [
            {
              name: 'main.go',
              type: 'file',
              path: '/cmd/server/main.go',
              size: 1024,
              language: 'go',
              content: `package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })
    
    log.Println("Server starting on :8080")
    r.Run(":8080")
}`,
            },
          ],
        },
      ],
    },
    {
      name: 'internal',
      type: 'folder',
      path: '/internal',
      children: [
        {
          name: 'config',
          type: 'folder',
          path: '/internal/config',
          children: [
            {
              name: 'config.go',
              type: 'file',
              path: '/internal/config/config.go',
              size: 512,
              language: 'go',
              content: `package config

type Config struct {
    Port     string
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    Database string
}`,
            },
          ],
        },
        {
          name: 'handlers',
          type: 'folder',
          path: '/internal/handlers',
          children: [
            {
              name: 'user.go',
              type: 'file',
              path: '/internal/handlers/user.go',
              size: 768,
              language: 'go',
            },
            {
              name: 'auth.go',
              type: 'file',
              path: '/internal/handlers/auth.go',
              size: 1536,
              language: 'go',
            },
          ],
        },
      ],
    },
    {
      name: 'go.mod',
      type: 'file',
      path: '/go.mod',
      size: 256,
      language: 'go',
      content: `module github.com/user/my-awesome-project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
)`,
    },
    {
      name: 'README.md',
      type: 'file',
      path: '/README.md',
      size: 1024,
      language: 'markdown',
      content: `# My Awesome Project

A Go web API built with Gin framework.

## Getting Started

\`\`\`bash
go run cmd/server/main.go
\`\`\``,
    },
    {
      name: 'Dockerfile',
      type: 'file',
      path: '/Dockerfile',
      size: 512,
      language: 'dockerfile',
    },
    {
      name: 'Makefile',
      type: 'file',
      path: '/Makefile',
      size: 384,
      language: 'makefile',
    },
  ]

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

  const selectFile = (file: FileNode) => {
    if (file.type === 'file') {
      setSelectedFile(file.path)
    }
  }

  const renderFileTree = (nodes: FileNode[], depth = 0) => {
    return nodes.map((node) => (
      <div key={node.path}>
        <div
          className={`flex items-center space-x-2 py-1 px-2 hover:bg-gray-100 cursor-pointer rounded ${
            selectedFile === node.path ? 'bg-blue-50 border-l-2 border-blue-500' : ''
          }`}
          style={{ paddingLeft: `${depth * 16 + 8}px` }}
          onClick={() => {
            if (node.type === 'folder') {
              toggleFolder(node.path)
            } else {
              selectFile(node)
            }
          }}
        >
          {node.type === 'folder' ? (
            <>
              {expandedFolders.has(node.path) ? (
                <ChevronDownIcon className="w-4 h-4 text-gray-400" />
              ) : (
                <ChevronRightIcon className="w-4 h-4 text-gray-400" />
              )}
              {expandedFolders.has(node.path) ? (
                <FolderOpenIcon className="w-4 h-4 text-blue-500" />
              ) : (
                <FolderIcon className="w-4 h-4 text-blue-500" />
              )}
            </>
          ) : (
            <>
              <div className="w-4" />
              <DocumentIcon className="w-4 h-4 text-gray-400" />
            </>
          )}
          <span className="text-sm text-gray-700 flex-1">{node.name}</span>
          {node.type === 'file' && node.size && (
            <span className="text-xs text-gray-400">{formatFileSize(node.size)}</span>
          )}
        </div>
        {node.type === 'folder' && expandedFolders.has(node.path) && node.children && (
          <div>{renderFileTree(node.children, depth + 1)}</div>
        )}
      </div>
    ))
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes}B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)}MB`
  }

  const getSelectedFileContent = () => {
    const findFile = (nodes: FileNode[], path: string): FileNode | null => {
      for (const node of nodes) {
        if (node.path === path) return node
        if (node.children) {
          const found = findFile(node.children, path)
          if (found) return found
        }
      }
      return null
    }

    if (!selectedFile) return null
    return findFile(fileStructure, selectedFile)
  }

  const selectedFileNode = getSelectedFileContent()

  return (
    <div className="card h-full flex flex-col">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-lg font-semibold text-gray-900">File Explorer</h2>
          <div className="flex space-x-2">
            <button
              className="p-1.5 text-gray-400 hover:text-gray-600 rounded hover:bg-gray-100"
              title="Download Project"
            >
              <ArrowDownTrayIcon className="w-4 h-4" />
            </button>
            <button
              className="p-1.5 text-gray-400 hover:text-gray-600 rounded hover:bg-gray-100"
              title="Copy Structure"
            >
              <ClipboardDocumentIcon className="w-4 h-4" />
            </button>
          </div>
        </div>
        <p className="text-sm text-gray-600">Browse generated project files</p>
      </div>

      <div className="flex-1 flex flex-col min-h-0">
        {/* File Tree */}
        <div className="flex-1 overflow-y-auto border rounded-lg bg-white">
          <div className="p-2">
            {renderFileTree(fileStructure)}
          </div>
        </div>

        {/* File Preview */}
        {selectedFileNode && selectedFileNode.content && (
          <div className="mt-4 flex-1 min-h-0">
            <div className="bg-gray-50 rounded-lg p-3 border">
              <div className="flex items-center justify-between mb-2">
                <h4 className="text-sm font-medium text-gray-900">{selectedFileNode.name}</h4>
                <span className="text-xs text-gray-500">{selectedFileNode.language}</span>
              </div>
              <div className="bg-white rounded border p-3 max-h-48 overflow-y-auto">
                <pre className="text-xs text-gray-700 font-mono whitespace-pre-wrap">
                  {selectedFileNode.content}
                </pre>
              </div>
            </div>
          </div>
        )}

        {/* File Count Summary */}
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">Total Files</span>
            <span className="font-medium text-gray-900">
              {fileStructure.reduce((count, node) => {
                const countFiles = (n: FileNode): number => {
                  if (n.type === 'file') return 1
                  return (n.children || []).reduce((acc, child) => acc + countFiles(child), 0)
                }
                return count + countFiles(node)
              }, 0)}
            </span>
          </div>
        </div>
      </div>
    </div>
  )
}