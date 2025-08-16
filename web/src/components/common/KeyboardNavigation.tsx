/**
 * Enhanced Keyboard Navigation System for WCAG 2.1 AA Compliance
 * Provides comprehensive keyboard shortcuts and navigation patterns
 */

import React, { useEffect, useCallback, useRef, useState } from 'react';
import { keyboardNavigation, ScreenReaderOnly } from './FocusManagement';

// Keyboard shortcut definitions
export interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  altKey?: boolean;
  shiftKey?: boolean;
  metaKey?: boolean;
  description: string;
  action: () => void;
  disabled?: boolean;
  scope?: 'global' | 'modal' | 'form';
}

// Global keyboard shortcuts manager
export class KeyboardShortcutManager {
  private shortcuts: Map<string, KeyboardShortcut> = new Map();
  private activeScope: string = 'global';

  register(id: string, shortcut: KeyboardShortcut) {
    this.shortcuts.set(id, shortcut);
  }

  unregister(id: string) {
    this.shortcuts.delete(id);
  }

  setScope(scope: string) {
    this.activeScope = scope;
  }

  private getShortcutKey(shortcut: KeyboardShortcut): string {
    const parts = [];
    if (shortcut.ctrlKey) parts.push('ctrl');
    if (shortcut.altKey) parts.push('alt');
    if (shortcut.shiftKey) parts.push('shift');
    if (shortcut.metaKey) parts.push('meta');
    parts.push(shortcut.key.toLowerCase());
    return parts.join('+');
  }

  handleKeyDown = (event: KeyboardEvent) => {
    const shortcutKey = [
      event.ctrlKey && 'ctrl',
      event.altKey && 'alt',
      event.shiftKey && 'shift',
      event.metaKey && 'meta',
      event.key.toLowerCase()
    ].filter(Boolean).join('+');

    for (const [id, shortcut] of this.shortcuts) {
      if (shortcut.disabled) continue;
      if (shortcut.scope && shortcut.scope !== this.activeScope) continue;

      const expectedKey = this.getShortcutKey(shortcut);
      if (expectedKey === shortcutKey) {
        event.preventDefault();
        shortcut.action();
        break;
      }
    }
  };

  getShortcuts(): Array<{ id: string; shortcut: KeyboardShortcut }> {
    return Array.from(this.shortcuts.entries()).map(([id, shortcut]) => ({
      id,
      shortcut
    }));
  }
}

// Global instance
export const globalShortcutManager = new KeyboardShortcutManager();

// Hook for using keyboard shortcuts
export const useKeyboardShortcuts = (
  shortcuts: Array<{ id: string; shortcut: KeyboardShortcut }>,
  scope: string = 'global'
) => {
  useEffect(() => {
    // Register shortcuts
    shortcuts.forEach(({ id, shortcut }) => {
      globalShortcutManager.register(id, { ...shortcut, scope: scope as 'global' | 'modal' | 'form' });
    });

    // Set scope
    globalShortcutManager.setScope(scope);

    // Add event listener
    document.addEventListener('keydown', globalShortcutManager.handleKeyDown);

    return () => {
      // Unregister shortcuts
      shortcuts.forEach(({ id }) => {
        globalShortcutManager.unregister(id);
      });

      // Remove event listener
      document.removeEventListener('keydown', globalShortcutManager.handleKeyDown);
    };
  }, [shortcuts, scope]);
};

// Common keyboard navigation patterns
export const KeyboardNavigationPatterns = {
  // Arrow key navigation for grids
  useGridNavigation: (
    rows: number,
    cols: number,
    onFocusChange: (row: number, col: number) => void
  ) => {
    const [currentRow, setCurrentRow] = useState(0);
    const [currentCol, setCurrentCol] = useState(0);

    const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
      let newRow = currentRow;
      let newCol = currentCol;

      switch (event.key) {
        case 'ArrowUp':
          event.preventDefault();
          newRow = Math.max(0, currentRow - 1);
          break;
        case 'ArrowDown':
          event.preventDefault();
          newRow = Math.min(rows - 1, currentRow + 1);
          break;
        case 'ArrowLeft':
          event.preventDefault();
          newCol = Math.max(0, currentCol - 1);
          break;
        case 'ArrowRight':
          event.preventDefault();
          newCol = Math.min(cols - 1, currentCol + 1);
          break;
        case 'Home':
          event.preventDefault();
          if (event.ctrlKey) {
            newRow = 0;
            newCol = 0;
          } else {
            newCol = 0;
          }
          break;
        case 'End':
          event.preventDefault();
          if (event.ctrlKey) {
            newRow = rows - 1;
            newCol = cols - 1;
          } else {
            newCol = cols - 1;
          }
          break;
        default:
          return;
      }

      if (newRow !== currentRow || newCol !== currentCol) {
        setCurrentRow(newRow);
        setCurrentCol(newCol);
        onFocusChange(newRow, newCol);
      }
    }, [currentRow, currentCol, rows, cols, onFocusChange]);

    return {
      currentRow,
      currentCol,
      handleKeyDown,
      setPosition: (row: number, col: number) => {
        setCurrentRow(row);
        setCurrentCol(col);
        onFocusChange(row, col);
      }
    };
  },

  // Tab panel navigation
  useTabNavigation: (
    tabIds: string[],
    onTabChange: (tabId: string) => void
  ) => {
    const [activeTabIndex, setActiveTabIndex] = useState(0);

    const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
      let newIndex = activeTabIndex;

      switch (event.key) {
        case 'ArrowLeft':
        case 'ArrowUp':
          event.preventDefault();
          newIndex = activeTabIndex > 0 ? activeTabIndex - 1 : tabIds.length - 1;
          break;
        case 'ArrowRight':
        case 'ArrowDown':
          event.preventDefault();
          newIndex = activeTabIndex < tabIds.length - 1 ? activeTabIndex + 1 : 0;
          break;
        case 'Home':
          event.preventDefault();
          newIndex = 0;
          break;
        case 'End':
          event.preventDefault();
          newIndex = tabIds.length - 1;
          break;
        default:
          return;
      }

      if (newIndex !== activeTabIndex) {
        setActiveTabIndex(newIndex);
        onTabChange(tabIds[newIndex]);
      }
    }, [activeTabIndex, tabIds, onTabChange]);

    return {
      activeTabIndex,
      handleKeyDown,
      setActiveTab: (tabId: string) => {
        const index = tabIds.indexOf(tabId);
        if (index >= 0) {
          setActiveTabIndex(index);
          onTabChange(tabId);
        }
      }
    };
  },

  // List navigation with type-ahead
  useListNavigation: (
    items: string[],
    onItemSelect: (item: string, index: number) => void
  ) => {
    const [focusedIndex, setFocusedIndex] = useState(0);
    const [typeAheadBuffer, setTypeAheadBuffer] = useState('');
    const typeAheadTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
      switch (event.key) {
        case 'ArrowUp':
          event.preventDefault();
          setFocusedIndex(prev => prev > 0 ? prev - 1 : items.length - 1);
          break;
        case 'ArrowDown':
          event.preventDefault();
          setFocusedIndex(prev => prev < items.length - 1 ? prev + 1 : 0);
          break;
        case 'Home':
          event.preventDefault();
          setFocusedIndex(0);
          break;
        case 'End':
          event.preventDefault();
          setFocusedIndex(items.length - 1);
          break;
        case 'Enter':
        case ' ':
          event.preventDefault();
          onItemSelect(items[focusedIndex], focusedIndex);
          break;
        default:
          // Type-ahead functionality
          if (event.key.length === 1 && !event.ctrlKey && !event.altKey && !event.metaKey) {
            event.preventDefault();
            
            const newBuffer = typeAheadBuffer + event.key.toLowerCase();
            setTypeAheadBuffer(newBuffer);
            
            // Find matching item
            const matchingIndex = items.findIndex((item, index) => 
              index >= focusedIndex && item.toLowerCase().startsWith(newBuffer)
            );
            
            if (matchingIndex >= 0) {
              setFocusedIndex(matchingIndex);
            } else {
              // Search from beginning
              const fallbackIndex = items.findIndex(item => 
                item.toLowerCase().startsWith(newBuffer)
              );
              if (fallbackIndex >= 0) {
                setFocusedIndex(fallbackIndex);
              }
            }
            
            // Clear buffer after delay
            if (typeAheadTimeoutRef.current) {
              clearTimeout(typeAheadTimeoutRef.current);
            }
            typeAheadTimeoutRef.current = setTimeout(() => {
              setTypeAheadBuffer('');
            }, 1000);
          }
          break;
      }
    }, [focusedIndex, items, onItemSelect, typeAheadBuffer]);

    useEffect(() => {
      return () => {
        if (typeAheadTimeoutRef.current) {
          clearTimeout(typeAheadTimeoutRef.current);
        }
      };
    }, []);

    return {
      focusedIndex,
      handleKeyDown,
      setFocusedIndex
    };
  }
};

// Keyboard shortcuts help overlay
export interface KeyboardShortcutsHelpProps {
  isVisible: boolean;
  onClose: () => void;
  className?: string;
}

export const KeyboardShortcutsHelp: React.FC<KeyboardShortcutsHelpProps> = ({
  isVisible,
  onClose,
  className = ''
}) => {
  const shortcuts = globalShortcutManager.getShortcuts();

  const formatShortcut = (shortcut: KeyboardShortcut): string => {
    const parts = [];
    if (shortcut.ctrlKey) parts.push('Ctrl');
    if (shortcut.altKey) parts.push('Alt');
    if (shortcut.shiftKey) parts.push('Shift');
    if (shortcut.metaKey) parts.push('Cmd');
    parts.push(shortcut.key);
    return parts.join(' + ');
  };

  useKeyboardShortcuts([
    {
      id: 'close-help',
      shortcut: {
        key: 'Escape',
        description: 'Close help',
        action: onClose
      }
    }
  ], 'modal');

  if (!isVisible) return null;

  return (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 ${className}`}
      role="dialog"
      aria-modal="true"
      aria-labelledby="shortcuts-title"
    >
      <div className="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 max-h-96 overflow-auto">
        <div className="flex items-center justify-between p-4 border-b">
          <h2 id="shortcuts-title" className="text-lg font-semibold">
            Keyboard Shortcuts
          </h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded"
            aria-label="Close keyboard shortcuts help"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </button>
        </div>
        
        <div className="p-4">
          {shortcuts.length === 0 ? (
            <p className="text-gray-500">No keyboard shortcuts available.</p>
          ) : (
            <div className="space-y-3">
              {shortcuts.map(({ id, shortcut }) => (
                <div key={id} className="flex justify-between items-center">
                  <span className="text-sm text-gray-700">{shortcut.description}</span>
                  <kbd className="px-2 py-1 bg-gray-100 border rounded text-xs font-mono">
                    {formatShortcut(shortcut)}
                  </kbd>
                </div>
              ))}
            </div>
          )}
          
          <div className="mt-6 pt-4 border-t">
            <h3 className="text-sm font-semibold mb-2">General Navigation</h3>
            <div className="space-y-2 text-sm text-gray-600">
              <div className="flex justify-between">
                <span>Tab / Shift+Tab</span>
                <span>Navigate between focusable elements</span>
              </div>
              <div className="flex justify-between">
                <span>Enter / Space</span>
                <span>Activate buttons and links</span>
              </div>
              <div className="flex justify-between">
                <span>Escape</span>
                <span>Close dialogs and menus</span>
              </div>
              <div className="flex justify-between">
                <span>Arrow Keys</span>
                <span>Navigate within components</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Hook for common form keyboard navigation
export const useFormKeyboardNavigation = () => {
  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    const target = event.target as HTMLElement;
    
    switch (event.key) {
      case 'Enter':
        // Submit form if on submit button or single-line input
        if (target.tagName === 'BUTTON' && target.getAttribute('type') !== 'button') {
          // Let the default behavior handle form submission
          return;
        }
        if (target.tagName === 'INPUT' && target.getAttribute('type') !== 'textarea') {
          const form = target.closest('form');
          if (form) {
            const submitButton = form.querySelector<HTMLButtonElement>('button[type="submit"], input[type="submit"]');
            if (submitButton) {
              event.preventDefault();
              submitButton.click();
            }
          }
        }
        break;
        
      case 'Escape':
        // Clear input or close form
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') {
          (target as HTMLInputElement).blur();
        }
        break;
    }
  }, []);

  return { handleKeyDown };
};

// Accessible button component with keyboard support
export interface AccessibleButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  children: React.ReactNode;
}

export const AccessibleButton: React.FC<AccessibleButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  children,
  className = '',
  onKeyDown,
  ...props
}) => {
  const handleKeyDown = useCallback((event: React.KeyboardEvent<HTMLButtonElement>) => {
    // Allow space key to activate button (in addition to Enter)
    if (event.key === ' ') {
      event.preventDefault();
      (event.target as HTMLButtonElement).click();
    }
    
    onKeyDown?.(event);
  }, [onKeyDown]);

  const baseClasses = 'inline-flex items-center justify-center font-medium rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';
  
  const variantClasses = {
    primary: 'bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500',
    secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900 focus:ring-gray-500',
    ghost: 'bg-transparent hover:bg-gray-100 text-gray-700 focus:ring-gray-500'
  };
  
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  };

  return (
    <button
      className={`${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      onKeyDown={handleKeyDown}
      aria-pressed={props['aria-pressed']}
      disabled={loading || props.disabled}
      {...props}
    >
      {loading && (
        <>
          <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <ScreenReaderOnly>Loading...</ScreenReaderOnly>
        </>
      )}
      {children}
    </button>
  );
};

// Quick access navigation component
export interface QuickAccessNavigationProps {
  className?: string;
}

export const QuickAccessNavigation: React.FC<QuickAccessNavigationProps> = ({
  className = ''
}) => {
  const quickLinks = [
    { href: '#main-content', label: 'Skip to main content' },
    { href: '#navigation', label: 'Skip to navigation' },
    { href: '#search', label: 'Skip to search' },
    { href: '#footer', label: 'Skip to footer' }
  ];

  return (
    <nav className={`sr-only focus-within:not-sr-only ${className}`} aria-label="Quick access navigation">
      <ul className="flex space-x-2 p-2 bg-blue-600">
        {quickLinks.map((link) => (
          <li key={link.href}>
            <a
              href={link.href}
              className="inline-block px-3 py-1 bg-white text-blue-600 rounded text-sm font-medium focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-blue-600"
            >
              {link.label}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
};