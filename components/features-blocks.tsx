export default function FeaturesBlocks() {
    return (
        <section
            className="app-store-blue relative inset-0 top-1/2 md:mt-24 lg:mt-0 pointer-events-none text-white"
            aria-hidden="true">
            <div className="relative max-w-6xl mx-auto px-4 sm:px-6 py-16">
                <div
                    id={"questions"}
                    className="max-w-3xl mx-auto text-center pb-12 md:pb-4 border-b-2 border-white border-opacity-70 pb-4">
                    <h2 className="h2 mb-4">Questions?</h2>
                </div>

                <div className="max-w-7xl mx-auto px-4 pt-16 sm:px-6 lg:px-8">
                    <dl className="space-y-10 md:space-y-10 md:grid md:grid-cols-2 md:grid-rows-2 md:gap-x-8">
                        <div>
                            <dt className="text-lg leading-6 font-medium text-white">What kind of data do you store?
                            </dt>
                            <dd className="mt-2 text-base text-indigo-200">
                                <p>
                                    We only store encrypted App Store & Slack authorization credentials. We do not
                                    store any requests that you make
                                </p>
                                <p className="mt-4">Embeds that fail to
                                    send are logged for
                                    24 hours to assist with debugging - this data is only used for fixing bugs and
                                    is
                                    not otherwise
                                    viewed.
                                </p>
                            </dd>
                        </div>
                        <div>
                            <dt className="text-lg leading-6 font-medium text-white">Want a native integration?</dt>
                            <dd className="mt-2 text-base text-indigo-200">Me too! There's an <a
                                className="underline text-white"
                                href="https://github.com/getsentry/sentry/issues/10925">open
                                issue on GitHub</a> that you can go and leave reactions on to help get it
                                prioritized. If official
                                support lands, this service will likely stop allowing new registrations but will
                                remain up so long
                                as webhooks are receiving events.
                            </dd>
                        </div>
                        <div>
                            <dt className="text-lg leading-6 font-medium text-white">Why doesn't the embed include
                                [thing]?
                            </dt>
                            <dd className="mt-2 text-base text-indigo-200">I've tried to add what I view as useful
                                information,
                                but if you think I've missed something please open an issue on GitHub!
                            </dd>
                        </div>
                        <div>
                            <dt className="text-lg leading-6 font-medium text-white">Have a feature request or want
                                to report a
                                bug?
                            </dt>
                            <dd className="mt-2 text-base text-indigo-200">Awesome! You can file an issue on the <a
                                className="underline text-white"
                                href="https://github.com/ianmitchell/sentrydiscord.dev">GitHub
                                repository</a></dd>
                        </div>
                    </dl>
                </div>
            </div>
        </section>
    )
}
