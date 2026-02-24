export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={`skeleton rounded-sm ${className}`} />;
}

export function SkeletonCard() {
  return (
    <div className="rounded-lg bg-surface-secondary p-6 shadow-elevation-1">
      <Skeleton className="mb-3 h-5 w-32" />
      <Skeleton className="mb-2 h-4 w-full" />
      <Skeleton className="mb-2 h-4 w-3/4" />
      <Skeleton className="h-4 w-1/2" />
    </div>
  );
}

export function SkeletonRow() {
  return (
    <div className="flex items-center gap-3 border-b border-separator px-4 py-3">
      <Skeleton className="h-10 w-10 !rounded-full" />
      <div className="flex-1">
        <Skeleton className="mb-2 h-4 w-40" />
        <Skeleton className="h-3 w-24" />
      </div>
    </div>
  );
}

export function SkeletonStatGrid() {
  return (
    <div className="grid gap-4 sm:grid-cols-3">
      {[0, 1, 2].map((i) => (
        <div key={i} className="rounded-lg bg-surface-secondary p-5 shadow-elevation-1">
          <Skeleton className="mb-2 h-3 w-20" />
          <Skeleton className="h-7 w-16" />
        </div>
      ))}
    </div>
  );
}
