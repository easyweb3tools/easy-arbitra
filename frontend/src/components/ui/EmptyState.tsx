import type { ReactNode } from "react";
import { Star, Brain, Search, CheckCircle, Clock, List } from "lucide-react";

type Preset =
  | "watchlist-empty"
  | "no-ai-report"
  | "no-results"
  | "all-clear"
  | "no-activity"
  | "feed-empty";

const presetConfig: Record<
  Preset,
  { icon: ReactNode; titleEn: string; titleZh: string; descEn: string; descZh: string }
> = {
  "watchlist-empty": {
    icon: <Star className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "Start Following Wallets",
    titleZh: "开始关注钱包",
    descEn: "Follow wallets to receive alerts on important changes.",
    descZh: "关注后你会收到这些钱包的重要变动通知。",
  },
  "no-ai-report": {
    icon: <Brain className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "AI Analysis Pending",
    titleZh: "AI 分析进行中",
    descEn: "The system will automatically generate an analysis report.",
    descZh: "系统将自动为该钱包生成分析报告。",
  },
  "no-results": {
    icon: <Search className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "No Results Found",
    titleZh: "没有找到匹配结果",
    descEn: "Try adjusting your filters or search terms.",
    descZh: "尝试调整筛选条件或搜索关键词。",
  },
  "all-clear": {
    icon: <CheckCircle className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "All Clear",
    titleZh: "一切正常",
    descEn: "No anomalies detected at this time.",
    descZh: "目前没有需要关注的异常信号。",
  },
  "no-activity": {
    icon: <Clock className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "No Recent Activity",
    titleZh: "暂无新动态",
    descEn: "Updates will appear here when wallets have new activity.",
    descZh: "当钱包有新交易或分析更新时会显示在这里。",
  },
  "feed-empty": {
    icon: <List className="h-12 w-12" strokeWidth={1.2} />,
    titleEn: "Feed Is Empty",
    titleZh: "事件流为空",
    descEn: "Events from followed wallets will appear here.",
    descZh: "你关注的钱包尚未产生新事件。",
  },
};

export function EmptyState({
  preset,
  locale = "en",
  action,
}: {
  preset: Preset;
  locale?: "en" | "zh";
  action?: ReactNode;
}) {
  const cfg = presetConfig[preset];
  return (
    <div className="flex flex-col items-center justify-center px-6 py-16 text-center animate-fade-in">
      <div className="mb-5 rounded-2xl bg-surface-tertiary/60 p-4 text-label-quaternary">{cfg.icon}</div>
      <h3 className="text-title-3 text-label-primary">
        {locale === "zh" ? cfg.titleZh : cfg.titleEn}
      </h3>
      <p className="mt-1.5 max-w-[280px] text-subheadline leading-relaxed text-label-tertiary">
        {locale === "zh" ? cfg.descZh : cfg.descEn}
      </p>
      {action && <div className="mt-6">{action}</div>}
    </div>
  );
}
