import { AlertTriangle, Shield } from "lucide-react";
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

  const riskColor = card.risk_level === "low" ? "text-tint-green" : card.risk_level === "high" ? "text-tint-red" : "text-tint-orange";
  const momentumColor = card.momentum === "heating" ? "text-tint-green" : card.momentum === "cooling" ? "text-tint-red" : "text-label-secondary";

  return (
    <article className="rounded-2xl border border-separator/50 bg-surface-secondary p-6 shadow-elevation-2 transition-all duration-300 ease-apple hover:shadow-[var(--card-hover-shadow)]">
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-tint-blue/10">
            <Shield className="h-4 w-4 text-tint-blue" />
          </div>
          <h2 className="text-title-3 text-label-primary">{text.title}</h2>
        </div>
        <TierBadge tier={card.pool_tier} locale={locale} />
      </div>
      <div className="grid gap-3 sm:grid-cols-2">
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium tracking-wide uppercase text-label-tertiary">{text.suitable}</span>
          <p className="mt-0.5 text-subheadline font-semibold text-label-primary">{card.suitable_for}</p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium tracking-wide uppercase text-label-tertiary">{text.risk}</span>
          <p className={`mt-0.5 text-subheadline font-semibold ${riskColor}`}>{card.risk_level}</p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium tracking-wide uppercase text-label-tertiary">{text.position}</span>
          <p className="mt-0.5 text-subheadline font-semibold text-label-primary">{card.suggested_position}</p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium tracking-wide uppercase text-label-tertiary">{text.momentum}</span>
          <p className={`mt-0.5 text-subheadline font-semibold ${momentumColor}`}>{card.momentum}</p>
        </div>
      </div>
      <div className="mt-4 rounded-xl bg-tint-blue/[0.06] p-4">
        <p className="text-callout leading-relaxed text-label-secondary">
          {locale === "zh" ? card.recommendation_zh : card.recommendation}
        </p>
      </div>
      <p className="mt-3 flex items-start gap-2 text-caption-1 text-tint-orange">
        <AlertTriangle className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        <span className="leading-relaxed">{locale === "zh" ? card.disclaimer_zh : card.disclaimer}</span>
      </p>
    </article>
  );
}
