const { defineConfig } = require('cypress');

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    specPattern: 'cypress/e2e/userExperience.cy.js',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
