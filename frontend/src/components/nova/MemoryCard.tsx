"use client";

import { LearningRecord } from "@/lib/types";
import { CheckCircle, XCircle, TrendingUp, TrendingDown } from "lucide-react";

interface MemoryCardProps {
  record: LearningRecord;
  locale: string;
}

export function MemoryCard({ record, locale }: MemoryCardProps) {
  const lesson = locale === "zh" ? record.lesson_learned_zh : record.lesson_learned;
  const date = new Date(record.validation_date).toLocaleDateString(
    locale === "zh" ? "zh-CN" : "en-US",
    { year: "numeric", month: "short", day: "numeric" }
  );

  return (
    <div
      className={`rounded-lg p-4 border ${
        record.is_success
          ? "bg-green-50 dark:bg-green-950/20 border-green-200 dark:border-green-800"
          : "bg-red-50 dark:bg-red-950/20 border-red-200 dark:border-red-800"
      }`}
    >
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0 mt-1">
          {record.is_success ? (
            <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
          ) : (
            <XCircle className="w-5 h-5 text-red-600 dark:text-red-400" />
          )}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-sm font-semibold text-gray-900 dark:text-white">
              {date}
            </span>
            <span
              className={`text-xs px-2 py-0.5 rounded ${
                record.is_success
                  ? "bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300"
                  : "bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300"
              }`}
            >
              {record.is_success
                ? locale === "zh"
                  ? "成功"
                  : "Success"
                : locale === "zh"
                ? "失败"
                : "Failed"}
            </span>
          </div>

          <div className="mb-2">
            <div className="text-xs text-gray-600 dark:text-gray-400 mb-1">
              {locale === "zh" ? "推荐钱包" : "Recommended Wallet"}
            </div>
            <div className="text-sm font-mono text-gray-900 dark:text-white">
              {record.wallet_address}
            </div>
          </div>

          {record.follow_pnl !== undefined && (
            <div className="flex items-center gap-2 mb-2">
              {record.follow_pnl >= 0 ? (
                <TrendingUp className="w-4 h-4 text-green-600 dark:text-green-400" />
              ) : (
                <TrendingDown className="w-4 h-4 text-red-600 dark:text-red-400" />
              )}
              <span
                className={`text-sm font-semibold ${
                  record.follow_pnl >= 0
                    ? "text-green-600 dark:text-green-400"
                    : "text-red-600 dark:text-red-400"
                }`}
              >
                {record.follow_pnl >= 0 ? "+" : ""}
                {record.follow_pnl.toFixed(2)}
              </span>
              <span className="text-xs text-gray-600 dark:text-gray-400">
                {locale === "zh" ? "次日跟单收益" : "Next-day Follow PnL"}
              </span>
            </div>
          )}

          <div className="bg-white dark:bg-gray-800 rounded p-3">
            <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 mb-1">
              💡 {locale === "zh" ? "Nova 学到了" : "Nova Learned"}
            </div>
            <p className="text-sm text-gray-700 dark:text-gray-300">{lesson}</p>
          </div>

          {record.strategy_adjustment && (
            <div className="mt-2 text-xs text-purple-600 dark:text-purple-400">
              🔧 {locale === "zh" ? "策略已调整" : "Strategy Adjusted"}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
