"use client";

import { ThinkingRound } from "@/lib/types";
import { CheckCircle, Search, Zap, AlertCircle } from "lucide-react";
import { useState } from "react";

interface ThinkingTimelineProps {
  rounds: ThinkingRound[];
  locale: string;
}

export function ThinkingTimeline({ rounds, locale }: ThinkingTimelineProps) {
  const [expandedRounds, setExpandedRounds] = useState<Set<number>>(new Set());

  const toggleRound = (round: number) => {
    const newExpanded = new Set(expandedRounds);
    if (newExpanded.has(round)) {
      newExpanded.delete(round);
    } else {
      newExpanded.add(round);
    }
    setExpandedRounds(newExpanded);
  };

  const getRoundIcon = (round: ThinkingRound) => {
    if (round.session.phase === "final") {
      return <CheckCircle className="w-5 h-5 text-green-500" />;
    }
    if (round.is_breakthrough) {
      return <Zap className="w-5 h-5 text-yellow-500" />;
    }
    if (round.is_hesitation) {
      return <AlertCircle className="w-5 h-5 text-orange-500" />;
    }
    return <Search className="w-5 h-5 text-blue-500" />;
  };

  const getRoundLabel = (round: ThinkingRound) => {
    if (round.session.phase === "final") return "Final Decision";
    if (round.is_breakthrough) return "Breakthrough";
    if (round.is_hesitation) return "Analyzing";
    return "Evaluating";
  };

  return (
    <div className="space-y-4">
      {rounds.map((round, idx) => {
        const isExpanded = expandedRounds.has(round.session.round);
        const summary = locale === "zh" ? round.session.nl_summary_zh : round.session.nl_summary;

        return (
          <div
            key={round.session.id}
            className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
          >
            <button
              onClick={() => toggleRound(round.session.round)}
              className="w-full p-4 flex items-start gap-4 hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors text-left"
            >
              <div className="flex-shrink-0 mt-1">{getRoundIcon(round)}</div>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-sm font-semibold text-gray-900 dark:text-white">
                    {new Date(round.session.created_at).toLocaleTimeString()} — Round {round.session.round}
                  </span>
                  <span className="text-xs px-2 py-0.5 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded">
                    {getRoundLabel(round)}
                  </span>
                  {round.session.confidence_score && (
                    <span className="text-xs text-gray-600 dark:text-gray-400">
                      Confidence: {Math.round(round.session.confidence_score)}%
                    </span>
                  )}
                  {round.confidence_change && (
                    <span
                      className={`text-xs ${
                        round.confidence_change > 0
                          ? "text-green-600 dark:text-green-400"
                          : "text-red-600 dark:text-red-400"
                      }`}
                    >
                      {round.confidence_change > 0 ? "+" : ""}
                      {round.confidence_change.toFixed(1)}%
                    </span>
                  )}
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
                  {summary || "Nova is analyzing candidates..."}
                </p>
              </div>

              <div className="flex-shrink-0">
                <svg
                  className={`w-5 h-5 text-gray-400 transition-transform ${
                    isExpanded ? "rotate-180" : ""
                  }`}
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              </div>
            </button>

            {isExpanded && (
              <div className="px-4 pb-4 border-t border-gray-200 dark:border-gray-700 pt-4">
                <div className="prose dark:prose-invert max-w-none">
                  <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
                    {summary || "No detailed analysis available."}
                  </p>
                </div>

                {round.session.focus_metrics && (
                  <div className="mt-4">
                    <h4 className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-2">
                      Focus Metrics:
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {JSON.parse(JSON.stringify(round.session.focus_metrics) || "[]").map(
                        (metric: string, i: number) => (
                          <span
                            key={i}
                            className="text-xs px-2 py-1 bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300 rounded"
                          >
                            {metric}
                          </span>
                        )
                      )}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
