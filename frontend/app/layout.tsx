import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import clsx from "clsx";

import { Providers } from "./providers";

import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";
import { Navbar } from "@/components/navbar";
import { ToastContainer } from "react-toastify";

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
        <html suppressHydrationWarning lang="en">
            <head />
            <body
                className={clsx(
                    "h-full bg-background font-sans antialiased flex flex-col flex-grow",
                    fontSans.variable
                )}
            >
                <div className="relative flex flex-col min-h-dvh justify-between">
                    <Providers
                        themeProps={{
                            attribute: "class",
                            defaultTheme: "dark",
                        }}
                    >
                        <main className="flex p-0 w-full min-h-full mx-auto flex-grow">
                            {children}
                            <ToastContainer />
                        </main>
                    </Providers>
                    <Navbar />
                </div>
            </body>
        </html>
    );
}
