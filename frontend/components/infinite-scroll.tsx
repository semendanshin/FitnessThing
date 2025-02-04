"use client";

import { useState } from "react";
import { useIntersectionObserver } from "@siberiacancode/reactuse";
import clsx from "clsx";

import { Loading } from "./loading";

export default function InfiniteScroll({
  fetchMore,
  hasMore,
  showLoading,
  showError,
  showEnd,
  children,
  className,
}: {
  fetchMore: () => Promise<void>;
  hasMore: boolean;
  showLoading?: boolean;
  showError?: boolean;
  showEnd?: boolean;
  children: React.ReactNode;
  className?: string;
} & React.HTMLAttributes<HTMLDivElement>) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isError, setIsError] = useState<boolean>(false);

  async function fetchData() {
    if (!hasMore) return;

    setIsLoading(true);
    try {
      await fetchMore();
      setIsError(false);
    } catch (error) {
      console.log(error);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  }

  const { ref } = useIntersectionObserver<HTMLDivElement>({
    threshold: 0,
    rootMargin: "0px",
    onChange: (entry) => {
      if (entry.isIntersecting && !isLoading && hasMore) {
        fetchData();
      }
    },
  });

  return (
    <>
      <div className={clsx("relative", className)}>
        {children}
        {showLoading && isLoading && <Loading />}
        {showError && isError && (
          <div className="text-red-500 text-center">
            Ошибка при загрузке данных
          </div>
        )}
        {showEnd && !hasMore && <div className="text-center">Конец!</div>}
        <div ref={ref} className="bottom-0 h-[100px] absolute z-[-1]" />
      </div>
    </>
  );
}

export function useInfiniteScroll() {
  const [hasMore, setHasMore] = useState<boolean>(true);

  return {
    hasMore,
    setHasMore,
  };
}
