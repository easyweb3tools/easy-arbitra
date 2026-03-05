import { getNovaStatus } from "@/lib/api";
import { BrainDashboard } from "@/components/nova/BrainDashboard";
import { Brain } from "lucide-react";

export const metadata = {
  title: "Nova AI Brain | Easy Arbitra",
  description: "Watch Nova's real-time thinking process",
};

export default async function NovaBrainPage() {
  let status;
  try {
    status = await getNovaStatus();
  } catch (err) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <Brain className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            Nova Brain Status Unavailable
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Unable to fetch Nova's current status. Please try again later.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Nova AI Brain
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Real-time view of Nova's analytical process
        </p>
      </div>

      <BrainDashboard status={status} />

      <div className="mt-8 bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          How Nova Works
        </h2>
        <div className="grid md:grid-cols-4 gap-6">
          <div>
            <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center mb-3">
              <span className="text-2xl font-bold text-blue-600 dark:text-blue-400">1</span>
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
              Data Collection
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Nova gathers trading data from top-performing wallets every hour
            </p>
          </div>

          <div>
            <div className="w-12 h-12 bg-purple-100 dark:bg-purple-900 rounded-lg flex items-center justify-center mb-3">
              <span className="text-2xl font-bold text-purple-600 dark:text-purple-400">2</span>
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
              Analysis & Evaluation
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Evaluates candidates based on win rate, stability, and risk control
            </p>
          </div>

          <div>
            <div className="w-12 h-12 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center mb-3">
              <span className="text-2xl font-bold text-green-600 dark:text-green-400">3</span>
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
              Memory Accumulation
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Builds context across multiple rounds to make informed decisions
            </p>
          </div>

          <div>
            <div className="w-12 h-12 bg-orange-100 dark:bg-orange-900 rounded-lg flex items-center justify-center mb-3">
              <span className="text-2xl font-bold text-orange-600 dark:text-orange-400">4</span>
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
              Decision Output
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Makes final recommendation when confidence threshold is reached
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
