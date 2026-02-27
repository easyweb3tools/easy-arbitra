"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/auth/AuthProvider";

export default function AuthCallbackPage() {
  const router = useRouter();
  const { refresh } = useAuth();

  useEffect(() => {
    refresh().then(() => {
      router.push("/");
    });
  }, [refresh, router]);

  return (
    <div className="flex items-center justify-center py-20">
      <p className="text-subheadline text-label-tertiary">Signing in...</p>
    </div>
  );
}
