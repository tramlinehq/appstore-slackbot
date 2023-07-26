import './css/style.css'

import { Inter } from 'next/font/google'

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap'
})
export const metadata = {
    title: 'App Store → Slackbot',
    description: 'Talk to the App Store directly from your Slack workspace.',
}

export default function RootLayout({ children, }: {
    children: React.ReactNode
}) {
    return (
        <html lang="en">
            <head>
                <meta property="og:title" content="App Store Slackbot" />
                <meta property="og:description" content="Talk to the App Store directly from your Slack workspace" />
                <meta property="og:image" content="https://storage.googleapis.com/tramline-public-assets/appstoreslackbot.png" />
                <meta name="twitter:title" content="App Store Slacbot" />
                <meta name="twitter:description" content="Talk to the App Store directly from your Slack workspace" />
                <meta name="twitter:image" content="https://storage.googleapis.com/tramline-public-assets/appstoreslackbot.png" />
                <link rel="preconnect" href="https://stijndv.com" />
                <link rel="stylesheet" href="https://stijndv.com/fonts/Eudoxus-Sans.css" />
                <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
                <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
                <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
                <link rel="manifest" href="/site.webmanifest" />
                <script defer data-domain="appstoreslackbot.com" src="https://plausible.io/js/script.tagged-events.js"></script>
            </head>
            <body
                className={`${inter.className} font-inter antialiased bg-white text-gray-900 cover-gradient`}>
                <div className="flex flex-col min-h-screen overflow-hidden supports-[overflow:clip]:overflow-clip">
                    {children}
                </div>
            </body>
        </html>
    )
}
