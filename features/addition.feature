Feature: Addition

  Scenario: Add two numbers
    Given I have a number 1
    And I have another number 1
    When I add the numbers
    Then the result should be 2
