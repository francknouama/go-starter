import { useEffect, useState } from 'react'
import { 
  WifiIcon, 
  ExclamationTriangleIcon, 
  ArrowPathIcon,
  XMarkIcon 
} from '@heroicons/react/24/outline'
import type { WSConnectionState } from '../../types'

interface ConnectionManagerProps {
  connectionState: WSConnectionState
  onRetry?: () => void
  onDismiss?: () => void
  className?: string
}

interface ConnectionAlert {
  id: string
  type: 'warning' | 'error' | 'info'
  message: string
  timestamp: number
  persistent?: boolean
}

export default function ConnectionManager({ 
  connectionState, 
  onRetry, 
  onDismiss,
  className = '' 
}: ConnectionManagerProps) {
  const [alerts, setAlerts] = useState<ConnectionAlert[]>([])
  const [dismissed, setDismissed] = useState<Set<string>>(new Set())
  
  // Generate alerts based on connection state
  useEffect(() => {
    const newAlerts: ConnectionAlert[] = []
    const now = Date.now()
    
    // Connection lost
    if (!connectionState.connected && !connectionState.connecting && connectionState.error) {
      newAlerts.push({
        id: 'connection-lost',
        type: 'error',
        message: 'Real-time connection lost. Preview updates are disabled.',
        timestamp: now,
        persistent: true
      })
    }
    
    // Reconnecting
    if (connectionState.connecting && connectionState.reconnectAttempts > 0) {
      newAlerts.push({
        id: 'reconnecting',
        type: 'warning',
        message: `Attempting to reconnect... (${connectionState.reconnectAttempts}/5)`,
        timestamp: now,
        persistent: true
      })
    }
    
    // Multiple reconnection attempts
    if (connectionState.reconnectAttempts >= 3) {
      newAlerts.push({
        id: 'connection-unstable',
        type: 'warning',
        message: 'Connection appears unstable. Check your network connection.',
        timestamp: now
      })
    }
    
    // Connection restored
    if (connectionState.connected && connectionState.reconnectAttempts > 0) {
      newAlerts.push({
        id: 'connection-restored',
        type: 'info',
        message: 'Real-time connection restored.',
        timestamp: now
      })
      
      // Auto-dismiss after 3 seconds
      setTimeout(() => {
        setDismissed(prev => new Set([...prev, 'connection-restored']))
      }, 3000)
    }
    
    setAlerts(newAlerts)
  }, [
    connectionState.connected, 
    connectionState.connecting, 
    connectionState.error, 
    connectionState.reconnectAttempts
  ])
  
  const dismissAlert = (alertId: string) => {
    setDismissed(prev => new Set([...prev, alertId]))
  }
  
  const getAlertStyles = (type: ConnectionAlert['type']) => {
    switch (type) {
      case 'error':
        return 'bg-red-50 border-red-200 text-red-800'
      case 'warning':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800'
      case 'info':
        return 'bg-blue-50 border-blue-200 text-blue-800'
      default:
        return 'bg-gray-50 border-gray-200 text-gray-800'
    }
  }
  
  const getAlertIcon = (type: ConnectionAlert['type']) => {
    switch (type) {
      case 'error':
        return <ExclamationTriangleIcon className="h-4 w-4 text-red-500" />
      case 'warning':
        return <ExclamationTriangleIcon className="h-4 w-4 text-yellow-500" />
      case 'info':
        return <WifiIcon className="h-4 w-4 text-blue-500" />
      default:
        return <WifiIcon className="h-4 w-4 text-gray-500" />
    }
  }
  
  // Filter out dismissed alerts
  const visibleAlerts = alerts.filter(alert => !dismissed.has(alert.id))
  
  if (visibleAlerts.length === 0) {
    return null
  }
  
  return (
    <div className={`space-y-2 ${className}`}>
      {visibleAlerts.map(alert => (
        <div
          key={alert.id}
          className={`flex items-center justify-between p-3 rounded-lg border ${getAlertStyles(alert.type)}`}
        >
          <div className="flex items-center gap-3">
            {getAlertIcon(alert.type)}
            <div>
              <div className="text-sm font-medium">{alert.message}</div>
              {alert.id === 'connection-lost' && onRetry && (
                <button
                  onClick={onRetry}
                  className="inline-flex items-center gap-1 mt-2 text-xs text-red-700 hover:text-red-800 underline"
                >
                  <ArrowPathIcon className="h-3 w-3" />
                  Retry Connection
                </button>
              )}
            </div>
          </div>
          
          {!alert.persistent && (
            <button
              onClick={() => dismissAlert(alert.id)}
              className="text-gray-500 hover:text-gray-700"
            >
              <XMarkIcon className="h-4 w-4" />
            </button>
          )}
        </div>
      ))}
    </div>
  )
}

// Connection Status Indicator Component
interface ConnectionStatusProps {
  connectionState: WSConnectionState
  showLabel?: boolean
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function ConnectionStatus({ 
  connectionState, 
  showLabel = true, 
  size = 'md',
  className = '' 
}: ConnectionStatusProps) {
  const getStatusColor = () => {
    if (connectionState.connected) return 'bg-green-500'
    if (connectionState.connecting) return 'bg-yellow-500'
    return 'bg-red-500'
  }
  
  const getStatusText = () => {
    if (connectionState.connected) return 'Connected'
    if (connectionState.connecting) {
      return connectionState.reconnectAttempts > 0 
        ? `Reconnecting (${connectionState.reconnectAttempts}/5)`
        : 'Connecting...'
    }
    return connectionState.error || 'Disconnected'
  }
  
  const getSizeClasses = () => {
    switch (size) {
      case 'sm':
        return 'h-2 w-2'
      case 'lg':
        return 'h-4 w-4'
      default:
        return 'h-3 w-3'
    }
  }
  
  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <div className={`rounded-full ${getSizeClasses()} ${getStatusColor()} ${connectionState.connecting ? 'animate-pulse' : ''}`} />
      {showLabel && (
        <span className={`text-sm ${
          connectionState.connected 
            ? 'text-green-600' 
            : connectionState.connecting 
              ? 'text-yellow-600' 
              : 'text-red-600'
        }`}>
          {getStatusText()}
        </span>
      )}
    </div>
  )
}

// Network Error Recovery Component
interface NetworkErrorRecoveryProps {
  error: string
  onRetry: () => void
  onDismiss?: () => void
  retryCount?: number
  maxRetries?: number
}

export function NetworkErrorRecovery({ 
  error, 
  onRetry, 
  onDismiss, 
  retryCount = 0, 
  maxRetries = 5 
}: NetworkErrorRecoveryProps) {
  const [isRetrying, setIsRetrying] = useState(false)
  
  const handleRetry = async () => {
    setIsRetrying(true)
    try {
      await onRetry()
    } finally {
      setTimeout(() => setIsRetrying(false), 1000)
    }
  }
  
  const canRetry = retryCount < maxRetries
  
  return (
    <div className="bg-red-50 border border-red-200 rounded-lg p-4">
      <div className="flex items-start gap-3">
        <ExclamationTriangleIcon className="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
        <div className="flex-1">
          <h3 className="text-sm font-medium text-red-800 mb-1">
            Connection Error
          </h3>
          <p className="text-sm text-red-700 mb-3">{error}</p>
          
          <div className="flex items-center gap-3">
            {canRetry ? (
              <button
                onClick={handleRetry}
                disabled={isRetrying}
                className="inline-flex items-center gap-2 px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isRetrying ? (
                  <>
                    <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-white" />
                    Retrying...
                  </>
                ) : (
                  <>
                    <ArrowPathIcon className="h-3 w-3" />
                    Retry ({retryCount + 1}/{maxRetries})
                  </>
                )}
              </button>
            ) : (
              <span className="text-sm text-red-700">
                Maximum retry attempts reached
              </span>
            )}
            
            {onDismiss && (
              <button
                onClick={onDismiss}
                className="text-sm text-red-600 hover:text-red-800 underline"
              >
                Dismiss
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}