"use client";

import { Badge } from "@/components/ui/badge";
import type { DecisionStep } from "@/lib/types";

interface DecisionLogProps {
  steps: DecisionStep[];
}

const toolColors: Record<string, string> = {
  resolve_wallet_target: "bg-blue-500/20 text-blue-300 border-blue-500/30",
  fetch_sports_trades: "bg-green-500/20 text-green-300 border-green-500/30",
  calculate_style_metrics: "bg-amber-500/20 text-amber-300 border-amber-500/30",
  build_report_payload: "bg-purple-500/20 text-purple-300 border-purple-500/30",
};

const toolLabels: Record<string, string> = {
  resolve_wallet_target: "Resolve Wallet",
  fetch_sports_trades: "Fetch Trades",
  calculate_style_metrics: "Calculate Metrics",
  build_report_payload: "Build Report",
};

export function DecisionLog({ steps }: DecisionLogProps) {
  return (
    <div className="space-y-0">
      {steps.map((step, idx) => (
        <div key={step.step} className="flex gap-4">
          {/* Timeline */}
          <div className="flex flex-col items-center">
            <div className="w-8 h-8 rounded-full bg-white/10 border border-white/20 flex items-center justify-center text-sm font-bold text-white">
              {step.step}
            </div>
            {idx < steps.length - 1 && (
              <div className="w-px h-full min-h-[40px] bg-white/10" />
            )}
          </div>

          {/* Content */}
          <div className="pb-6 flex-1">
            <div className="flex items-center gap-2 mb-1">
              <Badge
                variant="outline"
                className={toolColors[step.tool] || "bg-gray-500/20 text-gray-300"}
              >
                {toolLabels[step.tool] || step.tool}
              </Badge>
              <span className="text-xs text-white/40">
                {new Date(step.timestamp).toLocaleTimeString()}
              </span>
            </div>
            <p className="text-sm text-white/70">{step.reasoning}</p>
            {step.result_summary && (
              <p className="text-xs text-white/50 mt-1">
                {step.result_summary}
              </p>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
