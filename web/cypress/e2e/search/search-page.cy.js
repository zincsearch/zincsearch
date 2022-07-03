/// <reference types="cypress" />

describe("Search Page", () => {
    beforeEach(() => {
        cy.visit("http://localhost:8080");

        cy.get('[data-cy="login-user-id"]').type('admin')

        // login
        cy.get('[data-cy="login-password"]').type('admin')
        cy.get('[data-cy="login-sign-in"]').click()
        cy.contains('admin').should('be.visible')
    })

    it("should have basic fields", () => {
        // search bar is visible
        cy.get("[data-cy='search-bar-input']").should("be.visible");

        // syntax guide button is visible
        cy.get("[data-cy='syntax-guide-button']").should("be.visible");

        // date-time button is visible
        cy.get("[data-cy='date-time-button']").should("be.visible");

        // search bar refresh button is visible
        cy.get("[data-cy='search-bar-refresh-button']").should("be.visible");

        // search button is visible
        cy.get("[data-cy='search-bar-button-dropdown']").should("be.visible");

        // index drop down field is visible
        cy.get("[data-cy='index-dropdown']").should("be.visible");

        // index field search input is visible
        cy.get("[data-cy='index-field-search-input']").should("be.visible");

        // search result area is visible
        cy.get("[data-cy='search-result-area']").should("be.visible");
    })
})