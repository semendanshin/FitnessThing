"use client";

import type { ThemeProviderProps } from "next-themes";

import * as React from "react";
import { NextUIProvider } from "@nextui-org/system";
import { useRouter } from "next/navigation";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import { init, swipeBehavior } from "@telegram-apps/sdk-react";

export interface ProvidersProps {
  children: React.ReactNode;
  themeProps?: ThemeProviderProps;
}

declare module "@react-types/shared" {
  interface RouterConfig {
    routerOptions: NonNullable<
      Parameters<ReturnType<typeof useRouter>["push"]>[1]
    >;
  }
}

function TelegramMiniAppProvider({ children }: ProvidersProps) {
  init();

  if (swipeBehavior.disableVertical.isAvailable()) {
    swipeBehavior.disableVertical();
  } else {
    console.warn("Swipe behavior is not available");
  }

  return <>{children}</>;
}

export function Providers({ children, themeProps }: ProvidersProps) {
  const router = useRouter();

  return (
    <NextUIProvider
      className="flex flex-col flex-grow h-full"
      navigate={router.push}
    >
      <NextThemesProvider {...themeProps}>
        <TelegramMiniAppProvider>{children}</TelegramMiniAppProvider>
      </NextThemesProvider>
    </NextUIProvider>
  );
}
