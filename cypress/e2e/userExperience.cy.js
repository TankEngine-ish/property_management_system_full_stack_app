describe("User API and Frontend Integration E2E Test", () => {
  // we define the local environment
  const apiBaseUrl = "http://localhost:8000";
  const frontendUrl = "http://localhost:3000";
  

  const testUserName = `cypress-test-${Math.floor(Math.random() * 10000)}`;
  const testUserStatement = Math.floor(Math.random() * 1000);
  
  beforeEach(() => {
    cy.request("GET", `${apiBaseUrl}/api/repair/users`).then((response) => {
      const users = response.body;
      users.forEach(user => {
        if (user.name && user.name.startsWith('cypress-test-')) {
          cy.request("DELETE", `${apiBaseUrl}/api/repair/users/${user.id}`);
        }
      });
    });
    
    cy.request("GET", `${apiBaseUrl}/health`).then((response) => {
      expect(response.status).to.eq(200);
      expect(response.body).to.eq("OK");
    });
  });

  it("Creates a user manually through the UI and verifies it", () => {
    cy.visit(frontendUrl);
    cy.wait(3000); 
    
    cy.get('input[placeholder="Name"]').type(testUserName);
    cy.get('input[placeholder="Statement"]').type(testUserStatement.toString());
    cy.contains('button', 'Add User').click();
    cy.wait(3000);
    cy.contains(testUserName, { timeout: 10000 }).should("be.visible");
    cy.contains(testUserName)
      .parents('div.flex.items-center.justify-between')
      .find('button')
      .contains('Delete User')
      .click();
    cy.wait(2000);
    cy.contains(testUserName).should("not.exist");
  });

  it("Verifies basic API functionality", () => {
    cy.request("POST", `${apiBaseUrl}/api/repair/users`, {
      name: `${testUserName}-api`,
      statement: testUserStatement
    }).then((response) => {
      expect(response.status).to.eq(200);
      const userId = response.body.id;
      cy.request("GET", `${apiBaseUrl}/api/repair/users/${userId}`).then((getResponse) => {
        expect(getResponse.status).to.eq(200);
        expect(getResponse.body.name).to.eq(`${testUserName}-api`);
        
        // Clean up after yourself
        cy.request("DELETE", `${apiBaseUrl}/api/repair/users/${userId}`);
      });
    });
  });

  it("Tests the Update functionality through the UI", () => {
    cy.visit(frontendUrl);
    cy.wait(3000);
    
    const updateUserName = `cypress-update-${Math.floor(Math.random() * 10000)}`;
    cy.get('input[placeholder="Name"]').type(updateUserName);
    cy.get('input[placeholder="Statement"]').type('100');
    cy.contains('button', 'Add User').click();
    
    cy.wait(3000);
    
    cy.contains(updateUserName)
      .parents('div.flex.items-center.justify-between')
      .find('div.text-sm.text-gray-600')
      .invoke('text')
      .then((idText) => {
        const userId = idText.replace('Id: ', '').trim();
        cy.log(`Found user ID: ${userId} for user: ${updateUserName}`);
        
        cy.get('input[placeholder="User Id"]').type(userId);
        cy.get('input[placeholder="New Name"]').type(`${updateUserName}-updated`);
        cy.get('input[placeholder="New statement"]').type('200');
        cy.contains('button', 'Update User').click();
        
        cy.wait(3000);
        cy.contains(`${updateUserName}-updated`).should('be.visible');
        cy.contains(`${updateUserName}-updated`)
          .parents('div.flex.items-center.justify-between')
          .find('button')
          .contains('Delete User')
          .click();
      });
  });
});