"use client";
import { link as linkStyles } from "@nextui-org/theme";
import NextLink from "next/link";
import clsx from "clsx";

import { siteConfig } from "@/config/site";

export const Navbar = () => {
  return (
    <div className="h-16 flex items-center justify-between sticky bottom-0 left-0 w-full z-50 shadow-inner bg-background">
      <div className="mx-auto max-w-7xl px-2 flex items-center justify-around w-full z-49">
        {siteConfig.navItems.map((item, id) => (
          <NextLink
            key={id}
            className={clsx(
              linkStyles({ color: "foreground" }),
              "data-[active=true]:text-primary data-[active=true]:font-medium",
            )}
            color="foreground"
            href={item.href}
          >
            <div className="flex flex-col items-center justify-center gap-1">
              {item.icon}
              <p className="text-xs">{item.label}</p>
            </div>
          </NextLink>
          // {id !== siteConfig.navItems.length - 1 && (
          //     <Divider
          //         className="h-6"
          //         orientation="vertical"
          //     />
          // )}
        ))}
      </div>
    </div>
  );
};
