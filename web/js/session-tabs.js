/**
 * Session Tab Component for Dash Terminal
 * Handles JavaScript interactivity for tmux session tabs
 */

class SessionTabs {
    constructor() {
        this.sessions = new Map(); // sessionId -> session info
        this.activeSession = null;
        this.init();
    }

    init() {
        // Add click handlers to existing session tabs
        this.attachClickHandlers();

        // Initialize with default session
        this.initializeDefaultSession();

        console.log('SessionTabs: Initialized with click handlers');
    }

    attachClickHandlers() {
        const sessionTabsContainer = document.querySelector('.session-tabs');
        if (!sessionTabsContainer) {
            console.error('SessionTabs: .session-tabs container not found');
            return;
        }

        // Event delegation for session tabs
        sessionTabsContainer.addEventListener('click', (e) => {
            if (e.target.classList.contains('session-tab')) {
                this.handleSessionTabClick(e.target);
            } else if (e.target.classList.contains('session-add')) {
                this.handleAddSessionClick(e.target);
            }
        });
    }

    handleSessionTabClick(tabElement) {
        const sessionName = tabElement.textContent.trim();
        console.log(`SessionTabs: Clicked session tab: ${sessionName}`);

        // Remove active class from all session tabs
        document.querySelectorAll('.session-tab').forEach(tab => {
            tab.classList.remove('active');
        });

        // Add active class to clicked tab
        tabElement.classList.add('active');

        // Store active session
        this.activeSession = sessionName;

        // Trigger session selection event
        this.onSessionSelect(sessionName);
    }

    handleAddSessionClick(addElement) {
        console.log('SessionTabs: Add session clicked');

        // Generate new session name
        const sessionCount = document.querySelectorAll('.session-tab').length;
        const newSessionName = `Session ${sessionCount + 1}`;

        // Add new session tab
        this.addSessionTab(newSessionName);

        // Trigger session creation event
        this.onSessionCreate(newSessionName);
    }

    addSessionTab(sessionName) {
        const sessionTabsContainer = document.querySelector('.session-tabs');
        const addButton = sessionTabsContainer.querySelector('.session-add');

        // Create new session tab
        const newTab = document.createElement('div');
        newTab.className = 'session-tab';
        newTab.textContent = sessionName;

        // Insert before the add button
        sessionTabsContainer.insertBefore(newTab, addButton);

        // Select the new tab
        this.handleSessionTabClick(newTab);

        console.log(`SessionTabs: Added new session tab: ${sessionName}`);
    }

    onSessionSelect(sessionName) {
        // Emit custom event for session selection
        const event = new CustomEvent('sessionSelect', {
            detail: { sessionName }
        });
        document.dispatchEvent(event);

        console.log(`SessionTabs: Session selected: ${sessionName}`);
    }

    onSessionCreate(sessionName) {
        // Emit custom event for session creation
        const event = new CustomEvent('sessionCreate', {
            detail: { sessionName }
        });
        document.dispatchEvent(event);

        console.log(`SessionTabs: Session creation requested: ${sessionName}`);
    }

    initializeDefaultSession() {
        // Find the first session tab and make it active if none is active
        const firstTab = document.querySelector('.session-tab');
        const activeTab = document.querySelector('.session-tab.active');

        if (firstTab && !activeTab) {
            this.handleSessionTabClick(firstTab);
        } else if (activeTab) {
            this.activeSession = activeTab.textContent.trim();
        }
    }

    // Public methods for external control
    selectSession(sessionName) {
        const tabs = document.querySelectorAll('.session-tab');
        for (const tab of tabs) {
            if (tab.textContent.trim() === sessionName) {
                this.handleSessionTabClick(tab);
                return true;
            }
        }
        return false;
    }

    getCurrentSession() {
        return this.activeSession;
    }

    // Mobile-specific touch enhancements
    enableMobileOptimizations() {
        const sessionTabsContainer = document.querySelector('.session-tabs');

        // Enable smooth scrolling for tab overflow
        sessionTabsContainer.style.scrollBehavior = 'smooth';

        // Add touch feedback with preventDefault to avoid iOS click issues
        sessionTabsContainer.addEventListener('touchstart', (e) => {
            if (e.target.classList.contains('session-tab') || e.target.classList.contains('session-add')) {
                e.target.style.opacity = '0.7';
                console.log('TouchStart on:', e.target.textContent);
            }
        }, { passive: false });

        sessionTabsContainer.addEventListener('touchend', (e) => {
            if (e.target.classList.contains('session-tab') || e.target.classList.contains('session-add')) {
                e.target.style.opacity = '1';
                console.log('TouchEnd on:', e.target.textContent);

                // Trigger click manually on iOS
                setTimeout(() => {
                    e.target.click();
                }, 10);
            }
        }, { passive: false });

        // iOS Safari specific - prevent double tap zoom
        sessionTabsContainer.style.touchAction = 'manipulation';

        console.log('SessionTabs: Mobile optimizations enabled');
    }
}

// Export for global use
window.SessionTabs = SessionTabs;