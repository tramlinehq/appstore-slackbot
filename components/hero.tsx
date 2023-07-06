export default function Hero() {
    return (
        <section className="relative bg-white top-noise wave">
            <div className="max-w-6xl mx-auto px-4 sm:px-6">
                {/* Hero content */}
                <div className="pt-32 pb-12 md:pt-40 md:pb-20">
                    {/* Section header */}
                    <div className="text-center pb-12 md:pb-16">
                        <span className="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-pink-100 text-pink-800 mb-4">Beta</span>
                        <h1 className="text-2xl md:text-6xl font-extrabold leading-tighter tracking-tighter drop-shadow-sm" data-aos="zoom-y-out">
                            <span className="text-[#0043b6]">App Store</span>
                            <span className="pl-4 pr-4 bg-clip-text text-transparent bg-gradient-to-r from-[#0043b6] to-[#4A154B]">→</span>
                            <span className="text-[#4A154B]">Slackbot</span>
                        </h1>
                        <div className="max-w-3xl space-y-5 mx-auto">
                            <p className="max-w-md mt-5 mx-auto text-base text-black sm:text-lg md:my-8 md:text-2xl md:max-w-3xl drop-shadow-sm" data-aos="zoom-y-out" data-aos-delay="50">
                                Talk to the App Store directly from your Slack workspace
                            </p>
                            <div className="max-w-xs mx-auto sm:max-w-none sm:flex sm:justify-center"
                                 data-aos="zoom-y-out" data-aos-delay="100">
                                <div className="drop-shadow-lg">
                                    <a className="btn text-white bg-[#4A154B] opacity-90 hover:opacity-100 w-full mb-4 sm:w-auto sm:mb-0 font-bold"
                                       href="https://service.appstoreslackbot.com/">Get Started</a>
                                </div>
                                <div className="drop-shadow-lg">
                                    <a className="btn text-white bg-gray-900 hover:bg-gray-800 w-full sm:w-auto sm:ml-4 font-bold"
                                       href="#questions">Questions?</a>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    )
}
