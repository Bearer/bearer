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
    static get style() { return "*[data-bearer-button-popover], *[data-bearer-button-popover]::before, *[data-bearer-button-popover]::after {\n  -webkit-box-sizing: border-box;\n  box-sizing: border-box; }\n\nhtml[data-bearer-button-popover] {\n  font-family: sans-serif;\n  line-height: 1.15;\n  -webkit-text-size-adjust: 100%;\n  -ms-text-size-adjust: 100%;\n  -ms-overflow-style: scrollbar;\n  -webkit-tap-highlight-color: rgba(0, 0, 0, 0); }\n\n\@-ms-viewport {\n  width: device-width; }\n\narticle[data-bearer-button-popover], aside[data-bearer-button-popover], figcaption[data-bearer-button-popover], figure[data-bearer-button-popover], footer[data-bearer-button-popover], header[data-bearer-button-popover], hgroup[data-bearer-button-popover], main[data-bearer-button-popover], nav[data-bearer-button-popover], section[data-bearer-button-popover] {\n  display: block; }\n\nbody[data-bearer-button-popover] {\n  margin: 0;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\";\n  font-size: 1em;\n  font-weight: 400;\n  line-height: 1.5;\n  color: #212529;\n  text-align: left;\n  background-color: #fff; }\n\n[tabindex=\"-1\"][data-bearer-button-popover]:focus {\n  outline: 0 !important; }\n\nhr[data-bearer-button-popover] {\n  -webkit-box-sizing: content-box;\n  box-sizing: content-box;\n  height: 0;\n  overflow: visible; }\n\nh1[data-bearer-button-popover], h2[data-bearer-button-popover], h3[data-bearer-button-popover], h4[data-bearer-button-popover], h5[data-bearer-button-popover], h6[data-bearer-button-popover] {\n  margin-top: 0;\n  margin-bottom: 0.5rem; }\n\np[data-bearer-button-popover] {\n  margin-top: 0;\n  margin-bottom: 1rem; }\n\nabbr[title][data-bearer-button-popover], abbr[data-original-title][data-bearer-button-popover] {\n  text-decoration: underline;\n  -webkit-text-decoration: underline dotted;\n  text-decoration: underline dotted;\n  cursor: help;\n  border-bottom: 0; }\n\naddress[data-bearer-button-popover] {\n  margin-bottom: 1rem;\n  font-style: normal;\n  line-height: inherit; }\n\nol[data-bearer-button-popover], ul[data-bearer-button-popover], dl[data-bearer-button-popover] {\n  margin-top: 0;\n  margin-bottom: 1rem; }\n\nol[data-bearer-button-popover]   ol[data-bearer-button-popover], ul[data-bearer-button-popover]   ul[data-bearer-button-popover], ol[data-bearer-button-popover]   ul[data-bearer-button-popover], ul[data-bearer-button-popover]   ol[data-bearer-button-popover] {\n  margin-bottom: 0; }\n\ndt[data-bearer-button-popover] {\n  font-weight: 700; }\n\ndd[data-bearer-button-popover] {\n  margin-bottom: .5rem;\n  margin-left: 0; }\n\nblockquote[data-bearer-button-popover] {\n  margin: 0 0 1rem; }\n\ndfn[data-bearer-button-popover] {\n  font-style: italic; }\n\nb[data-bearer-button-popover], strong[data-bearer-button-popover] {\n  font-weight: bolder; }\n\nsmall[data-bearer-button-popover] {\n  font-size: 80%; }\n\nsub[data-bearer-button-popover], sup[data-bearer-button-popover] {\n  position: relative;\n  font-size: 75%;\n  line-height: 0;\n  vertical-align: baseline; }\n\nsub[data-bearer-button-popover] {\n  bottom: -.25em; }\n\nsup[data-bearer-button-popover] {\n  top: -.5em; }\n\na[data-bearer-button-popover] {\n  color: #007bff;\n  text-decoration: none;\n  background-color: transparent;\n  -webkit-text-decoration-skip: objects; }\n  a[data-bearer-button-popover]:hover {\n    color: #0056b3;\n    text-decoration: underline; }\n\na[data-bearer-button-popover]:not([href]):not([tabindex]) {\n  color: inherit;\n  text-decoration: none; }\n  a[data-bearer-button-popover]:not([href]):not([tabindex]):hover, a[data-bearer-button-popover]:not([href]):not([tabindex]):focus {\n    color: inherit;\n    text-decoration: none; }\n  a[data-bearer-button-popover]:not([href]):not([tabindex]):focus {\n    outline: 0; }\n\npre[data-bearer-button-popover], code[data-bearer-button-popover], kbd[data-bearer-button-popover], samp[data-bearer-button-popover] {\n  font-family: SFMono-Regular, Menlo, Monaco, Consolas, \"Liberation Mono\", \"Courier New\", monospace;\n  font-size: 1em; }\n\npre[data-bearer-button-popover] {\n  margin-top: 0;\n  margin-bottom: 1rem;\n  overflow: auto;\n  -ms-overflow-style: scrollbar; }\n\nfigure[data-bearer-button-popover] {\n  margin: 0 0 1rem; }\n\nimg[data-bearer-button-popover] {\n  vertical-align: middle;\n  border-style: none; }\n\nsvg[data-bearer-button-popover]:not(:root) {\n  overflow: hidden; }\n\ntable[data-bearer-button-popover] {\n  border-collapse: collapse; }\n\ncaption[data-bearer-button-popover] {\n  padding-top: 0.75rem;\n  padding-bottom: 0.75rem;\n  color: #6c757d;\n  text-align: left;\n  caption-side: bottom; }\n\nth[data-bearer-button-popover] {\n  text-align: inherit; }\n\nlabel[data-bearer-button-popover] {\n  display: inline-block;\n  margin-bottom: 0.5rem; }\n\nbutton[data-bearer-button-popover] {\n  border-radius: 0; }\n\nbutton[data-bearer-button-popover]:focus {\n  outline: 1px dotted;\n  outline: 5px auto -webkit-focus-ring-color; }\n\ninput[data-bearer-button-popover], button[data-bearer-button-popover], select[data-bearer-button-popover], optgroup[data-bearer-button-popover], textarea[data-bearer-button-popover] {\n  margin: 0;\n  font-family: inherit;\n  font-size: inherit;\n  line-height: inherit; }\n\nbutton[data-bearer-button-popover], input[data-bearer-button-popover] {\n  overflow: visible; }\n\nbutton[data-bearer-button-popover], select[data-bearer-button-popover] {\n  text-transform: none; }\n\nbutton[data-bearer-button-popover], html[data-bearer-button-popover]   [type=\"button\"][data-bearer-button-popover], [type=\"reset\"][data-bearer-button-popover], [type=\"submit\"][data-bearer-button-popover] {\n  -webkit-appearance: button; }\n\nbutton[data-bearer-button-popover]::-moz-focus-inner, [type=\"button\"][data-bearer-button-popover]::-moz-focus-inner, [type=\"reset\"][data-bearer-button-popover]::-moz-focus-inner, [type=\"submit\"][data-bearer-button-popover]::-moz-focus-inner {\n  padding: 0;\n  border-style: none; }\n\ninput[type=\"radio\"][data-bearer-button-popover], input[type=\"checkbox\"][data-bearer-button-popover] {\n  -webkit-box-sizing: border-box;\n  box-sizing: border-box;\n  padding: 0; }\n\ninput[type=\"date\"][data-bearer-button-popover], input[type=\"time\"][data-bearer-button-popover], input[type=\"datetime-local\"][data-bearer-button-popover], input[type=\"month\"][data-bearer-button-popover] {\n  -webkit-appearance: listbox; }\n\ntextarea[data-bearer-button-popover] {\n  overflow: auto;\n  resize: vertical; }\n\nfieldset[data-bearer-button-popover] {\n  min-width: 0;\n  padding: 0;\n  margin: 0;\n  border: 0; }\n\nlegend[data-bearer-button-popover] {\n  display: block;\n  width: 100%;\n  max-width: 100%;\n  padding: 0;\n  margin-bottom: .5rem;\n  font-size: 1.5rem;\n  line-height: inherit;\n  color: inherit;\n  white-space: normal; }\n\nprogress[data-bearer-button-popover] {\n  vertical-align: baseline; }\n\n[type=\"number\"][data-bearer-button-popover]::-webkit-inner-spin-button, [type=\"number\"][data-bearer-button-popover]::-webkit-outer-spin-button {\n  height: auto; }\n\n[type=\"search\"][data-bearer-button-popover] {\n  outline-offset: -2px;\n  -webkit-appearance: none; }\n\n[type=\"search\"][data-bearer-button-popover]::-webkit-search-cancel-button, [type=\"search\"][data-bearer-button-popover]::-webkit-search-decoration {\n  -webkit-appearance: none; }\n\n[data-bearer-button-popover]::-webkit-file-upload-button {\n  font: inherit;\n  -webkit-appearance: button; }\n\noutput[data-bearer-button-popover] {\n  display: inline-block; }\n\nsummary[data-bearer-button-popover] {\n  display: list-item;\n  cursor: pointer; }\n\ntemplate[data-bearer-button-popover] {\n  display: none; }\n\n[hidden][data-bearer-button-popover] {\n  display: none !important; }\n\n.popover[data-bearer-button-popover] {\n  position: absolute;\n  top: 0;\n  left: 0;\n  z-index: 1060;\n  display: block;\n  max-width: 276px;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, \"Helvetica Neue\", Arial, sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\";\n  font-style: normal;\n  font-weight: 400;\n  line-height: 1.5;\n  text-align: left;\n  text-align: start;\n  text-decoration: none;\n  text-shadow: none;\n  text-transform: none;\n  letter-spacing: normal;\n  word-break: normal;\n  word-spacing: normal;\n  white-space: normal;\n  line-break: auto;\n  font-size: 1em;\n  word-wrap: break-word;\n  background-color: #fff;\n  background-clip: padding-box;\n  border: 1px solid rgba(0, 0, 0, 0.2);\n  border-radius: 0.3rem; }\n  .popover[data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n    position: absolute;\n    display: block;\n    width: 1rem;\n    height: 0.5rem;\n    margin: 0 0.3rem; }\n    .popover[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .popover[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n      position: absolute;\n      display: block;\n      content: \"\";\n      border-color: transparent;\n      border-style: solid; }\n\n.bs-popover-top[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover] {\n  margin-bottom: 0.5rem; }\n  .bs-popover-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n    bottom: calc((0.5rem + 1px) * -1); }\n  .bs-popover-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    border-width: 0.5rem 0.5rem 0; }\n  .bs-popover-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before {\n    bottom: 0;\n    border-top-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"top\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    bottom: 1px;\n    border-top-color: #fff; }\n\n.bs-popover-right[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover] {\n  margin-left: 0.5rem; }\n  .bs-popover-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n    left: calc((0.5rem + 1px) * -1);\n    width: 0.5rem;\n    height: 1rem;\n    margin: 0.3rem 0; }\n  .bs-popover-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    border-width: 0.5rem 0.5rem 0.5rem 0; }\n  .bs-popover-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before {\n    left: 0;\n    border-right-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"right\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    left: 1px;\n    border-right-color: #fff; }\n\n.bs-popover-bottom[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover] {\n  margin-top: 0.5rem; }\n  .bs-popover-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n    top: calc((0.5rem + 1px) * -1); }\n  .bs-popover-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    border-width: 0 0.5rem 0.5rem 0.5rem; }\n  .bs-popover-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before {\n    top: 0;\n    border-bottom-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    top: 1px;\n    border-bottom-color: #fff; }\n  .bs-popover-bottom[data-bearer-button-popover]   .popover-header[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"bottom\"][data-bearer-button-popover]   .popover-header[data-bearer-button-popover]::before {\n    position: absolute;\n    top: 0;\n    left: 50%;\n    display: block;\n    width: 1rem;\n    margin-left: -0.5rem;\n    content: \"\";\n    border-bottom: 1px solid #f7f7f7; }\n\n.bs-popover-left[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover] {\n  margin-right: 0.5rem; }\n  .bs-popover-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover], .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n    right: calc((0.5rem + 1px) * -1);\n    width: 0.5rem;\n    height: 1rem;\n    margin: 0.3rem 0; }\n  .bs-popover-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    border-width: 0.5rem 0 0.5rem 0.5rem; }\n  .bs-popover-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before, .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::before {\n    right: 0;\n    border-left-color: rgba(0, 0, 0, 0.25); }\n  \n  .bs-popover-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after, .bs-popover-auto[x-placement^=\"left\"][data-bearer-button-popover]   .arrow[data-bearer-button-popover]::after {\n    right: 1px;\n    border-left-color: #fff; }\n\n.popover-header[data-bearer-button-popover] {\n  padding: 0.5rem 0.75rem;\n  margin-bottom: 0;\n  font-size: 1em;\n  color: inherit;\n  background-color: #f7f7f7;\n  border-bottom: 1px solid #ebebeb;\n  border-top-left-radius: calc(0.3rem - 1px);\n  border-top-right-radius: calc(0.3rem - 1px); }\n  .popover-header[data-bearer-button-popover]:empty {\n    display: none; }\n\n.popover-body[data-bearer-button-popover] {\n  padding: 0.5rem 0.75rem;\n  color: #212529; }\n\n.root[data-bearer-button-popover] {\n  position: relative;\n  display: inline-block; }\n\n.popover[data-bearer-button-popover] {\n  min-width: 250px;\n  min-height: 272px;\n  max-height: 80vh;\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  -webkit-box-orient: vertical;\n  -webkit-box-direction: normal;\n  -ms-flex-direction: column;\n  flex-direction: column; }\n  .popover.direction-bottom[data-bearer-button-popover] {\n    top: 100%;\n    -webkit-transform: translateX(-50%);\n    transform: translateX(-50%);\n    margin-left: 50%; }\n    .popover.direction-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n      left: 50%;\n      -webkit-transform: translateX(-50%);\n      transform: translateX(-50%);\n      margin: 0; }\n      .popover.direction-bottom[data-bearer-button-popover]   .arrow[data-bearer-button-popover]:after {\n        border-bottom-color: #f7f7f7; }\n  .popover.direction-top[data-bearer-button-popover] {\n    bottom: 100%;\n    left: 0;\n    top: initial;\n    -webkit-transform: translateX(-50%);\n    transform: translateX(-50%);\n    margin-left: 50%; }\n    .popover.direction-top[data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n      left: 50%;\n      -webkit-transform: translateX(-50%);\n      transform: translateX(-50%);\n      margin: 0; }\n  .popover.direction-left[data-bearer-button-popover] {\n    top: calc((-50px / 2) - (0.5rem / 2));\n    left: initial;\n    right: 100%; }\n    .popover.direction-left[data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n      -webkit-transform: translateY(-50%);\n      transform: translateY(-50%);\n      top: 50px;\n      margin: 0; }\n  .popover.direction-right[data-bearer-button-popover] {\n    left: 100%;\n    top: calc((-50px / 2) - (0.5rem / 2)); }\n    .popover.direction-right[data-bearer-button-popover]   .arrow[data-bearer-button-popover] {\n      -webkit-transform: translateY(-50%);\n      transform: translateY(-50%);\n      margin-top: 0;\n      top: 50px; }\n\n.popover-body[data-bearer-button-popover] {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  overflow: auto;\n  max-height: calc(80vh - 50px);\n  padding: 0;\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }\n\n.header[data-bearer-button-popover] {\n  padding-left: 25px; }\n\n.header-arrow[data-bearer-button-popover] {\n  position: absolute;\n  left: 5px; }"; }
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
