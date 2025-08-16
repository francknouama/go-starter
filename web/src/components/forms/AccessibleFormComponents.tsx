/**
 * Accessible Form Components for WCAG 2.1 AA Compliance
 * Provides comprehensive form accessibility with proper labels, descriptions, and error handling
 */

import React, { 
  forwardRef, 
  useState,
  useId,
  useRef,
  useEffect
} from 'react';
import type {
  InputHTMLAttributes, 
  SelectHTMLAttributes, 
  TextareaHTMLAttributes,
  ReactNode
} from 'react';
import { ScreenReaderOnly } from '../common/FocusManagement';
import { FormField, Description } from '../common/AriaComponents';
import { accessibility } from '../../styles/design-tokens';

// Base form control interface
interface BaseFormControlProps {
  label: string;
  description?: string;
  error?: string;
  required?: boolean;
  helpText?: string;
  className?: string;
  labelClassName?: string;
  controlClassName?: string;
}

// Enhanced Input Component
export interface AccessibleInputProps extends 
  Omit<InputHTMLAttributes<HTMLInputElement>, 'id'>,
  BaseFormControlProps {
  inputType?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search';
  leftAddon?: ReactNode;
  rightAddon?: ReactNode;
  characterLimit?: number;
  showCharacterCount?: boolean;
}

export const AccessibleInput = forwardRef<HTMLInputElement, AccessibleInputProps>(({
  label,
  description,
  error,
  required = false,
  helpText,
  className = '',
  labelClassName = '',
  controlClassName = '',
  inputType = 'text',
  leftAddon,
  rightAddon,
  characterLimit,
  showCharacterCount = false,
  value,
  onChange,
  ...props
}, ref) => {
  const inputId = useId();
  const descriptionId = useId();
  const errorId = useId();
  const helpId = useId();
  const characterCountId = useId();
  
  const [currentLength, setCurrentLength] = useState(
    typeof value === 'string' ? value.length : 0
  );

  const isHighContrast = accessibility.prefersHighContrast();

  // Update character count when value changes
  useEffect(() => {
    if (typeof value === 'string') {
      setCurrentLength(value.length);
    }
  }, [value]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setCurrentLength(event.target.value.length);
    onChange?.(event);
  };

  // Build describedby attribute
  const describedBy = [
    description && descriptionId,
    error && errorId,
    helpText && helpId,
    showCharacterCount && characterCountId
  ].filter(Boolean).join(' ');

  const isNearLimit = characterLimit && currentLength > characterLimit * 0.8;
  const isOverLimit = characterLimit && currentLength > characterLimit;

  return (
    <div className={className}>
      <label 
        htmlFor={inputId}
        className={`block text-sm font-medium mb-1 ${
          error ? 'text-red-700' : 'text-gray-700'
        } ${labelClassName}`}
      >
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-label="required">
            *
          </span>
        )}
      </label>
      
      {description && (
        <Description id={descriptionId} className="text-sm text-gray-600 mb-2">
          {description}
        </Description>
      )}
      
      <div className="relative">
        {leftAddon && (
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span className="text-gray-500 sm:text-sm" aria-hidden="true">
              {leftAddon}
            </span>
          </div>
        )}
        
        <input
          ref={ref}
          id={inputId}
          type={inputType}
          value={value}
          onChange={handleChange}
          aria-describedby={describedBy || undefined}
          aria-required={required}
          aria-invalid={!!error}
          maxLength={characterLimit}
          className={`
            block w-full rounded-md border-gray-300 shadow-sm
            focus:border-blue-500 focus:ring-blue-500 focus:ring-1
            disabled:bg-gray-50 disabled:text-gray-500
            ${error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}
            ${leftAddon ? 'pl-10' : 'pl-3'}
            ${rightAddon ? 'pr-10' : 'pr-3'}
            ${isHighContrast ? 'border-2' : 'border'}
            ${controlClassName}
          `}
          {...props}
        />
        
        {rightAddon && (
          <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span className="text-gray-500 sm:text-sm" aria-hidden="true">
              {rightAddon}
            </span>
          </div>
        )}
      </div>
      
      {showCharacterCount && characterLimit && (
        <div 
          id={characterCountId}
          className={`mt-1 text-sm ${
            isOverLimit ? 'text-red-600' : 
            isNearLimit ? 'text-yellow-600' : 'text-gray-500'
          }`}
          aria-live={isNearLimit ? 'polite' : 'off'}
        >
          {currentLength} / {characterLimit} characters
          {isOverLimit && (
            <span className="ml-1 font-medium">
              (Exceeds limit by {currentLength - characterLimit})
            </span>
          )}
        </div>
      )}
      
      {helpText && (
        <Description id={helpId} className="text-sm text-gray-600 mt-1">
          {helpText}
        </Description>
      )}
      
      {error && (
        <Description id={errorId} className="text-sm text-red-600 mt-1">
          <span role="alert">{error}</span>
        </Description>
      )}
    </div>
  );
});

AccessibleInput.displayName = 'AccessibleInput';

// Enhanced Select Component
export interface AccessibleSelectProps extends 
  Omit<SelectHTMLAttributes<HTMLSelectElement>, 'id'>,
  BaseFormControlProps {
  options: Array<{
    value: string;
    label: string;
    disabled?: boolean;
    group?: string;
  }>;
  placeholder?: string;
  emptyMessage?: string;
}

export const AccessibleSelect = forwardRef<HTMLSelectElement, AccessibleSelectProps>(({
  label,
  description,
  error,
  required = false,
  helpText,
  className = '',
  labelClassName = '',
  controlClassName = '',
  options,
  placeholder = 'Select an option',
  emptyMessage = 'No options available',
  ...props
}, ref) => {
  const selectId = useId();
  const descriptionId = useId();
  const errorId = useId();
  const helpId = useId();

  const isHighContrast = accessibility.prefersHighContrast();

  // Group options by group property
  const groupedOptions = options.reduce((acc, option) => {
    const group = option.group || 'default';
    if (!acc[group]) acc[group] = [];
    acc[group].push(option);
    return acc;
  }, {} as Record<string, typeof options>);

  const hasGroups = Object.keys(groupedOptions).length > 1 || 
                   !groupedOptions.default;

  // Build describedby attribute
  const describedBy = [
    description && descriptionId,
    error && errorId,
    helpText && helpId
  ].filter(Boolean).join(' ');

  return (
    <div className={className}>
      <label 
        htmlFor={selectId}
        className={`block text-sm font-medium mb-1 ${
          error ? 'text-red-700' : 'text-gray-700'
        } ${labelClassName}`}
      >
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-label="required">
            *
          </span>
        )}
      </label>
      
      {description && (
        <Description id={descriptionId} className="text-sm text-gray-600 mb-2">
          {description}
        </Description>
      )}
      
      <select
        ref={ref}
        id={selectId}
        aria-describedby={describedBy || undefined}
        aria-required={required}
        aria-invalid={!!error}
        className={`
          block w-full rounded-md border-gray-300 shadow-sm
          focus:border-blue-500 focus:ring-blue-500 focus:ring-1
          disabled:bg-gray-50 disabled:text-gray-500
          ${error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}
          ${isHighContrast ? 'border-2' : 'border'}
          ${controlClassName}
        `}
        {...props}
      >
        {!required && (
          <option value="">{placeholder}</option>
        )}
        
        {options.length === 0 ? (
          <option value="" disabled>{emptyMessage}</option>
        ) : hasGroups ? (
          Object.entries(groupedOptions).map(([group, groupOptions]) => (
            group === 'default' ? (
              groupOptions.map((option) => (
                <option 
                  key={option.value} 
                  value={option.value}
                  disabled={option.disabled}
                >
                  {option.label}
                </option>
              ))
            ) : (
              <optgroup key={group} label={group}>
                {groupOptions.map((option) => (
                  <option 
                    key={option.value} 
                    value={option.value}
                    disabled={option.disabled}
                  >
                    {option.label}
                  </option>
                ))}
              </optgroup>
            )
          ))
        ) : (
          options.map((option) => (
            <option 
              key={option.value} 
              value={option.value}
              disabled={option.disabled}
            >
              {option.label}
            </option>
          ))
        )}
      </select>
      
      {helpText && (
        <Description id={helpId} className="text-sm text-gray-600 mt-1">
          {helpText}
        </Description>
      )}
      
      {error && (
        <Description id={errorId} className="text-sm text-red-600 mt-1">
          <span role="alert">{error}</span>
        </Description>
      )}
    </div>
  );
});

AccessibleSelect.displayName = 'AccessibleSelect';

// Enhanced Textarea Component
export interface AccessibleTextareaProps extends 
  Omit<TextareaHTMLAttributes<HTMLTextAreaElement>, 'id'>,
  BaseFormControlProps {
  characterLimit?: number;
  showCharacterCount?: boolean;
  autoResize?: boolean;
}

export const AccessibleTextarea = forwardRef<HTMLTextAreaElement, AccessibleTextareaProps>(({
  label,
  description,
  error,
  required = false,
  helpText,
  className = '',
  labelClassName = '',
  controlClassName = '',
  characterLimit,
  showCharacterCount = false,
  autoResize = false,
  value,
  onChange,
  ...props
}, ref) => {
  const textareaId = useId();
  const descriptionId = useId();
  const errorId = useId();
  const helpId = useId();
  const characterCountId = useId();
  
  const [currentLength, setCurrentLength] = useState(
    typeof value === 'string' ? value.length : 0
  );

  const isHighContrast = accessibility.prefersHighContrast();

  // Update character count when value changes
  useEffect(() => {
    if (typeof value === 'string') {
      setCurrentLength(value.length);
    }
  }, [value]);

  const handleChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setCurrentLength(event.target.value.length);
    
    // Auto-resize functionality
    if (autoResize) {
      const textarea = event.target;
      textarea.style.height = 'auto';
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
    
    onChange?.(event);
  };

  // Build describedby attribute
  const describedBy = [
    description && descriptionId,
    error && errorId,
    helpText && helpId,
    showCharacterCount && characterCountId
  ].filter(Boolean).join(' ');

  const isNearLimit = characterLimit && currentLength > characterLimit * 0.8;
  const isOverLimit = characterLimit && currentLength > characterLimit;

  return (
    <div className={className}>
      <label 
        htmlFor={textareaId}
        className={`block text-sm font-medium mb-1 ${
          error ? 'text-red-700' : 'text-gray-700'
        } ${labelClassName}`}
      >
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-label="required">
            *
          </span>
        )}
      </label>
      
      {description && (
        <Description id={descriptionId} className="text-sm text-gray-600 mb-2">
          {description}
        </Description>
      )}
      
      <textarea
        ref={ref}
        id={textareaId}
        value={value}
        onChange={handleChange}
        aria-describedby={describedBy || undefined}
        aria-required={required}
        aria-invalid={!!error}
        maxLength={characterLimit}
        className={`
          block w-full rounded-md border-gray-300 shadow-sm
          focus:border-blue-500 focus:ring-blue-500 focus:ring-1
          disabled:bg-gray-50 disabled:text-gray-500
          ${error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}
          ${isHighContrast ? 'border-2' : 'border'}
          ${autoResize ? 'resize-none overflow-hidden' : 'resize-vertical'}
          ${controlClassName}
        `}
        {...props}
      />
      
      {showCharacterCount && characterLimit && (
        <div 
          id={characterCountId}
          className={`mt-1 text-sm ${
            isOverLimit ? 'text-red-600' : 
            isNearLimit ? 'text-yellow-600' : 'text-gray-500'
          }`}
          aria-live={isNearLimit ? 'polite' : 'off'}
        >
          {currentLength} / {characterLimit} characters
          {isOverLimit && (
            <span className="ml-1 font-medium">
              (Exceeds limit by {currentLength - characterLimit})
            </span>
          )}
        </div>
      )}
      
      {helpText && (
        <Description id={helpId} className="text-sm text-gray-600 mt-1">
          {helpText}
        </Description>
      )}
      
      {error && (
        <Description id={errorId} className="text-sm text-red-600 mt-1">
          <span role="alert">{error}</span>
        </Description>
      )}
    </div>
  );
});

AccessibleTextarea.displayName = 'AccessibleTextarea';

// Checkbox Component with enhanced accessibility
export interface AccessibleCheckboxProps extends 
  Omit<InputHTMLAttributes<HTMLInputElement>, 'id' | 'type'> {
  label: string;
  description?: string;
  error?: string;
  indeterminate?: boolean;
  className?: string;
  labelClassName?: string;
}

export const AccessibleCheckbox = forwardRef<HTMLInputElement, AccessibleCheckboxProps>(({
  label,
  description,
  error,
  indeterminate = false,
  className = '',
  labelClassName = '',
  ...props
}, ref) => {
  const checkboxId = useId();
  const descriptionId = useId();
  const errorId = useId();
  
  const checkboxRef = useRef<HTMLInputElement>(null);
  
  // Merge refs
  React.useImperativeHandle(ref, () => checkboxRef.current!, []);

  const isHighContrast = accessibility.prefersHighContrast();

  // Set indeterminate state
  useEffect(() => {
    if (checkboxRef.current) {
      checkboxRef.current.indeterminate = indeterminate;
    }
  }, [indeterminate]);

  // Build describedby attribute
  const describedBy = [
    description && descriptionId,
    error && errorId
  ].filter(Boolean).join(' ');

  return (
    <div className={className}>
      <div className="flex items-start">
        <div className="flex items-center h-5">
          <input
            ref={checkboxRef}
            id={checkboxId}
            type="checkbox"
            aria-describedby={describedBy || undefined}
            aria-invalid={!!error}
            className={`
              w-4 h-4 text-blue-600 border-gray-300 rounded
              focus:ring-blue-500 focus:ring-2 focus:ring-offset-0
              disabled:bg-gray-50 disabled:border-gray-300
              ${error ? 'border-red-300 focus:ring-red-500' : ''}
              ${isHighContrast ? 'border-2' : 'border'}
            `}
            {...props}
          />
        </div>
        
        <div className="ml-3 text-sm">
          <label 
            htmlFor={checkboxId}
            className={`font-medium cursor-pointer ${
              error ? 'text-red-700' : 'text-gray-700'
            } ${labelClassName}`}
          >
            {label}
          </label>
          
          {description && (
            <Description id={descriptionId} className="text-gray-600 mt-1">
              {description}
            </Description>
          )}
          
          {error && (
            <Description id={errorId} className="text-red-600 mt-1">
              <span role="alert">{error}</span>
            </Description>
          )}
        </div>
      </div>
    </div>
  );
});

AccessibleCheckbox.displayName = 'AccessibleCheckbox';

// Radio Group Component
export interface RadioOption {
  value: string;
  label: string;
  description?: string;
  disabled?: boolean;
}

export interface AccessibleRadioGroupProps {
  name: string;
  label: string;
  options: RadioOption[];
  value?: string;
  onChange?: (value: string) => void;
  error?: string;
  description?: string;
  required?: boolean;
  orientation?: 'horizontal' | 'vertical';
  className?: string;
  labelClassName?: string;
}

export const AccessibleRadioGroup: React.FC<AccessibleRadioGroupProps> = ({
  name,
  label,
  options,
  value,
  onChange,
  error,
  description,
  required = false,
  orientation = 'vertical',
  className = '',
  labelClassName = ''
}) => {
  const groupId = useId();
  const descriptionId = useId();
  const errorId = useId();

  const isHighContrast = accessibility.prefersHighContrast();

  // Build describedby attribute
  const describedBy = [
    description && descriptionId,
    error && errorId
  ].filter(Boolean).join(' ');

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange?.(event.target.value);
  };

  return (
    <fieldset className={className} aria-describedby={describedBy || undefined}>
      <legend className={`text-sm font-medium mb-2 ${
        error ? 'text-red-700' : 'text-gray-700'
      } ${labelClassName}`}>
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-label="required">
            *
          </span>
        )}
      </legend>
      
      {description && (
        <Description id={descriptionId} className="text-sm text-gray-600 mb-3">
          {description}
        </Description>
      )}
      
      <div className={`space-y-3 ${orientation === 'horizontal' ? 'sm:flex sm:space-y-0 sm:space-x-6' : ''}`}>
        {options.map((option, index) => {
          const radioId = `${groupId}-${index}`;
          const radioDescriptionId = option.description ? `${radioId}-desc` : undefined;
          
          return (
            <div key={option.value} className="flex items-start">
              <div className="flex items-center h-5">
                <input
                  id={radioId}
                  name={name}
                  type="radio"
                  value={option.value}
                  checked={value === option.value}
                  onChange={handleChange}
                  disabled={option.disabled}
                  aria-describedby={radioDescriptionId}
                  className={`
                    w-4 h-4 text-blue-600 border-gray-300
                    focus:ring-blue-500 focus:ring-2 focus:ring-offset-0
                    disabled:bg-gray-50 disabled:border-gray-300
                    ${error ? 'border-red-300 focus:ring-red-500' : ''}
                    ${isHighContrast ? 'border-2' : 'border'}
                  `}
                />
              </div>
              
              <div className="ml-3 text-sm">
                <label 
                  htmlFor={radioId}
                  className={`font-medium cursor-pointer ${
                    option.disabled ? 'text-gray-400' : 
                    error ? 'text-red-700' : 'text-gray-700'
                  }`}
                >
                  {option.label}
                </label>
                
                {option.description && (
                  <Description id={radioDescriptionId || ''} className="text-gray-600 mt-1">
                    {option.description}
                  </Description>
                )}
              </div>
            </div>
          );
        })}
      </div>
      
      {error && (
        <Description id={errorId} className="text-sm text-red-600 mt-3">
          <span role="alert">{error}</span>
        </Description>
      )}
    </fieldset>
  );
};

// Form Validation Hook
export interface ValidationRule {
  required?: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  email?: boolean;
  custom?: (value: any) => string | undefined;
}

export interface FormField {
  value: any;
  error?: string;
  touched?: boolean;
  rules?: ValidationRule;
}

export const useFormValidation = <T extends Record<string, FormField>>(
  initialState: T
) => {
  const [fields, setFields] = useState<T>(initialState);

  const validate = (fieldName: keyof T, value: any): string | undefined => {
    const rules = fields[fieldName].rules;
    if (!rules) return undefined;

    if (rules.required && (!value || value.toString().trim() === '')) {
      return 'This field is required';
    }

    if (value && rules.minLength && value.toString().length < rules.minLength) {
      return `Must be at least ${rules.minLength} characters`;
    }

    if (value && rules.maxLength && value.toString().length > rules.maxLength) {
      return `Must be no more than ${rules.maxLength} characters`;
    }

    if (value && rules.pattern && !rules.pattern.test(value.toString())) {
      return 'Invalid format';
    }

    if (value && rules.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value.toString())) {
      return 'Invalid email address';
    }

    if (rules.custom) {
      return rules.custom(value);
    }

    return undefined;
  };

  const updateField = (fieldName: keyof T, value: any, shouldValidate = true) => {
    setFields(prev => ({
      ...prev,
      [fieldName]: {
        ...prev[fieldName],
        value,
        error: shouldValidate ? validate(fieldName, value) : undefined,
        touched: true
      }
    }));
  };

  const validateAll = (): boolean => {
    let isValid = true;
    const updatedFields = { ...fields };

    Object.keys(fields).forEach(fieldName => {
      const field = fields[fieldName];
      const error = validate(fieldName, field.value);
      (updatedFields as any)[fieldName] = {
        ...field,
        error,
        touched: true
      };
      if (error) isValid = false;
    });

    setFields(updatedFields);
    return isValid;
  };

  const reset = () => {
    setFields(initialState);
  };

  const getErrors = (): Record<string, string> => {
    const errors: Record<string, string> = {};
    Object.entries(fields).forEach(([key, field]) => {
      if (field.error) {
        errors[key] = field.error;
      }
    });
    return errors;
  };

  return {
    fields,
    updateField,
    validateAll,
    reset,
    getErrors,
    isValid: Object.values(fields).every(field => !field.error)
  };
};