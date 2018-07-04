/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { d as process, b as Bearer } from './chunk-721a4283.js';

class BearerDropdownButton {
    constructor() {
        this.visible = true;
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
        if (this.innerListener) {
            Bearer.emitter.addListener(this.innerListener, () => {
                this.visible = false;
            });
        }
    }
    render() {
        return (h("div", { class: "root" },
            h("bearer-button", { onClick: this.toggleDisplay, kind: this.btnKind },
                h("slot", { name: "buttonText" }),
                h("span", { class: "symbol" }, "\u25BE")),
            this.visible && (h("div", { class: "dropdown-down" },
                h("slot", null)))));
    }
    static get is() { return "bearer-dropdown-button"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "btnKind": {
            "type": String,
            "attr": "btn-kind"
        },
        "innerListener": {
            "type": String,
            "attr": "inner-listener"
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
    static get style() { return ".root[data-bearer-dropdown-button] {\n  position: relative;\n  display: inline-block; }\n  .root[data-bearer-dropdown-button]   .dropdown-down[data-bearer-dropdown-button] {\n    position: absolute;\n    border: solid 1px silver;\n    background-color: white;\n    z-index: 1;\n    padding: 20px 10px 10px 10px;\n    min-width: 320px;\n    max-height: 80vh;\n    border-radius: 3px; }\n  .root[data-bearer-dropdown-button]   .symbol[data-bearer-dropdown-button] {\n    padding-left: 10px; }"; }
}

export { BearerDropdownButton };
