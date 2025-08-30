/**
 * Configuration State Management Hook
 * Manages project configuration with validation, persistence, and defaults
 */

import { useState, useEffect, useCallback, useRef } from 'react'
import type { ProjectConfig, Blueprint } from '../services/api'
import { useBlueprints, useDefaultConfig, useBlueprintValidation } from './useApi'

export interface ConfigurationState {
  config: ProjectConfig
  isValid: boolean
  errors: { [key: string]: string }
  isDirty: boolean
  lastSaved: Date | null
}

export interface ConfigurationActions {
  updateConfig: (updates: Partial<ProjectConfig>) => void
  resetConfig: () => void
  loadConfig: (config: ProjectConfig) => void
  validateConfig: () => Promise<boolean>
  saveConfig: () => void
  loadSavedConfig: () => boolean
}

interface UseConfigurationOptions {
  blueprint?: Blueprint | null
  enablePersistence?: boolean
  autoSave?: boolean
  autoSaveDelay?: number
  storageKey?: string
}

export function useConfiguration({
  blueprint,
  enablePersistence = true,
  autoSave = true,
  autoSaveDelay = 2000,
  storageKey = 'go-starter-config'
}: UseConfigurationOptions = {}) {
  // Default configuration
  const defaultConfig: ProjectConfig = {
    projectName: '',
    moduleUrl: '',
    goVersion: '1.21',
    projectType: 'web-api',
    framework: 'gin',
    architecture: 'standard',
    logger: 'slog'
  }

  // Hooks
  const { config: serverDefaultConfig } = useDefaultConfig()
  const { validateConfig: validateWithServer } = useBlueprintValidation()

  // State
  const [state, setState] = useState<ConfigurationState>({
    config: defaultConfig,
    isValid: false,
    errors: {},
    isDirty: false,
    lastSaved: null
  })

  // Refs for debouncing
  const autoSaveTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const validationTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  // Validation rules
  const validationRules = {
    projectName: {
      required: true,
      pattern: /^[a-z0-9-]+$/,
      minLength: 2,
      maxLength: 50,
      message: 'Project name must be lowercase alphanumeric with hyphens, 2-50 characters'
    },
    moduleUrl: {
      required: true,
      pattern: /^[a-z0-9.-]+\/[a-z0-9.-\/]+$/,
      message: 'Module URL must be a valid Go module path (e.g., github.com/user/project)'
    },
    projectType: {
      required: true,
      options: ['cli', 'web-api', 'library', 'lambda', 'microservice', 'monolith'],
      message: 'Please select a valid project type'
    },
    framework: {
      required: true,
      message: 'Please select a framework'
    },
    architecture: {
      required: true,
      options: ['standard', 'clean', 'ddd', 'hexagonal'],
      message: 'Please select an architecture pattern'
    },
    logger: {
      required: true,
      options: ['slog', 'zap', 'logrus', 'zerolog'],
      message: 'Please select a logger type'
    },
    goVersion: {
      required: true,
      pattern: /^1\.\d+$/,
      message: 'Go version must be in format 1.x'
    }
  }

  // Validate configuration
  const validateConfiguration = useCallback(async (config: ProjectConfig): Promise<{ isValid: boolean; errors: { [key: string]: string } }> => {
    const errors: { [key: string]: string } = {}

    // Client-side validation
    Object.entries(validationRules).forEach(([key, rules]) => {
      const value = config[key as keyof ProjectConfig] as any
      
      if (rules.required && (!value || value.trim() === '')) {
        errors[key] = `${key} is required`
        return
      }

      if (value && rules.pattern && !rules.pattern.test(value)) {
        errors[key] = rules.message
        return
      }

      if (value && rules.minLength && value.length < rules.minLength) {
        errors[key] = rules.message
        return
      }

      if (value && rules.maxLength && value.length > rules.maxLength) {
        errors[key] = rules.message
        return
      }

      if (value && rules.options && !rules.options.includes(value)) {
        errors[key] = rules.message
        return
      }
    })

    // Server-side validation (if blueprint is available)
    if (blueprint && Object.keys(errors).length === 0) {
      try {
        const serverValidation = await validateWithServer(blueprint.id, config)
        if (!serverValidation.valid && serverValidation.errors) {
          serverValidation.errors.forEach((error, index) => {
            errors[`server_${index}`] = error
          })
        }
      } catch (error) {
        console.warn('Server validation failed:', error)
        // Don't fail client-side validation if server validation fails
      }
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors
    }
  }, [blueprint, validateWithServer])

  // Save configuration to localStorage (moved before updateConfig)
  const saveConfig = useCallback(() => {
    if (!enablePersistence) return

    try {
      const saveData = {
        config: state.config,
        timestamp: Date.now(),
        version: '1.0'
      }
      localStorage.setItem(storageKey, JSON.stringify(saveData))
      
      setState(prev => ({
        ...prev,
        isDirty: false,
        lastSaved: new Date()
      }))
    } catch (error) {
      console.error('Failed to save configuration:', error)
    }
  }, [state.config, enablePersistence, storageKey])

  // Update configuration (prevent excessive re-renders)
  const updateConfig = useCallback((updates: Partial<ProjectConfig>) => {
    setState(prev => {
      const newConfig = { ...prev.config, ...updates }
      
      // Only update if something actually changed
      if (JSON.stringify(newConfig) === JSON.stringify(prev.config)) {
        return prev
      }
      
      // Clear validation timeout
      if (validationTimeoutRef.current) {
        clearTimeout(validationTimeoutRef.current)
      }

      // Schedule validation (debounced)
      validationTimeoutRef.current = setTimeout(async () => {
        try {
          const validation = await validateConfiguration(newConfig)
          setState(current => ({
            ...current,
            isValid: validation.isValid,
            errors: validation.errors
          }))
        } catch (error) {
          console.error('Validation error:', error)
        }
      }, 300)

      return {
        ...prev,
        config: newConfig,
        isDirty: true,
        errors: {} // Clear errors while typing
      }
    })

    // Auto-save (debounced)
    if (autoSave && enablePersistence) {
      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current)
      }
      
      autoSaveTimeoutRef.current = setTimeout(() => {
        try {
          saveConfig()
        } catch (error) {
          console.error('Auto-save error:', error)
        }
      }, autoSaveDelay)
    }
  }, [validateConfiguration, autoSave, enablePersistence, autoSaveDelay, saveConfig])

  // Reset configuration
  const resetConfig = useCallback(() => {
    const resetToConfig = serverDefaultConfig || defaultConfig
    setState({
      config: resetToConfig,
      isValid: false,
      errors: {},
      isDirty: false,
      lastSaved: null
    })
    
    if (enablePersistence) {
      localStorage.removeItem(storageKey)
    }
  }, [serverDefaultConfig, defaultConfig, enablePersistence, storageKey])

  // Load configuration
  const loadConfig = useCallback((config: ProjectConfig) => {
    setState(prev => ({
      ...prev,
      config: { ...defaultConfig, ...config },
      isDirty: false
    }))

    // Validate loaded config
    validateConfiguration(config).then(validation => {
      setState(current => ({
        ...current,
        isValid: validation.isValid,
        errors: validation.errors
      }))
    })
  }, [defaultConfig, validateConfiguration])

  // Validate configuration (manual trigger)
  const validateConfig = useCallback(async (): Promise<boolean> => {
    const validation = await validateConfiguration(state.config)
    setState(prev => ({
      ...prev,
      isValid: validation.isValid,
      errors: validation.errors
    }))
    return validation.isValid
  }, [state.config, validateConfiguration])


  // Load saved configuration from localStorage
  const loadSavedConfig = useCallback((): boolean => {
    if (!enablePersistence) return false

    try {
      const saved = localStorage.getItem(storageKey)
      if (!saved) return false

      const saveData = JSON.parse(saved)
      if (saveData.config && saveData.timestamp) {
        // Check if saved config is not too old (7 days)
        const maxAge = 7 * 24 * 60 * 60 * 1000
        if (Date.now() - saveData.timestamp > maxAge) {
          localStorage.removeItem(storageKey)
          return false
        }

        loadConfig(saveData.config)
        setState(prev => ({
          ...prev,
          lastSaved: new Date(saveData.timestamp),
          isDirty: false
        }))
        return true
      }
    } catch (error) {
      console.error('Failed to load saved configuration:', error)
      localStorage.removeItem(storageKey)
    }
    
    return false
  }, [enablePersistence, storageKey, loadConfig])

  // Initialize configuration
  useEffect(() => {
    // Try to load saved config first
    if (!loadSavedConfig()) {
      // Fall back to server default or local default
      if (serverDefaultConfig) {
        loadConfig(serverDefaultConfig)
      }
    }
  }, [serverDefaultConfig, loadSavedConfig, loadConfig])

  // Blueprint-specific configuration adjustments (prevent circular updates)
  const blueprintIdRef = useRef<string | null>(null)
  
  useEffect(() => {
    if (blueprint && blueprint.id !== blueprintIdRef.current) {
      blueprintIdRef.current = blueprint.id
      
      const updates: Partial<ProjectConfig> = {}
      const { framework, architecture } = state.config
      
      // Auto-adjust framework based on blueprint
      if (blueprint.frameworks && blueprint.frameworks.length > 0) {
        if (!blueprint.frameworks.includes(framework)) {
          updates.framework = blueprint.frameworks[0] as any
        }
      }

      // Auto-adjust architecture based on blueprint
      if (blueprint.architectures && blueprint.architectures.length > 0) {
        if (!blueprint.architectures.includes(architecture)) {
          updates.architecture = blueprint.architectures[0] as any
        }
      }

      // Apply updates directly to state without causing updateConfig loop
      if (Object.keys(updates).length > 0) {
        setState(prev => ({
          ...prev,
          config: { ...prev.config, ...updates },
          isDirty: true,
          isValid: false // Reset validation when config changes
        }))
      }
    }
  }, [blueprint?.id]) // Only depend on blueprint ID

  // Auto-generate module URL from project name (prevent circular updates)
  useEffect(() => {
    const { projectName, moduleUrl } = state.config
    if (projectName && 
        !moduleUrl.includes(projectName) && 
        (!moduleUrl || moduleUrl.includes('user') || moduleUrl.includes('project')) &&
        projectName.length >= 2) {
      const cleanProjectName = projectName.toLowerCase().replace(/[^a-z0-9-]/g, '-')
      const newModuleUrl = `github.com/user/${cleanProjectName}`
      
      // Only update if it's actually different
      if (newModuleUrl !== moduleUrl) {
        setState(prev => ({
          ...prev,
          config: { ...prev.config, moduleUrl: newModuleUrl },
          isDirty: true
        }))
      }
    }
  }, [state.config.projectName]) // Remove moduleUrl and updateConfig from dependencies

  // Cleanup timeouts on unmount
  useEffect(() => {
    return () => {
      if (autoSaveTimeoutRef.current) {
        clearTimeout(autoSaveTimeoutRef.current)
      }
      if (validationTimeoutRef.current) {
        clearTimeout(validationTimeoutRef.current)
      }
    }
  }, [])

  const actions: ConfigurationActions = {
    updateConfig,
    resetConfig,
    loadConfig,
    validateConfig,
    saveConfig,
    loadSavedConfig
  }

  return {
    ...state,
    ...actions
  }
}