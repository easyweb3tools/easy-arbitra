import Link from "next/link";

const tabs = [
  { key: "", label: "All" },
  { key: "market_maker", label: "Maker" },
  { key: "event_trader", label: "Event" },
  { key: "quant", label: "Quant" },
  { key: "arbitrage", label: "Arb" },
  { key: "lucky", label: "Lucky" },
];

export function StrategyTabs({
  current,
  makeHref,
}: {
  current?: string;
  makeHref: (strategyType: string) => string;
}) {
  return (
    <div className="flex gap-2 overflow-x-auto pb-1">
      {tabs.map((tab) => {
        const active = (current || "") === tab.key;
        return (
          <Link
            key={tab.key || "all"}
            href={makeHref(tab.key)}
            className={[
              "inline-flex h-8 items-center rounded-full px-3 text-caption-1 font-medium whitespace-nowrap",
              active
                ? "bg-tint-blue/15 text-tint-blue"
                : "bg-surface-tertiary text-label-tertiary hover:text-label-secondary"
            ].join(" ")}
          >
            {tab.label}
          </Link>
        );
      })}
    </div>
  );
}
