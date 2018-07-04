const path = require('path');

module.exports = (storybookBaseConfig, configType, defaultConfig) => {
  defaultConfig.module.rules.push({
    test: [/\.(js)$/],
    loaders: [require.resolve('@storybook/addon-storysource/loader')],
    enforce: 'pre',
  });

  defaultConfig.module.rules.push({
    test: /\.(ts|tsx)$/,
    loaders: [
      require.resolve('@storybook/addon-storysource/loader'),
    ],
    enforce: 'pre'
  });
  defaultConfig.resolve.extensions.push(".ts", ".tsx");

  return defaultConfig;
};
