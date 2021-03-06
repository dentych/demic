export default {
    server: {
        host: "0.0.0.0"
    },
    // head: {
    //     meta: [
    //         { charset: 'utf-8' },
    //         { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=1.0' },
    //
    //         // hid is used as unique identifier. Do not use `vmid` for it as it will not work
    //         { hid: 'description', name: 'description', content: 'Meta description' }
    //     ]
    // },
    buildModules: [
        "@nuxtjs/tailwindcss",
        "@nuxtjs/pwa",
    ],
    env: {
        apiBaseUrl: process.env.API_BASE_URL || "http://localhost:8080"
    }
}