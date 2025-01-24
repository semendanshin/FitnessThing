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
  function inner() {
    try {
      init();
      swipeBehavior.disableVertical();
    } catch (e) {
      console.log("Error while initializing Telegram Mini App SDK");
      console.log(e);
    }
  }

  React.useEffect(inner, []);

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
