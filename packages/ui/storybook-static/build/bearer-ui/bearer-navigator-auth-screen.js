/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { b as Bearer, e as Events$2 } from './chunk-721a4283.js';

class BearerNavigatorAuthScreen {
    constructor() {
        this.sessionInitialized = false;
        this.scenarioAuthorized = false;
        this.setupId = 'BEARER_SCENARIO_ID';
        this.authenticate = () => {
            Bearer.instance.askAuthorizations({
                scenarioId: this.setupId,
                setupId: this.setupId
            });
        };
        this.revoke = () => {
            Bearer.instance.revokeAuthorization(this.setupId);
        };
    }
    willAppear() {
        this.el.shadowRoot.querySelector('#screen')['willAppear']();
    }
    willDisappear() {
        this.el.shadowRoot.querySelector('#screen')['willDisappear']();
    }
    getTitle() {
        return 'Authentication';
    }
    componentDidLoad() {
        Bearer.instance.maybeInitialized.then(() => {
            this.sessionInitialized = true;
            this.scenarioAuthorized = Bearer.instance.hasAuthorized(this.setupId);
            this.listener = Bearer.emitter.addListener(Events$2.SCENARIO_AUTHORIZED,
            // TODO: we need to ensure the tokens are not confused
            ({ scenarioId, authorized }) => {
                if (this.setupId === scenarioId) {
                    this.scenarioAuthorized = authorized;
                    this.goNext();
                }
            });
            this.goNext();
        });
    }
    componentDidUnload() {
        if (this.listener) {
            this.listener.remove();
            this.listener = null;
        }
    }
    goNext() {
        if (this.scenarioAuthorized) {
            this.scenarioAuthenticate.emit();
            this.stepCompleted.emit();
        }
    }
    render() {
        return (h("bearer-navigator-screen", { id: "screen", navigationTitle: "Authentication", class: "in" }, this.sessionInitialized &&
            (this.scenarioAuthorized ? (h("bearer-button", { kind: "warning", onClick: this.revoke }, "Logout")) : (h("bearer-button", { kind: "primary", onClick: this.authenticate }, "Login")))));
    }
    static get is() { return "bearer-navigator-auth-screen"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "el": {
            "elementRef": true
        },
        "getTitle": {
            "method": true
        },
        "scenarioAuthorized": {
            "state": true
        },
        "sessionInitialized": {
            "state": true
        },
        "setupId": {
            "type": String,
            "attr": "setup-id"
        },
        "willAppear": {
            "method": true
        },
        "willDisappear": {
            "method": true
        }
    }; }
    static get events() { return [{
            "name": "scenarioAuthenticate",
            "method": "scenarioAuthenticate",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }, {
            "name": "stepCompleted",
            "method": "stepCompleted",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }]; }
    static get style() { return ":host {\n  display: none; }\n\n:host(.in) {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }\n\n:host(bearer-navigator-auth-screen.in) {\n  -ms-flex-item-align: center;\n  align-self: center;\n  width: 100%;\n  text-align: center;\n  -webkit-box-flex: initial;\n  -ms-flex: initial;\n  flex: initial; }\n\n:host(bearer-navigator-auth-screen.in) bearer-navigator-screen {\n  -webkit-box-pack: center;\n  -ms-flex-pack: center;\n  justify-content: center; }\n\n.screen {\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }"; }
}

class BearerNavigatorScreen {
    constructor() {
        this.visible = false;
        this.next = data => {
            const payload = this.name ? { [this.name]: data } : data;
            this.stepCompleted.emit(payload);
        };
        this.prev = () => {
            this.navigatorGoBack.emit();
        };
    }
    willAppear(data) {
        this.data = data;
        this.visible = true;
    }
    willDisappear() {
        this.visible = false;
    }
    getTitle() {
        if (typeof this.navigationTitle === 'string') {
            return this.navigationTitle;
        }
        return this.navigationTitle(this.data);
    }
    completeScreenHandler({ detail }) {
        this.next(detail);
    }
    render() {
        if (!this.visible) {
            return false;
        }
        return (h("div", { class: "screen" }, this.renderFunc ? (this.renderFunc({ data: this.data, next: this.next, prev: this.prev })) : (h("slot", null))));
    }
    static get is() { return "bearer-navigator-screen"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "data": {
            "state": true
        },
        "getTitle": {
            "method": true
        },
        "name": {
            "type": String,
            "attr": "name"
        },
        "navigationTitle": {
            "type": "Any",
            "attr": "navigation-title"
        },
        "renderFunc": {
            "type": "Any",
            "attr": "render-func"
        },
        "visible": {
            "state": true
        },
        "willAppear": {
            "method": true
        },
        "willDisappear": {
            "method": true
        }
    }; }
    static get events() { return [{
            "name": "stepCompleted",
            "method": "stepCompleted",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }, {
            "name": "navigatorGoBack",
            "method": "navigatorGoBack",
            "bubbles": true,
            "cancelable": true,
            "composed": true
        }]; }
    static get listeners() { return [{
            "name": "completeScreen",
            "method": "completeScreenHandler"
        }]; }
    static get style() { return ":host {\n  display: none; }\n\n:host(.in) {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }\n\n:host(bearer-navigator-auth-screen.in) {\n  -ms-flex-item-align: center;\n  align-self: center;\n  width: 100%;\n  text-align: center;\n  -webkit-box-flex: initial;\n  -ms-flex: initial;\n  flex: initial; }\n\n:host(bearer-navigator-auth-screen.in) bearer-navigator-screen {\n  -webkit-box-pack: center;\n  -ms-flex-pack: center;\n  justify-content: center; }\n\n.screen {\n  -webkit-box-flex: 1;\n  -ms-flex: 1;\n  flex: 1; }"; }
}

export { BearerNavigatorAuthScreen, BearerNavigatorScreen };
