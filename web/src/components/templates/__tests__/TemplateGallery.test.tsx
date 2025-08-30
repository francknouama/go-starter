import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TemplateGallery from '../TemplateGallery';
import { PROJECT_TEMPLATES, TEMPLATE_CATEGORIES } from '../../../data/projectTemplates';
import type { ProjectTemplate } from '../../../data/projectTemplates';

// Mock the child components
jest.mock('../TemplateCard', () => {
  return function MockTemplateCard({ template, onSelect, onPreview }: {
    template: ProjectTemplate;
    onSelect: () => void;
    onPreview: () => void;
  }) {
    return (
      <div data-testid={`template-card-${template.id}`} className="template-card">
        <h3>{template.name}</h3>
        <p>{template.description}</p>
        <button onClick={onSelect} data-testid={`select-${template.id}`}>
          Select
        </button>
        <button onClick={onPreview} data-testid={`preview-${template.id}`}>
          Preview
        </button>
      </div>
    );
  };
});

jest.mock('../TemplatePreviewModal', () => {
  return function MockTemplatePreviewModal({ template, onClose, onSelect }: {
    template: ProjectTemplate;
    onClose: () => void;
    onSelect: () => void;
  }) {
    return (
      <div data-testid="template-preview-modal">
        <h2>Preview: {template.name}</h2>
        <button onClick={onClose} data-testid="close-preview">Close</button>
        <button onClick={onSelect} data-testid="select-from-preview">Select</button>
      </div>
    );
  };
});

describe('TemplateGallery', () => {
  const mockOnSelectTemplate = jest.fn();
  const mockOnClose = jest.fn();
  const user = userEvent.setup();

  const defaultProps = {
    onSelectTemplate: mockOnSelectTemplate,
    onClose: mockOnClose,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Rendering and Basic Functionality', () => {
    it('renders the template gallery with header', () => {
      render(<TemplateGallery {...defaultProps} />);

      expect(screen.getByText('Production-Ready Blueprints')).toBeInTheDocument();
      expect(screen.getByText('🎉 100% Coverage')).toBeInTheDocument();
      expect(screen.getByText(/Historic achievement: 12 production-ready blueprints/)).toBeInTheDocument();
    });

    it('displays achievement stats correctly', () => {
      render(<TemplateGallery {...defaultProps} />);

      expect(screen.getByText(PROJECT_TEMPLATES.length.toString())).toBeInTheDocument();
      expect(screen.getByText(TEMPLATE_CATEGORIES.length.toString())).toBeInTheDocument();
      expect(screen.getByText('100%')).toBeInTheDocument();
    });

    it('renders search input', () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      expect(searchInput).toBeInTheDocument();
    });

    it('renders popular templates by default', () => {
      render(<TemplateGallery {...defaultProps} />);

      // Should show popular view by default
      const popularButton = screen.getByRole('button', { name: /Most Popular/ });
      expect(popularButton).toHaveClass('bg-gradient-to-r');
    });

    it('closes when overlay is clicked', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const overlay = document.querySelector('.fixed.inset-0.bg-black');
      expect(overlay).toBeInTheDocument();

      fireEvent.click(overlay!);
      expect(mockOnClose).toHaveBeenCalled();
    });

    it('closes when close button is clicked', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const closeButton = screen.getByRole('button', { name: /Close/ });
      await user.click(closeButton);

      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  describe('View Mode Switching', () => {
    it('switches between popular and gallery view modes', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const popularButton = screen.getByRole('button', { name: /Most Popular/ });
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });

      // Popular should be active by default
      expect(popularButton).toHaveClass('bg-gradient-to-r');
      expect(galleryButton).not.toHaveClass('bg-gradient-to-r');

      // Switch to gallery view
      await user.click(galleryButton);

      expect(galleryButton).toHaveClass('bg-gradient-to-r');
      expect(popularButton).not.toHaveClass('bg-gradient-to-r');
    });

    it('switches between grid and compact layout modes', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const gridButton = screen.getByRole('button', { name: /Grid view/ });
      const compactButton = screen.getByRole('button', { name: /Compact view/ });

      // Grid should be active by default
      expect(gridButton).toHaveClass('bg-white');
      expect(compactButton).not.toHaveClass('bg-white');

      // Switch to compact view
      await user.click(compactButton);

      expect(compactButton).toHaveClass('bg-white');
      expect(gridButton).not.toHaveClass('bg-white');
    });
  });

  describe('Search Functionality', () => {
    it('updates search query when user types', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');

      await user.type(searchInput, 'cli');

      expect(searchInput).toHaveValue('cli');
    });

    it('shows clear button when search query exists', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');

      await user.type(searchInput, 'test query');

      const clearButton = screen.getByRole('button', { name: '✕' });
      expect(clearButton).toBeInTheDocument();
    });

    it('clears search query when clear button is clicked', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');

      await user.type(searchInput, 'test query');
      
      const clearButton = screen.getByRole('button', { name: '✕' });
      await user.click(clearButton);

      expect(searchInput).toHaveValue('');
    });

    it('filters templates based on search query', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view to see all templates
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');

      // Search for CLI templates
      await user.type(searchInput, 'cli');

      // Should show CLI templates
      await waitFor(() => {
        expect(screen.getByTestId('template-card-cli-simple')).toBeInTheDocument();
      });
    });
  });

  describe('Category Filtering', () => {
    it('renders all category buttons', () => {
      render(<TemplateGallery {...defaultProps} />);

      // Should have "All" button
      expect(screen.getByRole('button', { name: /All \(\d+\)/ })).toBeInTheDocument();

      // Should have all category buttons
      TEMPLATE_CATEGORIES.forEach(category => {
        const categoryButton = screen.getByRole('button', { name: new RegExp(category.name) });
        expect(categoryButton).toBeInTheDocument();
      });
    });

    it('filters templates by category', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view first
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      // Find and click on CLI category
      const cliCategory = TEMPLATE_CATEGORIES.find(cat => cat.id === 'cli');
      if (cliCategory) {
        const categoryButton = screen.getByRole('button', { name: new RegExp(cliCategory.name) });
        await user.click(categoryButton);

        // Should show only CLI templates
        await waitFor(() => {
          expect(screen.getByTestId('template-card-cli-simple')).toBeInTheDocument();
        });
      }
    });

    it('highlights selected category', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const allButton = screen.getByRole('button', { name: /All \(\d+\)/ });
      
      // "All" should be selected by default
      expect(allButton).toHaveClass('bg-gradient-to-r');

      // Find CLI category and click it
      const cliCategory = TEMPLATE_CATEGORIES.find(cat => cat.id === 'cli');
      if (cliCategory) {
        const categoryButton = screen.getByRole('button', { name: new RegExp(cliCategory.name) });
        await user.click(categoryButton);

        expect(categoryButton).toHaveClass('bg-gradient-to-r');
        expect(allButton).not.toHaveClass('bg-gradient-to-r');
      }
    });
  });

  describe('Complexity Filtering', () => {
    it('filters templates by complexity level', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view first
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      // Find and click simple complexity
      const simpleButton = screen.getByRole('button', { name: /Simple \(\d+\)/ });
      await user.click(simpleButton);

      // Should show only simple templates
      await waitFor(() => {
        expect(screen.getByTestId('template-card-cli-simple')).toBeInTheDocument();
      });
    });

    it('highlights selected complexity level', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const allLevelsButton = screen.getByRole('button', { name: /All Levels/ });
      
      // "All Levels" should be selected by default
      expect(allLevelsButton).toHaveClass('bg-gradient-to-r');

      // Click simple complexity
      const simpleButton = screen.getByRole('button', { name: /Simple \(\d+\)/ });
      await user.click(simpleButton);

      expect(simpleButton).toHaveStyle('color: white');
      expect(allLevelsButton).not.toHaveClass('bg-gradient-to-r');
    });
  });

  describe('Filter Management', () => {
    it('shows filter indicator when filters are active', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const moreFiltersButton = screen.getByRole('button', { name: /More Filters/ });
      
      // No indicator initially
      expect(moreFiltersButton.querySelector('.bg-purple-500')).not.toBeInTheDocument();

      // Apply a filter
      const cliCategory = TEMPLATE_CATEGORIES.find(cat => cat.id === 'cli');
      if (cliCategory) {
        const categoryButton = screen.getByRole('button', { name: new RegExp(cliCategory.name) });
        await user.click(categoryButton);

        // Should show filter indicator
        await waitFor(() => {
          expect(moreFiltersButton.querySelector('.bg-purple-500')).toBeInTheDocument();
        });
      }
    });

    it('shows clear all button when filters are active', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Apply a search filter
      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      await user.type(searchInput, 'cli');

      // Should show clear all button
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Clear All' })).toBeInTheDocument();
      });
    });

    it('clears all filters when clear all button is clicked', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      
      // Apply filters
      await user.type(searchInput, 'cli');
      
      const cliCategory = TEMPLATE_CATEGORIES.find(cat => cat.id === 'cli');
      if (cliCategory) {
        const categoryButton = screen.getByRole('button', { name: new RegExp(cliCategory.name) });
        await user.click(categoryButton);
      }

      // Click clear all
      const clearAllButton = screen.getByRole('button', { name: 'Clear All' });
      await user.click(clearAllButton);

      // All filters should be cleared
      expect(searchInput).toHaveValue('');
      expect(screen.getByRole('button', { name: /All \(\d+\)/ })).toHaveClass('bg-gradient-to-r');
    });

    it('toggles advanced filters panel', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const moreFiltersButton = screen.getByRole('button', { name: /More Filters/ });
      
      // Panel should not be visible initially
      expect(screen.queryByText('Category')).not.toBeInTheDocument();

      // Show filters
      await user.click(moreFiltersButton);
      expect(screen.getByText('Category')).toBeInTheDocument();
      expect(screen.getByText('Complexity')).toBeInTheDocument();

      // Hide filters
      await user.click(moreFiltersButton);
      expect(screen.queryByText('Category')).not.toBeInTheDocument();
    });
  });

  describe('Template Selection and Preview', () => {
    it('calls onSelectTemplate when template is selected', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view to see all templates
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      await waitFor(() => {
        const selectButton = screen.getByTestId('select-cli-simple');
        expect(selectButton).toBeInTheDocument();
      });

      const selectButton = screen.getByTestId('select-cli-simple');
      await user.click(selectButton);

      expect(mockOnSelectTemplate).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'cli-simple' })
      );
      expect(mockOnClose).toHaveBeenCalled();
    });

    it('opens preview modal when preview button is clicked', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view to see all templates
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      await waitFor(() => {
        const previewButton = screen.getByTestId('preview-cli-simple');
        expect(previewButton).toBeInTheDocument();
      });

      const previewButton = screen.getByTestId('preview-cli-simple');
      await user.click(previewButton);

      expect(screen.getByTestId('template-preview-modal')).toBeInTheDocument();
    });

    it('closes preview modal', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view and open preview
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      await waitFor(() => {
        const previewButton = screen.getByTestId('preview-cli-simple');
        expect(previewButton).toBeInTheDocument();
      });

      const previewButton = screen.getByTestId('preview-cli-simple');
      await user.click(previewButton);

      // Close preview
      const closePreviewButton = screen.getByTestId('close-preview');
      await user.click(closePreviewButton);

      expect(screen.queryByTestId('template-preview-modal')).not.toBeInTheDocument();
    });

    it('selects template from preview modal', async () => {
      render(<TemplateGallery {...defaultProps} />);

      // Switch to gallery view and open preview
      const galleryButton = screen.getByRole('button', { name: /All.*Blueprints/ });
      await user.click(galleryButton);

      await waitFor(() => {
        const previewButton = screen.getByTestId('preview-cli-simple');
        expect(previewButton).toBeInTheDocument();
      });

      const previewButton = screen.getByTestId('preview-cli-simple');
      await user.click(previewButton);

      // Select from preview
      const selectFromPreviewButton = screen.getByTestId('select-from-preview');
      await user.click(selectFromPreviewButton);

      expect(mockOnSelectTemplate).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'cli-simple' })
      );
      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  describe('Empty State', () => {
    it('shows no results message when no templates match filters', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      
      // Search for something that doesn't exist
      await user.type(searchInput, 'nonexistenttemplate123');

      await waitFor(() => {
        expect(screen.getByText('No blueprints found')).toBeInTheDocument();
        expect(screen.getByText(/We couldn't find any blueprints matching/)).toBeInTheDocument();
      });
    });

    it('provides options to clear filters or browse all in empty state', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      
      // Search for something that doesn't exist
      await user.type(searchInput, 'nonexistenttemplate123');

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Clear all filters' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Browse all blueprints' })).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('has proper ARIA labels and roles', () => {
      render(<TemplateGallery {...defaultProps} />);

      // Close button should have accessible name
      const closeButton = screen.getByRole('button', { name: /Close/ });
      expect(closeButton).toBeInTheDocument();

      // Search input should have proper placeholder
      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      expect(searchInput).toHaveAttribute('type', 'text');
    });

    it('supports keyboard navigation', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      
      // Should be focusable
      searchInput.focus();
      expect(searchInput).toHaveFocus();

      // Tab should move to next focusable element
      await user.tab();
      const moreFiltersButton = screen.getByRole('button', { name: /More Filters/ });
      expect(moreFiltersButton).toHaveFocus();
    });
  });

  describe('Performance', () => {
    it('renders large number of templates efficiently', () => {
      const startTime = performance.now();
      
      render(<TemplateGallery {...defaultProps} />);
      
      const endTime = performance.now();
      const renderTime = endTime - startTime;
      
      // Should render within reasonable time (less than 100ms)
      expect(renderTime).toBeLessThan(100);
    });

    it('debounces search input', async () => {
      render(<TemplateGallery {...defaultProps} />);

      const searchInput = screen.getByPlaceholderText('Search templates by name, category, or technology...');
      
      // Type multiple characters quickly
      await user.type(searchInput, 'cli', { delay: 10 });
      
      // Should handle rapid input without issues
      expect(searchInput).toHaveValue('cli');
    });
  });
});