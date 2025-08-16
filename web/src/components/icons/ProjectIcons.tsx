import { 
  CommandLineIcon,
  ServerIcon,
  CubeIcon,
  BoltIcon,
  CloudIcon,
  CpuChipIcon,
  RectangleStackIcon,
  CircleStackIcon
} from '@heroicons/react/24/outline'

// Project Type Icons
export const ProjectTypeIcons = {
  'cli': CommandLineIcon,
  'web-api': ServerIcon,
  'library': CubeIcon,
  'lambda': BoltIcon,
  'lambda-proxy': CloudIcon,
  'microservice': CpuChipIcon,
  'monolith': RectangleStackIcon,
  'workspace': CircleStackIcon,
} as const

// Architecture Icons
export const ArchitectureIcons = {
  'standard': CubeIcon,
  'clean': RectangleStackIcon,
  'ddd': CircleStackIcon,
  'hexagonal': CpuChipIcon,
  'event-driven': BoltIcon,
} as const

// Framework Icons - Using text representations since framework logos aren't available
export const FrameworkIcons = {
  gin: () => <span className="font-bold text-lg">G</span>,
  echo: () => <span className="font-bold text-lg">E</span>,
  fiber: () => <span className="font-bold text-lg">F</span>,
  chi: () => <span className="font-bold text-lg">C</span>,
  cobra: () => <span className="font-bold text-lg">C</span>,
} as const

// Logger Icons - Using text representations
export const LoggerIcons = {
  slog: () => <span className="font-bold text-lg">S</span>,
  zap: () => <span className="font-bold text-lg">Z</span>,
  logrus: () => <span className="font-bold text-lg">L</span>,
  zerolog: () => <span className="font-bold text-lg">0</span>,
} as const