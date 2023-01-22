module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    'plugin:react/recommended',
    'airbnb',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  plugins: [
    'react',
    'typescript',
    '@typescript-eslint',
  ],
  settings: {
    'import/resolver': {
      node: {
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    },
  },
  rules: {
    "semi": ['error', 'never'],
    'react/jsx-filename-extension': 'off',
    "prefer-arrow-callback": [
      "error",
      { "allowNamedFunctions": true }
    ],
    "func-style": [
      "error",
      "expression",
      { "allowArrowFunctions": true }
    ]
  },
}