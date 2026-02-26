"use client";

import { useRouter, usePathname, useSearchParams } from "next/navigation";
import { useCallback } from "react";

type Option = {
  label: string;
  value: string;
};

export function SortToggle({
  options,
  paramName = "sort_by",
  defaultValue,
}: {
  options: Option[];
  paramName?: string;
  defaultValue?: string;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const current = searchParams.get(paramName) || defaultValue || options[0]?.value;

  const handleClick = useCallback(
    (value: string) => {
      const params = new URLSearchParams(searchParams.toString());
      params.set(paramName, value);
      params.set("page", "1");
      router.push(`${pathname}?${params.toString()}`);
    },
    [router, pathname, searchParams, paramName],
  );

  return (
    <div className="flex gap-1.5">
      {options.map((opt) => (
        <button
          key={opt.value}
          onClick={() => handleClick(opt.value)}
          className={[
            "rounded-full px-4 py-1.5 text-caption-1 font-semibold transition-all duration-200 ease-apple",
            current === opt.value
              ? "bg-label-primary text-surface-primary shadow-sm"
              : "bg-surface-tertiary/80 text-label-secondary hover:bg-surface-tertiary",
          ].join(" ")}
        >
          {opt.label}
        </button>
      ))}
    </div>
  );
}
