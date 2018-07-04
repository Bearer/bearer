/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import './chunk-721a4283.js';

class Typography {
    constructor() {
        this.as = 'p';
        this.kind = '';
    }
    render() {
        const Tag = this.as;
        return (h(Tag, { class: this.kind },
            h("slot", null)));
    }
    static get is() { return "bearer-typography"; }
    static get encapsulation() { return "shadow"; }
    static get properties() { return {
        "as": {
            "type": String,
            "attr": "as"
        },
        "kind": {
            "type": String,
            "attr": "kind"
        }
    }; }
    static get style() { return "h1[data-bearer-typography], h2[data-bearer-typography], h3[data-bearer-typography], h4[data-bearer-typography], h5[data-bearer-typography], h6[data-bearer-typography], .h1[data-bearer-typography], .h2[data-bearer-typography], .h3[data-bearer-typography], .h4[data-bearer-typography], .h5[data-bearer-typography], .h6[data-bearer-typography] {\n  margin-bottom: 0.5rem;\n  font-family: inherit;\n  font-weight: 500;\n  line-height: 1.2;\n  color: inherit; }\n\nh1[data-bearer-typography], .h1[data-bearer-typography] {\n  font-size: 2.5rem; }\n\nh2[data-bearer-typography], .h2[data-bearer-typography] {\n  font-size: 2rem; }\n\nh3[data-bearer-typography], .h3[data-bearer-typography] {\n  font-size: 1.75rem; }\n\nh4[data-bearer-typography], .h4[data-bearer-typography] {\n  font-size: 1.5rem; }\n\nh5[data-bearer-typography], .h5[data-bearer-typography] {\n  font-size: 1.25rem; }\n\nh6[data-bearer-typography], .h6[data-bearer-typography] {\n  font-size: 1rem; }\n\n.lead[data-bearer-typography] {\n  font-size: 1.25rem;\n  font-weight: 300; }\n\n.display-1[data-bearer-typography] {\n  font-size: 6rem;\n  font-weight: 300;\n  line-height: 1.2; }\n\n.display-2[data-bearer-typography] {\n  font-size: 5.5rem;\n  font-weight: 300;\n  line-height: 1.2; }\n\n.display-3[data-bearer-typography] {\n  font-size: 4.5rem;\n  font-weight: 300;\n  line-height: 1.2; }\n\n.display-4[data-bearer-typography] {\n  font-size: 3.5rem;\n  font-weight: 300;\n  line-height: 1.2; }\n\nhr[data-bearer-typography] {\n  margin-top: 1rem;\n  margin-bottom: 1rem;\n  border: 0;\n  border-top: 1px solid rgba(0, 0, 0, 0.1); }\n\nsmall[data-bearer-typography], .small[data-bearer-typography] {\n  font-size: 80%;\n  font-weight: 400; }\n\nmark[data-bearer-typography], .mark[data-bearer-typography] {\n  padding: 0.2em;\n  background-color: #fcf8e3; }\n\n.list-unstyled[data-bearer-typography] {\n  padding-left: 0;\n  list-style: none; }\n\n.list-inline[data-bearer-typography] {\n  padding-left: 0;\n  list-style: none; }\n\n.list-inline-item[data-bearer-typography] {\n  display: inline-block; }\n  .list-inline-item[data-bearer-typography]:not(:last-child) {\n    margin-right: 0.5rem; }\n\n.initialism[data-bearer-typography] {\n  font-size: 90%;\n  text-transform: uppercase; }\n\n.blockquote[data-bearer-typography] {\n  margin-bottom: 1rem;\n  font-size: 1.25rem; }\n\n.blockquote-footer[data-bearer-typography] {\n  display: block;\n  font-size: 80%;\n  color: #6c757d; }\n  .blockquote-footer[data-bearer-typography]::before {\n    content: \"\\2014 \\00A0\"; }"; }
}

export { Typography as BearerTypography };
