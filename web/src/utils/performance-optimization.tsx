/**
 * Performance Optimization Utilities for go-starter Web UI
 * Provides comprehensive performance monitoring, optimization, and user experience enhancements
 */

import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';

// Performance metrics and monitoring
export interface PerformanceMetrics {
  // Core Web Vitals
  lcp?: number; // Largest Contentful Paint
  fid?: number; // First Input Delay
  cls?: number; // Cumulative Layout Shift
  fcp?: number; // First Contentful Paint
  ttfb?: number; // Time to First Byte
  
  // Custom metrics
  componentRenderTime?: number;
  apiResponseTime?: number;
  bundleLoadTime?: number;
  memoryUsage?: number;
  
  // User experience metrics
  pageLoadTime?: number;
  interactionDelay?: number;
  scrollPerformance?: number;
}

export interface PerformanceThresholds {
  lcp: { good: number; needsImprovement: number };
  fid: { good: number; needsImprovement: number };
  cls: { good: number; needsImprovement: number };
  renderTime: { good: number; needsImprovement: number };
  apiResponse: { good: number; needsImprovement: number };
}

export const PERFORMANCE_THRESHOLDS: PerformanceThresholds = {
  lcp: { good: 2500, needsImprovement: 4000 },
  fid: { good: 100, needsImprovement: 300 },
  cls: { good: 0.1, needsImprovement: 0.25 },
  renderTime: { good: 16, needsImprovement: 100 }, // 60fps = 16ms
  apiResponse: { good: 1000, needsImprovement: 3000 }
};

// Performance monitoring class
export class PerformanceMonitor {
  private static instance: PerformanceMonitor;
  private metrics: PerformanceMetrics = {};
  private observers: Map<string, PerformanceObserver> = new Map();
  private listeners: Set<(metrics: PerformanceMetrics) => void> = new Set();

  static getInstance(): PerformanceMonitor {
    if (!PerformanceMonitor.instance) {
      PerformanceMonitor.instance = new PerformanceMonitor();
    }
    return PerformanceMonitor.instance;
  }

  initialize(): void {
    this.observeWebVitals();
    this.observeNavigation();
    this.observeResources();
    this.monitorMemory();
  }

  private observeWebVitals(): void {
    // Largest Contentful Paint
    if ('PerformanceObserver' in window) {
      const lcpObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const lastEntry = entries[entries.length - 1] as any;
        this.updateMetric('lcp', lastEntry.startTime);
      });
      
      try {
        lcpObserver.observe({ type: 'largest-contentful-paint', buffered: true });
        this.observers.set('lcp', lcpObserver);
      } catch (e) {
        console.warn('LCP observation not supported');
      }

      // First Input Delay
      const fidObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          this.updateMetric('fid', entry.processingStart - entry.startTime);
        });
      });

      try {
        fidObserver.observe({ type: 'first-input', buffered: true });
        this.observers.set('fid', fidObserver);
      } catch (e) {
        console.warn('FID observation not supported');
      }

      // Cumulative Layout Shift
      const clsObserver = new PerformanceObserver((list) => {
        let clsValue = 0;
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
          }
        });
        this.updateMetric('cls', clsValue);
      });

      try {
        clsObserver.observe({ type: 'layout-shift', buffered: true });
        this.observers.set('cls', clsObserver);
      } catch (e) {
        console.warn('CLS observation not supported');
      }
    }
  }

  private observeNavigation(): void {
    if ('PerformanceObserver' in window) {
      const navObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          this.updateMetric('fcp', entry.firstContentfulPaint);
          this.updateMetric('ttfb', entry.responseStart - entry.requestStart);
          this.updateMetric('pageLoadTime', entry.loadEventEnd - entry.fetchStart);
        });
      });

      try {
        navObserver.observe({ type: 'navigation', buffered: true });
        this.observers.set('navigation', navObserver);
      } catch (e) {
        console.warn('Navigation observation not supported');
      }
    }
  }

  private observeResources(): void {
    if ('PerformanceObserver' in window) {
      const resourceObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          if (entry.name.includes('.js') || entry.name.includes('.css')) {
            const loadTime = entry.responseEnd - entry.startTime;
            this.updateMetric('bundleLoadTime', loadTime);
          }
        });
      });

      try {
        resourceObserver.observe({ type: 'resource', buffered: true });
        this.observers.set('resource', resourceObserver);
      } catch (e) {
        console.warn('Resource observation not supported');
      }
    }
  }

  private monitorMemory(): void {
    if ('memory' in performance) {
      const updateMemory = () => {
        const memory = (performance as any).memory;
        this.updateMetric('memoryUsage', memory.usedJSHeapSize);
      };

      updateMemory();
      setInterval(updateMemory, 10000); // Every 10 seconds
    }
  }

  private updateMetric(key: keyof PerformanceMetrics, value: number): void {
    this.metrics = { ...this.metrics, [key]: value };
    this.notifyListeners();
  }

  private notifyListeners(): void {
    this.listeners.forEach(listener => listener(this.metrics));
  }

  getMetrics(): PerformanceMetrics {
    return { ...this.metrics };
  }

  subscribe(listener: (metrics: PerformanceMetrics) => void): () => void {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  measureComponentRender<T>(
    componentName: string,
    renderFunction: () => T
  ): T {
    const startTime = performance.now();
    const result = renderFunction();
    const endTime = performance.now();
    
    this.updateMetric('componentRenderTime', endTime - startTime);
    
    if (endTime - startTime > PERFORMANCE_THRESHOLDS.renderTime.needsImprovement) {
      console.warn(`Slow render detected in ${componentName}: ${endTime - startTime}ms`);
    }
    
    return result;
  }

  measureApiCall<T>(
    url: string,
    apiCall: () => Promise<T>
  ): Promise<T> {
    const startTime = performance.now();
    
    return apiCall().then(
      (result) => {
        const endTime = performance.now();
        this.updateMetric('apiResponseTime', endTime - startTime);
        return result;
      },
      (error) => {
        const endTime = performance.now();
        this.updateMetric('apiResponseTime', endTime - startTime);
        throw error;
      }
    );
  }

  getPerformanceScore(): number {
    const { lcp, fid, cls } = this.metrics;
    
    if (!lcp || !fid || cls === undefined) {
      return 0; // Can't calculate without core metrics
    }

    // Calculate individual scores (0-100)
    const lcpScore = lcp <= PERFORMANCE_THRESHOLDS.lcp.good ? 100 : 
                     lcp <= PERFORMANCE_THRESHOLDS.lcp.needsImprovement ? 50 : 0;
    
    const fidScore = fid <= PERFORMANCE_THRESHOLDS.fid.good ? 100 :
                     fid <= PERFORMANCE_THRESHOLDS.fid.needsImprovement ? 50 : 0;
    
    const clsScore = cls <= PERFORMANCE_THRESHOLDS.cls.good ? 100 :
                     cls <= PERFORMANCE_THRESHOLDS.cls.needsImprovement ? 50 : 0;

    // Weighted average (LCP and CLS are more important)
    return Math.round((lcpScore * 0.4 + fidScore * 0.3 + clsScore * 0.3));
  }

  dispose(): void {
    this.observers.forEach(observer => observer.disconnect());
    this.observers.clear();
    this.listeners.clear();
  }
}

// React performance hooks
export function usePerformanceMonitoring() {
  const [metrics, setMetrics] = useState<PerformanceMetrics>({});
  const monitor = useMemo(() => PerformanceMonitor.getInstance(), []);

  useEffect(() => {
    monitor.initialize();
    const unsubscribe = monitor.subscribe(setMetrics);
    setMetrics(monitor.getMetrics());
    
    return unsubscribe;
  }, [monitor]);

  return {
    metrics,
    score: monitor.getPerformanceScore(),
    measureRender: monitor.measureComponentRender.bind(monitor),
    measureApi: monitor.measureApiCall.bind(monitor)
  };
}

// Component performance optimization hooks
export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

export function useThrottle<T extends (...args: any[]) => any>(
  callback: T,
  delay: number
): T {
  const throttledRef = useRef<T | null>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const throttledCallback = useCallback((...args: Parameters<T>) => {
    if (!throttledRef.current) {
      throttledRef.current = callback;
      callback(...args);
      
      timeoutRef.current = setTimeout(() => {
        throttledRef.current = undefined;
      }, delay);
    }
  }, [callback, delay]) as T;

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return throttledCallback;
}

export function useVirtualization<T>(
  items: T[],
  itemHeight: number,
  containerHeight: number
): {
  visibleItems: T[];
  startIndex: number;
  endIndex: number;
  totalHeight: number;
  offsetY: number;
} {
  const [scrollTop, setScrollTop] = useState(0);

  const startIndex = Math.floor(scrollTop / itemHeight);
  const endIndex = Math.min(
    startIndex + Math.ceil(containerHeight / itemHeight) + 1,
    items.length - 1
  );

  const visibleItems = items.slice(startIndex, endIndex + 1);
  const totalHeight = items.length * itemHeight;
  const offsetY = startIndex * itemHeight;

  return {
    visibleItems,
    startIndex,
    endIndex,
    totalHeight,
    offsetY
  };
}

// Memoization utilities
export function useStableMemo<T>(
  factory: () => T,
  deps: React.DependencyList
): T {
  const ref = useRef<{ deps: React.DependencyList; value: T } | null>(null);

  if (!ref.current || !areEqual(ref.current.deps, deps)) {
    ref.current = {
      deps: [...deps],
      value: factory()
    };
  }

  return ref.current.value;
}

function areEqual(a: React.DependencyList, b: React.DependencyList): boolean {
  if (a.length !== b.length) return false;
  
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false;
  }
  
  return true;
}

export function useDeepMemo<T>(
  factory: () => T,
  deps: React.DependencyList
): T {
  const ref = useRef<{ deps: React.DependencyList; value: T } | null>(null);

  if (!ref.current || !deepEqual(ref.current.deps, deps)) {
    ref.current = {
      deps: JSON.parse(JSON.stringify(deps)),
      value: factory()
    };
  }

  return ref.current.value;
}

function deepEqual(a: any, b: any): boolean {
  return JSON.stringify(a) === JSON.stringify(b);
}

// Lazy loading utilities
export function useLazyLoading(
  threshold: number = 0.1
): {
  ref: React.RefObject<HTMLElement>;
  isVisible: boolean;
} {
  const ref = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsVisible(entry.isIntersecting);
      },
      { threshold }
    );

    observer.observe(element);

    return () => {
      observer.unobserve(element);
    };
  }, [threshold]);

  return { ref: ref as React.RefObject<HTMLElement>, isVisible };
}

// Image optimization
export interface OptimizedImageProps {
  src: string;
  alt: string;
  width?: number;
  height?: number;
  lazy?: boolean;
  placeholder?: string;
  className?: string;
  onLoad?: () => void;
  onError?: () => void;
}

export const OptimizedImage: React.FC<OptimizedImageProps> = ({
  src,
  alt,
  width,
  height,
  lazy = true,
  placeholder,
  className = '',
  onLoad,
  onError
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [hasError, setHasError] = useState(false);
  const { ref, isVisible } = useLazyLoading();

  const shouldLoad = !lazy || isVisible;

  const handleLoad = useCallback(() => {
    setIsLoaded(true);
    onLoad?.();
  }, [onLoad]);

  const handleError = useCallback(() => {
    setHasError(true);
    onError?.();
  }, [onError]);

  return (
    <div 
      ref={ref as React.Ref<HTMLDivElement>}
      className={`relative overflow-hidden ${className}`}
      style={{ width, height }}
    >
      {placeholder && !isLoaded && !hasError && (
        <div 
          className="absolute inset-0 bg-gray-200 animate-pulse"
          aria-hidden="true"
        />
      )}
      
      {shouldLoad && !hasError && (
        <img
          src={src}
          alt={alt}
          width={width}
          height={height}
          className={`transition-opacity duration-300 ${
            isLoaded ? 'opacity-100' : 'opacity-0'
          }`}
          onLoad={handleLoad}
          onError={handleError}
          loading={lazy ? 'lazy' : 'eager'}
        />
      )}
      
      {hasError && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100 text-gray-500">
          <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
      )}
    </div>
  );
};

// Code splitting utilities
export function createLazyComponent<T extends React.ComponentType<any>>(
  importFunc: () => Promise<{ default: T }>,
  fallback?: React.ComponentType
) {
  const LazyComponent = React.lazy(importFunc);
  
  return React.forwardRef<any, React.ComponentProps<T>>((props, ref) => (
    <React.Suspense fallback={fallback ? React.createElement(fallback) : <div>Loading...</div>}>
      <LazyComponent {...props as any} />
    </React.Suspense>
  ));
}

// Bundle analysis utilities
export class BundleAnalyzer {
  static analyzeBundleSize(): void {
    if (process.env.NODE_ENV === 'development') {
      console.log('Bundle analysis feature is disabled in Vite environment');
    }
  }

  static logChunkSizes(): void {
    if ('performance' in window && 'getEntriesByType' in performance) {
      const resources = performance.getEntriesByType('resource') as any[];
      const jsResources = resources.filter(resource => 
        resource.name.includes('.js') && !resource.name.includes('node_modules')
      );
      
      console.group('📦 JavaScript Bundle Sizes');
      jsResources.forEach(resource => {
        const sizeKB = resource.transferSize / 1024;
        console.log(`${resource.name}: ${sizeKB.toFixed(2)} KB`);
      });
      console.groupEnd();
    }
  }
}

// Memory management utilities
export class MemoryManager {
  private static cleanupTasks: Set<() => void> = new Set();

  static addCleanupTask(task: () => void): () => void {
    this.cleanupTasks.add(task);
    return () => this.cleanupTasks.delete(task);
  }

  static runCleanup(): void {
    this.cleanupTasks.forEach(task => {
      try {
        task();
      } catch (error) {
        console.warn('Cleanup task failed:', error);
      }
    });
    this.cleanupTasks.clear();
  }

  static monitorMemoryPressure(): void {
    if ('memory' in performance) {
      const checkMemory = () => {
        const memory = (performance as any).memory;
        const usagePercent = (memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100;
        
        if (usagePercent > 80) {
          console.warn(`High memory usage: ${usagePercent.toFixed(2)}%`);
          this.runCleanup();
        }
      };

      setInterval(checkMemory, 30000); // Check every 30 seconds
    }
  }
}

// Performance testing utilities
export class PerformanceTester {
  static async runLighthouseAudit(): Promise<void> {
    if (process.env.NODE_ENV === 'development') {
      console.log('Lighthouse audit feature is disabled in Vite environment');
    }
  }

  static measurePageSpeed(): PerformanceMetrics {
    const navigation = performance.getEntriesByType('navigation')[0] as any;
    
    return {
      pageLoadTime: navigation.loadEventEnd - navigation.fetchStart,
      ttfb: navigation.responseStart - navigation.requestStart,
      fcp: navigation.firstContentfulPaint,
    };
  }

  static startPerformanceTest(testName: string): () => void {
    const startTime = performance.now();
    const startMark = `${testName}-start`;
    const endMark = `${testName}-end`;
    const measureName = `${testName}-duration`;
    
    performance.mark(startMark);
    
    return () => {
      performance.mark(endMark);
      performance.measure(measureName, startMark, endMark);
      
      const measure = performance.getEntriesByName(measureName)[0];
      console.log(`⏱️ ${testName}: ${measure.duration.toFixed(2)}ms`);
    };
  }
}

// Initialize performance monitoring
if (typeof window !== 'undefined') {
  MemoryManager.monitorMemoryPressure();
  BundleAnalyzer.logChunkSizes();
}