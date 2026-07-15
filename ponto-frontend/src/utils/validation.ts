export const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
// Precisa ficar idêntico a cnpjRegex em ponto/internal/service/company_service.go —
// validar aqui só evita uma viagem à API para um erro que o backend rejeitaria de qualquer forma.
export const cnpjRegex = /^\d{2}\.\d{3}\.\d{3}\/\d{4}-\d{2}$/;

export function isValidEmail(value: string): boolean {
  return emailRegex.test(value);
}

export function isValidCNPJ(value: string): boolean {
  return cnpjRegex.test(value);
}

export function isValidPassword(value: string): boolean {
  return value.length >= 6;
}

export function formatCNPJ(value: string): string {
  const digits = value.replace(/\D/g, "").slice(0, 14);
  const blocks = [digits.slice(0, 2), digits.slice(2, 5), digits.slice(5, 8), digits.slice(8, 12), digits.slice(12, 14)];

  let result = blocks[0];
  if (blocks[1]) result += `.${blocks[1]}`;
  if (blocks[2]) result += `.${blocks[2]}`;
  if (blocks[3]) result += `/${blocks[3]}`;
  if (blocks[4]) result += `-${blocks[4]}`;
  return result;
}
