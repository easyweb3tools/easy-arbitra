"use client";

import { NovaStatus } from "@/lib/types";
import { Brain, Clock, Target, TrendingUp } from "lucide-react";
import { useEffect, useState } from "react";

interface BrainDashboardProps {
  status: NovaStatus;
}

export function BrainDashboard({ status: initialStatus }: BrainDashboardProps) {
  const [status, setStatus] = useState(initialStatus);

  useEffect(() => {
    if (!status.is_active) return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch("/api/v1/nova/status");
        const data = await res.json();
        setStatus(data.data);
      } catch (err) {
        console.error("Failed to fetch Nova status:", err);
      }
    }, 30000); // Poll every 30 seconds

    return () => clearInterval(interval);
  }, [status.is_active]);

  const progress = (status.current_round / status.total_rounds) * 100;

  return (
    <div className="bg-gradient-to-br from-purple-50 to-blue-50 dark:from-purple-950/20 dark:to-blue-950/20 rounded-lg p-6 border border-purple-200 dark:border-purple-800">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
          <Brain className="w-6 h-6 text-purple-600 dark:text-purple-400" />
        </div>
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            Nova AI Brain
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            {status.is_active ? "Thinking..." : "Idle"}
          </p>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="mb-6">
        <div className="flex justify-between text-sm mb-2">
          <span className="text-gray-600 dark:text-gray-400">
            Round {status.current_round} / {status.total_rounds}
          </span>
          <span className="text-gray-600 dark:text-gray-400">
            {Math.round(progress)}%
          </span>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div
            className="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full transition-all duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-1">
            <Target className="w-4 h-4 text-blue-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              Candidates
            </span>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {status.candidate_count}
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp className="w-4 h-4 text-green-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              Confidence
            </span>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {status.confidence_score ? `${Math.round(status.confidence_score)}%` : "—"}
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-1">
            <Clock className="w-4 h-4 text-orange-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              Phase
            </span>
          </div>
          <p className="text-sm font-semibold text-gray-900 dark:text-white capitalize">
            {status.phase}
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-1">
            <Brain className="w-4 h-4 text-purple-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              Focus
            </span>
          </div>
          <p className="text-xs font-medium text-gray-900 dark:text-white">
            {status.focus_metrics && status.focus_metrics.length > 0
              ? status.focus_metrics.slice(0, 2).join(", ")
              : "Analyzing..."}
          </p>
        </div>
      </div>

      {/* Next Round Timer */}
      {status.is_active && status.next_round_at && (
        <div className="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
          Next analysis round at{" "}
          {new Date(status.next_round_at).toLocaleTimeString()}
        </div>
      )}
    </div>
  );
}
