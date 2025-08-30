import { useEffect, useRef, useState, useCallback } from 'react'
import type { WSMessage, WSConnectionState, WSMessageType } from '../types/index'

interface UseWebSocketOptions {
  url?: string
  reconnectAttempts?: number
  reconnectInterval?: number
  onMessage?: (message: WSMessage) => void
  onConnect?: () => void
  onDisconnect?: () => void
  onError?: (error: Event) => void
}

interface UseWebSocketReturn {
  connectionState: WSConnectionState
  sendMessage: (type: WSMessageType, data: any, requestId?: string) => void
  connect: () => void
  disconnect: () => void
  subscribe: (messageType: WSMessageType, handler: (data: any) => void) => () => void
}

export function useWebSocket(options: UseWebSocketOptions = {}): UseWebSocketReturn {
  const {
    url = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/v1/ws`,
    reconnectAttempts = 5,
    reconnectInterval = 3000,
    onMessage,
    onConnect,
    onDisconnect,
    onError
  } = options

  const [connectionState, setConnectionState] = useState<WSConnectionState>({
    connected: false,
    connecting: false,
    reconnectAttempts: 0
  })

  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const messageHandlersRef = useRef<Map<WSMessageType, Set<(data: any) => void>>>(new Map())
  const reconnectAttemptsRef = useRef(0)

  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
  }, [])

  const updateConnectionState = useCallback((updates: Partial<WSConnectionState>) => {
    setConnectionState(prev => ({ ...prev, ...updates }))
  }, [])

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const message: WSMessage = JSON.parse(event.data)
      
      // Call global message handler
      onMessage?.(message)
      
      // Call type-specific handlers
      const handlers = messageHandlersRef.current.get(message.type)
      if (handlers) {
        handlers.forEach(handler => {
          try {
            handler(message.data)
          } catch (error) {
            console.error(`Error in WebSocket message handler for type ${message.type}:`, error)
          }
        })
      }
    } catch (error) {
      console.error('Error parsing WebSocket message:', error)
    }
  }, [onMessage])

  const handleOpen = useCallback(() => {
    console.log('WebSocket connected')
    reconnectAttemptsRef.current = 0
    updateConnectionState({
      connected: true,
      connecting: false,
      error: undefined,
      reconnectAttempts: 0
    })
    clearReconnectTimeout()
    onConnect?.()
  }, [updateConnectionState, clearReconnectTimeout, onConnect])

  const handleClose = useCallback(() => {
    console.log('WebSocket disconnected')
    updateConnectionState({
      connected: false,
      connecting: false
    })
    onDisconnect?.()
    
    // Attempt reconnection if we haven't exceeded max attempts
    if (reconnectAttemptsRef.current < reconnectAttempts) {
      const nextAttempt = reconnectAttemptsRef.current + 1
      reconnectAttemptsRef.current = nextAttempt
      
      updateConnectionState({
        reconnectAttempts: nextAttempt,
        lastReconnectAttempt: Date.now()
      })
      
      console.log(`Attempting to reconnect... (${nextAttempt}/${reconnectAttempts})`)
      
      reconnectTimeoutRef.current = setTimeout(() => {
        connect()
      }, reconnectInterval)
    } else {
      updateConnectionState({
        error: `Failed to reconnect after ${reconnectAttempts} attempts`
      })
    }
  }, [reconnectAttempts, reconnectInterval, updateConnectionState, onDisconnect])

  const handleError = useCallback((event: Event) => {
    console.error('WebSocket error:', event)
    updateConnectionState({
      connected: false,
      connecting: false,
      error: 'Connection error'
    })
    onError?.(event)
  }, [updateConnectionState, onError])

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return // Already connected
    }
    
    if (wsRef.current?.readyState === WebSocket.CONNECTING) {
      return // Already connecting
    }

    updateConnectionState({ connecting: true, error: undefined })
    
    try {
      const ws = new WebSocket(url)
      wsRef.current = ws
      
      ws.addEventListener('open', handleOpen)
      ws.addEventListener('message', handleMessage)
      ws.addEventListener('close', handleClose)
      ws.addEventListener('error', handleError)
    } catch (error) {
      console.error('Error creating WebSocket:', error)
      updateConnectionState({
        connecting: false,
        error: 'Failed to create WebSocket connection'
      })
    }
  }, [url, handleOpen, handleMessage, handleClose, handleError, updateConnectionState])

  const disconnect = useCallback(() => {
    clearReconnectTimeout()
    reconnectAttemptsRef.current = reconnectAttempts // Prevent reconnection
    
    if (wsRef.current) {
      // Remove event listeners to prevent handleClose from triggering reconnection
      wsRef.current.removeEventListener('open', handleOpen)
      wsRef.current.removeEventListener('message', handleMessage)
      wsRef.current.removeEventListener('close', handleClose)
      wsRef.current.removeEventListener('error', handleError)
      
      wsRef.current.close()
      wsRef.current = null
    }
    
    updateConnectionState({
      connected: false,
      connecting: false,
      error: undefined
    })
  }, [clearReconnectTimeout, reconnectAttempts, handleOpen, handleMessage, handleClose, handleError, updateConnectionState])

  const sendMessage = useCallback((type: WSMessageType, data: any, requestId?: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const message: WSMessage = {
        type,
        data,
        timestamp: new Date().toISOString(),
        ...(requestId && { requestId })
      }
      
      try {
        wsRef.current.send(JSON.stringify(message))
      } catch (error) {
        console.error('Error sending WebSocket message:', error)
      }
    } else {
      console.warn('Cannot send message: WebSocket is not connected')
    }
  }, [])

  const subscribe = useCallback((messageType: WSMessageType, handler: (data: any) => void) => {
    const handlers = messageHandlersRef.current.get(messageType) || new Set()
    handlers.add(handler)
    messageHandlersRef.current.set(messageType, handlers)
    
    // Return unsubscribe function
    return () => {
      const currentHandlers = messageHandlersRef.current.get(messageType)
      if (currentHandlers) {
        currentHandlers.delete(handler)
        if (currentHandlers.size === 0) {
          messageHandlersRef.current.delete(messageType)
        }
      }
    }
  }, [])

  // Auto-connect on mount
  useEffect(() => {
    connect()
    
    return () => {
      disconnect()
    }
  }, []) // Only run on mount/unmount

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      clearReconnectTimeout()
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [clearReconnectTimeout])

  return {
    connectionState,
    sendMessage,
    connect,
    disconnect,
    subscribe
  }
}