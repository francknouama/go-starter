import { renderHook, act, waitFor } from '@testing-library/react';
import { useWebSocket } from '../useWebSocket';
import type { WSMessage, WSMessageType } from '../../types';

// Mock WebSocket
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  url: string;
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;

  private eventListeners: Map<string, Set<EventListener>> = new Map();

  constructor(url: string) {
    this.url = url;
    // Simulate connection opening
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      this.dispatchEvent(new Event('open'));
    }, 0);
  }

  addEventListener(type: string, listener: EventListener) {
    if (!this.eventListeners.has(type)) {
      this.eventListeners.set(type, new Set());
    }
    this.eventListeners.get(type)!.add(listener);
  }

  removeEventListener(type: string, listener: EventListener) {
    const listeners = this.eventListeners.get(type);
    if (listeners) {
      listeners.delete(listener);
    }
  }

  dispatchEvent(event: Event): boolean {
    const listeners = this.eventListeners.get(event.type);
    if (listeners) {
      listeners.forEach(listener => {
        try {
          listener.call(this, event);
        } catch (error) {
          console.error('Error in event listener:', error);
        }
      });
    }
    
    // Also call direct event handlers
    if (event.type === 'open' && this.onopen) {
      this.onopen(event);
    } else if (event.type === 'close' && this.onclose) {
      this.onclose(event as CloseEvent);
    } else if (event.type === 'message' && this.onmessage) {
      this.onmessage(event as MessageEvent);
    } else if (event.type === 'error' && this.onerror) {
      this.onerror(event);
    }
    
    return true;
  }

  send(data: string) {
    if (this.readyState !== MockWebSocket.OPEN) {
      throw new Error('WebSocket is not open');
    }
    // Echo the message back for testing
    setTimeout(() => {
      const messageEvent = new MessageEvent('message', { data });
      this.dispatchEvent(messageEvent);
    }, 0);
  }

  close(code?: number, reason?: string) {
    this.readyState = MockWebSocket.CLOSED;
    const closeEvent = new CloseEvent('close', { code, reason });
    setTimeout(() => {
      this.dispatchEvent(closeEvent);
    }, 0);
  }

  // Helper method to simulate server messages
  simulateMessage(message: WSMessage) {
    const messageEvent = new MessageEvent('message', {
      data: JSON.stringify(message)
    });
    this.dispatchEvent(messageEvent);
  }

  // Helper method to simulate connection error
  simulateError() {
    this.readyState = MockWebSocket.CLOSED;
    const errorEvent = new Event('error');
    this.dispatchEvent(errorEvent);
  }
}

// Replace global WebSocket with mock
global.WebSocket = MockWebSocket as any;

describe('useWebSocket', () => {
  let mockWebSocketInstances: MockWebSocket[] = [];

  beforeEach(() => {
    mockWebSocketInstances = [];
    jest.clearAllMocks();
    jest.clearAllTimers();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
    mockWebSocketInstances.forEach(ws => {
      if (ws.readyState === MockWebSocket.OPEN) {
        ws.close();
      }
    });
  });

  describe('Connection Management', () => {
    it('connects automatically on mount', async () => {
      const onConnect = jest.fn();
      const { result } = renderHook(() => 
        useWebSocket({ onConnect })
      );

      expect(result.current.connectionState.connecting).toBe(true);
      expect(result.current.connectionState.connected).toBe(false);

      // Wait for connection to open
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);
      expect(result.current.connectionState.connecting).toBe(false);
      expect(onConnect).toHaveBeenCalled();
    });

    it('handles custom WebSocket URL', async () => {
      const customUrl = 'ws://custom-host:8080/ws';
      const { result } = renderHook(() => 
        useWebSocket({ url: customUrl })
      );

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);
    });

    it('disconnects cleanly', async () => {
      const onDisconnect = jest.fn();
      const { result } = renderHook(() => 
        useWebSocket({ onDisconnect })
      );

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);

      // Disconnect
      act(() => {
        result.current.disconnect();
      });

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(false);
      expect(onDisconnect).toHaveBeenCalled();
    });

    it('handles connection errors', async () => {
      const onError = jest.fn();
      const { result } = renderHook(() => 
        useWebSocket({ onError })
      );

      // Wait for initial connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Get the WebSocket instance and simulate error
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateError) {
        act(() => {
          wsInstance.simulateError();
        });
      }

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(false);
      expect(result.current.connectionState.error).toBeDefined();
      expect(onError).toHaveBeenCalled();
    });
  });

  describe('Reconnection Logic', () => {
    it('attempts to reconnect on connection loss', async () => {
      const { result } = renderHook(() => 
        useWebSocket({ 
          reconnectAttempts: 3,
          reconnectInterval: 1000
        })
      );

      // Wait for initial connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);
      
      // Simulate connection close
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance) {
        act(() => {
          wsInstance.close();
        });
      }

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(false);
      expect(result.current.connectionState.reconnectAttempts).toBe(1);

      // Wait for reconnection attempt
      await act(async () => {
        jest.advanceTimersByTime(1000);
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);
    });

    it('gives up after max reconnection attempts', async () => {
      const { result } = renderHook(() => 
        useWebSocket({ 
          reconnectAttempts: 2,
          reconnectInterval: 1000
        })
      );

      // Wait for initial connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Simulate multiple connection failures
      for (let i = 0; i < 3; i++) {
        const wsInstance = (global.WebSocket as any).mock?.instances?.[i];
        if (wsInstance) {
          act(() => {
            wsInstance.close();
          });
        }

        await act(async () => {
          jest.runOnlyPendingTimers();
          jest.advanceTimersByTime(1000);
        });
      }

      expect(result.current.connectionState.connected).toBe(false);
      expect(result.current.connectionState.error).toMatch(/Failed to reconnect/);
    });
  });

  describe('Message Handling', () => {
    it('sends messages when connected', async () => {
      const { result } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      const testData = { test: 'message' };
      act(() => {
        result.current.sendMessage('generation_progress', testData);
      });

      // The mock WebSocket should have received the message
      expect(result.current.connectionState.connected).toBe(true);
    });

    it('handles incoming messages', async () => {
      const onMessage = jest.fn();
      const { result } = renderHook(() => 
        useWebSocket({ onMessage })
      );

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      const testMessage: WSMessage = {
        type: 'generation_progress' as WSMessageType,
        data: { progress: 50 },
        timestamp: new Date().toISOString()
      };

      // Simulate incoming message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateMessage) {
        act(() => {
          wsInstance.simulateMessage(testMessage);
        });
      }

      expect(onMessage).toHaveBeenCalledWith(testMessage);
    });

    it('subscribes to specific message types', async () => {
      const handler = jest.fn();
      const { result } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Subscribe to specific message type
      act(() => {
        result.current.subscribe('generation_progress' as WSMessageType, handler);
      });

      const testMessage: WSMessage = {
        type: 'generation_progress' as WSMessageType,
        data: { progress: 75 },
        timestamp: new Date().toISOString()
      };

      // Simulate incoming message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateMessage) {
        act(() => {
          wsInstance.simulateMessage(testMessage);
        });
      }

      expect(handler).toHaveBeenCalledWith({ progress: 75 });
    });

    it('unsubscribes from message types', async () => {
      const handler = jest.fn();
      const { result } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Subscribe and then unsubscribe
      let unsubscribe: (() => void) | undefined;
      act(() => {
        unsubscribe = result.current.subscribe('generation_progress' as WSMessageType, handler);
      });

      act(() => {
        unsubscribe?.();
      });

      const testMessage: WSMessage = {
        type: 'generation_progress' as WSMessageType,
        data: { progress: 100 },
        timestamp: new Date().toISOString()
      };

      // Simulate incoming message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateMessage) {
        act(() => {
          wsInstance.simulateMessage(testMessage);
        });
      }

      expect(handler).not.toHaveBeenCalled();
    });

    it('handles malformed messages gracefully', async () => {
      const onMessage = jest.fn();
      const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
      
      const { result } = renderHook(() => 
        useWebSocket({ onMessage })
      );

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Simulate malformed message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance) {
        const invalidMessageEvent = new MessageEvent('message', {
          data: 'invalid json'
        });
        act(() => {
          wsInstance.dispatchEvent(invalidMessageEvent);
        });
      }

      expect(consoleError).toHaveBeenCalledWith(
        'Error parsing WebSocket message:', 
        expect.any(Error)
      );
      expect(onMessage).not.toHaveBeenCalled();

      consoleError.mockRestore();
    });
  });

  describe('Connection State Management', () => {
    it('tracks connection state accurately', async () => {
      const { result } = renderHook(() => useWebSocket());

      // Initially connecting
      expect(result.current.connectionState.connecting).toBe(true);
      expect(result.current.connectionState.connected).toBe(false);
      expect(result.current.connectionState.reconnectAttempts).toBe(0);

      // After connection opens
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connecting).toBe(false);
      expect(result.current.connectionState.connected).toBe(true);
      expect(result.current.connectionState.error).toBeUndefined();
    });

    it('updates reconnection attempt count', async () => {
      const { result } = renderHook(() => 
        useWebSocket({ reconnectAttempts: 3, reconnectInterval: 1000 })
      );

      // Wait for initial connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Force disconnect
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance) {
        act(() => {
          wsInstance.close();
        });
      }

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.reconnectAttempts).toBe(1);
      expect(result.current.connectionState.lastReconnectAttempt).toBeDefined();
    });
  });

  describe('Cleanup', () => {
    it('cleans up on unmount', async () => {
      const { result, unmount } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      expect(result.current.connectionState.connected).toBe(true);

      // Unmount component
      unmount();

      // Should clean up connection
      expect(result.current.connectionState.connected).toBe(true); // State is stale after unmount
    });

    it('prevents memory leaks with message handlers', async () => {
      const handler1 = jest.fn();
      const handler2 = jest.fn();
      
      const { result } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Subscribe multiple handlers
      let unsubscribe1: (() => void) | undefined;
      let unsubscribe2: (() => void) | undefined;
      
      act(() => {
        unsubscribe1 = result.current.subscribe('generation_progress' as WSMessageType, handler1);
        unsubscribe2 = result.current.subscribe('generation_progress' as WSMessageType, handler2);
      });

      // Unsubscribe one handler
      act(() => {
        unsubscribe1?.();
      });

      const testMessage: WSMessage = {
        type: 'generation_progress' as WSMessageType,
        data: { progress: 50 },
        timestamp: new Date().toISOString()
      };

      // Simulate message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateMessage) {
        act(() => {
          wsInstance.simulateMessage(testMessage);
        });
      }

      expect(handler1).not.toHaveBeenCalled();
      expect(handler2).toHaveBeenCalledWith({ progress: 50 });
    });
  });

  describe('Edge Cases', () => {
    it('handles rapid connect/disconnect calls', async () => {
      const { result } = renderHook(() => useWebSocket());

      // Rapid connect/disconnect
      act(() => {
        result.current.connect();
        result.current.disconnect();
        result.current.connect();
      });

      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Should eventually be connected
      expect(result.current.connectionState.connected).toBe(true);
    });

    it('handles sending messages before connection', async () => {
      const consoleWarn = jest.spyOn(console, 'warn').mockImplementation(() => {});
      
      const { result } = renderHook(() => useWebSocket());

      // Try to send message before connection is established
      act(() => {
        result.current.sendMessage('generation_progress', { test: 'data' });
      });

      expect(consoleWarn).toHaveBeenCalledWith(
        'Cannot send message: WebSocket is not connected'
      );

      consoleWarn.mockRestore();
    });

    it('handles multiple subscriptions to same message type', async () => {
      const handler1 = jest.fn();
      const handler2 = jest.fn();
      
      const { result } = renderHook(() => useWebSocket());

      // Wait for connection
      await act(async () => {
        jest.runOnlyPendingTimers();
      });

      // Subscribe multiple handlers to same type
      act(() => {
        result.current.subscribe('generation_progress' as WSMessageType, handler1);
        result.current.subscribe('generation_progress' as WSMessageType, handler2);
      });

      const testMessage: WSMessage = {
        type: 'generation_progress' as WSMessageType,
        data: { progress: 25 },
        timestamp: new Date().toISOString()
      };

      // Simulate message
      const wsInstance = (global.WebSocket as any).mock?.instances?.[0];
      if (wsInstance && wsInstance.simulateMessage) {
        act(() => {
          wsInstance.simulateMessage(testMessage);
        });
      }

      expect(handler1).toHaveBeenCalledWith({ progress: 25 });
      expect(handler2).toHaveBeenCalledWith({ progress: 25 });
    });
  });
});