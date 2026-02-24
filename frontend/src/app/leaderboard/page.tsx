import { redirect } from "next/navigation";

export default function LeaderboardPage() {
  redirect("/wallets?sort_by=smart_score&order=desc");
}
