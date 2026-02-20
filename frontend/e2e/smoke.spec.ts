import { expect, test } from "@playwright/test";

test("home page loads", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("Easy Arbitra")).toBeVisible();
});

test("core pages load", async ({ page }) => {
  await page.goto("/wallets");
  await expect(page.getByText("Wallets")).toBeVisible();

  await page.goto("/markets");
  await expect(page.getByText("Markets")).toBeVisible();

  await page.goto("/leaderboard");
  await expect(page.getByText("Leaderboard")).toBeVisible();

  await page.goto("/anomalies");
  await expect(page.getByText("Anomaly Feed")).toBeVisible();
});
