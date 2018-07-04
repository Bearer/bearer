/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import './chunk-721a4283.js';

class AuthConfig {
    handleSubmit() { }
    render() {
        return (h("form", { onSubmit: () => this.handleSubmit() },
            h("bearer-input", { type: "text", label: "hello", controlName: "hello", value: "Hello", hint: "hello" }),
            h("bearer-input", { type: "submit", onSubmit: () => this.handleSubmit() })));
    }
    static get is() { return "auth-config"; }
    static get encapsulation() { return "shadow"; }
    static get events() { return [{
            "name": "submit",
            "method": "submit",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }]; }
    static get style() { return ""; }
}

export { AuthConfig };
