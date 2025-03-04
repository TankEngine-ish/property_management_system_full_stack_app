describe("User API and Frontend Integration E2E Test", () => {
  const apiBaseUrl = "http://localhost:8000";
  const frontendUrl = "http://localhost:3000";

  it("Creates, fetches, validates, and deletes a user dynamically via API and validates on the frontend", () => {
    
    cy.request("POST", `${apiBaseUrl}/api/repair/users`, {
      name: "Geralt of Rivia",
      statement: 300,
    }).then((response) => {
      expect(response.status).to.eq(200); // creates Geralt
      const userId = response.body.id;

     
      cy.request("GET", `${apiBaseUrl}/api/repair/users/${userId}`).then((fetchResponse) => {
        expect(fetchResponse.status).to.eq(200); // fetches Geralt
        expect(fetchResponse.body).to.have.property("name", "Geralt of Rivia");
        expect(fetchResponse.body).to.have.property("statement", 300);
      });

      
      cy.visit(frontendUrl);
      cy.contains(`Id: ${userId}`); 
      cy.contains("Geralt of Rivia"); 
      cy.contains("300"); 

      
      cy.request("DELETE", `${apiBaseUrl}/api/repair/users/${userId}`).then((deleteResponse) => {
        expect(deleteResponse.status).to.eq(200); //deletes Geralt

        
        cy.request({
          method: "GET",
          url: `${apiBaseUrl}/api/repair/users/${userId}`,
          failOnStatusCode: false, 
        }).then((verifyResponse) => {
          expect(verifyResponse.status).to.eq(404); // making sure Geralt is gone
        });

        // finally we refresh the frontend and validate Geralt is no longer displayed
        cy.visit(frontendUrl);
        cy.contains(`Id: ${userId}`).should("not.exist");
        cy.contains("Geralt of Rivia").should("not.exist");
      });
    });
  });
});
