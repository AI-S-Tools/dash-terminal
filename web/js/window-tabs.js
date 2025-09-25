/**
 * Window Tab Component for Dash Terminal
 * Handles JavaScript interactivity for tmux window tabs
 * Based on SessionTabs pattern with unified touch/click handling
 */

class WindowTabs {
    constructor() {
        this.windows = new Map(); // windowId -> window info
        this.activeWindow = null;
        this.currentSession = null;
        this.init();
    }

    init() {
        // Add unified click/touch handlers
        this.attachInteractionHandlers();

        // Initialize with default window
        this.initializeDefaultWindow();

        console.log('WindowTabs: Initialized with interaction handlers');
    }

    attachInteractionHandlers() {
        const windowTabsContainer = document.querySelector('.window-tabs');
        if (!windowTabsContainer) {
            console.error('WindowTabs: .window-tabs container not found');
            return;
        }

        // Unified interaction handler for both click and touch
        const handleInteraction = (e) => {
            // Prevent default behavior and event bubbling
            e.preventDefault();
            e.stopPropagation();

            const target = e.target;
            console.log(`WindowTabs: ${e.type} on ${target.className}:`, target.textContent);

            if (target.classList.contains('window-tab')) {
                this.handleWindowTabClick(target);
            } else if (target.classList.contains('window-add')) {
                this.handleAddWindowClick(target);
            }
        };

        // Use touchstart for mobile (fires before click) and click as fallback
        windowTabsContainer.addEventListener('touchstart', handleInteraction, { passive: false });
        windowTabsContainer.addEventListener('click', (e) => {
            // Only handle click if it wasn't already handled by touch
            if (e.type === 'click' && !e.isTrusted) return;
            handleInteraction(e);
        });
    }

    handleWindowTabClick(tabElement) {
        const windowName = tabElement.textContent.trim();
        console.log(`WindowTabs: Clicked window tab: ${windowName}`);

        // Remove active class from all window tabs
        document.querySelectorAll('.window-tab').forEach(tab => {
            tab.classList.remove('active');
        });

        // Add active class to clicked tab
        tabElement.classList.add('active');

        // Store active window
        this.activeWindow = windowName;

        // Trigger window selection event
        this.onWindowSelect(windowName);
    }

    handleAddWindowClick(addElement) {
        console.log('WindowTabs: Add window clicked');

        // Generate new window name
        const windowCount = document.querySelectorAll('.window-tab').length;
        const defaultNames = ['bash', 'htop', 'vim', 'logs', 'ssh', 'tail', 'top'];
        let newWindowName;

        if (windowCount < defaultNames.length) {
            newWindowName = defaultNames[windowCount];
        } else {
            newWindowName = `window-${windowCount + 1}`;
        }

        // Add new window tab
        this.addWindowTab(newWindowName);

        // Trigger window creation event
        this.onWindowCreate(newWindowName);
    }

    addWindowTab(windowName) {
        const windowTabsContainer = document.querySelector('.window-tabs');
        const addButton = windowTabsContainer.querySelector('.window-add');

        // Create new window tab
        const newTab = document.createElement('div');
        newTab.className = 'window-tab';
        newTab.textContent = windowName;

        // Insert before the add button
        windowTabsContainer.insertBefore(newTab, addButton);

        // Select the new tab
        this.handleWindowTabClick(newTab);

        console.log(`WindowTabs: Added new window tab: ${windowName}`);
    }

    onWindowSelect(windowName) {
        // Emit custom event for window selection
        const event = new CustomEvent('windowSelect', {
            detail: {
                windowName,
                sessionName: this.currentSession
            }
        });
        document.dispatchEvent(event);

        console.log(`WindowTabs: Window selected: ${windowName} in session: ${this.currentSession}`);
    }

    onWindowCreate(windowName) {
        // Emit custom event for window creation
        const event = new CustomEvent('windowCreate', {
            detail: {
                windowName,
                sessionName: this.currentSession
            }
        });
        document.dispatchEvent(event);

        console.log(`WindowTabs: Window creation requested: ${windowName} in session: ${this.currentSession}`);
    }

    initializeDefaultWindow() {
        // Find the first window tab and make it active if none is active
        const firstTab = document.querySelector('.window-tab');
        const activeTab = document.querySelector('.window-tab.active');

        if (firstTab && !activeTab) {
            this.handleWindowTabClick(firstTab);
        } else if (activeTab) {
            this.activeWindow = activeTab.textContent.trim();
        }
    }

    // Public methods for external control
    selectWindow(windowName) {
        const tabs = document.querySelectorAll('.window-tab');
        for (const tab of tabs) {
            if (tab.textContent.trim() === windowName) {
                this.handleWindowTabClick(tab);
                return true;
            }
        }
        return false;
    }

    getCurrentWindow() {
        return this.activeWindow;
    }

    // Integration with session tabs
    setCurrentSession(sessionName) {
        this.currentSession = sessionName;
        console.log(`WindowTabs: Current session set to: ${sessionName}`);
    }

    // Reset windows for new session (optional - for future use)
    resetWindowsForSession(sessionName) {
        // Clear all windows except the first one
        const windowTabs = document.querySelectorAll('.window-tab');
        windowTabs.forEach((tab, index) => {
            if (index > 0) {
                tab.remove();
            }
        });

        // Reset to first window
        if (windowTabs.length > 0) {
            this.handleWindowTabClick(windowTabs[0]);
        }

        this.setCurrentSession(sessionName);
        console.log(`WindowTabs: Reset windows for session: ${sessionName}`);
    }

    // Mobile-specific touch enhancements
    enableMobileOptimizations() {
        const windowTabsContainer = document.querySelector('.window-tabs');

        // Enable smooth scrolling for tab overflow
        windowTabsContainer.style.scrollBehavior = 'smooth';

        // Simple touch feedback without interfering with main touch handler
        windowTabsContainer.addEventListener('touchend', (e) => {
            if (e.target.classList.contains('window-tab') || e.target.classList.contains('window-add')) {
                // Visual feedback
                e.target.style.opacity = '0.7';
                setTimeout(() => {
                    e.target.style.opacity = '1';
                }, 150);
                console.log('Window touch feedback on:', e.target.textContent);
            }
        });

        // iOS Safari specific - prevent double tap zoom
        windowTabsContainer.style.touchAction = 'manipulation';

        console.log('WindowTabs: Mobile optimizations enabled');
    }
}

// Export for global use
window.WindowTabs = WindowTabs;