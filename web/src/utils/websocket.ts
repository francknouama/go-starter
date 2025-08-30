import type { WSMessage, WSMessageType } from '../types/index'

/**
 * WebSocket utilities for the Go Starter Web UI
 */

export interface WebSocketConfig {
  url?: string
  reconnectAttempts?: number
  reconnectInterval?: number
  heartbeatInterval?: number
  debug?: boolean
}

export interface WebSocketEventHandlers {
  onOpen?: () => void
  onClose?: (event: CloseEvent) => void
  onError?: (error: Event) => void
  onMessage?: (message: WSMessage) => void
  onReconnect?: (attempt: number) => void
  onReconnectFailed?: () => void
}

export class WebSocketManager {
  private ws: WebSocket | null = null
  private config: Required<WebSocketConfig>
  private handlers: WebSocketEventHandlers
  private reconnectAttempts: number = 0
  private reconnectTimeout: NodeJS.Timeout | null = null
  private heartbeatInterval: NodeJS.Timeout | null = null
  private messageQueue: WSMessage[] = []
  private isManualClose: boolean = false

  constructor(config: WebSocketConfig = {}, handlers: WebSocketEventHandlers = {}) {
    this.config = {
      url: config.url || this.getDefaultUrl(),
      reconnectAttempts: config.reconnectAttempts || 5,
      reconnectInterval: config.reconnectInterval || 3000,
      heartbeatInterval: config.heartbeatInterval || 30000,
      debug: config.debug || false
    }
    this.handlers = handlers
  }

  private getDefaultUrl(): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    return `${protocol}//${window.location.host}/api/v1/ws`
  }

  private log(message: string, ...args: any[]): void {
    if (this.config.debug) {
      console.log(`[WebSocketManager] ${message}`, ...args)
    }
  }

  private setupEventHandlers(): void {
    if (!this.ws) return

    this.ws.onopen = () => {
      this.log('Connected')
      this.reconnectAttempts = 0
      this.isManualClose = false
      this.startHeartbeat()
      this.flushMessageQueue()
      this.handlers.onOpen?.()
    }

    this.ws.onclose = (event) => {
      this.log('Disconnected', event.code, event.reason)
      this.stopHeartbeat()
      this.handlers.onClose?.(event)

      if (!this.isManualClose && this.shouldReconnect()) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      this.log('Error', error)
      this.handlers.onError?.(error)
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data)
        this.log('Received message', message.type, message)
        
        // Handle built-in message types
        if (message.type === 'ping') {
          this.send('pong', { timestamp: Date.now() })
        } else {
          this.handlers.onMessage?.(message)
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
  }

  private shouldReconnect(): boolean {
    return this.reconnectAttempts < this.config.reconnectAttempts
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
    }

    const attempt = this.reconnectAttempts + 1
    const delay = Math.min(this.config.reconnectInterval * Math.pow(1.5, attempt - 1), 30000)

    this.log(`Scheduling reconnect attempt ${attempt}/${this.config.reconnectAttempts} in ${delay}ms`)

    this.reconnectTimeout = setTimeout(() => {
      this.reconnectAttempts = attempt
      this.handlers.onReconnect?.(attempt)
      this.connect()
    }, delay)
  }

  private startHeartbeat(): void {
    this.stopHeartbeat()
    
    this.heartbeatInterval = setInterval(() => {
      if (this.isConnected()) {
        this.send('ping', { timestamp: Date.now() })
      }
    }, this.config.heartbeatInterval)
  }

  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  private flushMessageQueue(): void {
    while (this.messageQueue.length > 0 && this.isConnected()) {
      const message = this.messageQueue.shift()
      if (message) {
        this.sendImmediate(message)
      }
    }
  }

  private sendImmediate(message: WSMessage): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return false
    }

    try {
      this.ws.send(JSON.stringify(message))
      this.log('Sent message', message.type, message)
      return true
    } catch (error) {
      console.error('Error sending WebSocket message:', error)
      return false
    }
  }

  public connect(): void {
    if (this.isConnected()) {
      this.log('Already connected')
      return
    }

    if (this.isConnecting()) {
      this.log('Already connecting')
      return
    }

    this.log('Connecting to', this.config.url)

    try {
      this.ws = new WebSocket(this.config.url)
      this.setupEventHandlers()
    } catch (error) {
      console.error('Error creating WebSocket:', error)
      this.handlers.onError?.(new Event('connection-failed'))
    }
  }

  public disconnect(): void {
    this.log('Disconnecting')
    this.isManualClose = true
    this.stopHeartbeat()

    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.messageQueue = []
    this.reconnectAttempts = 0
  }

  public send(type: WSMessageType, data: any, requestId?: string): void {
    const message: WSMessage = {
      type,
      data,
      timestamp: new Date().toISOString(),
      ...(requestId && { requestId })
    }

    if (this.isConnected()) {
      this.sendImmediate(message)
    } else {
      // Queue message for when connection is established
      this.messageQueue.push(message)
      this.log('Queued message', type, 'Total queued:', this.messageQueue.length)

      // Attempt to connect if not connected
      if (!this.isConnecting()) {
        this.connect()
      }
    }
  }

  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  public isConnecting(): boolean {
    return this.ws?.readyState === WebSocket.CONNECTING
  }

  public getConnectionState(): 'connecting' | 'connected' | 'disconnected' | 'error' {
    if (!this.ws) return 'disconnected'
    
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting'
      case WebSocket.OPEN:
        return 'connected'
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
        return 'disconnected'
      default:
        return 'error'
    }
  }

  public getReconnectAttempts(): number {
    return this.reconnectAttempts
  }

  public getMaxReconnectAttempts(): number {
    return this.config.reconnectAttempts
  }

  public getQueuedMessageCount(): number {
    return this.messageQueue.length
  }
}

/**
 * Create a WebSocket connection with automatic reconnection
 */
export function createWebSocket(
  config: WebSocketConfig = {},
  handlers: WebSocketEventHandlers = {}
): WebSocketManager {
  return new WebSocketManager(config, handlers)
}

/**
 * WebSocket connection status indicator
 */
export function getConnectionStatusInfo(manager: WebSocketManager) {
  const state = manager.getConnectionState()
  const attempts = manager.getReconnectAttempts()
  const maxAttempts = manager.getMaxReconnectAttempts()
  const queuedMessages = manager.getQueuedMessageCount()

  return {
    state,
    attempts,
    maxAttempts,
    queuedMessages,
    isConnected: state === 'connected',
    isConnecting: state === 'connecting',
    hasQueuedMessages: queuedMessages > 0,
    reconnectExhausted: attempts >= maxAttempts,
    statusText: getStatusText(state, attempts, maxAttempts)
  }
}

function getStatusText(
  state: 'connecting' | 'connected' | 'disconnected' | 'error',
  attempts: number,
  maxAttempts: number
): string {
  switch (state) {
    case 'connected':
      return 'Connected'
    case 'connecting':
      return attempts > 0 ? `Reconnecting... (${attempts}/${maxAttempts})` : 'Connecting...'
    case 'disconnected':
      return attempts >= maxAttempts ? 'Connection failed' : 'Disconnected'
    case 'error':
      return 'Connection error'
    default:
      return 'Unknown'
  }
}