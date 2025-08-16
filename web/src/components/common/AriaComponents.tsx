/**
 * Comprehensive ARIA Components for WCAG 2.1 AA Screen Reader Optimization
 * Provides semantic HTML elements with proper ARIA roles, states, and properties
 */

import React, { useState, useEffect, useRef } from 'react';
import type { ReactNode } from 'react';

// Simple stub for ScreenReaderOnly component
const ScreenReaderOnly: React.FC<{ children: ReactNode }> = ({ children }) => (
  <span className="sr-only">{children}</span>
);

// ARIA Landmarks for better navigation
export interface LandmarkProps {
  children: ReactNode;
  className?: string;
  label?: string;
  describedBy?: string;
  id?: string;
}

export const MainLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label,
  describedBy,
  id 
}) => (
  <main 
    id={id}
    className={className}
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </main>
);

export const NavigationLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label = 'Main navigation',
  describedBy,
  id 
}) => (
  <nav 
    id={id}
    className={className}
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </nav>
);

export const BannerLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label,
  describedBy,
  id 
}) => (
  <header 
    id={id}
    className={className}
    role="banner"
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </header>
);

export const ContentInfoLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label,
  describedBy,
  id 
}) => (
  <footer 
    id={id}
    className={className}
    role="contentinfo"
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </footer>
);

export const ComplementaryLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label,
  describedBy,
  id 
}) => (
  <aside 
    id={id}
    className={className}
    role="complementary"
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </aside>
);

export const SearchLandmark: React.FC<LandmarkProps> = ({ 
  children, 
  className = '', 
  label = 'Search',
  describedBy 
}) => (
  <section 
    className={className}
    role="search"
    aria-label={label}
    aria-describedby={describedBy}
  >
    {children}
  </section>
);

// ARIA Status and Alert Components
export interface StatusProps {
  children: ReactNode;
  className?: string;
  type?: 'status' | 'alert' | 'log' | 'marquee';
  polite?: boolean;
  atomic?: boolean;
}

export const StatusRegion: React.FC<StatusProps> = ({
  children,
  className = '',
  type = 'status',
  polite = true,
  atomic = true
}) => (
  <div
    className={className}
    role={type}
    aria-live={polite ? 'polite' : 'assertive'}
    aria-atomic={atomic}
  >
    {children}
  </div>
);

export const AlertRegion: React.FC<StatusProps> = ({
  children,
  className = '',
  atomic = true
}) => (
  <div
    className={className}
    role="alert"
    aria-live="assertive"
    aria-atomic={atomic}
  >
    {children}
  </div>
);

// ARIA Describedby Helper Component
export interface DescribedByProps {
  id: string;
  children: ReactNode;
  className?: string;
}

export const Description: React.FC<DescribedByProps> = ({
  id,
  children,
  className = ''
}) => (
  <div id={id} className={`${className}`}>
    {children}
  </div>
);

// Enhanced Form Components with ARIA support
export interface FormFieldProps {
  id: string;
  label: string;
  children: ReactNode;
  description?: string;
  error?: string;
  required?: boolean;
  className?: string;
}

export const FormField: React.FC<FormFieldProps> = ({
  id,
  label,
  children,
  description,
  error,
  required = false,
  className = ''
}) => {
  const descriptionId = description ? `${id}-description` : undefined;
  const errorId = error ? `${id}-error` : undefined;
  const describedBy = [descriptionId, errorId].filter(Boolean).join(' ');

  return (
    <div className={className}>
      <label 
        htmlFor={id}
        className="block text-sm font-medium text-gray-700 mb-1"
      >
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-label="required">
            *
          </span>
        )}
      </label>
      
      {description && (
        <Description id={descriptionId!} className="text-sm text-gray-600 mb-2">
          {description}
        </Description>
      )}
      
      {React.cloneElement(children as React.ReactElement<any>, {
        'aria-describedby': describedBy || undefined,
        'aria-required': required,
        'aria-invalid': !!error,
        className: `${(children as React.ReactElement<any>).props.className || ''} ${
          error ? 'border-red-500 focus:ring-red-500' : ''
        }`
      })}
      
      {error && (
        <Description id={errorId!} className="text-sm text-red-600 mt-1">
          <span role="alert">{error}</span>
        </Description>
      )}
    </div>
  );
};

// ARIA Expanded/Collapsed Components
export interface CollapsibleProps {
  id: string;
  trigger: ReactNode;
  children: ReactNode;
  defaultExpanded?: boolean;
  className?: string;
  triggerClassName?: string;
  contentClassName?: string;
}

export const Collapsible: React.FC<CollapsibleProps> = ({
  id,
  trigger,
  children,
  defaultExpanded = false,
  className = '',
  triggerClassName = '',
  contentClassName = ''
}) => {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);
  const contentId = `${id}-content`;

  return (
    <div className={className}>
      <button
        id={id}
        className={triggerClassName}
        aria-expanded={isExpanded}
        aria-controls={contentId}
        onClick={() => setIsExpanded(!isExpanded)}
      >
        {trigger}
        <ScreenReaderOnly>
          {isExpanded ? 'Collapse' : 'Expand'} section
        </ScreenReaderOnly>
      </button>
      
      <div
        id={contentId}
        className={contentClassName}
        role="region"
        aria-labelledby={id}
        hidden={!isExpanded}
      >
        {children}
      </div>
    </div>
  );
};

// ARIA Tab Components
export interface TabsProps {
  tabs: Array<{
    id: string;
    label: string;
    content: ReactNode;
    disabled?: boolean;
  }>;
  defaultTab?: string;
  className?: string;
  tabListClassName?: string;
  tabClassName?: string;
  activeTabClassName?: string;
  panelClassName?: string;
}

export const Tabs: React.FC<TabsProps> = ({
  tabs,
  defaultTab,
  className = '',
  tabListClassName = '',
  tabClassName = '',
  activeTabClassName = '',
  panelClassName = ''
}) => {
  const [activeTab, setActiveTab] = useState(defaultTab || tabs[0]?.id);

  return (
    <div className={className}>
      <div
        role="tablist"
        className={tabListClassName}
        aria-label="Configuration tabs"
      >
        {tabs.map((tab) => (
          <button
            key={tab.id}
            id={`tab-${tab.id}`}
            role="tab"
            aria-selected={activeTab === tab.id}
            aria-controls={`panel-${tab.id}`}
            disabled={tab.disabled}
            tabIndex={activeTab === tab.id ? 0 : -1}
            className={`${tabClassName} ${
              activeTab === tab.id ? activeTabClassName : ''
            }`}
            onClick={() => setActiveTab(tab.id)}
            onKeyDown={(e) => {
              const currentIndex = tabs.findIndex(t => t.id === activeTab);
              let newIndex = currentIndex;

              switch (e.key) {
                case 'ArrowLeft':
                  e.preventDefault();
                  newIndex = currentIndex > 0 ? currentIndex - 1 : tabs.length - 1;
                  break;
                case 'ArrowRight':
                  e.preventDefault();
                  newIndex = currentIndex < tabs.length - 1 ? currentIndex + 1 : 0;
                  break;
                case 'Home':
                  e.preventDefault();
                  newIndex = 0;
                  break;
                case 'End':
                  e.preventDefault();
                  newIndex = tabs.length - 1;
                  break;
                default:
                  return;
              }

              const newTab = tabs[newIndex];
              if (!newTab.disabled) {
                setActiveTab(newTab.id);
                document.getElementById(`tab-${newTab.id}`)?.focus();
              }
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {tabs.map((tab) => (
        <div
          key={tab.id}
          id={`panel-${tab.id}`}
          role="tabpanel"
          aria-labelledby={`tab-${tab.id}`}
          hidden={activeTab !== tab.id}
          className={panelClassName}
          tabIndex={0}
        >
          {tab.content}
        </div>
      ))}
    </div>
  );
};

// ARIA List Components
export interface ListProps {
  items: ReactNode[];
  className?: string;
  itemClassName?: string;
  label?: string;
}

export const List: React.FC<ListProps> = ({
  items,
  className = '',
  itemClassName = '',
  label
}) => (
  <ul className={className} aria-label={label}>
    {items.map((item, index) => (
      <li key={index} className={itemClassName}>
        {item}
      </li>
    ))}
  </ul>
);

// ARIA Menu Components
export interface MenuProps {
  trigger: ReactNode;
  items: Array<{
    id: string;
    label: string;
    onClick: () => void;
    disabled?: boolean;
    separator?: boolean;
  }>;
  className?: string;
  triggerClassName?: string;
  menuClassName?: string;
  itemClassName?: string;
}

export const Menu: React.FC<MenuProps> = ({
  trigger,
  items,
  className = '',
  triggerClassName = '',
  menuClassName = '',
  itemClassName = ''
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [focusedIndex, setFocusedIndex] = useState(-1);
  const menuRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setFocusedIndex((prev) => {
          const nextIndex = prev < items.length - 1 ? prev + 1 : 0;
          return items[nextIndex].disabled ? 
            (nextIndex < items.length - 1 ? nextIndex + 1 : 0) : nextIndex;
        });
        break;
      case 'ArrowUp':
        e.preventDefault();
        setFocusedIndex((prev) => {
          const nextIndex = prev > 0 ? prev - 1 : items.length - 1;
          return items[nextIndex].disabled ? 
            (nextIndex > 0 ? nextIndex - 1 : items.length - 1) : nextIndex;
        });
        break;
      case 'Enter':
      case ' ':
        e.preventDefault();
        if (focusedIndex >= 0 && !items[focusedIndex].disabled) {
          items[focusedIndex].onClick();
          setIsOpen(false);
          triggerRef.current?.focus();
        }
        break;
      case 'Escape':
        e.preventDefault();
        setIsOpen(false);
        triggerRef.current?.focus();
        break;
    }
  };

  useEffect(() => {
    if (isOpen && focusedIndex >= 0) {
      const menuItem = menuRef.current?.children[focusedIndex] as HTMLElement;
      menuItem?.focus();
    }
  }, [focusedIndex, isOpen]);

  return (
    <div className={className}>
      <button
        ref={triggerRef}
        className={triggerClassName}
        aria-expanded={isOpen}
        aria-haspopup="menu"
        onClick={() => {
          setIsOpen(!isOpen);
          setFocusedIndex(0);
        }}
      >
        {trigger}
      </button>

      {isOpen && (
        <div
          ref={menuRef}
          role="menu"
          className={menuClassName}
          onKeyDown={handleKeyDown}
        >
          {items.map((item, index) => {
            if (item.separator) {
              return <hr key={index} role="separator" className="my-1" />;
            }

            return (
              <button
                key={item.id}
                role="menuitem"
                className={itemClassName}
                disabled={item.disabled}
                tabIndex={-1}
                onClick={() => {
                  item.onClick();
                  setIsOpen(false);
                  triggerRef.current?.focus();
                }}
              >
                {item.label}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};

// ARIA Progress Components
export interface ProgressProps {
  value: number;
  max?: number;
  label?: string;
  description?: string;
  className?: string;
  showPercentage?: boolean;
}

export const Progress: React.FC<ProgressProps> = ({
  value,
  max = 100,
  label,
  description,
  className = '',
  showPercentage = true
}) => {
  const percentage = Math.round((value / max) * 100);

  return (
    <div className={className}>
      {label && (
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm font-medium text-gray-700">{label}</span>
          {showPercentage && (
            <span className="text-sm text-gray-500">{percentage}%</span>
          )}
        </div>
      )}
      
      <div
        role="progressbar"
        aria-valuenow={value}
        aria-valuemin={0}
        aria-valuemax={max}
        aria-label={label}
        aria-describedby={description ? `${label}-desc` : undefined}
        className="w-full bg-gray-200 rounded-full h-2"
      >
        <div
          className="bg-blue-600 h-2 rounded-full transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>
      
      {description && (
        <Description id={`${label}-desc`} className="text-sm text-gray-600 mt-1">
          {description}
        </Description>
      )}
    </div>
  );
};

// ARIA Loading Component with Spinner
export interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  label?: string;
  className?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  label = 'Loading',
  className = ''
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8'
  };

  return (
    <div className={`inline-flex items-center ${className}`}>
      <svg
        className={`animate-spin ${sizeClasses[size]} text-gray-500`}
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        role="img"
        aria-label={label}
      >
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
        />
        <path
          className="opacity-75"
          fill="currentColor"
          d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        />
      </svg>
      <ScreenReaderOnly>{label}</ScreenReaderOnly>
    </div>
  );
};

// ARIA Tooltip Component
export interface TooltipProps {
  content: string;
  children: ReactNode;
  position?: 'top' | 'bottom' | 'left' | 'right';
  className?: string;
}

export const Tooltip: React.FC<TooltipProps> = ({
  content,
  children,
  position = 'top',
  className = ''
}) => {
  const [isVisible, setIsVisible] = useState(false);
  const tooltipId = useRef(`tooltip-${Math.random().toString(36).substr(2, 9)}`);

  const positionClasses = {
    top: 'bottom-full left-1/2 transform -translate-x-1/2 mb-2',
    bottom: 'top-full left-1/2 transform -translate-x-1/2 mt-2',
    left: 'right-full top-1/2 transform -translate-y-1/2 mr-2',
    right: 'left-full top-1/2 transform -translate-y-1/2 ml-2'
  };

  return (
    <div className={`relative inline-block ${className}`}>
      <div
        aria-describedby={isVisible ? tooltipId.current : undefined}
        onMouseEnter={() => setIsVisible(true)}
        onMouseLeave={() => setIsVisible(false)}
        onFocus={() => setIsVisible(true)}
        onBlur={() => setIsVisible(false)}
      >
        {children}
      </div>
      
      {isVisible && (
        <div
          id={tooltipId.current}
          role="tooltip"
          className={`absolute z-50 px-2 py-1 text-sm text-white bg-gray-900 rounded shadow-lg ${positionClasses[position]}`}
        >
          {content}
        </div>
      )}
    </div>
  );
};

// ARIA Dialog/Modal Component
export interface DialogProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  className?: string;
  overlayClassName?: string;
}

export const Dialog: React.FC<DialogProps> = ({
  isOpen,
  onClose,
  title,
  children,
  className = '',
  overlayClassName = ''
}) => {
  const dialogRef = useRef<HTMLDivElement>(null);
  const titleId = useRef(`dialog-title-${Math.random().toString(36).substr(2, 9)}`);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      dialogRef.current?.focus();
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center ${overlayClassName}`}
      role="dialog"
      aria-modal="true"
      aria-labelledby={titleId.current}
    >
      <div
        className="fixed inset-0 bg-black bg-opacity-50"
        onClick={onClose}
        aria-hidden="true"
      />
      
      <div
        ref={dialogRef}
        className={`relative bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 ${className}`}
        tabIndex={-1}
      >
        <div className="flex items-center justify-between p-4 border-b">
          <h2 id={titleId.current} className="text-lg font-semibold">
            {title}
          </h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded"
            aria-label="Close dialog"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </button>
        </div>
        
        <div className="p-4">
          {children}
        </div>
      </div>
    </div>
  );
};