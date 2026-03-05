"use client";

import { CandidateScore } from "@/lib/types";
import { Star, TrendingUp, Activity, Shield } from "lucide-react";
import Link from "next/link";

interface CandidateScoreCardProps {
  candidate: CandidateScore;
  locale: string;
}

export function CandidateScoreCard({ candidate, locale }: CandidateScoreCardProps) {
  const comment = locale === "zh" ? candidate.nova_comment_zh : candidate.nova_comment;

  const getStars = (value: number) => {
    const stars = Math.round((value / 100) * 5);
    return Array.from({ length: 5 }, (_, i) => (
      <Star
        key={i}
        className={`w-4 h-4 ${
          i < stars
            ? "fill-yellow-400 text-yellow-400"
            : "text-gray-300 dark:text-gray-600"
        }`}
      />
    ));
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 hover:border-purple-300 dark:hover:border-purple-700 transition-colors">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900 rounded-full flex items-center justify-center">
            <span className="text-sm font-bold text-purple-600 dark:text-purple-400">
              #{candidate.rank}
            </span>
          </div>
          <div>
            <Link
              href={`/wallets/${candidate.wallet_id}`}
              className="text-sm font-mono font-semibold text-blue-600 dark:text-blue-400 hover:underline"
            >
              {candidate.address
                ? `${candidate.address.slice(0, 6)}...${candidate.address.slice(-4)}`
                : `#${candidate.wallet_id}`}
            </Link>
            {candidate.pseudonym && (
              <div className="text-xs text-gray-600 dark:text-gray-400">
                {candidate.pseudonym}
              </div>
            )}
          </div>
        </div>
        <div className="text-right">
          <div className="text-xs text-gray-600 dark:text-gray-400">
            {locale === "zh" ? "Nova 评分" : "Nova Score"}
          </div>
          <div className="text-2xl font-bold text-purple-600 dark:text-purple-400">
            {candidate.nova_score}
          </div>
        </div>
      </div>

      <div className="space-y-2 mb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <TrendingUp className="w-4 h-4 text-green-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "胜率" : "Win Rate"}
            </span>
          </div>
          <div className="flex items-center gap-1">
            {getStars(candidate.win_rate)}
            <span className="text-xs text-gray-600 dark:text-gray-400 ml-1">
              {candidate.win_rate.toFixed(0)}%
            </span>
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Shield className="w-4 h-4 text-blue-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "稳定性" : "Stability"}
            </span>
          </div>
          <div className="flex items-center gap-1">
            {getStars(candidate.stability)}
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Activity className="w-4 h-4 text-orange-500" />
            <span className="text-xs text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "活跃度" : "Activity"}
            </span>
          </div>
          <div className="flex items-center gap-1">
            {getStars(candidate.activity)}
          </div>
        </div>
      </div>

      {comment && (
        <div className="bg-purple-50 dark:bg-purple-950/30 rounded p-3">
          <div className="text-xs font-semibold text-purple-600 dark:text-purple-400 mb-1">
            💬 {locale === "zh" ? "Nova 点评" : "Nova's Comment"}
          </div>
          <p className="text-xs text-gray-700 dark:text-gray-300">{comment}</p>
        </div>
      )}
    </div>
  );
}
