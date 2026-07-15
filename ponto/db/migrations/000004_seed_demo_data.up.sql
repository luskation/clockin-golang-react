-- Dados de demonstração para dev/teste local. Remova esta migration (ou o
-- conteúdo dela) antes de rodar contra um banco de produção de verdade.

INSERT INTO companies (id, name, cnpj)
VALUES (gen_random_uuid(), 'Comp Júnior', '12.345.678/0001-90');

INSERT INTO users (id, company_id, name, email, password, role)
VALUES (
    gen_random_uuid(),
    (SELECT id FROM companies WHERE cnpj = '12.345.678/0001-90'),
    'Admin',
    'admin@compjunior.com',
    '$2a$10$NpeQFaqhmz4wXA/8pMaEbeUjTbrP7u6dE1kk8XPJVYb2/Oik5a6.a', -- admin123
    'admin'
);

INSERT INTO users (id, company_id, name, email, password, role)
VALUES (
    gen_random_uuid(),
    (SELECT id FROM companies WHERE cnpj = '12.345.678/0001-90'),
    'Colaborador Teste',
    'colaborador@compjunior.com',
    '$2a$10$CJmQsbEKGlP0fEkExWcGOO7MdA96reAq0FsMZIR1tTZLWw9/rTxN6', -- colab123
    'employee'
);
