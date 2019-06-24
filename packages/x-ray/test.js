const axios = require('axios')
const capture = require('./lib').captureHttps
capture()

const https = require('https')

axios.get('https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY').then()
axios.post('https://enj66527mhry9.x.pipedream.net/').then()
// https.get('https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY', (resp) => {
//   let data = '';

//   // A chunk of data has been recieved.
//   resp.on('data', (chunk) => {
//     data += chunk;
//   });

//   // The whole response has been received. Print out the result.
//   resp.on('end', () => {
//     console.log(JSON.parse(data).explanation);
//   });

// }).on("error", (err) => {
//   console.log("Error: " + err.message);
// });

// // const options = {
// //   hostname: 'encrypted.google.com',
// //   // port: 443,
// //   path: '/',
// //   method: 'GET'
// // };

// // https.request(options, (resp) => {
// //   // let data = '';

// //   // // A chunk of data has been recieved.
// //   // resp.on('data', (chunk) => {
// //   //   data += chunk;
// //   // });

// //   // // The whole response has been received. Print out the result.
// //   // resp.on('end', () => {
// //   //   console.log(JSON.parse(data).explanation);
// //   // });

// // }).on("error", (err) => {
// //   console.log("Error: " + err.message);
// // }).end()
