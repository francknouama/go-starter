/**
 * API Service Layer for go-starter Web UI
 * Provides integration with the backend API endpoints
 */

// Base configuration
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

// API Response Types
export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  timestamp: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
  uptime: string;
  checks: {
    database: string;
    memory: string;
    disk: string;
  };
}

export interface Blueprint {
  id: string;
  name: string;
  description: string;
  category: string;
  difficulty: string;
  features: string[];
  architectures: string[];
  frameworks: string[];
  estimatedFiles: number;
  tags: string[];
}

// Import ProjectConfig from types
export type { ProjectConfig } from '../types'

export interface PreviewResponse {
  fileStructure: FileNode[];
  estimatedSize: string;
  fileCount: number;
}

export interface FileNode {
  name: string;
  type: 'file' | 'directory';
  path: string;
  size?: number;
  children?: FileNode[];
  content?: string;
}

export interface GenerationRequest extends ProjectConfig {
  outputFormat: 'zip' | 'preview';
  includeTests: boolean;
  includeDocs: boolean;
}

export interface GenerationResponse {
  projectId: string;
  downloadUrl?: string;
  previewUrl?: string;
  fileCount: number;
  estimatedSize: string;
}

// HTTP Client with error handling
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...defaultHeaders,
          ...options.headers,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error);
      throw error;
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  async downloadBlob(endpoint: string): Promise<Blob> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new Error(`Download failed: ${response.status} ${response.statusText}`);
    }
    
    return response.blob();
  }
}

// API Service Instance
const apiClient = new ApiClient();

// Health & System APIs
export const healthApi = {
  async getHealth(): Promise<HealthResponse> {
    const response = await apiClient.get<HealthResponse>('/health');
    return response.data!;
  },

  async getSimpleHealth(): Promise<{ status: string }> {
    const response = await apiClient.get<{ status: string }>('/health/simple');
    return response.data!;
  },

  async getMetrics(): Promise<Record<string, unknown>> {
    const response = await apiClient.get<Record<string, unknown>>('/metrics');
    return response.data!;
  },

  async getStatus(): Promise<{ status: string }> {
    const response = await apiClient.get<{ status: string }>('/status');
    return response.data!;
  },
};

// Configuration APIs
export const configApi = {
  async getDefaultConfig(): Promise<ProjectConfig> {
    const response = await apiClient.get<ProjectConfig>('/config');
    return response.data!;
  },

  async getProjectTypeDetails(type: string): Promise<{ type: string; description: string; features?: string[] }> {
    const response = await apiClient.get<{ type: string; description: string; features?: string[] }>(`/config/types/${type}`);
    return response.data!;
  },

  async getFrameworks(): Promise<string[]> {
    const response = await apiClient.get<string[]>('/config/frameworks');
    return response.data!;
  },

  async getArchitectures(): Promise<string[]> {
    const response = await apiClient.get<string[]>('/config/architectures');
    return response.data!;
  },
};

// Blueprint APIs
export const blueprintApi = {
  async getBlueprints(filters?: {
    category?: string;
    difficulty?: string;
    framework?: string;
    architecture?: string;
  }): Promise<Blueprint[]> {
    let endpoint = '/blueprints';
    
    if (filters) {
      const params = new URLSearchParams();
      Object.entries(filters).forEach(([key, value]) => {
        if (value) params.append(key, value);
      });
      
      if (params.toString()) {
        endpoint += `?${params.toString()}`;
      }
    }
    
    const response = await apiClient.get<Blueprint[]>(endpoint);
    return response.data!;
  },

  async getBlueprintById(id: string): Promise<Blueprint> {
    const response = await apiClient.get<Blueprint>(`/blueprints/${id}`);
    return response.data!;
  },

  async getBlueprintsByCategory(category: string): Promise<Blueprint[]> {
    const response = await apiClient.get<Blueprint[]>(`/blueprints/category/${category}`);
    return response.data!;
  },

  async validateBlueprintConfig(id: string, config: ProjectConfig): Promise<{ valid: boolean; errors?: string[] }> {
    const response = await apiClient.post<{ valid: boolean; errors?: string[] }>(`/blueprints/${id}/validate`, config);
    return response.data || { valid: true };
  },
};

// Project Generation APIs
export const projectApi = {
  async generatePreview(config: ProjectConfig): Promise<PreviewResponse> {
    const response = await apiClient.post<PreviewResponse>('/preview', config);
    return response.data!;
  },

  async generateProject(request: GenerationRequest): Promise<GenerationResponse> {
    const response = await apiClient.post<GenerationResponse>('/generate', request);
    return response.data!;
  },

  async generateAndDownload(request: GenerationRequest): Promise<Blob> {
    return apiClient.downloadBlob('/generate/download');
  },

  async downloadProject(token: string): Promise<Blob> {
    return apiClient.downloadBlob(`/download/${token}`);
  },

  async getDownloadStatus(token: string): Promise<{ status: string; progress?: number }> {
    const response = await apiClient.get<{ status: string; progress?: number }>(`/download/${token}/status`);
    return response.data || { status: 'completed' };
  },
};

// WebSocket Message Types
export interface WebSocketMessage {
  type: string;
  data: unknown;
}

// WebSocket Service
export class WebSocketService {
  private ws: WebSocket | null = null;
  private listeners: Map<string, Set<(data: unknown) => void>> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/v1/ws`;
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          this.reconnectDelay = 1000;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            this.handleMessage(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          this.handleReconnect();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  send(message: WebSocketMessage): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket not connected, message not sent:', message);
    }
  }

  subscribe(event: string, callback: (data: unknown) => void): () => void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event)!.add(callback);

    // Return unsubscribe function
    return () => {
      const eventListeners = this.listeners.get(event);
      if (eventListeners) {
        eventListeners.delete(callback);
        if (eventListeners.size === 0) {
          this.listeners.delete(event);
        }
      }
    };
  }

  private handleMessage(message: WebSocketMessage): void {
    const { type, data } = message;
    const eventListeners = this.listeners.get(type);
    
    if (eventListeners) {
      eventListeners.forEach(callback => {
        try {
          callback(data);
        } catch (error) {
          console.error('Error in WebSocket message handler:', error);
        }
      });
    }
  }

  private handleReconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      
      setTimeout(() => {
        console.log(`Attempting to reconnect WebSocket (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        this.connect().catch(error => {
          console.error('WebSocket reconnection failed:', error);
        });
      }, this.reconnectDelay);

      // Exponential backoff
      this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
    } else {
      console.error('Max WebSocket reconnection attempts reached');
    }
  }
}

// Export singleton WebSocket service
export const wsService = new WebSocketService();

// Convenience function to initialize WebSocket connection
export const connectWebSocket = () => wsService.connect();

// Export all API groups
export const api = {
  health: healthApi,
  config: configApi,
  blueprints: blueprintApi,
  projects: projectApi,
  ws: wsService,
};

export default api;