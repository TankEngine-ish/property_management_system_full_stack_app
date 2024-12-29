describe('User Interface Tests', () => {
  it('should display at least one user', () => {
    cy.visit('cypress/integration/cypress/e2e/userExperience.spec.js'); // Adjust the path if needed

    // Dynamically check for user elements instead of hardcoding an ID
    cy.get('.user-interface .card').first().within(() => {
      cy.get('.text-sm').invoke('text').should('match', /^Id: \d+$/); // Check if an ID exists
      cy.get('.text-lg').should('not.be.empty'); // Check that the name is not empty
      cy.get('.text-md').should('not.be.empty'); // Check that the statement is not empty
    });
  });

  it('should create a new user and verify its presence', () => {
    const newUserName = `Test User ${Date.now()}`;
    const newUserStatement = Math.floor(Math.random() * 100);

    cy.visit('/');

    // Fill in the form to create a new user
    cy.get('form').eq(0).within(() => {
      cy.get('input[placeholder="Name"]').type(newUserName);
      cy.get('input[placeholder="Statement"]').type(newUserStatement.toString());
      cy.get('button[type="submit"]').click();
    });

    // Verify that the new user appears in the list
    cy.contains(newUserName).should('exist');
    cy.contains(newUserStatement).should('exist');
  });

  it('should delete a user and verify its absence', () => {
    cy.visit('/');

    // Dynamically select the first user and delete it
    cy.get('.user-interface .card').first().within(() => {
      cy.get('button').contains('Delete').click();
    });

    // Verify the user no longer exists
    cy.get('.user-interface .card').should('have.length.lessThan', 1);
  });
});
