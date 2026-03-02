import {
  getWalletPositions,
  getWalletProfile,
  getWalletTrades,
} from "@/lib/api";
import { Card } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { AccountInfoCard } from "@/components/wallet/AccountInfoCard";
import { PositionsList } from "@/components/wallet/PositionsList";
import { TradeHistory } from "@/components/wallet/TradeHistory";

export const dynamic = "force-dynamic";

export default async function WalletProfilePage({
  params,
  searchParams,
}: {
  params: { id: string };
  searchParams: { trades_page?: string };
}) {
  const locale = await getLocaleFromCookies();
  const tradesPage = Math.max(1, parseInt(searchParams.trades_page || "1", 10) || 1);

  const [profile, positions, trades] = await Promise.all([
    getWalletProfile(params.id),
    getWalletPositions(params.id).catch(() => []),
    getWalletTrades(params.id, new URLSearchParams({ page: String(tradesPage), page_size: "20" })).catch(() => ({
      items: [],
      pagination: { page: tradesPage, page_size: 20, total: 0 },
    })),
  ]);

  return (
    <section className="space-y-8 animate-fade-in">
      {/* ── 1. Account Info Card ── */}
      <AccountInfoCard profile={profile} decisionCard={null} locale={locale} />

      {/* ── 2. Key Metrics ── */}
      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <StatCell
            label={t(locale, "profile.realizedPnl")}
            value={`${profile.layer1_facts.realized_pnl >= 0 ? "+" : ""}${profile.layer1_facts.realized_pnl.toFixed(2)}`}
            color={profile.layer1_facts.realized_pnl >= 0 ? "positive" : "negative"}
          />
        </Card>
        <Card>
          <StatCell label={t(locale, "profile.totalTrades")} value={String(profile.layer1_facts.total_trades)} />
        </Card>
        <Card>
          <StatCell label={t(locale, "profile.feesPaid")} value={profile.layer1_facts.fees_paid.toFixed(2)} />
        </Card>
        <Card>
          <StatCell label={t(locale, "profile.volume30d")} value={profile.layer1_facts.volume_30d.toFixed(2)} />
        </Card>
      </div>

      {/* ── 3. Positions ── */}
      <PositionsList positions={positions} locale={locale} />

      {/* ── 4. Trade History ── */}
      <TradeHistory
        items={trades.items}
        locale={locale}
        page={trades.pagination.page}
        pageSize={trades.pagination.page_size}
        total={trades.pagination.total}
        walletId={params.id}
      />

      {/* ── 5. Layer 3 Info-Edge ── */}
      <Card>
        <h3 className="mb-3 text-subheadline font-semibold text-label-secondary">{t(locale, "profile.layer3")}</h3>
        <div className="grid gap-4 sm:grid-cols-4">
          <StatCell label={t(locale, "profile.label")} value={profile.layer3_info_edge.label} />
          <StatCell label={t(locale, "profile.meanDt")} value={`${profile.layer3_info_edge.mean_delta_minutes.toFixed(2)} min`} />
          <StatCell label={t(locale, "profile.samples")} value={String(profile.layer3_info_edge.samples)} />
          <StatCell label={t(locale, "profile.pvalue")} value={profile.layer3_info_edge.p_value.toFixed(4)} />
        </div>
      </Card>

      {/* ── 6. Disclosures ── */}
      <div className="rounded-2xl bg-tint-orange/[0.05] p-5">
        <div className="space-y-1.5">
          {profile.meta.disclosures.map((d) => (
            <p key={d} className="text-caption-1 leading-relaxed text-tint-orange">{d}</p>
          ))}
        </div>
      </div>
    </section>
  );
}
