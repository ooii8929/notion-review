const { Client } = require('@notionhq/client');
require('dotenv').config()

const notion = new Client({ auth: process.env.NOTION_API_KEY });
console.log(process.env.NOTION_API_KEY);
(async () => {
  const databaseId = '3b88a7c6e2834561a304043557a47aee';
  const response = await notion.databases.query({
    database_id: databaseId,
    filter: {
      or: [
        {
          property: 'In stock',
          checkbox: {
            equals: true,
          },
        },
        {
          property: 'Cost of next trip',
          number: {
            greater_than_or_equal_to: 2,
          },
        },
      ],
    },
    sorts: [
      {
        property: 'Last ordered',
        direction: 'ascending',
      },
    ],
  });
  console.log(response);
})();

// const notion = new Client({ auth: process.env.NOTION_API_KEY });

// console.log(process.env.NOTION_ACCESS_TOKEN);
// (async () => {
//   const response = await notion.search({
//     query: 'External tasks',
//     filter: {
//       value: 'database',
//       property: 'object'
//     },
//     sort: {
//       direction: 'ascending',
//       timestamp: 'last_edited_time'
//     },
//   });
//   console.log(response);
// })();