export default function FeaturesBlocks() {
    return (
        <section
            className="app-store-blue relative inset-0 top-1/2 md:mt-24 lg:mt-0 text-white"
            aria-hidden="true">
            <div className="relative max-w-6xl mx-auto px-4 sm:px-6 py-16">
                <div
                    id={"questions"}
                    className="max-w-3xl mx-auto text-center pb-12 md:pb-4 border-b-2 border-white border-opacity-70 pb-4">
                    <h2 className="h2 mb-4">Questions?</h2>
                </div>

                <div className="max-w-8xl mx-auto px-4 pt-16 sm:px-6 lg:px-8">
                    <dl className="space-y-10 md:gap-y-10 md:grid md:grid-cols-2 md:grid-rows-2 md:gap-x-8">
                        <div className="mt-10">
                            <dt className="text-2xl leading-6 font-semibold text-white">How does this bot work?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">
                                <p>
                                    This bot is a thin wrapper over <a className="underline text-white" target="_blank" href="https://github.com/tramlinehq/applelink">Applelink</a>
                                - a stateless App Store API service built for Tramline. App Store Slackbot is open-sourced under the MIT license.
                                </p>
                            </dd>
                        </div>

                        <div>
                            <dt className="text-2xl leading-6 font-semibold text-white">What data do you store?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">
                                <p>
                                    We only store encrypted App Store & Slack authorization credentials. We do not
                                    log or store any requests that you make through Slack. You can wipe all your credentials and data from the portal whenever you like.
                                </p>
                            </dd>
                        </div>
                        <div>
                            <dt className="text-2xl leading-6 font-semibold text-white">What all can this bot do?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">We have added some practical commands that you could use to keep track of your app releases. Run <code>/appstoreslackbot help</code> in your Slack channel to list them all.
                            </dd>
                        </div>
                        <div>
                            <dt className="text-2xl leading-6 font-semibold text-white">Is Google Play Store supported as well?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">Not yet! If this is something that you'd like, please react to <a className="underline text-white" target="_blank" href="https://github.com/tramlinehq/appstore-slackbot/issues/2">this issue on GitHub</a> and we'll prioritize if there's enough interest.
                            </dd>
                        </div>
                        <div>
                            <dt className="text-2xl leading-6 font-semibold text-white">What kind of access does this bot need?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">
                                <p>
                                    It needs to authorize the Slack app with your Slack workspace and it requires configuring App Store Connect API key from your Developer Account.
                                </p>
                            </dd>
                        </div>
                        <div>
                            <dt className="text-2xl leading-6 font-semibold text-white">Have a command that doesn't exist?</dt>
                            <dd className="text-xl mt-2 text-base text-indigo-200">Commands to prepare a new version for release, send it to testers, submit it for review, and add marketing materials to the release are coming soon. Need more? You can file an issue on the <a className="underline text-white" target="_blank" href="https://github.com/tramlinehq/appstore-slackbot">GitHub repository.</a>
                            </dd>
                        </div>
                    </dl>
                </div>
            </div>
        </section>
    )
}
