import type { VueTSConfig } from "nuxt/schema"
import vueCssModule from "vite-plugin-vue-css-module"

function fixTsConfig(tsConfig: VueTSConfig) {
  // убрать лишние стандартные алиасы, чтобы остался ровно один ~, который и будем использовать
  for (const p of ["~~", "@@", "@"]) {
    delete tsConfig.compilerOptions!.paths[p]
    delete tsConfig.compilerOptions!.paths[p + "/*"]
  }
}

export default defineNuxtConfig({
  compatibilityDate: "2026-04-14",

  //
  // Build, imports, language
  //

  ssr: false,

  devtools: false,

  imports: {
    dirs: ["stores"],
  },

  typescript: {
    tsConfig: {
      vueCompilerOptions: {
        plugins: ["@vue/language-plugin-pug"],
      },
    },
  },

  vite: {
    plugins: [
      // Почему-то тут ошибка типов без as Plugin, скорее всего из-за разных версий vite.
      // TODO: Разобраться и пофиксить.
      vueCssModule({ attrName: "mclass", pugClassLiterals: true }) as Plugin,
    ],
  },

  css: [
    "~/styles/index.css",
  ],

  hooks: {
    "prepare:types": function ({ tsConfig, nodeTsConfig, sharedTsConfig }) {
      [tsConfig, nodeTsConfig, sharedTsConfig].forEach(fixTsConfig)
    },
  },

  experimental: {
    typescriptPlugin: true,
    serverAppConfig: false,
  },

  //
  // Modules & config
  //

  modules: [
    "@nuxt/icon",
    "@nuxt/ui",
    "nuxt-open-fetch",
  ],

  openFetch: {
    clients: {
      api: {
        schema: "../docs/api/openapi.yaml",
      },
    },
  },

  router: {
    options: {
      sensitive: true,
    },
  },

  //
  // Runtime config
  //

  runtimeConfig: {
    public: {
      // NUXT_PUBLIC_BASE_URL, без слеша на конце (пример: https://foodstock.org)
      baseUrl: "",
    },
  },

  $development: {
    vite: {
      server: {
        // Для отладки колбеков через туннельные сервисы.
        allowedHosts: true,
      },
    },
    runtimeConfig: {
      public: {
        baseUrl: "https://localhost:3151",
      },
    },
  },
})
