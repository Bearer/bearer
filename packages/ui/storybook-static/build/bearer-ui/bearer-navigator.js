/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import './chunk-721a4283.js';

class BearerNavigator {
    constructor() {
        this.screens = [];
        this.screenData = {};
        this.next = e => {
            if (e) {
                e.preventDefault();
                e.stopPropagation();
            }
            if (this.hasNext()) {
                this.visibleScreen = this.visibleScreen + 1;
            }
        };
        this.hasNext = () => this.visibleScreen < this.screens.length - 1;
        this.hasPrevious = () => this.visibleScreen > 0;
    }
    set visibleScreen(index) {
        if (this._visibleScreen >= 0) {
            const currentScreen = this.screens[this._visibleScreen];
            if (currentScreen) {
                currentScreen.willDisappear();
                currentScreen.classList.remove('in');
            }
        }
        this._visibleScreen = index;
        const newScreen = this.screens[this._visibleScreen];
        if (newScreen) {
            newScreen.willAppear(this.screenData);
            this.navigationTitle = newScreen.getTitle();
            newScreen.classList.add('in');
        }
    }
    get visibleScreen() {
        return this._visibleScreen;
    }
    scenarioCompletedHandler() {
        this.screenData = {};
        this.visibleScreen = 0;
    }
    stepCompletedHandler(event) {
        event.preventDefault();
        event.stopImmediatePropagation();
        this.screenData = Object.assign({}, this.screenData, event.detail);
        this.next(null);
    }
    prev(e) {
        if (e) {
            e.preventDefault();
            e.stopPropagation();
        }
        if (this.hasPrevious()) {
            this.visibleScreen = this.visibleScreen - 1;
        }
    }
    componentDidLoad() {
        if (this.el.shadowRoot) {
            this.screens = this.el.shadowRoot
                .querySelector('slot:not([name])')['assignedNodes']()
                .filter(node => node.willAppear);
        }
        this.visibleScreen = 0;
    }
    render() {
        return (h("div", null,
            h("h3", { class: "title" },
                h("bearer-navigator-back", { disabled: !this.hasPrevious() }),
                h("slot", { name: "header-name" }, this.navigationTitle)),
            h("slot", null)));
    }
    static get is() { return "bearer-navigator"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "_visibleScreen": {
            "state": true
        },
        "el": {
            "elementRef": true
        },
        "navigationTitle": {
            "state": true
        },
        "screenData": {
            "state": true
        },
        "screens": {
            "state": true
        }
    }; }
    static get listeners() { return [{
            "name": "scenarioCompleted",
            "method": "scenarioCompletedHandler"
        }, {
            "name": "stepCompleted",
            "method": "stepCompletedHandler"
        }, {
            "name": "navigatorGoBack",
            "method": "prev"
        }]; }
    static get style() { return ".title {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  margin-top: 0;\n  margin-left: -0.5rem; }"; }
}

export { BearerNavigator };
