/**
 * Accessible Layout Component with Skip Links and ARIA Landmarks
 * Provides WCAG 2.1 AA compliant page structure and navigation
 */

import React, { ReactNode } from 'react';
import { 
  SkipLink, 
  ScreenReaderOnly, 
  LiveRegion, 
  AnnouncementRegion,
  useHighContrast,
  useReducedMotion 
} from '../common/FocusManagement';
import { 
  MainLandmark, 
  NavigationLandmark, 
  BannerLandmark, 
  ContentInfoLandmark,
  ComplementaryLandmark 
} from '../common/AriaComponents';
import { QuickAccessNavigation } from '../common/KeyboardNavigation';

export interface AccessibleLayoutProps {
  children: ReactNode;
  header?: ReactNode;
  navigation?: ReactNode;
  sidebar?: ReactNode;
  footer?: ReactNode;
  className?: string;
  announcements?: string[];
}

export const AccessibleLayout: React.FC<AccessibleLayoutProps> = ({
  children,
  header,
  navigation,
  sidebar,
  footer,
  className = '',
  announcements = []
}) => {
  const isHighContrast = useHighContrast();
  const prefersReducedMotion = useReducedMotion();

  return (
    <div 
      className={`min-h-screen bg-white ${className}`}
      data-high-contrast={isHighContrast}
      data-reduced-motion={prefersReducedMotion}
    >
      {/* Skip Links - Always first for keyboard users */}
      <QuickAccessNavigation className="relative z-50" />
      
      {/* Additional skip links */}
      <div className="sr-only focus-within:not-sr-only relative z-50">
        <SkipLink href="#main-content">
          Skip to main content
        </SkipLink>
        {navigation && (
          <SkipLink href="#primary-navigation">
            Skip to navigation
          </SkipLink>
        )}
        {sidebar && (
          <SkipLink href="#sidebar-content">
            Skip to sidebar
          </SkipLink>
        )}
        {footer && (
          <SkipLink href="#footer-content">
            Skip to footer
          </SkipLink>
        )}
      </div>

      {/* Live regions for announcements */}
      {announcements.map((message, index) => (
        <AnnouncementRegion 
          key={index}
          message={message}
          visible={false}
        />
      ))}

      {/* Page structure with proper landmarks */}
      <div className="flex flex-col min-h-screen">
        {/* Header/Banner */}
        {header && (
          <BannerLandmark 
            className="w-full"
            label="Site header"
          >
            {header}
          </BannerLandmark>
        )}

        {/* Main content area */}
        <div className="flex flex-1">
          {/* Primary Navigation */}
          {navigation && (
            <NavigationLandmark 
              id="primary-navigation"
              className="flex-shrink-0"
              label="Primary navigation"
            >
              {navigation}
            </NavigationLandmark>
          )}

          {/* Content wrapper */}
          <div className="flex flex-1">
            {/* Main content */}
            <MainLandmark 
              id="main-content"
              className="flex-1 focus:outline-none"
              label="Main content"
            >
              {/* Focus target for skip links */}
              <div 
                id="main-content-start" 
                tabIndex={-1}
                className="sr-only"
              >
                <ScreenReaderOnly>Main content begins</ScreenReaderOnly>
              </div>
              
              {children}
            </MainLandmark>

            {/* Sidebar/Complementary content */}
            {sidebar && (
              <ComplementaryLandmark 
                id="sidebar-content"
                className="flex-shrink-0"
                label="Sidebar content"
              >
                {sidebar}
              </ComplementaryLandmark>
            )}
          </div>
        </div>

        {/* Footer */}
        {footer && (
          <ContentInfoLandmark 
            id="footer-content"
            className="w-full"
            label="Site footer"
          >
            {footer}
          </ContentInfoLandmark>
        )}
      </div>
    </div>
  );
};

// Specialized layouts for different page types
export interface DashboardLayoutProps extends AccessibleLayoutProps {
  title: string;
  breadcrumbs?: ReactNode;
  actions?: ReactNode;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  title,
  breadcrumbs,
  actions,
  children,
  ...props
}) => {
  return (
    <AccessibleLayout {...props}>
      {/* Page header with title and actions */}
      <header className="bg-white border-b border-gray-200 px-4 py-6 sm:px-6 lg:px-8">
        {breadcrumbs && (
          <nav aria-label="Breadcrumb" className="mb-4">
            {breadcrumbs}
          </nav>
        )}
        
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {title}
            </h1>
            <ScreenReaderOnly>
              You are on the {title} page
            </ScreenReaderOnly>
          </div>
          
          {actions && (
            <div className="flex items-center space-x-4">
              {actions}
            </div>
          )}
        </div>
      </header>

      {/* Page content */}
      <div className="flex-1 p-4 sm:p-6 lg:p-8">
        {children}
      </div>
    </AccessibleLayout>
  );
};

// Form layout with enhanced accessibility
export interface FormLayoutProps extends AccessibleLayoutProps {
  title: string;
  description?: string;
  onSubmit?: (event: React.FormEvent) => void;
  submitLabel?: string;
  cancelLabel?: string;
  onCancel?: () => void;
  isSubmitting?: boolean;
  errors?: Record<string, string>;
}

export const FormLayout: React.FC<FormLayoutProps> = ({
  title,
  description,
  onSubmit,
  submitLabel = 'Submit',
  cancelLabel = 'Cancel',
  onCancel,
  isSubmitting = false,
  errors = {},
  children,
  ...props
}) => {
  const errorCount = Object.keys(errors).length;

  return (
    <AccessibleLayout 
      {...props}
      announcements={[
        ...(props.announcements || []),
        ...(errorCount > 0 ? [`Form has ${errorCount} error${errorCount === 1 ? '' : 's'}`] : [])
      ]}
    >
      <div className="max-w-2xl mx-auto p-6">
        {/* Form header */}
        <header className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            {title}
          </h1>
          {description && (
            <p className="text-gray-600">
              {description}
            </p>
          )}
        </header>

        {/* Error summary */}
        {errorCount > 0 && (
          <div 
            role="alert"
            aria-labelledby="error-summary-title"
            className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md"
          >
            <h2 id="error-summary-title" className="text-lg font-semibold text-red-800 mb-2">
              Please correct the following errors:
            </h2>
            <ul className="list-disc list-inside space-y-1">
              {Object.entries(errors).map(([field, error]) => (
                <li key={field} className="text-red-700">
                  <a 
                    href={`#${field}`}
                    className="underline hover:no-underline focus:outline-none focus:ring-2 focus:ring-red-500"
                  >
                    {error}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Form content */}
        <form 
          onSubmit={onSubmit}
          noValidate
          aria-describedby={description ? 'form-description' : undefined}
        >
          {description && (
            <div id="form-description" className="sr-only">
              {description}
            </div>
          )}
          
          <div className="space-y-6">
            {children}
          </div>

          {/* Form actions */}
          <div className="mt-8 flex items-center justify-end space-x-4">
            {onCancel && (
              <button
                type="button"
                onClick={onCancel}
                className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                disabled={isSubmitting}
              >
                {cancelLabel}
              </button>
            )}
            
            <button
              type="submit"
              disabled={isSubmitting}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              aria-describedby={isSubmitting ? 'submit-loading' : undefined}
            >
              {isSubmitting ? 'Submitting...' : submitLabel}
            </button>
            
            {isSubmitting && (
              <ScreenReaderOnly id="submit-loading">
                Form is being submitted, please wait
              </ScreenReaderOnly>
            )}
          </div>
        </form>
      </div>
    </AccessibleLayout>
  );
};

// Modal layout with proper focus management
export interface ModalLayoutProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: ReactNode;
  className?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  closeOnOverlayClick?: boolean;
  closeOnEscape?: boolean;
}

export const ModalLayout: React.FC<ModalLayoutProps> = ({
  isOpen,
  onClose,
  title,
  description,
  children,
  className = '',
  size = 'md',
  closeOnOverlayClick = true,
  closeOnEscape = true
}) => {
  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl'
  };

  // Prevent body scroll when modal is open
  React.useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      // Announce modal opening
      const announcement = document.createElement('div');
      announcement.setAttribute('aria-live', 'assertive');
      announcement.setAttribute('aria-atomic', 'true');
      announcement.className = 'sr-only';
      announcement.textContent = `${title} dialog opened`;
      document.body.appendChild(announcement);
      
      setTimeout(() => {
        document.body.removeChild(announcement);
      }, 1000);
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen, title]);

  // Handle escape key
  React.useEffect(() => {
    if (!isOpen || !closeOnEscape) return;

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, closeOnEscape, onClose]);

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 z-50 overflow-y-auto"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      aria-describedby={description ? 'modal-description' : undefined}
    >
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={closeOnOverlayClick ? onClose : undefined}
        aria-hidden="true"
      />

      {/* Modal container */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div 
          className={`relative bg-white rounded-lg shadow-xl w-full ${sizeClasses[size]} ${className}`}
        >
          {/* Modal header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <div>
              <h2 id="modal-title" className="text-lg font-semibold text-gray-900">
                {title}
              </h2>
              {description && (
                <p id="modal-description" className="mt-1 text-sm text-gray-600">
                  {description}
                </p>
              )}
            </div>
            
            <button
              onClick={onClose}
              className="p-1 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded"
              aria-label={`Close ${title} dialog`}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Modal content */}
          <div className="p-6">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
};

// Breadcrumb navigation component
export interface BreadcrumbItem {
  label: string;
  href?: string;
  current?: boolean;
}

export interface BreadcrumbNavigationProps {
  items: BreadcrumbItem[];
  className?: string;
}

export const BreadcrumbNavigation: React.FC<BreadcrumbNavigationProps> = ({
  items,
  className = ''
}) => {
  return (
    <nav aria-label="Breadcrumb" className={className}>
      <ol className="flex items-center space-x-2 text-sm">
        {items.map((item, index) => (
          <li key={index} className="flex items-center">
            {index > 0 && (
              <svg 
                className="w-4 h-4 text-gray-400 mr-2" 
                fill="currentColor" 
                viewBox="0 0 20 20"
                aria-hidden="true"
              >
                <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
              </svg>
            )}
            
            {item.current ? (
              <span 
                className="text-gray-900 font-medium"
                aria-current="page"
              >
                {item.label}
              </span>
            ) : (
              <a 
                href={item.href}
                className="text-gray-500 hover:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded"
              >
                {item.label}
              </a>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
};