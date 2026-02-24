import { AlertTriangle } from "lucide-react";
import type { WalletDecisionCard } from "@/lib/types";
import { TierBadge } from "@/components/ui/Badge";

export function DecisionCard({
  card,
  locale,
}: {
  card: WalletDecisionCard;
  locale: "en" | "zh";
}) {
  const text = locale === "zh" ? {
    title: "跟单决策卡",
    suitable: "适合人群",
    risk: "风险等级",
    position: "建议仓位",
    momentum: "近7天状态",
  } : {
    title: "Follow Decision",
    suitable: "Suitable for",
    risk: "Risk",
    position: "Position",
    momentum: "7D status",
  };

  return (
    <article className="rounded-lg border border-separator bg-surface-secondary p-5 shadow-elevation-2">
      <div className="mb-3 flex items-center justify-between">
        <h2 className="text-title-3 text-label-primary">{text.title}</h2>
        <TierBadge tier={card.pool_tier} locale={locale} />
      </div>
      <div className="grid gap-2 sm:grid-cols-2">
        <p className="text-footnote text-label-tertiary">{text.suitable}: <span className="font-medium text-label-secondary">{card.suitable_for}</span></p>
        <p className="text-footnote text-label-tertiary">{text.risk}: <span className="font-medium text-label-secondary">{card.risk_level}</span></p>
        <p className="text-footnote text-label-tertiary">{text.position}: <span className="font-medium text-label-secondary">{card.suggested_position}</span></p>
        <p className="text-footnote text-label-tertiary">{text.momentum}: <span className="font-medium text-label-secondary">{card.momentum}</span></p>
      </div>
      <p className="mt-3 rounded-md bg-surface-tertiary p-3 text-callout text-label-secondary">
        {locale === "zh" ? card.recommendation_zh : card.recommendation}
      </p>
      <p className="mt-3 flex items-start gap-1.5 text-caption-1 text-tint-orange">
        <AlertTriangle className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        {locale === "zh" ? card.disclaimer_zh : card.disclaimer}
      </p>
    </article>
  );
}
