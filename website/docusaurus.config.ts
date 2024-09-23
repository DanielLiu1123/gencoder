import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Gencoder',
  tagline: 'A code generator for any language or framework that preserves your changes during regeneration, powered by Handlebars.',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://danielliu1123.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/gencoder/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'danielliu1123', // Usually your GitHub org/user name.
  projectName: 'gencoder', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/DanielLiu1123/gencoder/tree/main/website',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    navbar: {
      title: 'Home',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        // {
        //   type: 'localeDropdown',
        //   position: 'right',
        // },
        {
          label: 'GitHub',
          href: 'https://github.com/danielliu1123/gencoder',
          position: 'right',
        },
      ],
    },
    // algolia: {
    //   // The application ID provided by Algolia
    //   appId: 'B5TXVPY7SP',
    //   // Public API key: it is safe to commit it
    //   apiKey: 'd1d0662c9bcd12936152178add34706d',
    //   indexName: 'danielliu1123io',
    //   // Optional: see doc section below
    //   contextualSearch: true,
    //   // Optional: path for search page that enabled by default (`false` to disable it)
    //   searchPagePath: 'search',
    // },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      // additionalLanguages: ['java'],
      // magicComments: [
      //   // Remember to extend the default highlight class name as well!
      //   {
      //     className: 'theme-code-block-highlighted-line',
      //     line: 'highlight-next-line',
      //     block: {start: 'highlight-start', end: 'highlight-end'},
      //   },
      //   // Customized
      //   {
      //     className: 'code-line-deleted',
      //     line: 'highlight-next-line-as-deleted',
      //     block: {start: 'highlight-deleted-start', end: 'highlight-deleted-end'},
      //   },
      //   {
      //     className: 'code-line-added',
      //     line: 'highlight-next-line-as-added',
      //     block: {start: 'highlight-added-start', end: 'highlight-added-end'},
      //   },
      // ],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
