export const metadata = {
    title: 'App Store Slackbot',
    description: 'Talk to the App Store directly from your Slack workspace.',
}

import Hero from '@/components/hero'
import Features from '@/components/features'
import FeaturesBlocks from '@/components/features-blocks'

export default function Home() {
    return (
        <>
            <Hero/>
            <div className="cover"><Features/></div>
            <div className="cover"><FeaturesBlocks/></div>
        </>
    )
}
