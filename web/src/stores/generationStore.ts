/**
 * Global Generation State Store
 * Zustand-based state management for the generation workflow
 */

import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'
import type { Blueprint, ProjectConfig } from '../services/api'
import type { WorkflowSession, GenerationResult } from '../hooks/useGenerationWorkflow'

export type AppView = 
  | 'blueprint-selection'  // User is selecting a blueprint
  | 'configuration'        // User is configuring the project
  | 'generation'          // Project is being generated
  | 'success'             // Generation completed successfully
  | 'error'               // Generation failed

export interface UserPreferences {
  theme: 'light' | 'dark' | 'system'
  disclosureMode: 'basic' | 'advanced'
  autoSave: boolean
  showHelpTooltips: boolean
  preferredFramework?: string
  preferredLogger?: string
  defaultGoVersion: string
}

export interface GenerationStore {
  // Current view and navigation
  currentView: AppView
  previousView: AppView | null
  navigationHistory: AppView[]
  
  // Blueprint selection
  selectedBlueprint: Blueprint | null
  availableBlueprints: Blueprint[]
  blueprintFilters: {
    category?: string
    difficulty?: string
    search?: string
  }
  
  // Configuration
  projectConfig: ProjectConfig
  configErrors: { [key: string]: string }
  isConfigValid: boolean
  
  // Generation workflow
  currentSession: WorkflowSession | null
  sessionHistory: WorkflowSession[]
  generationResult: GenerationResult | null
  
  // UI state
  isLoading: boolean
  isSidebarOpen: boolean
  showPreview: boolean
  showAdvancedOptions: boolean
  
  // User preferences
  preferences: UserPreferences
  
  // Actions
  setCurrentView: (view: AppView) => void
  navigateTo: (view: AppView) => void
  goBack: () => void
  
  setSelectedBlueprint: (blueprint: Blueprint | null) => void
  setAvailableBlueprints: (blueprints: Blueprint[]) => void
  setBlueprintFilters: (filters: Partial<GenerationStore['blueprintFilters']>) => void
  
  updateProjectConfig: (config: Partial<ProjectConfig>) => void
  setConfigErrors: (errors: { [key: string]: string }) => void
  setConfigValid: (valid: boolean) => void
  resetConfig: () => void
  
  setCurrentSession: (session: WorkflowSession | null) => void
  addToSessionHistory: (session: WorkflowSession) => void
  setGenerationResult: (result: GenerationResult | null) => void
  
  setLoading: (loading: boolean) => void
  toggleSidebar: () => void
  setShowPreview: (show: boolean) => void
  setShowAdvancedOptions: (show: boolean) => void
  
  updatePreferences: (preferences: Partial<UserPreferences>) => void
  
  // Utility actions
  resetApp: () => void
  clearHistory: () => void
}

const defaultConfig: ProjectConfig = {
  projectName: '',
  moduleUrl: '',
  goVersion: '1.21',
  projectType: 'web-api',
  framework: 'gin',
  architecture: 'standard',
  logger: 'slog'
}

const defaultPreferences: UserPreferences = {
  theme: 'system',
  disclosureMode: 'basic',
  autoSave: true,
  showHelpTooltips: true,
  defaultGoVersion: '1.21'
}

export const useGenerationStore = create<GenerationStore>()(
  persist(
    immer((set, get) => ({
      // Initial state
      currentView: 'blueprint-selection',
      previousView: null,
      navigationHistory: [],
      
      selectedBlueprint: null,
      availableBlueprints: [],
      blueprintFilters: {},
      
      projectConfig: { ...defaultConfig },
      configErrors: {},
      isConfigValid: false,
      
      currentSession: null,
      sessionHistory: [],
      generationResult: null,
      
      isLoading: false,
      isSidebarOpen: false,
      showPreview: false,
      showAdvancedOptions: false,
      
      preferences: { ...defaultPreferences },
      
      // Navigation actions
      setCurrentView: (view: AppView) =>
        set((state) => {
          state.previousView = state.currentView
          state.currentView = view
          state.navigationHistory.push(view)
          
          // Keep history limited
          if (state.navigationHistory.length > 10) {
            state.navigationHistory = state.navigationHistory.slice(-10)
          }
        }),
      
      navigateTo: (view: AppView) => {
        const { setCurrentView } = get()
        setCurrentView(view)
      },
      
      goBack: () =>
        set((state) => {
          if (state.previousView) {
            const temp = state.currentView
            state.currentView = state.previousView
            state.previousView = temp
          } else if (state.navigationHistory.length > 1) {
            // Go to previous view in history
            state.navigationHistory.pop() // Remove current
            const previousView = state.navigationHistory[state.navigationHistory.length - 1]
            if (previousView) {
              state.previousView = state.currentView
              state.currentView = previousView
            }
          }
        }),
      
      // Blueprint actions
      setSelectedBlueprint: (blueprint: Blueprint | null) =>
        set((state) => {
          // Only update if blueprint actually changed
          if (state.selectedBlueprint?.id === blueprint?.id) {
            return
          }
          
          state.selectedBlueprint = blueprint
          
          // Auto-adjust config based on blueprint (only if blueprint is different)
          if (blueprint) {
            let configChanged = false
            
            // Check framework compatibility
            if (blueprint.frameworks?.[0] && !blueprint.frameworks.includes(state.projectConfig.framework)) {
              state.projectConfig.framework = blueprint.frameworks[0] as any
              configChanged = true
            }
            
            // Check architecture compatibility
            if (blueprint.architectures?.[0] && !blueprint.architectures.includes(state.projectConfig.architecture)) {
              state.projectConfig.architecture = blueprint.architectures[0] as any
              configChanged = true
            }
            
            // Only reset validation if we changed config
            if (configChanged) {
              state.isConfigValid = false
            }
          }
        }),
      
      setAvailableBlueprints: (blueprints: Blueprint[]) =>
        set((state) => {
          state.availableBlueprints = blueprints
        }),
      
      setBlueprintFilters: (filters) =>
        set((state) => {
          state.blueprintFilters = { ...state.blueprintFilters, ...filters }
        }),
      
      // Configuration actions
      updateProjectConfig: (config) =>
        set((state) => {
          const newConfig = { ...state.projectConfig, ...config }
          
          // Only update if something actually changed
          if (JSON.stringify(newConfig) === JSON.stringify(state.projectConfig)) {
            return
          }
          
          state.projectConfig = newConfig
          
          // Auto-generate module URL if project name changes and doesn't already have a custom one
          if (config.projectName && (!state.projectConfig.moduleUrl || state.projectConfig.moduleUrl.includes('user'))) {
            const cleanName = config.projectName.toLowerCase().replace(/[^a-z0-9-]/g, '-')
            state.projectConfig.moduleUrl = `github.com/user/${cleanName}`
          }
          
          // Clear errors when config changes
          state.configErrors = {}
          state.isConfigValid = false
        }),
      
      setConfigErrors: (errors) =>
        set((state) => {
          state.configErrors = errors
          state.isConfigValid = Object.keys(errors).length === 0
        }),
      
      setConfigValid: (valid) =>
        set((state) => {
          state.isConfigValid = valid
        }),
      
      resetConfig: () =>
        set((state) => {
          state.projectConfig = { ...defaultConfig }
          state.configErrors = {}
          state.isConfigValid = false
        }),
      
      // Session actions
      setCurrentSession: (session) =>
        set((state) => {
          state.currentSession = session
        }),
      
      addToSessionHistory: (session) =>
        set((state) => {
          const existingIndex = state.sessionHistory.findIndex(s => s.id === session.id)
          if (existingIndex >= 0) {
            state.sessionHistory[existingIndex] = session
          } else {
            state.sessionHistory.push(session)
            // Keep history limited
            if (state.sessionHistory.length > 20) {
              state.sessionHistory = state.sessionHistory.slice(-20)
            }
          }
        }),
      
      setGenerationResult: (result) =>
        set((state) => {
          state.generationResult = result
        }),
      
      // UI actions
      setLoading: (loading) =>
        set((state) => {
          state.isLoading = loading
        }),
      
      toggleSidebar: () =>
        set((state) => {
          state.isSidebarOpen = !state.isSidebarOpen
        }),
      
      setShowPreview: (show) =>
        set((state) => {
          state.showPreview = show
        }),
      
      setShowAdvancedOptions: (show) =>
        set((state) => {
          state.showAdvancedOptions = show
        }),
      
      // Preferences actions
      updatePreferences: (preferences) =>
        set((state) => {
          const newPreferences = { ...state.preferences, ...preferences }
          
          // Only update if preferences actually changed
          if (JSON.stringify(newPreferences) === JSON.stringify(state.preferences)) {
            return
          }
          
          state.preferences = newPreferences
          
          // Apply preference changes to config if needed (only if not already set)
          if (preferences.defaultGoVersion && (!state.projectConfig.goVersion || state.projectConfig.goVersion === '1.21')) {
            state.projectConfig.goVersion = preferences.defaultGoVersion
          }
        }),
      
      // Utility actions
      resetApp: () =>
        set((state) => {
          state.currentView = 'blueprint-selection'
          state.previousView = null
          state.navigationHistory = []
          state.selectedBlueprint = null
          state.projectConfig = { ...defaultConfig }
          state.configErrors = {}
          state.isConfigValid = false
          state.currentSession = null
          state.generationResult = null
          state.isLoading = false
          state.showPreview = false
          state.showAdvancedOptions = false
        }),
      
      clearHistory: () =>
        set((state) => {
          state.sessionHistory = []
          state.navigationHistory = []
        })
    })),
    {
      name: 'go-starter-store',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        // Only persist these parts of the state
        selectedBlueprint: state.selectedBlueprint,
        projectConfig: state.projectConfig,
        sessionHistory: state.sessionHistory.slice(-5), // Keep only last 5 sessions
        preferences: state.preferences,
        blueprintFilters: state.blueprintFilters
      }),
      version: 1,
      migrate: (persistedState: any, version: number) => {
        // Handle migrations between versions
        if (version === 0) {
          // Migration from version 0 to 1
          return {
            ...persistedState,
            preferences: {
              ...defaultPreferences,
              ...persistedState.preferences
            }
          }
        }
        return persistedState
      }
    }
  )
)

// Selectors for common state access patterns
export const selectCurrentBlueprint = (state: GenerationStore) => state.selectedBlueprint
export const selectProjectConfig = (state: GenerationStore) => state.projectConfig
export const selectIsConfigValid = (state: GenerationStore) => state.isConfigValid
export const selectCurrentSession = (state: GenerationStore) => state.currentSession
export const selectGenerationResult = (state: GenerationStore) => state.generationResult
export const selectPreferences = (state: GenerationStore) => state.preferences
export const selectIsLoading = (state: GenerationStore) => state.isLoading

// Hook for common combinations
export const useCurrentWorkflow = () => {
  return useGenerationStore((state) => ({
    currentView: state.currentView,
    selectedBlueprint: state.selectedBlueprint,
    projectConfig: state.projectConfig,
    isConfigValid: state.isConfigValid,
    currentSession: state.currentSession,
    generationResult: state.generationResult,
    isLoading: state.isLoading
  }))
}

export const useBlueprintSelection = () => {
  return useGenerationStore((state) => ({
    selectedBlueprint: state.selectedBlueprint,
    availableBlueprints: state.availableBlueprints,
    blueprintFilters: state.blueprintFilters,
    setSelectedBlueprint: state.setSelectedBlueprint,
    setAvailableBlueprints: state.setAvailableBlueprints,
    setBlueprintFilters: state.setBlueprintFilters
  }))
}

export const useProjectConfiguration = () => {
  return useGenerationStore((state) => ({
    projectConfig: state.projectConfig,
    configErrors: state.configErrors,
    isConfigValid: state.isConfigValid,
    selectedBlueprint: state.selectedBlueprint,
    updateProjectConfig: state.updateProjectConfig,
    setConfigErrors: state.setConfigErrors,
    setConfigValid: state.setConfigValid,
    resetConfig: state.resetConfig
  }))
}