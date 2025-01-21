"use client";

import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { useIntersectionObserver } from "@siberiacancode/reactuse";
import { Loading } from "./loading";

export default function InfiniteScroll({
    fetchMore,
    offset,
    setOffset,
    hasMore,
    limit,
    showLoading,
    showError,
    showEnd,
    children,
    ...props
}: {
    fetchMore: (offset: number, limit: number) => Promise<void>;
    offset: number;
    setOffset: Dispatch<SetStateAction<number>>;
    hasMore: boolean;
    children: React.ReactNode;
    limit: number;
    showLoading?: boolean;
    showError?: boolean;
    showEnd?: boolean;
} & React.HTMLAttributes<HTMLDivElement>) {
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [isError, setIsError] = useState<boolean>(false);

    async function fetchData() {
        if (!hasMore) return;

        setIsLoading(true);
        try {
            await fetchMore(offset, limit);
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
        rootMargin: "200px",
        onChange: (entry) => {
            if (entry.isIntersecting && !isLoading && hasMore) {
                setOffset((prev) => prev + limit);
                fetchData();
            }
        },
    });

    return (
        <div {...props}>
            {children}
            {showLoading && isLoading && <Loading />}
            {showError && isError && (
                <div className="text-red-500 text-center">
                    Ошибка при загрузке данных
                </div>
            )}
            {showEnd && !hasMore && <div className="text-center">Конец!</div>}
            <div ref={ref} />
        </div>
    );
}

export function useInfiniteScroll() {
    const [offset, setOffset] = useState<number>(0);
    const [hasMore, setHasMore] = useState<boolean>(true);

    return {
        offset,
        setOffset,
        hasMore,
        setHasMore,
    };
}
