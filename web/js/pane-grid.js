
class PaneGrid {
    constructor(container) {
        this.container = container;
        this.panes = new Map();
        this.activePane = null;
        this.gridTemplate = {
            columns: '1fr',
            rows: '1fr',
            areas: [['main']]
        };
        this.render();
    }

    addPane(id, terminal) {
        const pane = { id, terminal, element: document.createElement('div') };
        pane.element.classList.add('pane');
        pane.element.dataset.paneId = id;
        pane.element.appendChild(terminal.element);
        this.panes.set(id, pane);
        this.container.appendChild(pane.element);
        this.setActivePane(id);
        this.updateLayout();
    }

    removePane(id) {
        const pane = this.panes.get(id);
        if (pane) {
            pane.element.remove();
            this.panes.delete(id);
            this.updateLayout();
        }
    }

    setActivePane(id) {
        if (this.activePane) {
            this.activePane.element.classList.remove('active');
        }
        const pane = this.panes.get(id);
        if (pane) {
            pane.element.classList.add('active');
            this.activePane = pane;
            pane.terminal.focus();
        }
    }

    updateLayout(layout) {
        // For now, we'll just re-render with the current grid template.
        // In the future, this will take a layout object to dynamically
        // change the grid.
        this.render();
    }

    render() {
        this.container.style.gridTemplateColumns = this.gridTemplate.columns;
        this.container.style.gridTemplateRows = this.gridTemplate.rows;
        this.container.style.gridTemplateAreas = this.gridTemplate.areas
            .map(row => `"${row.join(' ')}"`)
            .join(' ');

        this.panes.forEach((pane, id) => {
            pane.element.style.gridArea = this.findGridArea(id);
        });
    }

    findGridArea(id) {
        // This is a placeholder. A real implementation would have a mapping
        // from pane ID to grid area. For now, we assume a single pane 'main'.
        return 'main';
    }
}
