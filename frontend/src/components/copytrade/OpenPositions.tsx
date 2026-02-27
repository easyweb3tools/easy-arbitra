"use client";

import { useState } from "react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { closeCopyTradePosition } from "@/lib/api";
import type { CopyTradeDecision } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

function OpenPositionRow({
  position,
  locale,
  onClosed,
}: {
  position: CopyTradeDecision;
  locale: Locale;
  onClosed?: () => void;
}) {
  const [loading, setLoading] = useState(false);

  async function handleClose() {
    setLoading(true);
    try {
      await closeCopyTradePosition(position.id);
      onClosed?.();
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex items-center justify-between gap-3 px-4 py-3 border-b border-separator/60 last:border-b-0">
      <div className="flex-1 min-w-0">
        <p className="text-subheadline font-medium text-label-primary truncate">
          {position.market_title}
        </p>
        <div className="flex items-center gap-2 mt-0.5 text-caption-1 text-label-tertiary">
          <span>{position.outcome} {position.action}</span>
          <span>${position.size_usdc.toFixed(0)} @ {position.price.toFixed(3)}</span>
          <span>{new Date(position.created_at).toLocaleDateString()}</span>
        </div>
      </div>
      <Button variant="gray" size="mini" loading={loading} onClick={handleClose}>
        {loading ? t(locale, "copyTrade.closing") : t(locale, "copyTrade.close")}
      </Button>
    </div>
  );
}

export function OpenPositions({
  positions,
  locale,
  onRefresh,
}: {
  positions: CopyTradeDecision[];
  locale: Locale;
  onRefresh?: () => void;
}) {
  if (positions.length === 0) {
    return (
      <Card variant="flat">
        <p className="py-6 text-center text-subheadline text-label-tertiary">
          {t(locale, "copyTrade.noPositions")}
        </p>
      </Card>
    );
  }

  return (
    <Card padding={false}>
      {positions.map((pos) => (
        <OpenPositionRow
          key={pos.id}
          position={pos}
          locale={locale}
          onClosed={onRefresh}
        />
      ))}
    </Card>
  );
}
