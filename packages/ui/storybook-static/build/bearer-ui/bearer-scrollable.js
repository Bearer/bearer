/*! Built with http://stenciljs.com */
const { h } = window.BearerUi;

import './chunk-721a4283.js';

class BearerScrollable {
    constructor() {
        this.perPage = 5;
        this.hasMore = true;
        this.page = 1;
        this.fetching = false;
        this.collection = [];
        this._renderCollection = () => {
            if (this.fetching && !this.collection.length) {
                return null;
            }
            return this.renderCollection(this.collection);
        };
        this._renderFetching = () => this.renderFetching ? this.renderFetching() : h("bearer-loading", null);
        this.onScroll = () => {
            if (!this.fetching) {
                if (this.content.scrollTop + this.content.clientHeight >=
                    this.content.scrollHeight) {
                    this.fetchNext();
                }
            }
        };
    }
    fetchNext() {
        if (this.store && this.hasMore) {
            const redux = this.store.getState()[this.reducer];
            this.fetching = true;
            this.hasMore = redux.length > this.collection.length;
            this.collection = [
                ...this.collection,
                ...redux.splice((this.page - 1) * this.perPage, this.page * this.perPage)
            ];
            this.fetching = false;
            this.page = this.page + 1;
        }
        if (!this.store && this.hasMore) {
            this.fetching = true;
            this.fetcher({ page: this.page })
                .then(({ items }) => {
                this.hasMore = items.length === this.perPage;
                this.collection = [...this.collection, ...items];
                this.fetching = false;
                this.page = this.page + 1;
            })
                .catch(() => {
                this.fetching = false;
            });
        }
    }
    componentDidLoad() {
        if (this.element) {
            this.content = this.element.querySelector('.scrollable-list');
            this.content.addEventListener('scroll', this.onScroll);
            this.fetchNext();
        }
    }
    reset() {
        this.hasMore = true;
        this.collection = [];
        this.fetchNext();
    }
    render() {
        return (h("div", { class: "scrollable-list" },
            this._renderCollection(),
            this.fetching && this._renderFetching()));
    }
    static get is() { return "bearer-scrollable"; }
    static get properties() { return {
        "collection": {
            "state": true
        },
        "content": {
            "state": true
        },
        "element": {
            "elementRef": true
        },
        "fetcher": {
            "type": "Any",
            "attr": "fetcher"
        },
        "fetching": {
            "state": true
        },
        "hasMore": {
            "state": true
        },
        "page": {
            "state": true
        },
        "perPage": {
            "type": Number,
            "attr": "per-page"
        },
        "reducer": {
            "type": String,
            "attr": "reducer"
        },
        "renderCollection": {
            "type": "Any",
            "attr": "render-collection"
        },
        "renderFetching": {
            "type": "Any",
            "attr": "render-fetching"
        },
        "reset": {
            "method": true
        },
        "store": {
            "type": "Any",
            "attr": "store"
        }
    }; }
    static get listeners() { return [{
            "name": "BearerScrollableNext",
            "method": "fetchNext"
        }]; }
    static get style() { return ".scrollable-list {\n  height: 300px;\n  overflow-x: hidden;\n  overflow-y: scroll; }\n\n.scrollable-list::-webkit-scrollbar {\n  display: none; }"; }
}

export { BearerScrollable };
