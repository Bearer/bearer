/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import { b as Bearer } from './chunk-721a4283.js';

class Button {
    render() {
        const { loadingComponent } = Bearer.config;
        if (loadingComponent) {
            const Tag = loadingComponent;
            return h(Tag, null);
        }
        return (h("div", { id: "root" },
            h("div", { id: "loader" },
                h("div", { id: "d1" }),
                h("div", { id: "d2" }),
                h("div", { id: "d3" }),
                h("div", { id: "d4" }),
                h("div", { id: "d5" }))));
    }
    static get is() { return "bearer-loading"; }
    static get encapsulation() { return "shadow"; }
    static get style() { return "[data-bearer-loading-host] {\n  --dot-size: 15px;\n  --dot-stretch-width: 30px;\n  --dot-stretch-height: 12px;\n  --dot-count: 5; }\n\n#root[data-bearer-loading] {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  -webkit-box-pack: center;\n  -ms-flex-pack: center;\n  justify-content: center;\n  margin-top: 50px;\n  margin-bottom: 50px; }\n\n#loader[data-bearer-loading] {\n  position: relative;\n  height: var(--dot-size);\n  width: calc(var(--dot-size) * var(--dot-count));\n  margin: var(--dot-stretch-height) 0; }\n\n#loader[data-bearer-loading]   div[data-bearer-loading] {\n  width: var(--dot-size);\n  height: var(--dot-size);\n  border-radius: 50%;\n  position: absolute; }\n\n#d1[data-bearer-loading] {\n  background: #3179eb;\n  -webkit-animation: animate 3s linear infinite;\n  animation: animate 3s linear infinite; }\n\n#d2[data-bearer-loading] {\n  background: #4260e2;\n  -webkit-animation: animate 3s linear infinite -0.6s;\n  animation: animate 3s linear infinite -0.6s; }\n\n#d3[data-bearer-loading] {\n  background: #575fe7;\n  -webkit-animation: animate 3s linear infinite -1.2s;\n  animation: animate 3s linear infinite -1.2s; }\n\n#d4[data-bearer-loading] {\n  background: #7651e7;\n  -webkit-animation: animate 3s linear infinite -1.8s;\n  animation: animate 3s linear infinite -1.8s; }\n\n#d5[data-bearer-loading] {\n  background: #8d62ea;\n  -webkit-animation: animate 3s linear infinite -2.4s;\n  animation: animate 3s linear infinite -2.4s; }\n\n\@-webkit-keyframes animate {\n  0% {\n    left: calc(var(--dot-size) * var(--dot-count));\n    top: 0; }\n  80% {\n    left: 0;\n    top: 0; }\n  85% {\n    left: 0;\n    width: var(--dot-size);\n    height: var(--dot-size); }\n  90% {\n    width: var(--dot-stretch-width);\n    height: var(--dot-stretch-height); }\n  95% {\n    width: 0;\n    left: calc(var(--dot-size) * var(--dot-count)); }\n  100% {\n    left: calc(var(--dot-size) * var(--dot-count));\n    top: 0; } }"; }
}

export { Button as BearerLoading };
