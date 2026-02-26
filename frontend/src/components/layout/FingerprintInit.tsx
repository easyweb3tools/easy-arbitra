"use client";

import { useEffect } from "react";
import { ensureFingerprint } from "@/lib/fingerprint";

export function FingerprintInit() {
  useEffect(() => {
    ensureFingerprint();
  }, []);
  return null;
}
