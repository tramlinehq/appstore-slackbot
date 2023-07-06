import TramlineLogo from '@/public/images/tramline.png'
import Image from "next/image";

export default function Footer() {
    return (
        <footer className="pt-16 bg-white">
            <div className="max-w-6xl mx-auto px-4 sm:px-6">
                <Image className="opacity-70 hover:opacity-100" src={TramlineLogo} width={100} alt="Features bg"/>

                {/* Bottom area */}
                <div className="md:flex md:items-center md:justify-between pt-2 md:pt-4 pb-6 md:pb-8 text-grey-400 ">
                    {/* Social as */}
                    <div className="flex mb-4 md:order-1 md:ml-4 md:mb-0">
                        <a href="https://twitter.com/tramlinehq"
                           className="flex justify-center items-center text-gray-500 hover:text-gray-900 bg-white hover:bg-white-100 rounded-full shadow transition duration-150 ease-in-out"
                           aria-label="Twitter">
                            <svg className="w-8 h-8 fill-current" viewBox="0 0 32 32"
                                 xmlns="http://www.w3.org/2000/svg">
                                <path
                                    d="M24 11.5c-.6.3-1.2.4-1.9.5.7-.4 1.2-1 1.4-1.8-.6.4-1.3.6-2.1.8-.6-.6-1.5-1-2.4-1-1.7 0-3.2 1.5-3.2 3.3 0 .3 0 .5.1.7-2.7-.1-5.2-1.4-6.8-3.4-.3.5-.4 1-.4 1.7 0 1.1.6 2.1 1.5 2.7-.5 0-1-.2-1.5-.4 0 1.6 1.1 2.9 2.6 3.2-.3.1-.6.1-.9.1-.2 0-.4 0-.6-.1.4 1.3 1.6 2.3 3.1 2.3-1.1.9-2.5 1.4-4.1 1.4H8c1.5.9 3.2 1.5 5 1.5 6 0 9.3-5 9.3-9.3v-.4c.7-.5 1.3-1.1 1.7-1.8z"/>
                            </svg>
                        </a>
                        <a href="https://github.com/tramlinehq/appstore-slackbot"
                           className="ml-4 flex justify-center items-center text-gray-500 hover:text-gray-900 bg-white hover:bg-white-100 rounded-full shadow transition duration-150 ease-in-out"
                           aria-label="Github">
                            <svg className="w-8 h-8 fill-current" viewBox="0 0 32 32"
                                 xmlns="http://www.w3.org/2000/svg">
                                <path
                                    d="M16 8.2c-4.4 0-8 3.6-8 8 0 3.5 2.3 6.5 5.5 7.6.4.1.5-.2.5-.4V22c-2.2.5-2.7-1-2.7-1-.4-.9-.9-1.2-.9-1.2-.7-.5.1-.5.1-.5.8.1 1.2.8 1.2.8.7 1.3 1.9.9 2.3.7.1-.5.3-.9.5-1.1-1.8-.2-3.6-.9-3.6-4 0-.9.3-1.6.8-2.1-.1-.2-.4-1 .1-2.1 0 0 .7-.2 2.2.8.6-.2 1.3-.3 2-.3s1.4.1 2 .3c1.5-1 2.2-.8 2.2-.8.4 1.1.2 1.9.1 2.1.5.6.8 1.3.8 2.1 0 3.1-1.9 3.7-3.7 3.9.3.4.6.9.6 1.6v2.2c0 .2.1.5.6.4 3.2-1.1 5.5-4.1 5.5-7.6-.1-4.4-3.7-8-8.1-8z"/>
                            </svg>
                        </a>
                    </div>
                    {/* Copyrights note */}
                    <div className="text-medium text-gray-600 mr-4">
                        <div>Built and powered by <a href="https://tramline.app" className="underline">Tramline Inc.</a></div>
                        <div>Since 2023.</div>
                    </div>
                </div>
            </div>
        </footer>
    )
}
