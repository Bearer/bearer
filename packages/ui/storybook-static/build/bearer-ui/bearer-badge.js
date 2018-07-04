/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { b as Bearer } from './chunk-721a4283.js';

class BearerBadge {
    render() {
        return (h("span", { class: `badge badge-${this.kind}` },
            h("slot", null)));
    }
    static get is() { return "bearer-badge"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "kind": {
            "type": String,
            "attr": "kind"
        }
    }; }
    static get style() { return ".badge {\n  display: inline-block;\n  padding: 0.25em 0.4em;\n  font-size: 75%;\n  font-weight: 700;\n  line-height: 1;\n  text-align: center;\n  white-space: nowrap;\n  vertical-align: baseline;\n  border-radius: 0.25rem; }\n  .badge:empty {\n    display: none; }\n\n.btn .badge {\n  position: relative;\n  top: -1px; }\n\n.badge-pill {\n  padding-right: 0.6em;\n  padding-left: 0.6em;\n  border-radius: 10rem; }\n\n.badge-primary {\n  color: #fff;\n  background-color: #007bff; }\n  .badge-primary[href]:hover, .badge-primary[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #0062cc; }\n\n.badge-secondary {\n  color: #fff;\n  background-color: #6c757d; }\n  .badge-secondary[href]:hover, .badge-secondary[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #545b62; }\n\n.badge-success {\n  color: #fff;\n  background-color: #28a745; }\n  .badge-success[href]:hover, .badge-success[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #1e7e34; }\n\n.badge-info {\n  color: #fff;\n  background-color: #17a2b8; }\n  .badge-info[href]:hover, .badge-info[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #117a8b; }\n\n.badge-warning {\n  color: #212529;\n  background-color: #ffc107; }\n  .badge-warning[href]:hover, .badge-warning[href]:focus {\n    color: #212529;\n    text-decoration: none;\n    background-color: #d39e00; }\n\n.badge-danger {\n  color: #fff;\n  background-color: #dc3545; }\n  .badge-danger[href]:hover, .badge-danger[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #bd2130; }\n\n.badge-light {\n  color: #212529;\n  background-color: #f8f9fa; }\n  .badge-light[href]:hover, .badge-light[href]:focus {\n    color: #212529;\n    text-decoration: none;\n    background-color: #dae0e5; }\n\n.badge-dark {\n  color: #fff;\n  background-color: #343a40; }\n  .badge-dark[href]:hover, .badge-dark[href]:focus {\n    color: #fff;\n    text-decoration: none;\n    background-color: #1d2124; }"; }
}

class BearerSetupDisplay {
    constructor() {
        this.scenarioId = '';
        this.isSetup = false;
        this.setupId = '';
    }
    componentDidLoad() {
        Bearer.emitter.addListener(`setup_success:${this.scenarioId}`, data => {
            this.setupId = data.referenceID;
            this.isSetup = true;
        });
    }
    render() {
        if (this.isSetup) {
            return (h("p", null,
                "Scenario is currently setup with Setup ID:\u00A0",
                h("bearer-badge", { kind: "info" }, this.setupId)));
        }
        else {
            return (h("p", null,
                h("p", null, "Scenario hasn't been setup yet")));
        }
    }
    static get is() { return "bearer-setup-display"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "isSetup": {
            "state": true
        },
        "scenarioId": {
            "type": String,
            "attr": "scenario-id"
        },
        "setupId": {
            "state": true
        }
    }; }
}

export { BearerBadge, BearerSetupDisplay };
