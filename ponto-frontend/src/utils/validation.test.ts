import { formatCNPJ, isValidCNPJ, isValidEmail, isValidPassword } from "./validation";

describe("isValidEmail", () => {
  it("aceita e-mails válidos", () => {
    expect(isValidEmail("ana@example.com")).toBe(true);
  });

  it("rejeita e-mails sem @", () => {
    expect(isValidEmail("ana-example.com")).toBe(false);
  });

  it("rejeita e-mails com espaço", () => {
    expect(isValidEmail("ana @example.com")).toBe(false);
  });
});

describe("isValidCNPJ", () => {
  it("aceita o formato mascarado", () => {
    expect(isValidCNPJ("12.345.678/0001-90")).toBe(true);
  });

  it("rejeita CNPJ sem máscara", () => {
    expect(isValidCNPJ("12345678000190")).toBe(false);
  });
});

describe("isValidPassword", () => {
  it("exige ao menos 6 caracteres", () => {
    expect(isValidPassword("12345")).toBe(false);
    expect(isValidPassword("123456")).toBe(true);
  });
});

describe("formatCNPJ", () => {
  it("aplica a máscara progressivamente enquanto digita", () => {
    expect(formatCNPJ("12")).toBe("12");
    expect(formatCNPJ("12345")).toBe("12.345");
    expect(formatCNPJ("12345678")).toBe("12.345.678");
    expect(formatCNPJ("123456780001")).toBe("12.345.678/0001");
    expect(formatCNPJ("12345678000190")).toBe("12.345.678/0001-90");
  });

  it("ignora caracteres não numéricos e limita a 14 dígitos", () => {
    expect(formatCNPJ("12.345.678/0001-90extra")).toBe("12.345.678/0001-90");
  });
});
