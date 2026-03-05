import { getNovaMemory } from "@/lib/api";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { MemoryCard } from "@/components/nova/MemoryCard";
import { Brain, TrendingUp, Target, Lightbulb } from "lucide-react";

export const metadata = {
  title: "Nova Memory | Easy Arbitra",
  description: "Nova's learning history and strategy evolution",
};

export const dynamic = "force-dynamic";

export default async function NovaMemoryPage() {
  const locale = await getLocaleFromCookies();

  let memory;
  try {
    memory = await getNovaMemory(30);
  } catch (err) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <Brain className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            {locale === "zh" ? "记忆库暂不可用" : "Memory Unavailable"}
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            {locale === "zh"
              ? "无法获取 Nova 的学习记录，请稍后再试。"
              : "Unable to fetch Nova's learning records. Please try again later."}
          </p>
        </div>
      </div>
    );
  }

  const { summary, history } = memory;
  const lessons = locale === "zh" ? summary.recent_lessons_zh : summary.recent_lessons;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          {locale === "zh" ? "Nova 记忆库" : "Nova's Memory"}
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          {locale === "zh"
            ? "Nova 从历史验证中学习和改进"
            : "How Nova learns and improves from validation results"}
        </p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <Target className="w-5 h-5 text-blue-500" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "总验证次数" : "Total Validations"}
            </span>
          </div>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">
            {summary.total_validations}
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-5 h-5 text-green-500" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "成功率" : "Success Rate"}
            </span>
          </div>
          <p className="text-3xl font-bold text-green-600 dark:text-green-400">
            {summary.success_rate.toFixed(1)}%
          </p>
          <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
            {summary.success_count} / {summary.total_validations}
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-5 h-5 text-purple-500" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "本周成功率" : "Weekly Success"}
            </span>
          </div>
          <p className="text-3xl font-bold text-purple-600 dark:text-purple-400">
            {summary.weekly_success_rate.toFixed(1)}%
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <Brain className="w-5 h-5 text-orange-500" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {locale === "zh" ? "策略调整" : "Strategy Changes"}
            </span>
          </div>
          <p className="text-3xl font-bold text-gray-900 dark:text-white">
            {summary.strategy_evolution.length}
          </p>
        </div>
      </div>

      {/* Recent Lessons */}
      <div className="bg-gradient-to-br from-purple-50 to-blue-50 dark:from-purple-950/20 dark:to-blue-950/20 rounded-lg p-6 border border-purple-200 dark:border-purple-800 mb-8">
        <div className="flex items-center gap-2 mb-4">
          <Lightbulb className="w-5 h-5 text-yellow-500" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            {locale === "zh" ? "最近学到的经验" : "Recent Lessons Learned"}
          </h2>
        </div>
        <ul className="space-y-2">
          {lessons.map((lesson, idx) => (
            <li
              key={idx}
              className="flex items-start gap-2 text-gray-700 dark:text-gray-300"
            >
              <span className="text-purple-500 mt-1">•</span>
              <span>{lesson}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Strategy Evolution */}
      {summary.strategy_evolution.length > 0 && (
        <div className="mb-8">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            {locale === "zh" ? "策略演进" : "Strategy Evolution"}
          </h2>
          <div className="space-y-3">
            {summary.strategy_evolution.map((change, idx) => (
              <div
                key={idx}
                className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700"
              >
                <div className="flex items-center gap-2 mb-2">
                  <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                  <span className="text-sm text-gray-600 dark:text-gray-400">
                    {new Date(change.date).toLocaleDateString(
                      locale === "zh" ? "zh-CN" : "en-US",
                      { year: "numeric", month: "short", day: "numeric" }
                    )}
                  </span>
                </div>
                <p className="text-gray-900 dark:text-white">
                  {locale === "zh" ? change.description_zh : change.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Learning History */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          {locale === "zh" ? "验证历史" : "Validation History"}
        </h2>
        {history.length > 0 ? (
          <div className="space-y-4">
            {history.map((record, idx) => (
              <MemoryCard key={idx} record={record} locale={locale} />
            ))}
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 rounded-lg p-8 border border-gray-200 dark:border-gray-700 text-center">
            <p className="text-gray-600 dark:text-gray-400">
              {locale === "zh"
                ? "暂无验证记录"
                : "No validation records yet"}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
