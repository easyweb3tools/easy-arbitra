"use client";

import { DecisionExplanation } from "@/lib/types";
import { CheckCircle, XCircle } from "lucide-react";

interface DecisionExplainerProps {
  explanation: DecisionExplanation;
  locale: string;
}

export function DecisionExplainer({ explanation, locale }: DecisionExplainerProps) {
  const reasons = locale === "zh" ? explanation.key_reasons_zh : explanation.key_reasons;

  return (
    <div className="space-y-6">
      {/* Weight Distribution */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          {locale === "zh" ? "决策权重分布" : "Decision Weight Distribution"}
        </h3>
        <div className="space-y-3">
          {Object.entries(explanation.weight_distribution).map(([key, value]) => (
            <div key={key}>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-700 dark:text-gray-300 capitalize">
                  {key.replace(/_/g, " ")}
                </span>
                <span className="font-semibold text-gray-900 dark:text-white">
                  {value}%
                </span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-gradient-to-r from-purple-500 to-blue-500 h-2 rounded-full"
                  style={{ width: `${value}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Metric Comparison */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          {locale === "zh" ? "关键指标对比" : "Key Metrics Comparison"}
        </h3>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="text-left py-2 px-3 text-gray-600 dark:text-gray-400 font-medium">
                  {locale === "zh" ? "指标" : "Metric"}
                </th>
                <th className="text-right py-2 px-3 text-gray-600 dark:text-gray-400 font-medium">
                  {locale === "zh" ? "推荐钱包" : "Pick Value"}
                </th>
                <th className="text-right py-2 px-3 text-gray-600 dark:text-gray-400 font-medium">
                  {locale === "zh" ? "平均水平" : "Average"}
                </th>
                <th className="text-right py-2 px-3 text-gray-600 dark:text-gray-400 font-medium">
                  {locale === "zh" ? "Nova 标准" : "Standard"}
                </th>
                <th className="text-center py-2 px-3 text-gray-600 dark:text-gray-400 font-medium">
                  {locale === "zh" ? "达标" : "Pass"}
                </th>
              </tr>
            </thead>
            <tbody>
              {explanation.metric_comparison.map((metric, idx) => (
                <tr
                  key={idx}
                  className="border-b border-gray-100 dark:border-gray-800 last:border-0"
                >
                  <td className="py-3 px-3 text-gray-900 dark:text-white font-medium">
                    {locale === "zh" ? metric.metric_zh : metric.metric}
                  </td>
                  <td className="py-3 px-3 text-right text-gray-900 dark:text-white font-semibold">
                    {metric.pick_value}
                    {metric.metric.includes("Rate") ? "%" : ""}
                  </td>
                  <td className="py-3 px-3 text-right text-gray-600 dark:text-gray-400">
                    {metric.average_value}
                    {metric.metric.includes("Rate") ? "%" : ""}
                  </td>
                  <td className="py-3 px-3 text-right text-gray-600 dark:text-gray-400">
                    {metric.nova_standard}
                  </td>
                  <td className="py-3 px-3 text-center">
                    {metric.passed ? (
                      <CheckCircle className="w-5 h-5 text-green-500 inline" />
                    ) : (
                      <XCircle className="w-5 h-5 text-red-500 inline" />
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Key Reasons */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          {locale === "zh" ? "关键理由" : "Key Reasons"}
        </h3>
        <ul className="space-y-2">
          {reasons.map((reason, idx) => (
            <li
              key={idx}
              className="flex items-start gap-2 text-gray-700 dark:text-gray-300"
            >
              <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" />
              <span>{reason}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
