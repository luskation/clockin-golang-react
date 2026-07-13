import { formatDate, formatTime } from "./datetime";

describe("formatTime", () => {
  it("formata a hora em 24h com dois dígitos", () => {
    const date = new Date(2026, 0, 5, 9, 5, 3);
    expect(formatTime(date)).toBe("09:05:03");
  });
});

describe("formatDate", () => {
  it("capitaliza o dia da semana e inclui dia/mês", () => {
    const date = new Date(2026, 6, 13);

    const formatted = formatDate(date);

    expect(formatted[0]).toBe(formatted[0].toUpperCase());
    expect(formatted).toContain("13");
    expect(formatted).toContain("julho");
  });
});
