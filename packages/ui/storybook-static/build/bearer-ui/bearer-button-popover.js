/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { d as process } from './chunk-721a4283.js';

class BearerButtonPopover {
    constructor() {
        this.visible = true;
        this.direction = 'top';
        this.arrow = true;
        this.btnKind = 'primary';
        this.toggleDisplay = e => {
            e.preventDefault();
            this.visible = !this.visible;
        };
    }
    clickOutsideHandler() {
        this.visible = false;
    }
    clickInsideHandler(ev) {
        ev.stopImmediatePropagation();
    }
    toggle(opened) {
        this.visible = opened;
    }
    componentDidLoad() {
        if (this.opened === false) {
            this.visible = false;
        }
    }
    render() {
        return (h("div", { class: "root" },
            h("bearer-button", { onClick: this.toggleDisplay, kind: this.btnKind },
                h("slot", { name: "buttonText" })),
            this.visible && (h("div", { class: `popover fade show bs-popover-${this.direction} direction-${this.direction}` },
                h("h3", { class: "popover-header" },
                    this.backNav && h("bearer-navigator-back", { class: "header-arrow" }),
                    h("span", { class: "header" }, this.header)),
                h("div", { class: "popover-body" },
                    h("slot", null)),
                this.arrow && h("div", { class: "arrow" })))));
    }
    static get is() { return "bearer-button-popover"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "arrow": {
            "type": Boolean,
            "attr": "arrow"
        },
        "backNav": {
            "type": Boolean,
            "attr": "back-nav"
        },
        "btnKind": {
            "type": String,
            "attr": "btn-kind"
        },
        "direction": {
            "type": String,
            "attr": "direction"
        },
        "header": {
            "type": String,
            "attr": "header"
        },
        "opened": {
            "type": Boolean,
            "attr": "opened"
        },
        "toggle": {
            "method": true
        },
        "visible": {
            "state": true
        }
    }; }
    static get listeners() { return [{
            "name": "body:click",
            "method": "clickOutsideHandler"
        }, {
            "name": "click",
            "method": "clickInsideHandler"
        }]; }
    static get style() { return "*,\n*::before,\n*::after {\n  -webkit-box-sizing: border-box;\n  box-sizing: border-box; }\n\nhtml {\n  font-family: sans-serif;\n  line-height: 1.15;\n  -webkit-text-size-adjust: 100%;\n  -ms-text-size-adjust: 100%;\n  -ms-overflow-style: scrollbar;\n  -webkit-tap-highlight-color: rgba(0, 0, 0, 0); }\n\n\@-ms-viewport {\n  width: device-width; }\n\narticle, aside, figcaption, figure, footer, header, hgroup, main, nav, section {\n  display: block; }\n\nbody {\n  margin: 0;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\";\n  font-size: 1em;\n  font-weight: 400;\n  line-height: 1.5;\n  color: #212529;\n  text-align: left;\n  background-color: #fff; }\n\n[tabindex=\"-1\"]:focus {\n  outline: 0 !important; }\n\nhr {\n  -webkit-box-sizing: content-box;\n  box-sizing: content-box;\n  height: 0;\n  overflow: visible; }\n\nh1, h2, h3, h4, h5, h6 {\n  margin-top: 0;\n  margin-bottom: 0.5rem; }\n\np {\n  margin-top: 0;\n  margin-bottom: 1rem; }\n\nabbr[title],\nabbr[data-original-title] {\n  text-decoration: underline;\n  -webkit-text-decoration: underline dotted;\n  text-decoration: underline dotted;\n  cursor: help;\n  border-bottom: 0; }\n\naddress {\n  margin-bottom: 1rem;\n  font-style: normal;\n  line-height: inherit; }\n\nol,\nul,\ndl {\n  margin-top: 0;\n  margin-bottom: 1rem; }\n\nol ol,\nul ul,\nol ul,\nul ol {\n  margin-bottom: 0; }\n\ndt {\n  font-weight: 700; }\n\ndd {\n  margin-bottom: .5rem;\n  margin-left: 0; }\n\nblockquote {\n  margin: 0 0 1rem; }\n\ndfn {\n  font-style: italic; }\n\nb,\nstrong {\n  font-weight: bolder; }\n\nsmall {\n  font-size: 80%; }\n\nsub,\nsup {\n  position: relative;\n  font-size: 75%;\n  line-height: 0;\n  vertical-align: baseline; }\n\nsub {\n  bottom: -.25em; }\n\nsup {\n  top: -.5em; }\n\na {\n  color: #007bff;\n  text-decoration: none;\n  background-color: transparent;\n  -webkit-text-decoration-skip: objects; }\n  a:hover {\n    color: #0056b3;\n    text-decoration: underline; }\n\na:not([href]):not([tabindex]) {\n  color: inherit;\n  text-decoration: none; }\n  a:not([href]):not([tabindex]):hover, a:not([href]):not([tabindex]):focus {\n    color: inherit;\n    text-decoration: none; }\n  a:not([href]):not([tabindex]):focus {\n    outline: 0; }\n\npre,\ncode,\nkbd,\nsamp {\n  font-family: SFMono-Regular, Menlo, Monaco, Consolas, \"Liberation Mono\", \"Courier New\", monospace;\n  font-size: 1em; }\n\npre {\n  margin-top: 0;\n  margin-bottom: 1rem;\n  overflow: auto;\n  -ms-overflow-style: scrollbar; }\n\nfigure {\n  margin: 0 0 1rem; }\n\nimg {\n  vertical-align: middle;\n  border-style: none; }\n\nsvg:not(:root) {\n  overflow: hidden; }\n\ntable {\n  border-collapse: collapse; }\n\ncaption {\n  padding-top: 0.75rem;\n  padding-bottom: 0.75rem;\n  color: #6c757d;\n  text-align: left;\n  caption-side: bottom; }\n\nth {\n  text-align: inherit; }\n\nlabel {\n  display: inline-block;\n  margin-bottom: 0.5rem; }\n\nbutton {\n  border-radius: 0; }\n\nbutton:focus {\n  outline: 1px dotted;\n  outline: 5px auto -webkit-focus-ring-color; }\n\ninput,\nbutton,\nselect,\noptgroup,\ntextarea {\n  margin: 0;\n  font-family: inherit;\n  font-size: inherit;\n  line-height: inherit; }\n\nbutton,\ninput {\n  overflow: visible; }\n\nbutton,\nselect {\n  text-transform: none; }\n\nbutton,\nhtml [type=\"button\"],\n[type=\"reset\"],\n[type=\"submit\"] {\n  -webkit-appearance: button; }\n\nbutton::-moz-focus-inner,\n[type=\"button\"]::-moz-focus-inner,\n[type=\"reset\"]::-moz-focus-inner,\n[type=\"submit\"]::-moz-focus-inner {\n  padding: 0;\n  border-style: none; }\n\ninput[type=\"radio\"],\ninput[type=\"checkbox\"] {\n  -webkit-box-sizing: border-box;\n  box-sizing: border-box;\n  padding: 0; }\n\ninput[type=\"date\"],\ninput[type=\"time\"],\ninput[type=\"datetime-local\"],\ninput[type=\"month\"] {\n  -webkit-appearance: listbox; }\n\ntextarea {\n  overflow: auto;\n  resize: vertical; }\n\nfieldset {\n  min-width: 0;\n  padding: 0;\n  margin: 0;\n  border: 0; }\n\nlegend {\n  display: block;\n  width: 100%;\n  max-width: 100%;\n  padding: 0;\n  margin-bottom: .5rem;\n  font-size: 1.5rem;\n  line-height: inherit;\n  color: inherit;\n  white-space: normal; }\n\nprogress {\n  vertical-align: baseline; }\n\n[type=\"number\"]::-webkit-inner-spin-button,\n[type=\"number\"]::-webkit-outer-spin-button {\n  height: auto; }\n\n[type=\"search\"] {\n  outline-offset: -2px;\n  -webkit-appearance: none; }\n\n[type=\"search\"]::-webkit-search-cancel-button,\n[type=\"search\"]::-webkit-search-decoration {\n  -webkit-appearance: none; }\n\n::-webkit-file-upload-button {\n  font: inherit;\n  -webkit-appearance: button; }\n\noutput {\n  display: inline-block; }\n\nsummary {\n  display: list-item;\n  cursor: pointer; }\n\ntemplate {\n  display: none; }\n\n[hidden] {\n  display: none !important; }\n\n.popover {\n  position: absolute;\n  top: 0;\n  left: 0;\n  z-index: 1060;\n  display: block;\n  max-width: 276px;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\";\n  font-style: normal;\n  font-weight: 400;\n  line-height: 1.5;\n  text-align: left;\n  text-align: start;\n  text-decoration: none;\n  text-shadow: none;\n  text-transform: none;\n  letter-spacing: normal;\n  word-break: normal;\n  word-spacing: normal;\n  white-space: normal;\n  line-break: auto;\n  font-size: 1em;\n  word-wrap: break-word;\n  background-color: #fff;\n  background-clip: padding-box;\n  border: 1px solid rgba(0, 0, 0, 0.2);\n  border-radius: 0.3rem; }\n  .popover .arrow {\n    position: absolute;\n    display: block;\n    width: 1rem;\n    height: 0.5rem;\n    margin: 0 0.3rem; }\n    .popover .arrow::before, .popover .arrow::after {\n      position: absolute;\n      display: block;\n      content: \"\";\n      border-color: transparent;\n      border-style: solid; }\n\n.bs-popover-top, .bs-popover-auto[x-placement^=\"top\"] {\n  margin-bottom: 0.5rem; }\n  .bs-popover-top .arrow, .bs-popover-auto[x-placement^=\"top\"] .arrow {\n    bottom: calc((0.5rem + 1px) * -1); }\n  .bs-popover-top .arrow::before, .bs-popover-auto[x-placement^=\"top\"] .arrow::before,\n  .bs-popover-top .arrow::after,\n  .bs-popover-auto[x-placement^=\"top\"] .arrow::after {\n    border-width: 0.5rem 0.5rem 0; }\n  .bs-popover-top .arrow::before, .bs-popover-auto[x-placement^=\"top\"] .arrow::before {\n    bottom: 0;\n    border-top-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-top .arrow::after,\n  .bs-popover-auto[x-placement^=\"top\"] .arrow::after {\n    bottom: 1px;\n    border-top-color: #fff; }\n\n.bs-popover-right, .bs-popover-auto[x-placement^=\"right\"] {\n  margin-left: 0.5rem; }\n  .bs-popover-right .arrow, .bs-popover-auto[x-placement^=\"right\"] .arrow {\n    left: calc((0.5rem + 1px) * -1);\n    width: 0.5rem;\n    height: 1rem;\n    margin: 0.3rem 0; }\n  .bs-popover-right .arrow::before, .bs-popover-auto[x-placement^=\"right\"] .arrow::before,\n  .bs-popover-right .arrow::after,\n  .bs-popover-auto[x-placement^=\"right\"] .arrow::after {\n    border-width: 0.5rem 0.5rem 0.5rem 0; }\n  .bs-popover-right .arrow::before, .bs-popover-auto[x-placement^=\"right\"] .arrow::before {\n    left: 0;\n    border-right-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-right .arrow::after,\n  .bs-popover-auto[x-placement^=\"right\"] .arrow::after {\n    left: 1px;\n    border-right-color: #fff; }\n\n.bs-popover-bottom, .bs-popover-auto[x-placement^=\"bottom\"] {\n  margin-top: 0.5rem; }\n  .bs-popover-bottom .arrow, .bs-popover-auto[x-placement^=\"bottom\"] .arrow {\n    top: calc((0.5rem + 1px) * -1); }\n  .bs-popover-bottom .arrow::before, .bs-popover-auto[x-placement^=\"bottom\"] .arrow::before,\n  .bs-popover-bottom .arrow::after,\n  .bs-popover-auto[x-placement^=\"bottom\"] .arrow::after {\n    border-width: 0 0.5rem 0.5rem 0.5rem; }\n  .bs-popover-bottom .arrow::before, .bs-popover-auto[x-placement^=\"bottom\"] .arrow::before {\n    top: 0;\n    border-bottom-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-bottom .arrow::after,\n  .bs-popover-auto[x-placement^=\"bottom\"] .arrow::after {\n    top: 1px;\n    border-bottom-color: #fff; }\n  .bs-popover-bottom .popover-header::before, .bs-popover-auto[x-placement^=\"bottom\"] .popover-header::before {\n    position: absolute;\n    top: 0;\n    left: 50%;\n    display: block;\n    width: 1rem;\n    margin-left: -0.5rem;\n    content: \"\";\n    border-bottom: 1px solid #f7f7f7; }\n\n.bs-popover-left, .bs-popover-auto[x-placement^=\"left\"] {\n  margin-right: 0.5rem; }\n  .bs-popover-left .arrow, .bs-popover-auto[x-placement^=\"left\"] .arrow {\n    right: calc((0.5rem + 1px) * -1);\n    width: 0.5rem;\n    height: 1rem;\n    margin: 0.3rem 0; }\n  .bs-popover-left .arrow::before, .bs-popover-auto[x-placement^=\"left\"] .arrow::before,\n  .bs-popover-left .arrow::after,\n  .bs-popover-auto[x-placement^=\"left\"] .arrow::after {\n    border-width: 0.5rem 0 0.5rem 0.5rem; }\n  .bs-popover-left .arrow::before, .bs-popover-auto[x-placement^=\"left\"] .arrow::before {\n    right: 0;\n    border-left-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-left .arrow::after,\n  .bs-popover-auto[x-placement^=\"left\"] .arrow::after {\n    right: 1px;\n    border-left-color: #fff; }\n\n.popover-header {\n  padding: 0.5rem 0.75rem;\n  margin-bottom: 0;\n  font-size: 1em;\n  color: inherit;\n  background-color: #f7f7f7;\n  border-bottom: 1px solid #ebebeb;\n  border-top-left-radius: calc(0.3rem - 1px);\n  border-top-right-radius: calc(0.3rem - 1px); }\n  .popover-header:empty {\n    display: none; }\n\n.popover-body {\n  padding: 0.5rem 0.75rem;\n  color: #212529; }\n\n.root {\n  position: relative;\n  display: inline-block; }\n\n.popover {\n  min-width: 250px;\n  min-height: 272px;\n  max-height: 80vh;\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  -webkit-box-orient: vertical;\n  -webkit-box-direction: normal;\n  -ms-flex-direction: column;\n  flex-direction: column; }\n  .popover.direction-bottom {\n    top: 100%;\n    -webkit-transform: translateX(-50%);\n    transform: translateX(-50%);\n    margin-left: 50%; }\n    .popover.direction-bottom .arrow {\n      left: 50%;\n      -webkit-transform: translateX(-50%);\n      transform: translateX(-50%);\n      margin: 0; }\n      .popover.direction-bottom .arrow:after {\n        border-bottom-color: #f7f7f7; }\n  .popover.direction-top {\n    bottom: 100%;\n    left: 0;\n    top: initial;\n    -webkit-transform: translateX(-50%);\n    transform: translateX(-50%);\n    margin-left: 50%; }\n    .popover.direction-top .arrow {\n      left: 50%;\n      -webkit-transform: translateX(-50%);\n      transform: translateX(-50%);\n      margin: 0; }\n  .popover.direction-left {\n    top: calc((-50px / 2) - (0.5rem / 2));\n    left: initial;\n    right: 100%; }\n    .popover.direction-left .arrow {\n      -webkit-transform: translateY(-50%);\n      transform: translateY(-50%);\n      top: 50px;\n      margin: 0; }\n  .popover.direction-right {\n    left: 100%;\n    top: calc((-50px / 2) - (0.5rem / 2)); }\n    .popover.direction-right .arrow {\n      -webkit-transform: translateY(-50%);\n      transform: translateY(-50%);\n      margin-top: 0;\n      top: 50px; }\n\n.popover-body {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  overflow: auto;\n  max-height: calc(80vh - 50px);\n  padding: 0;\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }\n\n.header {\n  padding-left: 25px; }\n\n.header-arrow {\n  position: absolute;\n  left: 5px; }"; }
}

const NAVIGATOR_AUTH_SCREEN_NAME = 'BEARER-NAVIGATOR-AUTH-SCREEN';
class BearerPopoverNavigator {
    constructor() {
        this.screens = [];
        this.screenData = {};
        this.direction = 'top';
        this.button = 'Activate';
        this.header = 'Popover Navigator';
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
        this.hasAuthScreen = () => this.screenNodes.filter(node => node['tagName'] === NAVIGATOR_AUTH_SCREEN_NAME).length;
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
        this.visibleScreen = this.hasAuthScreen() ? 1 : 0;
        this.el.shadowRoot.querySelector('#button')['toggle'](false);
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
    get screenNodes() {
        return this.el.shadowRoot
            ? this.el.shadowRoot.querySelector('slot:not([name])')['assignedNodes']()
            : [];
    }
    componentDidLoad() {
        this.screens = this.screenNodes;
        this.visibleScreen = 0;
    }
    render() {
        return (h("bearer-button-popover", { id: "button", direction: this.direction, header: this.navigationTitle, backNav: this.hasPrevious() },
            h("span", { slot: "buttonText" }, this.button),
            h("slot", null)));
    }
    static get is() { return "bearer-popover-navigator"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "_visibleScreen": {
            "state": true
        },
        "button": {
            "type": String,
            "attr": "button"
        },
        "direction": {
            "type": String,
            "attr": "direction"
        },
        "el": {
            "elementRef": true
        },
        "header": {
            "state": true
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
    static get style() { return ""; }
}

export { BearerButtonPopover, BearerPopoverNavigator };
