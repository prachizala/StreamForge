import "./globals.css";
import React from "react";

export const metadata = {
  title: "StreamForge Studio — High-Performance Media Pipeline",
  description: "Real-time media processing engine built with Next.js, Go, and Rust.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-[#0B0F17] text-gray-100 antialiased">
        {children}
      </body>
    </html>
  );
}
