module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(ts|tsx)$': 'babel-jest',
  },
  moduleNameMapper: {
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
  },
  setupFilesAfterEnv: ['@testing-library/jest-dom/extend-expect'], // Custom matchers for DOM testing
  testPathIgnorePatterns: ['/node_modules/', '/dist/'], // Ignore these folders
  transformIgnorePatterns: ['node_modules/(?!@babel/runtime)'], // Allow Babel to transform specific dependencies
};
