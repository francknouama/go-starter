/**
 * React Hooks for API Integration
 * Provides easy-to-use hooks for all backend API operations
 */

import { useState, useEffect, useCallback } from 'react';
import { 
  Blueprint, 
  ProjectConfig, 
  PreviewResponse, 
  GenerationRequest,
  GenerationResponse,
  HealthResponse 
} from '../services/api';
import { mockApi as api } from '../services/mock-api';
// import { useErrorReporting } from '../components/common/ErrorBoundary';

// Simple error reporting stub
const useErrorReporting = () => ({
  reportNetworkError: (context: string, data?: unknown, message?: string) => {
    console.error(`[${context}] ${message}`, data);
  }
});

// Generic API hook for handling async operations
export function useAsyncOperation<T>() {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { reportNetworkError } = useErrorReporting();

  const execute = useCallback(async (operation: () => Promise<T>) => {
    try {
      setLoading(true);
      setError(null);
      const result = await operation();
      setData(result);
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'An error occurred';
      setError(errorMessage);
      reportNetworkError('API operation', undefined, errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [reportNetworkError]);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  return { data, loading, error, execute, reset };
}

// Health API hooks
export function useHealth() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const checkHealth = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const healthData = await api.health.getHealth();
      setHealth(healthData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Health check failed');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    checkHealth();
    
    // Check health every 30 seconds
    const interval = setInterval(checkHealth, 30000);
    return () => clearInterval(interval);
  }, [checkHealth]);

  return { health, loading, error, refetch: checkHealth };
}

// Configuration API hooks
export function useDefaultConfig() {
  const { data, loading, error, execute } = useAsyncOperation<ProjectConfig>();

  useEffect(() => {
    execute(() => api.config.getDefaultConfig());
  }, [execute]);

  return { config: data, loading, error };
}

export function useFrameworks() {
  const { data, loading, error, execute } = useAsyncOperation<string[]>();

  useEffect(() => {
    execute(() => api.config.getFrameworks());
  }, [execute]);

  return { frameworks: data || [], loading, error };
}

export function useArchitectures() {
  const { data, loading, error, execute } = useAsyncOperation<string[]>();

  useEffect(() => {
    execute(() => api.config.getArchitectures());
  }, [execute]);

  return { architectures: data || [], loading, error };
}

// Blueprint API hooks
export function useBlueprints(filters?: {
  category?: string;
  difficulty?: string;
  framework?: string;
  architecture?: string;
}) {
  const { data, loading, error, execute } = useAsyncOperation<Blueprint[]>();

  const fetchBlueprints = useCallback(() => {
    return execute(() => api.blueprints.getBlueprints(filters));
  }, [execute, filters]);

  useEffect(() => {
    fetchBlueprints();
  }, [fetchBlueprints]);

  return { 
    blueprints: data || [], 
    loading, 
    error, 
    refetch: fetchBlueprints 
  };
}

export function useBlueprint(id: string | null) {
  const { data, loading, error, execute } = useAsyncOperation<Blueprint>();

  useEffect(() => {
    if (id) {
      execute(() => api.blueprints.getBlueprintById(id));
    }
  }, [id, execute]);

  return { blueprint: data, loading, error };
}

export function useBlueprintValidation() {
  const { data, loading, error, execute } = useAsyncOperation<{ valid: boolean; errors?: string[] }>();

  const validateConfig = useCallback((blueprintId: string, config: ProjectConfig) => {
    return execute(() => api.blueprints.validateBlueprintConfig(blueprintId, config));
  }, [execute]);

  return { 
    validation: data, 
    loading, 
    error, 
    validateConfig 
  };
}

// Project generation hooks
export function useProjectPreview() {
  const { data, loading, error, execute } = useAsyncOperation<PreviewResponse>();

  const generatePreview = useCallback((config: ProjectConfig) => {
    return execute(() => api.projects.generatePreview(config));
  }, [execute]);

  return { 
    preview: data, 
    loading, 
    error, 
    generatePreview 
  };
}

export function useProjectGeneration() {
  const { data, loading, error, execute } = useAsyncOperation<GenerationResponse>();

  const generateProject = useCallback((request: GenerationRequest) => {
    return execute(() => api.projects.generateProject(request));
  }, [execute]);

  return { 
    generation: data, 
    loading, 
    error, 
    generateProject 
  };
}

export function useProjectDownload() {
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { reportNetworkError } = useErrorReporting();

  const downloadProject = useCallback(async (request: GenerationRequest, filename?: string) => {
    try {
      setDownloading(true);
      setError(null);
      
      const blob = await api.projects.generateAndDownload(request);
      
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename || `${request.projectName || 'project'}.zip`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Download failed';
      setError(errorMessage);
      reportNetworkError('project download', undefined, errorMessage);
    } finally {
      setDownloading(false);
    }
  }, [reportNetworkError]);

  const downloadByToken = useCallback(async (token: string, filename?: string) => {
    try {
      setDownloading(true);
      setError(null);
      
      const blob = await api.projects.downloadProject(token);
      
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename || 'project.zip';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Download failed';
      setError(errorMessage);
      reportNetworkError('token download', undefined, errorMessage);
    } finally {
      setDownloading(false);
    }
  }, [reportNetworkError]);

  return { 
    downloading, 
    error, 
    downloadProject, 
    downloadByToken 
  };
}

// WebSocket hooks
export function useWebSocket() {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const connect = async () => {
      try {
        await api.ws.connect();
        if (mounted) {
          setConnected(true);
          setError(null);
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err.message : 'WebSocket connection failed');
          setConnected(false);
        }
      }
    };

    connect();

    return () => {
      mounted = false;
      api.ws.disconnect();
      setConnected(false);
    };
  }, []);

  const subscribe = useCallback((event: string, callback: (data: unknown) => void) => {
    return api.ws.subscribe(event, callback);
  }, []);

  const send = useCallback((message: unknown) => {
    api.ws.send(message);
  }, []);

  return { connected, error, subscribe, send };
}

export function useWebSocketEvent<T>(event: string) {
  const [data, setData] = useState<T | null>(null);
  const { subscribe } = useWebSocket();

  useEffect(() => {
    const unsubscribe = subscribe(event, (eventData: T) => {
      setData(eventData);
    });

    return unsubscribe;
  }, [event, subscribe]);

  return data;
}

// Real-time project generation status
export function useGenerationProgress() {
  const progressData = useWebSocketEvent<{
    projectId: string;
    progress: number;
    stage: string;
    message: string;
    completed: boolean;
    error?: string;
  }>('generation_progress');

  return progressData;
}

// Real-time system notifications
export function useSystemNotifications() {
  const notification = useWebSocketEvent<{
    type: 'info' | 'warning' | 'error' | 'success';
    title: string;
    message: string;
    timestamp: string;
  }>('system_notification');

  return notification;
}

// Combined hook for complete project workflow
export function useProjectWorkflow() {
  const { config: defaultConfig, loading: configLoading } = useDefaultConfig();
  const { blueprints, loading: blueprintsLoading } = useBlueprints();
  const { generatePreview, preview, loading: previewLoading } = useProjectPreview();
  const { generateProject, generation, loading: generationLoading } = useProjectGeneration();
  const { downloadProject, downloading } = useProjectDownload();
  const generationProgress = useGenerationProgress();

  const isLoading = configLoading || blueprintsLoading || previewLoading || generationLoading || downloading;

  return {
    // Data
    defaultConfig,
    blueprints,
    preview,
    generation,
    generationProgress,
    
    // Loading states
    isLoading,
    configLoading,
    blueprintsLoading,
    previewLoading,
    generationLoading,
    downloading,
    
    // Actions
    generatePreview,
    generateProject,
    downloadProject,
  };
}

// Export WebSocket connection function
export const connectWebSocket = () => api.ws.connect()