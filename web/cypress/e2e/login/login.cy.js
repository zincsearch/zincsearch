/// <reference types="cypress" />

// Welcome to Cypress!
//
// This spec file contains a variety of sample tests
// for a todo list app that are designed to demonstrate
// the power of writing tests in Cypress.
//
// To learn more about how Cypress works and
// what makes it such an awesome testing tool,
// please read our getting started guide:
// https://on.cypress.io/introduction-to-cypress

describe('login', () => {
  beforeEach(() => {
    // Cypress starts out with a blank slate for each test
    // so we must tell it to visit our website with the `cy.visit()` command.
    // Since we want to visit the same URL at the start of all our tests,
    // we include it in our beforeEach function so that it runs before each test
    cy.visit('http://localhost:8080')
  })

  it('should have user id and password fields along with sign in button', () => {

    // user id field is visible
    cy.get('[data-cy="login-user-id"]').should('be.visible')

    // password field is visible
    cy.get('[data-cy="login-password"]').should('be.visible')

    // Sign In button is visible
    cy.get('[data-cy="login-sign-in"]').should('be.visible')

  })

  it('should be able to login', () => {
    cy.login();
  })
})
