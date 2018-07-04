import {
    configure
} from '@storybook/html';

const req = require.context("../src/components", true, /.stories.tsx$/);

function loadStories() {
    req.keys().forEach(filename => req(filename));
}

configure(loadStories, module);
