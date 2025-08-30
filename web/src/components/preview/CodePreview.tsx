import { useState, useEffect, useMemo } from 'react'
import { ClipboardDocumentIcon, CheckIcon } from '@heroicons/react/24/outline'
import type { WSFileContent } from '../../types'

interface CodePreviewProps {
  file: WSFileContent | null
  isLoading?: boolean
  className?: string
}

// Simple syntax highlighting for common Go file types
function getLanguageFromPath(path: string): string {
  const ext = path.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'go':
      return 'go'
    case 'mod':
    case 'sum':
      return 'gomod'
    case 'yaml':
    case 'yml':
      return 'yaml'
    case 'json':
      return 'json'
    case 'md':
      return 'markdown'
    case 'dockerfile':
      return 'dockerfile'
    case 'sh':
      return 'bash'
    case 'sql':
      return 'sql'
    default:
      return 'text'
  }
}

// Basic syntax highlighting patterns
const syntaxHighlightPatterns = {
  go: [
    { pattern: /\b(package|import|func|var|const|type|struct|interface|if|else|for|range|return|defer|go|chan|select|case|default|break|continue|fallthrough|switch|map|slice)\b/g, class: 'text-purple-600 font-semibold' },
    { pattern: /\b(string|int|int32|int64|uint|uint32|uint64|float32|float64|bool|byte|rune|error|nil|true|false)\b/g, class: 'text-blue-600' },
    { pattern: /\/\/.*$/gm, class: 'text-gray-500 italic' },
    { pattern: /\/\*[\s\S]*?\*\//g, class: 'text-gray-500 italic' },
    { pattern: /"([^"\\\\]|\\\\.)*"/g, class: 'text-green-600' },
    { pattern: /`[^`]*`/g, class: 'text-green-600' },
    { pattern: /\b\d+(\.\d+)?\b/g, class: 'text-orange-600' }
  ],
  yaml: [
    { pattern: /^[\w\-]+:/gm, class: 'text-blue-600 font-semibold' },
    { pattern: /#.*$/gm, class: 'text-gray-500 italic' },
    { pattern: /"([^"\\\\]|\\\\.)*"/g, class: 'text-green-600' },
    { pattern: /'([^'\\\\]|\\\\.)*'/g, class: 'text-green-600' },
    { pattern: /\btrue\b|\bfalse\b|\bnull\b/g, class: 'text-purple-600' }
  ],
  json: [
    { pattern: /"[\w\-]+"\s*:/g, class: 'text-blue-600 font-semibold' },
    { pattern: /"([^"\\\\]|\\\\.)*"/g, class: 'text-green-600' },
    { pattern: /\btrue\b|\bfalse\b|\bnull\b/g, class: 'text-purple-600' },
    { pattern: /\b\d+(\.\d+)?\b/g, class: 'text-orange-600' }
  ],
  markdown: [
    { pattern: /^#{1,6}\s.*$/gm, class: 'text-blue-600 font-bold' },
    { pattern: /\*\*([^*]+)\*\*/g, class: 'font-bold' },
    { pattern: /\*([^*]+)\*/g, class: 'italic' },
    { pattern: /`([^`]+)`/g, class: 'bg-gray-100 px-1 rounded text-red-600 font-mono' },
    { pattern: /\[([^\]]+)\]\(([^)]+)\)/g, class: 'text-blue-600 underline' }
  ]
}

function applySyntaxHighlighting(content: string, language: string): string {
  const patterns = syntaxHighlightPatterns[language as keyof typeof syntaxHighlightPatterns]
  if (!patterns) return content

  let highlighted = content
  patterns.forEach(({ pattern, class: className }) => {
    highlighted = highlighted.replace(pattern, (match) => {
      return `<span class="${className}">${match}</span>`
    })
  })

  return highlighted
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

export default function CodePreview({ file, isLoading = false, className = '' }: CodePreviewProps) {
  const [copied, setCopied] = useState(false)
  const [showLineNumbers, setShowLineNumbers] = useState(true)
  
  const language = useMemo(() => {
    return file ? getLanguageFromPath(file.path) : 'text'
  }, [file?.path])
  
  const highlightedContent = useMemo(() => {
    if (!file?.content) return ''
    return applySyntaxHighlighting(file.content, language)
  }, [file?.content, language])
  
  const lineCount = useMemo(() => {
    if (!file?.content) return 0
    return file.content.split('\n').length
  }, [file?.content])
  
  const copyToClipboard = async () => {
    if (!file?.content) return
    
    try {
      await navigator.clipboard.writeText(file.content)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy to clipboard:', err)
      // Fallback for older browsers
      const textArea = document.createElement('textarea')
      textArea.value = file.content
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }
  
  if (isLoading) {
    return (
      <div className={`flex items-center justify-center h-full bg-gray-50 ${className}`}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
          <div className="text-gray-500 text-sm">Loading file content...</div>
        </div>
      </div>
    )
  }
  
  if (!file) {
    return (
      <div className={`flex items-center justify-center h-full bg-gray-50 ${className}`}>
        <div className="text-center text-gray-500">
          <div className="text-lg mb-2">No file selected</div>
          <div className="text-sm">Select a file from the tree to preview its content</div>
        </div>
      </div>
    )
  }
  
  if (file.isDir) {
    return (
      <div className={`flex items-center justify-center h-full bg-gray-50 ${className}`}>
        <div className="text-center text-gray-500">
          <div className="text-lg mb-2">Directory Selected</div>
          <div className="text-sm">Select a file to view its content</div>
        </div>
      </div>
    )
  }
  
  return (
    <div className={`flex flex-col h-full ${className}`}>
      {/* File Header */}
      <div className="flex items-center justify-between p-3 border-b bg-gray-50">
        <div className="flex items-center gap-3">
          <div>
            <div className="font-medium text-gray-900">{file.path.split('/').pop()}</div>
            <div className="text-xs text-gray-500 flex items-center gap-2">
              <span>{file.path}</span>
              <span>•</span>
              <span>{formatFileSize(file.size)}</span>
              <span>•</span>
              <span>{lineCount} lines</span>
              <span>•</span>
              <span className="uppercase">{language}</span>
            </div>
          </div>
        </div>
        
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowLineNumbers(!showLineNumbers)}
            className="px-2 py-1 text-xs bg-white border rounded hover:bg-gray-50"
            title="Toggle line numbers"
          >
            {showLineNumbers ? 'Hide' : 'Show'} Lines
          </button>
          
          <button
            onClick={copyToClipboard}
            className="flex items-center gap-1 px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
            title="Copy to clipboard"
          >
            {copied ? (
              <>
                <CheckIcon className="h-3 w-3" />
                Copied!
              </>
            ) : (
              <>
                <ClipboardDocumentIcon className="h-3 w-3" />
                Copy
              </>
            )}
          </button>
        </div>
      </div>
      
      {/* Code Content */}
      <div className="flex-1 overflow-auto bg-white">
        <div className="flex">
          {/* Line Numbers */}
          {showLineNumbers && (
            <div className="flex-shrink-0 p-4 pr-2 bg-gray-50 border-r text-xs text-gray-400 font-mono select-none">
              {Array.from({ length: lineCount }, (_, i) => (
                <div key={i + 1} className="leading-5">
                  {i + 1}
                </div>
              ))}
            </div>
          )}
          
          {/* Code */}
          <div className="flex-1 p-4">
            {language === 'text' ? (
              <pre className="text-sm text-gray-700 whitespace-pre-wrap font-mono leading-5">
                {file.content}
              </pre>
            ) : (
              <pre 
                className="text-sm text-gray-700 whitespace-pre-wrap font-mono leading-5"
                dangerouslySetInnerHTML={{ __html: highlightedContent }}
              />
            )}
          </div>
        </div>
      </div>
      
      {/* File Stats */}
      <div className="border-t bg-gray-50 px-4 py-2 text-xs text-gray-500">
        <div className="flex justify-between items-center">
          <div>
            {file.modTime && (
              <span>Modified: {new Date(file.modTime).toLocaleString()}</span>
            )}
          </div>
          <div className="flex items-center gap-4">
            {file.mode && <span>Mode: {file.mode}</span>}
            <span>Size: {formatFileSize(file.size)}</span>
          </div>
        </div>
      </div>
    </div>
  )
}