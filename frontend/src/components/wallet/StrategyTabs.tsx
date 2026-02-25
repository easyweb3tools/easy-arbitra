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
    <div className="flex gap-2 overflow-x-auto pb-1 -mx-1 px-1">
      {tabs.map((tab) => {
        const active = (current || "") === tab.key;
        return (
          <Link
            key={tab.key || "all"}
            href={makeHref(tab.key)}
            className={[
              "inline-flex h-8 items-center rounded-full px-4 text-caption-1 font-semibold whitespace-nowrap",
              "transition-all duration-200 ease-apple",
              active
                ? "bg-tint-blue text-white shadow-[0_1px_4px_rgba(0,122,255,0.25)]"
                : "bg-surface-tertiary/70 text-label-tertiary hover:bg-surface-tertiary hover:text-label-secondary"
            ].join(" ")}
          >
            {tab.label}
          </Link>
        );
      })}
    </div>
  );
}
