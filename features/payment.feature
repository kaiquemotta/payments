Feature: Criar Pagamento

  Scenario: Criar pagamento com sucesso
    Given que tenho um pagamento válido
    When eu tentar criar o pagamento
    Then o pagamento deve ser criado com sucesso

