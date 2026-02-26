import { Card } from "@/components/ui/Card";
import { CategoryTag } from "@/components/ui/Badge";
import type { CopyTradeDecision } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

function fmt$(v: number) {
  return v >= 0 ? `+$${v.toFixed(2)}` : `-$${Math.abs(v).toFixed(2)}`;
}

function fmtPct(v: number) {
  return `${(v * 100).toFixed(0)}%`;
}

export function DecisionCard({
  decision,
  locale,
}: {
  decision: CopyTradeDecision;
  locale: Locale;
}) {
  const isCopy = decision.decision === "copy";
  const reasoning = locale === "zh" ? decision.reasoning : decision.reasoning_en;

  return (
    <Card variant="flat" className="!p-4">
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-2">
          <CategoryTag color={isCopy ? "green" : "gray"}>
            {isCopy ? t(locale, "copyTrade.copy") : t(locale, "copyTrade.skip")}
          </CategoryTag>
          <span className="text-subheadline font-medium text-label-primary line-clamp-1">
            {decision.market_title}
          </span>
        </div>
        <span className="text-caption-1 text-label-quaternary whitespace-nowrap">
          {new Date(decision.created_at).toLocaleDateString()}
        </span>
      </div>

      {isCopy && (
        <div className="mt-2 flex items-center gap-3 text-footnote text-label-secondary">
          <span>${decision.size_usdc.toFixed(0)} @ {decision.price.toFixed(3)}</span>
          <span>{decision.outcome} {decision.action}</span>
          <span>{t(locale, "copyTrade.confidence")} {fmtPct(decision.confidence)}</span>
          {decision.realized_pnl != null && (
            <span className={decision.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}>
              {fmt$(decision.realized_pnl)}
            </span>
          )}
        </div>
      )}

      <p className="mt-2 text-footnote text-label-tertiary leading-relaxed line-clamp-2">
        {reasoning}
      </p>

      {decision.risk_notes && decision.risk_notes.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1.5">
          {decision.risk_notes.map((note, i) => (
            <span key={i} className="inline-flex h-5 items-center rounded px-1.5 text-caption-2 bg-tint-orange/10 text-tint-orange">
              {note}
            </span>
          ))}
        </div>
      )}
    </Card>
  );
}

export function DecisionFeed({
  decisions,
  locale,
}: {
  decisions: CopyTradeDecision[];
  locale: Locale;
}) {
  if (decisions.length === 0) {
    return (
      <p className="py-8 text-center text-subheadline text-label-tertiary">
        {t(locale, "copyTrade.noDecisions")}
      </p>
    );
  }

  return (
    <div className="space-y-3">
      {decisions.map((d) => (
        <DecisionCard key={d.id} decision={d} locale={locale} />
      ))}
    </div>
  );
}
