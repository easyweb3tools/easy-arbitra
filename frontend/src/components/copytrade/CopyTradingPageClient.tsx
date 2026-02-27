"use client";

import { useEffect, useState } from "react";
import { getCopyTradeDashboard, getCopyTradePositions } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { SkeletonRow } from "@/components/ui/Skeleton";
import { DashboardStats } from "@/components/copytrade/DashboardStats";
import { ConfigList } from "@/components/copytrade/ConfigList";
import { DecisionFeed } from "@/components/copytrade/DecisionFeed";
import { OpenPositions } from "@/components/copytrade/OpenPositions";
import type { CopyTradeDashboard, CopyTradeDecision } from "@/lib/types";
import type { Locale } from "@/lib/i18n";

type Labels = {
  title: string;
  activeConfigs: string;
  positions: string;
  recentDecisions: string;
  noConfigs: string;
};

export function CopyTradingPageClient({
  locale,
  labels,
}: {
  locale: Locale;
  labels: Labels;
}) {
  const [loading, setLoading] = useState(true);
  const [dashboard, setDashboard] = useState<CopyTradeDashboard | null>(null);
  const [positions, setPositions] = useState<CopyTradeDecision[]>([]);

  function fetchData() {
    let cancelled = false;
    setLoading(true);
    Promise.all([
      getCopyTradeDashboard(),
      getCopyTradePositions(),
    ])
      .then(([d, p]) => {
        if (!cancelled) {
          setDashboard(d);
          setPositions(p);
        }
      })
      .catch(() => {
        // will show empty state
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }

  useEffect(fetchData, []);

  if (loading) {
    return (
      <section className="space-y-6">
        <h1 className="text-large-title font-bold text-label-primary">
          {labels.title}
        </h1>
        <Card padding={false}>
          <SkeletonRow />
          <SkeletonRow />
          <SkeletonRow />
        </Card>
      </section>
    );
  }

  return (
    <section className="space-y-8 animate-fade-in">
      <h1 className="text-large-title font-bold text-label-primary">
        {labels.title}
      </h1>

      {dashboard ? (
        <>
          <DashboardStats data={dashboard} locale={locale} />

          <div>
            <SectionHeader title={labels.activeConfigs} />
            <ConfigList configs={dashboard.configs} locale={locale} />
          </div>

          <div>
            <SectionHeader title={labels.positions} />
            <OpenPositions
              positions={positions}
              locale={locale}
              onRefresh={fetchData}
            />
          </div>

          <div>
            <SectionHeader title={labels.recentDecisions} />
            <DecisionFeed decisions={dashboard.recent_decisions} locale={locale} />
          </div>
        </>
      ) : (
        <Card>
          <p className="py-12 text-center text-subheadline text-label-tertiary">
            {labels.noConfigs}
          </p>
        </Card>
      )}
    </section>
  );
}
