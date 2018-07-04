/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import './chunk-721a4283.js';

const Icon = () => (h("svg", { fill: "currentColor", height: "24", viewBox: "0 0 24 24", width: "24", xmlns: "http://www.w3.org/2000/svg" },
    h("path", { d: "M15.41 16.09l-4.58-4.59 4.58-4.59L14 5.5l-6 6 6 6z" }),
    h("path", { d: "M0-.5h24v24H0z", fill: "none" })));

class BearerNavigatorBack {
    constructor() {
        this.disabled = false;
        this.back = () => {
            this.navigatorGoBack.emit();
        };
    }
    render() {
        return (h("button", { onClick: this.back, disabled: this.disabled },
            h(Icon, null)));
    }
    static get is() { return "bearer-navigator-back"; }
    static get properties() { return {
        "disabled": {
            "type": Boolean,
            "attr": "disabled"
        }
    }; }
    static get events() { return [{
            "name": "navigatorGoBack",
            "method": "navigatorGoBack",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }]; }
    static get style() { return "button {\n  background: none;\n  border: none;\n  padding: 0; }\n  button:not([disabled]) {\n    cursor: pointer; }"; }
}

export { BearerNavigatorBack };
