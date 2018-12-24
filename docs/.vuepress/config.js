module.exports = {
    title: "QStars Network",
    description: "Documentation for the QStars Network.",
    dest: "./dist/docs",
    base: "/qstars/",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        lastUpdated: "Last Updated",
        nav: [{ text: "Back to QStars", link: "https://qstar.network" }],
        sidebar: [
            {
                title: "Introduction",
                collapsable: false,
                children: [
                    "/",
                ]
            },
            {
                title: "Getting Started",
                collapsable: false,
                children: [
                    "/getting-started/installation.md",
                    "/getting-started/full-node.md",
                    "/getting-started/create-testnet.md",
                    "/getting-started/commands.md",
                ]
            },
            {
                title: "SDK Usage",
                collapsable: false,
                children: [
                    "/sdk/kvstore.md",
                    "/sdk/qstarsclicmd.md",
                    "/sdk/qstarsclirestful.md",
                    "/sdk/transactionmessage.md",
                    "/sdk/app-init.md",
                ]
            }
            ,
            {
                title: "library",
                collapsable: false,
                children: [
                    "/library/instruction.md",
                    "/library/introduction.md",
                ]
            }
            ,
            {
                title: "jianqian app",
                collapsable: false,
                children: [
                    "/jianqian/javaniolayer.md",
                ]
            }
            // {
            //     title: "QStars Code",
            //     collapsable: false,
            //     children: [
            //         ["/sdk/overview", "Overview"],
            //         ["/sdk/core/intro", "Core"],
            //         "/sdk/core/app1",
            //         "/sdk/core/app2",
            //         "/sdk/core/app3",
            //         "/sdk/core/app4",
            //         "/sdk/core/app5",
            //         // "/sdk/modules",
            //         "/sdk/clients"
            //     ]
            // },
            // {
            //   title: "Specifications",
            //   collapsable: false,
            //   children: [
            //     ["/specs/overview", "Overview"],
            //     "/specs/governance",
            //     "/specs/ibc",
            //     "/specs/staking",
            //     "/specs/icts",
            //   ]
            // },
            // {
            //     title: "Lotion JS",
            //     collapsable: false,
            //     children: [["/lotion/overview", "Overview"], "/lotion/building-an-app"]
            // },
            // {
            //     title: "Validators",
            //     collapsable: false,
            //     children: [
            //         ["/validators/overview", "Overview"],
            //         ["/validators/security", "Security"],
            //         ["/validators/validator-setup", "Validator Setup"],
            //         "/validators/validator-faq"
            //     ]
            // },
            // {
            //     title: "Resources",
            //     collapsable: false,
            //     children: [
            //         // ["/resources/faq" "General"],
            //         "/resources/delegator-faq",
            //         ["/resources/whitepaper", "Whitepaper - English"],
            //         ["/resources/whitepaper-ko", "Whitepaper - 한국어"],
            //         ["/resources/whitepaper-zh-CN", "Whitepaper - 中文"],
            //         ["/resources/whitepaper-pt", "Whitepaper - Português"]
            //     ]
            // }
        ]
    }
}
