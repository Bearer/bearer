/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { a as classnames, b as Bearer, c as BearerState } from './chunk-721a4283.js';

class Alert {
    constructor() {
        this.kind = 'primary';
    }
    render() {
        const classes = classnames({
            alert: true,
            [`alert-${this.kind}`]: true,
            'alert-dismissible': this.onDismiss
        });
        return (h("div", { class: classes },
            h("slot", null),
            this.onDismiss && (h("button", { type: "button", class: "close", "data-dismiss": "alert", "aria-label": "Close", onClick: this.onDismiss },
                h("span", { "aria-hidden": "true" }, "\u00D7")))));
    }
    static get is() { return "bearer-alert"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "kind": {
            "type": String,
            "attr": "kind"
        },
        "onDismiss": {
            "type": "Any",
            "attr": "on-dismiss"
        }
    }; }
    static get style() { return "[data-bearer-alert-slot] > hr {\n  margin-top: 1rem;\n  margin-bottom: 1rem;\n  border: 0;\n  border-top: 1px solid rgba(0, 0, 0, 0.1); }\n\n.alert[data-bearer-alert] {\n  position: relative;\n  padding: 0.75rem 1.25rem;\n  margin-bottom: 1rem;\n  border: 1px solid transparent;\n  border-radius: 0.25rem; }\n\n.alert-heading[data-bearer-alert] {\n  color: inherit; }\n\n.alert-link[data-bearer-alert] {\n  font-weight: 700; }\n\n.alert-dismissible[data-bearer-alert] {\n  padding-right: 4rem; }\n  .alert-dismissible[data-bearer-alert]   .close[data-bearer-alert] {\n    position: absolute;\n    top: 0;\n    right: 0;\n    padding: 0.75rem 1.25rem;\n    color: inherit; }\n\n.alert-primary[data-bearer-alert] {\n  color: #004085;\n  background-color: #cce5ff;\n  border-color: #b8daff; }\n  .alert-primary [data-bearer-alert-slot] > hr {\n    border-top-color: #9fcdff; }\n  .alert-primary [data-bearer-alert-slot] > .alert-link {\n    color: #002752; }\n\n.alert-secondary[data-bearer-alert] {\n  color: #383d41;\n  background-color: #e2e3e5;\n  border-color: #d6d8db; }\n  .alert-secondary [data-bearer-alert-slot] > hr {\n    border-top-color: #c8cbcf; }\n  .alert-secondary [data-bearer-alert-slot] > .alert-link {\n    color: #202326; }\n\n.alert-success[data-bearer-alert] {\n  color: #155724;\n  background-color: #d4edda;\n  border-color: #c3e6cb; }\n  .alert-success [data-bearer-alert-slot] > hr {\n    border-top-color: #b1dfbb; }\n  .alert-success [data-bearer-alert-slot] > .alert-link {\n    color: #0b2e13; }\n\n.alert-info[data-bearer-alert] {\n  color: #0c5460;\n  background-color: #d1ecf1;\n  border-color: #bee5eb; }\n  .alert-info [data-bearer-alert-slot] > hr {\n    border-top-color: #abdde5; }\n  .alert-info [data-bearer-alert-slot] > .alert-link {\n    color: #062c33; }\n\n.alert-warning[data-bearer-alert] {\n  color: #856404;\n  background-color: #fff3cd;\n  border-color: #ffeeba; }\n  .alert-warning [data-bearer-alert-slot] > hr {\n    border-top-color: #ffe8a1; }\n  .alert-warning [data-bearer-alert-slot] > .alert-link {\n    color: #533f03; }\n\n.alert-danger[data-bearer-alert] {\n  color: #721c24;\n  background-color: #f8d7da;\n  border-color: #f5c6cb; }\n  .alert-danger [data-bearer-alert-slot] > hr {\n    border-top-color: #f1b0b7; }\n  .alert-danger [data-bearer-alert-slot] > .alert-link {\n    color: #491217; }\n\n.alert-light[data-bearer-alert] {\n  color: #818182;\n  background-color: #fefefe;\n  border-color: #fdfdfe; }\n  .alert-light [data-bearer-alert-slot] > hr {\n    border-top-color: #ececf6; }\n  .alert-light [data-bearer-alert-slot] > .alert-link {\n    color: #686868; }\n\n.alert-dark[data-bearer-alert] {\n  color: #1b1e21;\n  background-color: #d6d8d9;\n  border-color: #c6c8ca; }\n  .alert-dark [data-bearer-alert-slot] > hr {\n    border-top-color: #b9bbbe; }\n  .alert-dark [data-bearer-alert-slot] > .alert-link {\n    color: #040505; }\n\n.close[data-bearer-alert] {\n  float: right;\n  font-size: 1.5rem;\n  font-weight: 700;\n  line-height: 1;\n  color: #000;\n  text-shadow: 0 1px 0 #fff;\n  opacity: .5; }\n  .close[data-bearer-alert]:hover, .close[data-bearer-alert]:focus {\n    color: #000;\n    text-decoration: none;\n    opacity: .75; }\n  .close[data-bearer-alert]:not(:disabled):not(.disabled) {\n    cursor: pointer; }\n\nbutton.close[data-bearer-alert] {\n  padding: 0;\n  background-color: transparent;\n  border: 0;\n  -webkit-appearance: none; }"; }
}

class FieldSet {
    constructor(set) {
        this.set = set;
    }
    get(controlName) {
        return this.set.find(el => el.controlName === controlName);
    }
    getValue(controlName) {
        return this.get(controlName).value;
    }
    setValue(controlName, value) {
        this.set.map(el => {
            if (el.controlName === controlName) {
                el.value = value;
                return el;
            }
            return el;
        });
    }
    map(func) {
        return this.set.map(func);
    }
    reduce(func) {
        return this.set.reduce(func);
    }
    filter(func) {
        return this.set.filter(func);
    }
}

const OAuth2SetupType = [
    {
        type: 'text',
        label: 'Client ID',
        controlName: 'clientId',
        required: true
    },
    {
        type: 'password',
        label: 'Client Secret',
        controlName: 'clientSecret',
        required: true
    }
];

const EmailSetupType = [
    {
        type: 'email',
        label: 'Email',
        controlName: 'email',
        required: true
    },
    {
        type: 'password',
        label: 'Password',
        controlName: 'password',
        required: true
    }
];

const KeySetupType = [
    {
        type: 'text',
        label: 'Enter your Key',
        controlName: 'key',
        required: true
    }
];

class BearerSetup {
    constructor() {
        this.fields = [];
        this.error = false;
        this.loading = false;
        this.handleSubmit = (e) => {
            e.preventDefault();
            this.loading = true;
            const secretSet = this.fieldSet.map(el => {
                return { key: el.controlName, value: el.value };
            });
            const publicSet = this.fieldSet
                .filter(el => el.type !== 'password')
                .map(el => {
                return { key: el.controlName, value: el.value };
            });
            const setupId = BearerState.generateUniqueId(30);
            BearerState.storeSecret(setupId, secretSet.reduce((acc, obj) => (Object.assign({}, acc, { [obj['key']]: obj['value'] })), {}))
                .then(() => {
                this.error = false;
                return BearerState.storeData(`${setupId}setup`, publicSet.reduce((acc, obj) => (Object.assign({}, acc, { [obj['key']]: obj['value'] })), {}));
            })
                .then(_ => {
                this.loading = false;
                Bearer.emitter.emit(`setup_success:${this.scenarioId}`, {
                    // clientID: this.inputs.clientID,
                    referenceID: setupId
                });
            })
                .catch(() => {
                this.error = true;
                this.loading = false;
                Bearer.emitter.emit(`setup_error:${this.scenarioId}`, {});
            });
        };
    }
    componentWillLoad() {
        if (typeof this.fields !== 'string') {
            this.fieldSet = new FieldSet(this.fields);
            return;
        }
        switch (this.fields) {
            case 'email':
                this.fieldSet = new FieldSet(EmailSetupType);
                break;
            case 'type':
                this.fieldSet = new FieldSet(KeySetupType);
                break;
            case 'oauth2':
            default:
                this.fieldSet = new FieldSet(OAuth2SetupType);
        }
    }
    componentDidLoad() {
        const form = this.element.shadowRoot.querySelector('bearer-form');
        if (this.referenceId) {
            BearerState.getData(`${this.referenceId}setup`)
                .then((data) => {
                Object.keys(data.Item).forEach(key => {
                    if (data.Item.hasOwnProperty(key) &&
                        key !== 'ReadAllowed' &&
                        key !== 'referenceId') {
                        this.fieldSet.setValue(key, data.Item[key]);
                    }
                });
                form.updateFieldSet(this.fieldSet);
                console.debug('[BEARER]', 'get_setup_success', data);
                Bearer.emitter.emit(`setup_success:${this.scenarioId}`, {
                    referenceID: this.referenceId
                });
            })
                .catch(e => {
                console.error('[BEARER]', 'get_setup_error', e);
            });
        }
    }
    render() {
        return [
            this.error && (h("bearer-alert", { kind: "danger" }, "[Error] Unable to store the credentials")),
            this.loading ? (h("bearer-loading", null)) : (h("bearer-form", { fields: this.fieldSet, clearOnInput: true, onSubmit: this.handleSubmit }))
        ];
    }
    static get is() { return "bearer-setup"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "element": {
            "elementRef": true
        },
        "error": {
            "state": true
        },
        "fields": {
            "type": String,
            "attr": "fields"
        },
        "fieldSet": {
            "state": true
        },
        "loading": {
            "state": true
        },
        "referenceId": {
            "type": String,
            "attr": "reference-id"
        },
        "scenarioId": {
            "type": String,
            "attr": "scenario-id"
        }
    }; }
    static get events() { return [{
            "name": "stepCompleted",
            "method": "stepCompleted",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }]; }
    static get style() { return ""; }
}

export { Alert as BearerAlert, BearerSetup };
