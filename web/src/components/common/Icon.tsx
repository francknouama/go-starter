import React from 'react'
import {
  CodeBracketIcon,
  CommandLineIcon,
  BookOpenIcon,
  CloudIcon,
  Cog6ToothIcon,
  CubeIcon,
  ServerIcon,
  // BuildingOfficeIcon,
  FolderIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
  XCircleIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  PlusIcon,
  MinusIcon,
  XMarkIcon,
  MagnifyingGlassIcon,
  AdjustmentsHorizontalIcon,
  EyeIcon,
  EyeSlashIcon,
  DocumentIcon,
  DocumentTextIcon,
  ClipboardDocumentIcon,
  ArrowDownTrayIcon,
  ShareIcon,
  HeartIcon,
  StarIcon,
  LightBulbIcon,
  QuestionMarkCircleIcon,
  SunIcon,
  MoonIcon,
  ComputerDesktopIcon,
} from '@heroicons/react/24/outline'

import {
  CheckCircleIcon as CheckCircleIconSolid,
  ExclamationCircleIcon as ExclamationCircleIconSolid,
  InformationCircleIcon as InformationCircleIconSolid,
  XCircleIcon as XCircleIconSolid,
  StarIcon as StarIconSolid,
  HeartIcon as HeartIconSolid,
} from '@heroicons/react/24/solid'

// Icon mapping for project types
const projectTypeIcons = {
  'web-api': CodeBracketIcon,
  'cli': CommandLineIcon,
  'library': BookOpenIcon,
  'lambda': CloudIcon,
  'lambda-proxy': CloudIcon,
  'microservice': CubeIcon,
  'monolith': ServerIcon,
  'workspace': FolderIcon,
  'event-driven': Cog6ToothIcon,
} as const;

// Icon mapping for common UI actions
const actionIcons = {
  // Navigation
  'chevron-down': ChevronDownIcon,
  'chevron-up': ChevronUpIcon,
  'chevron-left': ChevronLeftIcon,
  'chevron-right': ChevronRightIcon,
  
  // Actions
  'plus': PlusIcon,
  'minus': MinusIcon,
  'close': XMarkIcon,
  'search': MagnifyingGlassIcon,
  'settings': AdjustmentsHorizontalIcon,
  'view': EyeIcon,
  'hide': EyeSlashIcon,
  
  // Documents
  'document': DocumentIcon,
  'document-text': DocumentTextIcon,
  'clipboard': ClipboardDocumentIcon,
  'download': ArrowDownTrayIcon,
  'share': ShareIcon,
  
  // Feedback
  'heart': HeartIcon,
  'heart-solid': HeartIconSolid,
  'star': StarIcon,
  'star-solid': StarIconSolid,
  'lightbulb': LightBulbIcon,
  'help': QuestionMarkCircleIcon,
  
  // Theme
  'sun': SunIcon,
  'moon': MoonIcon,
  'desktop': ComputerDesktopIcon,
} as const;

// Icon mapping for status indicators
const statusIcons = {
  'success': CheckCircleIcon,
  'success-solid': CheckCircleIconSolid,
  'warning': ExclamationCircleIcon,
  'warning-solid': ExclamationCircleIconSolid,
  'info': InformationCircleIcon,
  'info-solid': InformationCircleIconSolid,
  'error': XCircleIcon,
  'error-solid': XCircleIconSolid,
} as const;

// Combined icon registry
const iconRegistry = {
  ...projectTypeIcons,
  ...actionIcons,
  ...statusIcons,
} as const;

export type IconName = keyof typeof iconRegistry;
export type ProjectTypeIconName = keyof typeof projectTypeIcons;
export type ActionIconName = keyof typeof actionIcons;
export type StatusIconName = keyof typeof statusIcons;

export interface IconProps {
  name: IconName;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  className?: string;
  'aria-hidden'?: boolean;
}

const Icon: React.FC<IconProps> = ({
  name,
  size = 'md',
  className = '',
  'aria-hidden': ariaHidden = true,
}) => {
  const IconComponent = iconRegistry[name];
  
  if (!IconComponent) {
    console.warn(`Icon "${name}" not found in registry`);
    return null;
  }
  
  const sizeClasses = {
    xs: 'w-3 h-3',
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
    xl: 'w-8 h-8',
    '2xl': 'w-12 h-12',
  };
  
  return (
    <IconComponent
      className={`${sizeClasses[size]} ${className}`}
      aria-hidden={ariaHidden}
    />
  );
};

export default Icon;

// Project Type Icon Component
export interface ProjectTypeIconProps {
  type: ProjectTypeIconName;
  size?: IconProps['size'];
  className?: string;
  selected?: boolean;
}

export const ProjectTypeIcon: React.FC<ProjectTypeIconProps> = ({
  type,
  size = 'md',
  className = '',
  selected = false,
}) => {
  const baseClasses = selected 
    ? 'text-primary-600' 
    : 'text-gray-600';
    
  return (
    <Icon
      name={type}
      size={size}
      className={`${baseClasses} ${className}`}
    />
  );
};

// Status Icon Component with semantic colors
export interface StatusIconProps {
  status: 'success' | 'warning' | 'info' | 'error';
  variant?: 'outline' | 'solid';
  size?: IconProps['size'];
  className?: string;
}

export const StatusIcon: React.FC<StatusIconProps> = ({
  status,
  variant = 'outline',
  size = 'md',
  className = '',
}) => {
  const iconName = `${status}${variant === 'solid' ? '-solid' : ''}` as StatusIconName;
  
  const statusColors = {
    success: 'text-success-600',
    warning: 'text-warning-600',
    info: 'text-primary-600',
    error: 'text-error-600',
  };
  
  return (
    <Icon
      name={iconName}
      size={size}
      className={`${statusColors[status]} ${className}`}
    />
  );
};

// Export icon registry for external use
export { iconRegistry, projectTypeIcons, actionIcons, statusIcons };