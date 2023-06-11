import './css/style.css'

import {Inter} from 'next/font/google'

const inter = Inter({
    subsets: ['latin'],
    variable: '--font-inter',
    display: 'swap'
})

export const metadata = {
    title: 'App Store Slackbot',
    description: 'Talk to the App Store directly from your Slack workspace.',
}

export default function RootLayout({
                                       children,
                                   }: {
    children: React.ReactNode
}) {
    return (
        <html lang="en">
        <head>
            <link rel="preconnect" href="https://fonts.googleapis.com"/>
            <link rel="preconnect" href="https://fonts.gstatic.com"/>
            <link
                href="https://fonts.googleapis.com/css2?family=Sarabun:ital,wght@0,300;0,400;0,500;0,600;1,300;1,400;1,500&display=swap"
                rel="stylesheet"/>
        </head>
        <body
            className={`${inter.variable} font-inter antialiased bg-white text-gray-900 tracking-tight cover-gradient`}>
        <div className="flex flex-col min-h-screen overflow-hidden supports-[overflow:clip]:overflow-clip">
            {children}
        </div>
        </body>
        </html>
    )
}
