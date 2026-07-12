export interface MenuItem {
  label: string;
  path: string;
}

export interface MenuGroup {
  label: string;
  items: MenuItem[];
}

export const menuGroups: MenuGroup[] = [
  {
    label: "Overview",
    items: [
      { label: "Dashboard", path: "/dashboard" },
      { label: "Agenda", path: "/agenda" },
    ],
  },
  {
    label: "System",
    items: [
      { label: "System Info", path: "/system/info" },
      { label: "Command Runner", path: "/system/command" },
      { label: "Scheduler", path: "/system/scheduler" },
    ],
  },
  {
    label: "Settings",
    items: [{ label: "Appearance", path: "/settings/appearance" }],
  },
  {
    label: "Product",
    items: [{ label: "About", path: "/about" }],
  },
];

export function findMenuItem(pathname: string) {
  return menuGroups.flatMap((group) => group.items).find((item) => item.path === pathname);
}
