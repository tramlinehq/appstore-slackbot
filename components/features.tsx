'use client'

import { useState, useRef, useEffect } from 'react'
import { Transition } from '@headlessui/react'
import Image from 'next/image'
import Help from '@/public/images/help.png'
import ControlRelease from '@/public/images/control_release.png'
import StoreInfo from '@/public/images/store_info.png'

export default function Features() {

    const [tab, setTab] = useState<number>(1)

    const tabs = useRef<HTMLDivElement>(null)

    return (
        <section className="relative inset-0 slack-purple">
            {/* Section background (needs .relative class on parent and next sibling elements) */}
            <div className="relative max-w-6xl mx-auto px-4 sm:px-6 py-16">
                {/* Section header */}
                <div
                    className="max-w-3xl mx-auto text-center pb-12 md:pb-4 border-b-2 border-white border-opacity-70 pb-4 text-white">
                    <h2 className="h2 mb-4">Examples</h2>
                </div>
                {/* Section content */}
                <div className="md:grid md:grid-cols-12 md:gap-6 px-4 sm:px-6 lg:px-8 space-y-10 md:space-y-10 py-16">
                    {/* Content */}
                    <div className="max-w-xl md:max-w-none md:w-full mx-auto md:col-span-5 lg:col-span-5 md:mt-6"
                        data-aos="fade-right">
                        <div className="mb-8 text-white">
                            <p className="text-xl"><a className="underline plausible-event-name=Get+started+Login" href="https://service.appstoreslackbot.com">Login</a> to connect your Slack workspace and add an App Store Connect API key to operate on your app right from your Slack channel.</p>
                        </div>
                        {/* Tabs buttons */}
                        <div className="mb-8 md:mb-0">
                            <a
                                className={`flex items-center text-lg p-5 bg-[#25262B] shadow-xs rounded border transition duration-300 ease-in-out mb-3 ${tab !== 1 ? 'hover:opacity-95 hover:shadow-lg border-transparent' : 'border-gray-200'}`}
                                href="#0"
                                onClick={(e) => {
                                    e.preventDefault();
                                    setTab(1);
                                }}
                            >
                                <div class="text-white leading-snug tracking-tight">
                                    <div className="font-bold leading-snug mb-1">Practical commands
                                    </div>
                                    <div>Commands to check current review status, list test groups, see what build is distributed where, and so much more
                                    </div>
                                </div>
                                <div
                                    className="flex justify-center items-center w-8 h-8 bg-white rounded-full shadow flex-shrink-0 ml-3">
                                    <svg className="w-3 h-3 fill-current" viewBox="0 0 12 12"
                                        xmlns="http://www.w3.org/2000/svg">
                                        <path
                                            d="M11.953 4.29a.5.5 0 00-.454-.292H6.14L6.984.62A.5.5 0 006.12.173l-6 7a.5.5 0 00.379.825h5.359l-.844 3.38a.5.5 0 00.864.445l6-7a.5.5 0 00.075-.534z" />
                                    </svg>
                                </div>
                            </a>
                            <a
                                className={`flex items-center text-lg p-5 bg-[#25262B] shadow-xs rounded border transition duration-300 ease-in-out mb-3 ${tab !== 2 ? 'hover:opacity-95 hover:shadow-lg border-transparent' : 'border-gray-200'}`}
                                href="#0"
                                onClick={(e) => {
                                    e.preventDefault();
                                    setTab(2);
                                }}
                            >
                                <div class="text-white leading-snug tracking-tight">
                                    <div className="font-bold leading-snug mb-1">Control your current release
                                    </div>
                                    <div>View your live release status, pause the phased release, resume it or release to all users
                                    </div>
                                </div>
                                <div
                                    className="flex justify-center items-center w-8 h-8 bg-white rounded-full shadow flex-shrink-0 ml-3">
                                    <svg className="w-3 h-3 fill-current" viewBox="0 0 12 12"
                                        xmlns="http://www.w3.org/2000/svg">
                                        <path
                                            d="M11.854.146a.5.5 0 00-.525-.116l-11 4a.5.5 0 00-.015.934l4.8 1.921 1.921 4.8A.5.5 0 007.5 12h.008a.5.5 0 00.462-.329l4-11a.5.5 0 00-.116-.525z"
                                            fillRule="nonzero" />
                                    </svg>
                                </div>
                            </a>
                            <a
                                className={`flex items-center text-lg p-5 bg-[#25262B] shadow-xs rounded border transition duration-300 ease-in-out mb-3 ${tab !== 3 ? 'hover:opacity-95 hover:shadow-lg border-transparent' : 'border-gray-200'}`}
                                href="#0"
                                onClick={(e) => {
                                    e.preventDefault();
                                    setTab(3);
                                }}
                            >
                                <div class="text-white leading-snug tracking-tight">
                                    <div className="font-bold leading-snug mb-1">Current status
                                    </div>
                                    <div>Check the overall status of your app across all beta channels and production
                                    </div>
                                </div>
                                <div
                                    className="flex justify-center items-center w-8 h-8 bg-white rounded-full shadow flex-shrink-0 ml-3">
                                    <svg className="w-3 h-3 fill-current" viewBox="0 0 12 12"
                                        xmlns="http://www.w3.org/2000/svg">
                                        <path
                                            d="M11.334 8.06a.5.5 0 00-.421-.237 6.023 6.023 0 01-5.905-6c0-.41.042-.82.125-1.221a.5.5 0 00-.614-.586 6 6 0 106.832 8.529.5.5 0 00-.017-.485z"
                                            fill="#191919" fillRule="nonzero" />
                                    </svg>
                                </div>
                            </a>
                        </div>
                    </div>

                    {/* Tabs items */}
                    <div
                        className="max-w-xl md:max-w-none md:w-full mx-auto md:col-span-7 lg:col-span-7 mb-8 md:mb-0 md:order-1"
                        data-aos="zoom-y-out" ref={tabs}>
                        <div className="relative flex flex-col text-center lg:text-right">
                            {/* Item 1 */}
                            <Transition
                                show={tab === 1}
                                appear={true}
                                className="w-full"
                                enter="transition ease-in-out duration-700 transform order-first"
                                enterFrom="opacity-0 translate-y-16"
                                enterTo="opacity-100 translate-y-0"
                                leave="transition ease-in-out duration-300 transform absolute"
                                leaveFrom="opacity-100 translate-y-0"
                                leaveTo="opacity-0 -translate-y-16"
                            >
                                <div className="relative inline-flex flex-col">
                                    <Image className="md:max-w-none mx-auto rounded transform w-full" src={Help} alt="Usage guide" />
                                </div>
                            </Transition>
                            {/* Item 2 */}
                            <Transition
                                show={tab === 2}
                                appear={true}
                                className="w-full"
                                enter="transition ease-in-out duration-700 transform order-first"
                                enterFrom="opacity-0 translate-y-16"
                                enterTo="opacity-100 translate-y-0"
                                leave="transition ease-in-out duration-300 transform absolute"
                                leaveFrom="opacity-100 translate-y-0"
                                leaveTo="opacity-0 -translate-y-16"
                            >
                                <div className="relative inline-flex flex-col">
                                    <Image className="md:max-w-none mx-auto rounded transform w-full" src={ControlRelease} alt="Control Release" />
                                </div>
                            </Transition>
                            {/* Item 3 */}
                            <Transition
                                show={tab === 3}
                                appear={true}
                                className="w-full"
                                enter="transition ease-in-out duration-700 transform order-first"
                                enterFrom="opacity-0 translate-y-16"
                                enterTo="opacity-100 translate-y-0"
                                leave="transition ease-in-out duration-300 transform absolute"
                                leaveFrom="opacity-100 translate-y-0"
                                leaveTo="opacity-0 -translate-y-16"
                            >
                                <div className="relative inline-flex flex-col">
                                    <Image className="md:max-w-none mx-auto rounded transform w-full" src={StoreInfo} alt="Store Info" />
                                </div>
                            </Transition>
                        </div>
                    </div>

                </div>
            </div>
        </section>
    )
}