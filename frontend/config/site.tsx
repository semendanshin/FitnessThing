import { HomeIcon, ProfileIcon } from "./icons";

export type SiteConfig = typeof siteConfig;

export const siteConfig = {
  name: "Fitness App",
  description: "",
  navItems: [
    {
      label: "Главная",
      href: "/",
      icon: <HomeIcon className="w-4 h-4" />,
    },
    // {
    //   label: "Рутины",
    //   href: "/routines",
    //   icon: <RoutinesIcon className="w-4 h-4" />,
    // },
    {
      label: "Профиль",
      href: "/profile",
      icon: <ProfileIcon className="w-4 h-4" />,
    },
  ],
};
