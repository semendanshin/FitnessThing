import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import clsx from "clsx";
import { ToastContainer } from "react-toastify";

import { Providers } from "./providers";

import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";
import { Navbar } from "@/components/navbar";

export const metadata: Metadata = {
  title: {
    default: siteConfig.name,
    template: `%s - ${siteConfig.name}`,
  },
  description: siteConfig.description,
  icons: {
    icon: "/favicon.ico",
  },
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "white" },
    { media: "(prefers-color-scheme: dark)", color: "black" },
  ],
  userScalable: false,
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html suppressHydrationWarning className="overflow-hidden" lang="ru">
      <head />
      <body
        className={clsx(
          "h-full w-full overflow-hidden font-sans antialiased ",
          fontSans.variable,
        )}
      >
        <div className="flex flex-col justify-between h-dvh overflow-y-scroll max-w-full">
          <Providers
            themeProps={{
              attribute: "class",
              defaultTheme: "dark",
              enableSystem: true,
            }}
          >
            <main className="flex mx-auto flex-grow overflow-y-auto mb-[4rem] w-full h-full">
              <div className="flex flex-grow max-h-full flex-col w-full">
                {children}
              </div>
              <ToastContainer />
              <Navbar />
            </main>
          </Providers>
        </div>
      </body>
    </html>
  );
}
