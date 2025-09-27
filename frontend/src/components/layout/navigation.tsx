"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

const navigationItems = [
  {
    name: "お知らせ",
    href: "/announcements",
    icon: "📢",
  },
  {
    name: "チャット", 
    href: "/chat",
    icon: "💬",
  },
];

export function Navigation() {
  const pathname = usePathname();

  return (
    <nav className="flex items-center space-x-1">
      {navigationItems.map((item) => {
        const isActive = pathname === item.href;
        
        return (
          <Button
            key={item.href}
            variant={isActive ? "default" : "ghost"}
            asChild
            className={cn(
              "flex items-center gap-2 px-4 py-2",
              isActive && "bg-primary text-primary-foreground"
            )}
          >
            <Link href={item.href}>
              <span className="text-sm">{item.icon}</span>
              <span>{item.name}</span>
            </Link>
          </Button>
        );
      })}
    </nav>
  );
}
