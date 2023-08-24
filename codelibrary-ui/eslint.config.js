const path = require('path')
const standard = require('eslint-config-standard')
const tsPlugin = require('@typescript-eslint/eslint-plugin')
const tsParser = require('@typescript-eslint/parser')
const importPlugin = require('eslint-plugin-import')
const nPlugin = require('eslint-plugin-n')
const promisePlugin = require('eslint-plugin-promise')

module.exports = [
  {
    plugins: {
      '@typescript-eslint': tsPlugin,
      import: importPlugin,
      n: nPlugin,
      promise: promisePlugin,
    },
    files: [
      '**/*.js',
      '**/*.jsx',
      '**/*.ts',
      '**/*.tsx',
    ],
    languageOptions: {
      parser: tsParser,
    },
    rules: {
      ...standard.rules,
      /**
       * Rules that go against StandardJS are here.
       *
       * StandardJS only gets a few things wrong.
       */

      // Go against StandardJS, and require trailing commas.
      // Trailing commas reduce the size of diffs.
      'comma-dangle': ['warn', {
        arrays: 'always-multiline',
        objects: 'always-multiline',
        imports: 'always-multiline',
        exports: 'always-multiline',
        // Trailing commas for function calls are annoying.
        functions: 'only-multiline',
      }],
      // Don't allow spaces inside braces.
      'object-curly-spacing': ['warn', 'never'],
      // Require NO spaces before the function parentheses.
      'space-before-function-paren': ['warn', 'never'],
      // Complaining about unused variables is too annoying.
      'no-unused-vars': 'off',

      /**
       * Rules that expand on StandardJS are here.
       *
       * These rules enforce even stricter checks.
       */

      'padding-line-between-statements': [
        'warn',
        // Require blank lines before control statements.
        {blankLine: 'always', prev: '*', next: 'block'},
        {blankLine: 'always', prev: '*', next: 'block-like'},
        {blankLine: 'always', prev: '*', next: 'break'},
        {blankLine: 'always', prev: '*', next: 'class'},
        {blankLine: 'always', prev: '*', next: 'continue'},
        {blankLine: 'always', prev: '*', next: 'do'},
        {blankLine: 'always', prev: '*', next: 'for'},
        {blankLine: 'always', prev: '*', next: 'if'},
        {blankLine: 'always', prev: '*', next: 'return'},
        {blankLine: 'always', prev: '*', next: 'switch'},
        {blankLine: 'always', prev: '*', next: 'throw'},
        {blankLine: 'always', prev: '*', next: 'try'},
        {blankLine: 'always', prev: '*', next: 'while'},
        // Do not allow blank lines before switch statement labels.
        {blankLine: 'never', prev: '*', next: 'case'},
        {blankLine: 'never', prev: '*', next: 'default'},
      ],
      // Prefer arrow functions.
      'prefer-arrow-callback': 'warn',
      // Require parentheses to make arrow function bodies less confusing.
      'no-confusing-arrow': ['error', {allowParens: true}],
      // Treat function parameters as if they are `const`.
      'no-param-reassign': 'error',
      // Prefer no quotes for properties.
      'quote-props': ['warn', 'as-needed', {
        keywords: false,
        numbers: false,
      }],

      // Disable ESLint warnings about undefined variables.
      // TypeScript does a better job of checking this.
      'no-undef': 'off',
    },
  },
]
