import type { Metadata } from 'next'
import React from 'react'

export const metadata: Metadata = {
  title: '{{.ProjectName}}',
  description: 'Scaffolded by Open Workbench CLI',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}


